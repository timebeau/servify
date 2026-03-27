package delivery

import (
	"context"

	"servify/apps/server/internal/models"
)

// HandlerService is the only agent contract that HTTP handlers should depend on.
type HandlerService interface {
	CreateAgent(ctx context.Context, req *AgentCreateRequest) (*models.Agent, error)
	GetAgentByUserID(ctx context.Context, userID uint) (*models.Agent, error)
	ListAgents(ctx context.Context, limit int) ([]models.Agent, error)
	AgentGoOnline(ctx context.Context, userID uint) error
	AgentGoOffline(ctx context.Context, userID uint) error
	UpdateAgentStatus(ctx context.Context, userID uint, status string) error
	RevokeAgentTokens(ctx context.Context, userID uint) (int, error)
	AssignSessionToAgent(ctx context.Context, sessionID string, agentID uint) error
	ReleaseSessionFromAgent(ctx context.Context, sessionID string, agentID uint) error
	FindAvailableAgent(ctx context.Context, skills []string, priority string) (*AgentInfo, error)
	GetOnlineAgents(ctx context.Context) []*AgentInfo
	GetAgentStats(ctx context.Context, agentID *uint) (*AgentStats, error)
}
