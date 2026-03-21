package services

import (
	"context"
	"time"

	agentapp "servify/apps/server/internal/modules/agent/application"

	"github.com/sirupsen/logrus"
)

type agentRuntimeMaintenance struct {
	logger   *logrus.Logger
	module   *agentapp.Service
	onSynced func(context.Context)
}

func newAgentRuntimeMaintenance(logger *logrus.Logger, module *agentapp.Service, onSynced func(context.Context)) *agentRuntimeMaintenance {
	return &agentRuntimeMaintenance{
		logger:   logger,
		module:   module,
		onSynced: onSynced,
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
		if time.Since(item.LastActivity) > timeout {
			m.logger.Warnf("Agent %d appears inactive, marking as away", item.UserID)
			_ = m.module.MarkAway(ctx, item.UserID)
		}
	}
	if m.onSynced != nil {
		m.onSynced(ctx)
	}
}

func (m *agentRuntimeMaintenance) updateAgentMetrics() {}
