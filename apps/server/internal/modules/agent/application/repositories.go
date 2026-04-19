package application

import (
	"context"
	"time"

	"servify/apps/server/internal/models"
	agentdomain "servify/apps/server/internal/modules/agent/domain"
)

type Repository interface {
	CreateAgent(ctx context.Context, userID uint, department string, skills []string, maxChatConcurrency int) (*agentdomain.AgentProfile, error)
	GetAgentByUserID(ctx context.Context, userID uint) (*agentdomain.AgentProfile, *models.Agent, error)
	GetAgentRuntimeByUserID(ctx context.Context, userID uint) (*AgentRuntimeDTO, error)
	ListActiveAgentRuntimes(ctx context.Context) ([]AgentRuntimeDTO, error)
	ListAgents(ctx context.Context, limit int) ([]models.Agent, error)
	UpdatePresenceStatus(ctx context.Context, userID uint, status agentdomain.PresenceStatus) error
	UpdateChatLoad(ctx context.Context, userID uint, currentLoad int) error
	// Persisted runtime metadata methods (replace in-memory registry)
	UpdateLastActivity(ctx context.Context, userID uint) error
	SetConnectedTime(ctx context.Context, userID uint) error
	ClearConnectedTime(ctx context.Context, userID uint) error
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	AssignSession(ctx context.Context, sessionID string, agentUserID uint) error
	ReleaseSession(ctx context.Context, sessionID string, agentUserID uint) error
	GetStats(ctx context.Context, agentUserID *uint) (*AgentStatsDTO, error)
	RevokeUserTokens(ctx context.Context, userID uint, revokeAt time.Time) (int, error)
}

type RuntimeRegistry interface {
	GoOnline(profile agentdomain.AgentProfile) (AgentRuntimeDTO, error)
	GoOffline(userID uint)
	UpdateStatus(userID uint, status agentdomain.PresenceStatus)
	AssignSession(userID uint, session *models.Session) (AgentRuntimeDTO, error)
	ReleaseSession(userID uint, sessionID string) (AgentRuntimeDTO, bool)
	ApplyTransfer(sessionID string, fromAgentID *uint, toAgentID uint)
	Get(userID uint) (AgentRuntimeDTO, bool)
	List() []AgentRuntimeDTO
}
