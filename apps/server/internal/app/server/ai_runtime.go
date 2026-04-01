package server

import (
	"context"
	"fmt"
	"time"

	"servify/apps/server/internal/config"
	aidelivery "servify/apps/server/internal/modules/ai/delivery"
	"servify/apps/server/internal/platform/configscope"
	"servify/apps/server/internal/platform/knowledgeprovider"
	weknorakp "servify/apps/server/internal/platform/knowledgeprovider/weknora"
	"servify/apps/server/internal/platform/llm/openai"
	"servify/apps/server/internal/services"
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
	WeKnoraClient   weknora.WeKnoraInterface
	WeKnoraHealthy  bool
	KnowledgeBaseID string
}

func (a *AIAssembly) KnowledgeProvider(cfg *config.Config) knowledgeprovider.KnowledgeProvider {
	if a == nil || !a.WeKnoraHealthy || a.WeKnoraClient == nil {
		return nil
	}
	return weknorakp.NewProvider(a.WeKnoraClient, a.KnowledgeBaseID)
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
		nil,
		"",
		logger,
	)

	assembly := &AIAssembly{
		Service:         aidelivery.NewHandlerServiceAdapter(defaultService),
		RuntimeService:  defaultService,
		KnowledgeBaseID: weKnoraConfig.KnowledgeBaseID,
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

	timeout := opts.HealthCheckTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
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

	enhanced := services.NewOrchestratedEnhancedAIService(
		baseAI,
		openai.NewProvider(openAIConfig.APIKey, openAIConfig.BaseURL),
		weknorakp.NewProvider(client, weKnoraConfig.KnowledgeBaseID),
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
