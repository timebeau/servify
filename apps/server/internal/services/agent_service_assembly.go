package services

import (
	"context"

	agentapp "servify/apps/server/internal/modules/agent/application"
	agentinfra "servify/apps/server/internal/modules/agent/infra"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AgentServiceAssembly struct {
	Service       *AgentService
	Maintenance   *agentRuntimeMaintenance
	LegacyRuntime *agentLegacyRuntimeAdapter
}

type AgentServiceDependencies struct {
	DB            *gorm.DB
	Logger        *logrus.Logger
	Module        *agentapp.Service
	LegacyRuntime *agentLegacyRuntimeAdapter
}

func BuildAgentServiceAssembly(db *gorm.DB, logger *logrus.Logger) *AgentServiceAssembly {
	if logger == nil {
		logger = logrus.New()
	}

	repo := agentinfra.NewGormRepository(db)
	registry := agentinfra.NewInMemoryRegistry()
	module := agentapp.NewService(repo, registry)
	cache := &agentRuntimeCache{}
	legacyRuntime := newAgentLegacyRuntimeAdapter(cache)
	service := NewAgentServiceWithDependencies(AgentServiceDependencies{
		DB:            db,
		Logger:        logger,
		Module:        module,
		LegacyRuntime: legacyRuntime,
	})

	return &AgentServiceAssembly{
		Service:       service,
		Maintenance:   newAgentRuntimeMaintenance(logger, module, func(ctx context.Context) { legacyRuntime.Sync(ctx, module) }),
		LegacyRuntime: legacyRuntime,
	}
}
