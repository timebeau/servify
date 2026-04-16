package services

import (
	"context"
	"time"

	agentapp "servify/apps/server/internal/modules/agent/application"

	"github.com/sirupsen/logrus"
)

type agentRuntimeMaintenance struct {
	logger *logrus.Logger
	module *agentapp.Service
}

func newAgentRuntimeMaintenance(logger *logrus.Logger, module *agentapp.Service) *agentRuntimeMaintenance {
	return &agentRuntimeMaintenance{
		logger: logger,
		module: module,
	}
}

func (m *agentRuntimeMaintenance) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.cleanupInactiveAgents(context.Background(), 5*time.Minute)
		m.updateAgentMetrics()
	}
}

func (m *agentRuntimeMaintenance) cleanupInactiveAgents(ctx context.Context, timeout time.Duration) {
	runtimes := m.module.GetOnlineAgents(ctx)
	for _, item := range runtimes {
		if item.LastActivity.IsZero() {
			continue
		}
		if time.Since(item.LastActivity) > timeout {
			m.logger.Warnf("Agent %d appears inactive, marking as away", item.UserID)
			_ = m.module.MarkAway(ctx, item.UserID)
		}
	}
}

func (m *agentRuntimeMaintenance) updateAgentMetrics() {}
