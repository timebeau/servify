package server

import (
	"context"
	"fmt"
	"time"

	"servify/apps/server/internal/config"
	aidelivery "servify/apps/server/internal/modules/ai/delivery"
	"servify/apps/server/internal/platform/configscope"
	"servify/apps/server/internal/platform/knowledgeprovider"
	difykp "servify/apps/server/internal/platform/knowledgeprovider/dify"
	weknorakp "servify/apps/server/internal/platform/knowledgeprovider/weknora"
	"servify/apps/server/internal/platform/llm/openai"
	"servify/apps/server/internal/services"
	"servify/apps/server/pkg/dify"
	"servify/apps/server/pkg/weknora"

	"github.com/sirupsen/logrus"
)

type AIAssemblyOptions struct {
	RequireWeKnoraHealthy bool
	SyncKnowledgeBase     bool
	HealthCheckTimeout    time.Duration
}

type AIAssembly struct {
	Service         aidelivery.HandlerService
	RuntimeService  aidelivery.RuntimeService
	KnowledgeDriver knowledgeprovider.KnowledgeProvider
	KnowledgeProviderID string
	DifyHealthy     bool
	DifyDatasetID   string
	WeKnoraClient   weknora.WeKnoraInterface
	WeKnoraHealthy  bool
	KnowledgeBaseID string
}

func (a *AIAssembly) KnowledgeProvider(cfg *config.Config) knowledgeprovider.KnowledgeProvider {
	if a == nil {
		return nil
	}
	return a.KnowledgeDriver
}

func BuildAIAssembly(cfg *config.Config, logger *logrus.Logger, opts AIAssemblyOptions) (*AIAssembly, error) {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	resolver := configscope.NewResolver(cfg)
	openAIConfig := resolver.ResolveOpenAI(context.Background(), nil)
	weKnoraConfig := resolver.ResolveWeKnora(context.Background(), nil)

	baseAI := services.NewAIService(openAIConfig.APIKey, openAIConfig.BaseURL)
	baseAI.InitializeKnowledgeBase()
	defaultService := services.NewOrchestratedEnhancedAIService(
		baseAI,
		openai.NewProvider(openAIConfig.APIKey, openAIConfig.BaseURL),
		nil,
		"",
		nil,
		"",
		logger,
	)

	assembly := &AIAssembly{
		Service:         aidelivery.NewHandlerServiceAdapter(defaultService),
		RuntimeService:  defaultService,
		KnowledgeBaseID: weKnoraConfig.KnowledgeBaseID,
	}

	if cfg.Dify.Enabled {
		difyClient := dify.NewClient(&dify.Config{
			BaseURL: cfg.Dify.BaseURL,
			APIKey:  cfg.Dify.APIKey,
			Timeout: cfg.Dify.Timeout,
		})
		ctx, cancel := context.WithTimeout(context.Background(), timeoutForHealthCheck(opts))
		defer cancel()
		if err := difyClient.HealthCheck(ctx, cfg.Dify.DatasetID); err != nil {
			logger.Warnf("Dify health check failed: %v", err)
			if !weKnoraConfig.Enabled && opts.RequireWeKnoraHealthy {
				return nil, fmt.Errorf("dify health check failed: %w", err)
			}
		} else {
			assembly.DifyHealthy = true
			assembly.DifyDatasetID = cfg.Dify.DatasetID
			assembly.KnowledgeProviderID = "dify"
			assembly.KnowledgeDriver = difykp.NewProvider(difyClient, cfg.Dify.DatasetID, difykp.SearchConfig{
				TopK:            cfg.Dify.Search.TopK,
				ScoreThreshold:  cfg.Dify.Search.ScoreThreshold,
				SearchMethod:    cfg.Dify.Search.SearchMethod,
				RerankingEnable: cfg.Dify.Search.RerankingEnable,
			})
			enhanced := services.NewOrchestratedEnhancedAIService(
				baseAI,
				openai.NewProvider(openAIConfig.APIKey, openAIConfig.BaseURL),
				assembly.KnowledgeDriver,
				"dify",
				nil,
				cfg.Dify.DatasetID,
				logger,
			)
			assembly.Service = aidelivery.NewHandlerServiceAdapter(enhanced)
			assembly.RuntimeService = enhanced
			return assembly, nil
		}
	}

	if !weKnoraConfig.Enabled {
		return assembly, nil
	}

	client := weknora.NewClient(&weknora.Config{
		BaseURL:    weKnoraConfig.BaseURL,
		APIKey:     weKnoraConfig.APIKey,
		TenantID:   weKnoraConfig.TenantID,
		Timeout:    weKnoraConfig.Timeout,
		MaxRetries: weKnoraConfig.MaxRetries,
	}, logger)
	assembly.WeKnoraClient = client

	timeout := timeoutForHealthCheck(opts)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := client.HealthCheck(ctx); err != nil {
		logger.Warnf("WeKnora health check failed: %v", err)
		if opts.RequireWeKnoraHealthy {
			return nil, fmt.Errorf("weknora health check failed: %w", err)
		}
		if !cfg.Fallback.Enabled {
			return nil, fmt.Errorf("weknora unavailable and fallback disabled: %w", err)
		}
		return assembly, nil
	}
	assembly.WeKnoraHealthy = true
	assembly.KnowledgeProviderID = "weknora"
	assembly.KnowledgeDriver = weknorakp.NewProvider(client, weKnoraConfig.KnowledgeBaseID)

	enhanced := services.NewOrchestratedEnhancedAIService(
		baseAI,
		openai.NewProvider(openAIConfig.APIKey, openAIConfig.BaseURL),
		assembly.KnowledgeDriver,
		"weknora",
		client,
		weKnoraConfig.KnowledgeBaseID,
		logger,
	)
	if opts.SyncKnowledgeBase {
		syncCtx, syncCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer syncCancel()
		if err := enhanced.SyncKnowledgeBase(syncCtx); err != nil {
			logger.Warnf("Knowledge base sync failed: %v", err)
		}
	}
	assembly.Service = aidelivery.NewHandlerServiceAdapter(enhanced)
	assembly.RuntimeService = enhanced
	return assembly, nil
}

func timeoutForHealthCheck(opts AIAssemblyOptions) time.Duration {
	if opts.HealthCheckTimeout > 0 {
		return opts.HealthCheckTimeout
	}
	return 10 * time.Second
}
