package infra

import (
	"testing"

	"servify/apps/server/internal/models"
	agentapp "servify/apps/server/internal/modules/agent/application"
)

func TestRedisRegistryHandleStatusChangeIgnoresMalformedPayload(t *testing.T) {
	registry := &RedisRegistry{
		localCache: NewInMemoryRegistry(),
	}

	registry.handleStatusChange("invalid-payload")
	if got := registry.localCache.List(); len(got) != 0 {
		t.Fatalf("expected no agents after malformed payload, got %d", len(got))
	}
}

func TestRedisRegistryHandleStatusChangeOfflineRemovesAgent(t *testing.T) {
	registry := &RedisRegistry{
		localCache: NewInMemoryRegistry(),
	}
	registry.localCache.agents[7] = &runtimeAgent{
		AgentRuntimeDTO: agentapp.AgentRuntimeDTO{UserID: 7, Status: "online"},
		Sessions:        map[string]*models.Session{},
	}

	registry.handleStatusChange("offline:7")
	if _, ok := registry.localCache.Get(7); ok {
		t.Fatal("expected agent to be removed from local cache")
	}
}
