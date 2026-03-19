package server

import (
	"context"
	"fmt"
	"time"

	"servify/apps/server/internal/config"
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
	Service        services.AIServiceInterface
	LegacyService  *services.AIService
	WeKnoraClient  weknora.WeKnoraInterface
	WeKnoraHealthy bool
}

func BuildAIAssembly(cfg *config.Config, logger *logrus.Logger, opts AIAssemblyOptions) (*AIAssembly, error) {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	baseAI := services.NewAIService(cfg.AI.OpenAI.APIKey, cfg.AI.OpenAI.BaseURL)
	baseAI.InitializeKnowledgeBase()

	assembly := &AIAssembly{
		Service:       services.NewOrchestratedAIService(openai.NewProvider(cfg.AI.OpenAI.APIKey, cfg.AI.OpenAI.BaseURL), nil),
		LegacyService: baseAI,
	}
	if !cfg.WeKnora.Enabled {
		return assembly, nil
	}

	client := weknora.NewClient(&weknora.Config{
		BaseURL:    cfg.WeKnora.BaseURL,
		APIKey:     cfg.WeKnora.APIKey,
		TenantID:   cfg.WeKnora.TenantID,
		Timeout:    cfg.WeKnora.Timeout,
		MaxRetries: cfg.WeKnora.MaxRetries,
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
		openai.NewProvider(cfg.AI.OpenAI.APIKey, cfg.AI.OpenAI.BaseURL),
		weknorakp.NewProvider(client, cfg.WeKnora.KnowledgeBaseID),
		client,
		cfg.WeKnora.KnowledgeBaseID,
		logger,
	)
	if opts.SyncKnowledgeBase {
		syncCtx, syncCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer syncCancel()
		if err := enhanced.SyncKnowledgeBase(syncCtx); err != nil {
			logger.Warnf("Knowledge base sync failed: %v", err)
		}
	}
	assembly.Service = enhanced
	return assembly, nil
}
