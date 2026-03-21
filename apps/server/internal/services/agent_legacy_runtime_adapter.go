package services

import (
	"context"

	agentapp "servify/apps/server/internal/modules/agent/application"
)

type agentLegacyRuntimeAdapter struct {
	cache *agentRuntimeCache
}

func newAgentLegacyRuntimeAdapter(cache *agentRuntimeCache) *agentLegacyRuntimeAdapter {
	return &agentLegacyRuntimeAdapter{cache: cache}
}

func (a *agentLegacyRuntimeAdapter) GetOnlineAgent(userID uint) (*AgentInfo, bool) {
	return a.cache.Load(userID)
}

func (a *agentLegacyRuntimeAdapter) Sync(ctx context.Context, module *agentapp.Service) {
	runtimes := module.GetOnlineAgents(ctx)
	active := make(map[uint]struct{}, len(runtimes))
	for _, runtime := range runtimes {
		active[runtime.UserID] = struct{}{}
		a.cache.Store(runtime.UserID, mapRuntimeToLegacy(&runtime))
	}
	stale := a.cache.CollectStale(active)
	for _, userID := range stale {
		a.cache.Delete(userID)
	}
}
