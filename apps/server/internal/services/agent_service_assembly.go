package services

import (
	agentapp "servify/apps/server/internal/modules/agent/application"
	agentinfra "servify/apps/server/internal/modules/agent/infra"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AgentServiceAssembly struct {
	Service     *AgentService
	Maintenance *agentRuntimeMaintenance
}

type AgentServiceDependencies struct {
	DB     *gorm.DB
	Logger *logrus.Logger
	Module *agentapp.Service
}

func BuildAgentServiceAssembly(db *gorm.DB, logger *logrus.Logger, redisClient *redis.Client) *AgentServiceAssembly {
	if logger == nil {
		logger = logrus.New()
	}

	repo := agentinfra.NewGormRepository(db)

	var registry agentapp.RuntimeRegistry
	if redisClient != nil {
		logger.Info("using redis-backed agent registry for multi-instance support")
		registry = agentinfra.NewRedisRegistry(redisClient, db, logger)
	} else {
		logger.Warn("using in-memory agent registry - not suitable for multi-instance deployment")
		registry = agentinfra.NewInMemoryRegistry()
	}

	module := agentapp.NewService(repo, registry)
	service := NewAgentServiceWithDependencies(AgentServiceDependencies{
		DB:     db,
		Logger: logger,
		Module: module,
	})

	return &AgentServiceAssembly{
		Service:     service,
		Maintenance: newAgentRuntimeMaintenance(logger, module),
	}
}
