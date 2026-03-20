package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/services"
)

// HandlerService is the only agent contract that HTTP handlers should depend on.
type HandlerService interface {
	CreateAgent(ctx context.Context, req *services.AgentCreateRequest) (*models.Agent, error)
	GetAgentByUserID(ctx context.Context, userID uint) (*models.Agent, error)
	ListAgents(ctx context.Context, limit int) ([]models.Agent, error)
	AgentGoOnline(ctx context.Context, userID uint) error
	AgentGoOffline(ctx context.Context, userID uint) error
	UpdateAgentStatus(ctx context.Context, userID uint, status string) error
	AssignSessionToAgent(ctx context.Context, sessionID string, agentID uint) error
	ReleaseSessionFromAgent(ctx context.Context, sessionID string, agentID uint) error
	FindAvailableAgent(ctx context.Context, skills []string, priority string) (*services.AgentInfo, error)
	GetOnlineAgents(ctx context.Context) []*services.AgentInfo
	GetAgentStats(ctx context.Context, agentID *uint) (*services.AgentStats, error)
}
