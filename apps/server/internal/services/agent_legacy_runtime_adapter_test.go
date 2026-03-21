package services

import (
	"context"
	"testing"
	"time"

	agentapp "servify/apps/server/internal/modules/agent/application"
	agentdomain "servify/apps/server/internal/modules/agent/domain"

	"servify/apps/server/internal/models"
)

type legacyAdapterRepo struct{}

func (r *legacyAdapterRepo) CreateAgent(ctx context.Context, userID uint, department string, skills []string, maxChatConcurrency int) (*agentdomain.AgentProfile, error) {
	return nil, nil
}

func (r *legacyAdapterRepo) GetAgentByUserID(ctx context.Context, userID uint) (*agentdomain.AgentProfile, *models.Agent, error) {
	return nil, nil, nil
}

func (r *legacyAdapterRepo) ListAgents(ctx context.Context, limit int) ([]models.Agent, error) {
	return nil, nil
}

func (r *legacyAdapterRepo) UpdatePresenceStatus(ctx context.Context, userID uint, status agentdomain.PresenceStatus) error {
	return nil
}

func (r *legacyAdapterRepo) UpdateChatLoad(ctx context.Context, userID uint, currentLoad int) error {
	return nil
}

func (r *legacyAdapterRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	return nil, nil
}

func (r *legacyAdapterRepo) AssignSession(ctx context.Context, sessionID string, agentUserID uint) error {
	return nil
}

func (r *legacyAdapterRepo) ReleaseSession(ctx context.Context, sessionID string, agentUserID uint) error {
	return nil
}

func (r *legacyAdapterRepo) GetStats(ctx context.Context, agentUserID *uint) (*agentapp.AgentStatsDTO, error) {
	return &agentapp.AgentStatsDTO{}, nil
}

type legacyAdapterRegistry struct {
	items map[uint]agentapp.AgentRuntimeDTO
}

func (r *legacyAdapterRegistry) GoOnline(profile agentdomain.AgentProfile) (agentapp.AgentRuntimeDTO, error) {
	return agentapp.AgentRuntimeDTO{}, nil
}

func (r *legacyAdapterRegistry) GoOffline(userID uint) {}

func (r *legacyAdapterRegistry) UpdateStatus(userID uint, status agentdomain.PresenceStatus) {}

func (r *legacyAdapterRegistry) AssignSession(userID uint, session *models.Session) (agentapp.AgentRuntimeDTO, error) {
	return agentapp.AgentRuntimeDTO{}, nil
}

func (r *legacyAdapterRegistry) ReleaseSession(userID uint, sessionID string) (agentapp.AgentRuntimeDTO, bool) {
	return agentapp.AgentRuntimeDTO{}, false
}

func (r *legacyAdapterRegistry) ApplyTransfer(sessionID string, fromAgentID *uint, toAgentID uint) {}

func (r *legacyAdapterRegistry) Get(userID uint) (agentapp.AgentRuntimeDTO, bool) {
	item, ok := r.items[userID]
	return item, ok
}

func (r *legacyAdapterRegistry) List() []agentapp.AgentRuntimeDTO {
	out := make([]agentapp.AgentRuntimeDTO, 0, len(r.items))
	for _, item := range r.items {
		out = append(out, item)
	}
	return out
}

func TestAgentLegacyRuntimeAdapter_Sync(t *testing.T) {
	cache := &agentRuntimeCache{}
	adapter := newAgentLegacyRuntimeAdapter(cache)
	module := agentapp.NewService(&legacyAdapterRepo{}, &legacyAdapterRegistry{
		items: map[uint]agentapp.AgentRuntimeDTO{
			1: {
				UserID:             1,
				Username:           "agent1",
				Name:               "Agent One",
				Department:         "support",
				Status:             string(agentdomain.PresenceStatusOnline),
				MaxChatConcurrency: 3,
				CurrentChatLoad:    1,
				LastActivity:       time.Now(),
			},
		},
	})

	adapter.Sync(context.Background(), module)

	agent, ok := adapter.GetOnlineAgent(1)
	if !ok {
		t.Fatal("expected online agent in legacy cache")
	}
	if agent.UserID != 1 || agent.CurrentLoad != 1 || agent.MaxConcurrent != 3 {
		t.Fatalf("unexpected cached agent: %+v", agent)
	}
}
