package infra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"servify/apps/server/internal/models"
	agentapp "servify/apps/server/internal/modules/agent/application"
	agentdomain "servify/apps/server/internal/modules/agent/domain"
	platformauth "servify/apps/server/internal/platform/auth"
	"servify/apps/server/internal/platform/usersecurity"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateAgent(ctx context.Context, userID uint, department string, skills []string, maxChatConcurrency int) (*agentdomain.AgentProfile, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	var existing models.Agent
	if err := applyAgentScope(r.db.WithContext(ctx), ctx).Where("user_id = ?", userID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("user is already an agent")
	}
	agent := &models.Agent{
		UserID:        userID,
		Department:    department,
		Skills:        strings.Join(skills, ","),
		Status:        string(agentdomain.PresenceStatusOffline),
		MaxConcurrent: maxChatConcurrency,
		Rating:        5,
	}
	applyAgentScopeFields(ctx, agent)
	if err := r.db.WithContext(ctx).Create(agent).Error; err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}
	_ = r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("role", "agent").Error
	return mapProfile(user, *agent), nil
}

func (r *GormRepository) GetAgentByUserID(ctx context.Context, userID uint) (*agentdomain.AgentProfile, *models.Agent, error) {
	var agent models.Agent
	if err := applyAgentScope(r.db.WithContext(ctx), ctx).Preload("User").Preload("Tickets", func(db *gorm.DB) *gorm.DB {
		db = db.Where("status NOT IN ?", []string{"closed"}).Order("created_at DESC")
		if tenantID := platformauth.TenantIDFromContext(ctx); tenantID != "" {
			db = db.Where("tenant_id = ?", tenantID)
		}
		if workspaceID := platformauth.WorkspaceIDFromContext(ctx); workspaceID != "" {
			db = db.Where("workspace_id = ?", workspaceID)
		}
		return db
	}).Where("user_id = ?", userID).First(&agent).Error; err != nil {
		return nil, nil, fmt.Errorf("agent not found: %w", err)
	}
	profile := mapProfile(agent.User, agent)
	return profile, &agent, nil
}

func (r *GormRepository) ListAgents(ctx context.Context, limit int) ([]models.Agent, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	var agents []models.Agent
	if err := applyAgentScope(r.db.WithContext(ctx), ctx).Preload("User").Order("updated_at DESC").Limit(limit).Find(&agents).Error; err != nil {
		return nil, fmt.Errorf("failed to list agents: %w", err)
	}
	return agents, nil
}

func (r *GormRepository) GetAgentRuntimeByUserID(ctx context.Context, userID uint) (*agentapp.AgentRuntimeDTO, error) {
	var agent models.Agent
	if err := applyAgentScope(r.db.WithContext(ctx), ctx).
		Preload("User").
		Where("user_id = ? AND status <> ?", userID, string(agentdomain.PresenceStatusOffline)).
		First(&agent).Error; err != nil {
		return nil, fmt.Errorf("agent runtime not found: %w", err)
	}
	runtime := mapRuntime(agent.User, agent)
	return &runtime, nil
}

func (r *GormRepository) ListActiveAgentRuntimes(ctx context.Context) ([]agentapp.AgentRuntimeDTO, error) {
	var agents []models.Agent
	if err := applyAgentScope(r.db.WithContext(ctx), ctx).
		Preload("User").
		Where("status <> ?", string(agentdomain.PresenceStatusOffline)).
		Order("user_id ASC").
		Find(&agents).Error; err != nil {
		return nil, fmt.Errorf("failed to list active agent runtimes: %w", err)
	}
	runtimes := make([]agentapp.AgentRuntimeDTO, 0, len(agents))
	for _, agent := range agents {
		runtimes = append(runtimes, mapRuntime(agent.User, agent))
	}
	return runtimes, nil
}

func (r *GormRepository) UpdatePresenceStatus(ctx context.Context, userID uint, status agentdomain.PresenceStatus) error {
	if err := applyAgentScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Where("user_id = ?", userID).Update("status", string(status)).Error; err != nil {
		return fmt.Errorf("failed to update agent status: %w", err)
	}
	return nil
}

func (r *GormRepository) UpdateChatLoad(ctx context.Context, userID uint, currentLoad int) error {
	return applyAgentScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Where("user_id = ?", userID).Update("current_load", currentLoad).Error
}

func (r *GormRepository) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error; err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	return &session, nil
}

func (r *GormRepository) AssignSession(ctx context.Context, sessionID string, agentUserID uint) error {
	if err := r.db.WithContext(ctx).Model(&models.Session{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"agent_id": agentUserID,
		"status":   "active",
		"ended_at": nil,
	}).Error; err != nil {
		return fmt.Errorf("failed to assign session: %w", err)
	}
	return nil
}

func (r *GormRepository) ReleaseSession(ctx context.Context, sessionID string, agentUserID uint) error {
	if err := r.db.WithContext(ctx).Model(&models.Session{}).Where("id = ? AND agent_id = ?", sessionID, agentUserID).Updates(map[string]interface{}{
		"status":   "ended",
		"ended_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to release session: %w", err)
	}
	return nil
}

func (r *GormRepository) GetStats(ctx context.Context, agentUserID *uint) (*agentapp.AgentStatsDTO, error) {
	stats := &agentapp.AgentStatsDTO{}
	query := applyAgentScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx)
	if agentUserID != nil {
		query = query.Where("user_id = ?", *agentUserID)
	}
	query.Count(&stats.Total)
	var avgResponseTime float64
	applyAgentScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Select("AVG(avg_response_time)").Row().Scan(&avgResponseTime)
	stats.AvgResponseTime = int64(avgResponseTime)
	var avgRating float64
	applyAgentScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Select("AVG(rating)").Row().Scan(&avgRating)
	stats.AvgRating = avgRating
	return stats, nil
}

func (r *GormRepository) RevokeUserTokens(ctx context.Context, userID uint, revokeAt time.Time) (int, error) {
	return usersecurity.RevokeUserTokens(ctx, r.db, userID, revokeAt)
}

func applyAgentScope(db *gorm.DB, ctx context.Context) *gorm.DB {
	if tenantID := platformauth.TenantIDFromContext(ctx); tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	if workspaceID := platformauth.WorkspaceIDFromContext(ctx); workspaceID != "" {
		db = db.Where("workspace_id = ?", workspaceID)
	}
	return db
}

func applyAgentScopeFields(ctx context.Context, agent *models.Agent) {
	if agent == nil {
		return
	}
	if tenantID := platformauth.TenantIDFromContext(ctx); tenantID != "" {
		agent.TenantID = tenantID
	}
	if workspaceID := platformauth.WorkspaceIDFromContext(ctx); workspaceID != "" {
		agent.WorkspaceID = workspaceID
	}
}

func mapProfile(user models.User, agent models.Agent) *agentdomain.AgentProfile {
	return &agentdomain.AgentProfile{
		UserID:              agent.UserID,
		Username:            user.Username,
		Name:                user.Name,
		Department:          agent.Department,
		Skills:              splitSkills(agent.Skills),
		MaxChatConcurrency:  agent.MaxConcurrent,
		MaxVoiceConcurrency: 1,
		CurrentChatLoad:     agent.CurrentLoad,
		CurrentVoiceLoad:    0,
		Rating:              agent.Rating,
		AvgResponseTime:     agent.AvgResponseTime,
		TotalTickets:        agent.TotalTickets,
	}
}

func mapRuntime(user models.User, agent models.Agent) agentapp.AgentRuntimeDTO {
	return agentapp.AgentRuntimeDTO{
		UserID:              agent.UserID,
		Username:            user.Username,
		Name:                user.Name,
		Department:          agent.Department,
		Skills:              splitSkills(agent.Skills),
		Status:              agent.Status,
		MaxChatConcurrency:  agent.MaxConcurrent,
		MaxVoiceConcurrency: 1,
		CurrentChatLoad:     agent.CurrentLoad,
		CurrentVoiceLoad:    0,
		Rating:              agent.Rating,
		AvgResponseTime:     agent.AvgResponseTime,
	}
}

func splitSkills(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
