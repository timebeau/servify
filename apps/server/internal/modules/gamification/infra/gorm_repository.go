package infra

import (
	"context"
	"strings"

	"servify/apps/server/internal/models"
	gamificationapp "servify/apps/server/internal/modules/gamification/application"
	platformauth "servify/apps/server/internal/platform/auth"

	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) ListAgentProfiles(ctx context.Context, department string) ([]gamificationapp.AgentProfile, error) {
	q := applyScopeFilter(r.db.WithContext(ctx).
		Model(&models.Agent{}).
		Select("agents.user_id as user_id, users.username as username, users.name as name, agents.department as department, agents.avg_response_time as avg_response_time").
		Joins("LEFT JOIN users ON users.id = agents.user_id"), ctx)
	if department != "" {
		q = q.Where("agents.department = ?", department)
	}
	var profiles []gamificationapp.AgentProfile
	if err := q.Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

func (r *GormRepository) ListResolvedCounts(ctx context.Context, startDate, endDate string) ([]gamificationapp.AgentResolvedCount, error) {
	var resolved []gamificationapp.AgentResolvedCount
	if err := applyScopeFilter(r.db.WithContext(ctx).
		Model(&models.Ticket{}).
		Select("agent_id as agent_id, COUNT(*) as count").
		Where("agent_id IS NOT NULL").
		Where("resolved_at IS NOT NULL").
		Where("resolved_at >= ? AND resolved_at <= ?", startDate, endDate).
		Where("status IN ?", []string{"resolved", "closed"}).
		Group("agent_id"), ctx).
		Scan(&resolved).Error; err != nil {
		return nil, err
	}
	return resolved, nil
}

func (r *GormRepository) ListCSATStats(ctx context.Context, startDate, endDate string) ([]gamificationapp.AgentCSAT, error) {
	var csats []gamificationapp.AgentCSAT
	if err := applyScopeFilter(r.db.WithContext(ctx).
		Model(&models.CustomerSatisfaction{}).
		Select("agent_id as agent_id, AVG(rating) as avg, COUNT(*) as count").
		Where("agent_id IS NOT NULL").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("agent_id"), ctx).
		Scan(&csats).Error; err != nil {
		return nil, err
	}
	return csats, nil
}

func applyScopeFilter(tx *gorm.DB, ctx context.Context) *gorm.DB {
	if tenantID := strings.TrimSpace(platformauth.TenantIDFromContext(ctx)); tenantID != "" {
		tx = tx.Where("tenant_id = ?", tenantID)
	}
	if workspaceID := strings.TrimSpace(platformauth.WorkspaceIDFromContext(ctx)); workspaceID != "" {
		tx = tx.Where("workspace_id = ?", workspaceID)
	}
	return tx
}
