package infra

import (
	"context"

	"servify/apps/server/internal/models"
	suggestionapp "servify/apps/server/internal/modules/suggestion/application"
	platformauth "servify/apps/server/internal/platform/auth"

	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindTicketCandidates(ctx context.Context, tokens []string, candidateMax int) ([]suggestionapp.TicketCandidate, error) {
	q := applyScopeFilter(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).
		Select("id, title, description, status, category, priority, created_at").
		Order("created_at DESC")

	where, args := suggestionapp.BuildLikeWhereTokens([]string{"title", "description"}, tokens, 3)
	if where != "" {
		q = q.Where(where, args...)
	}

	var rows []suggestionapp.TicketCandidate
	if err := q.Limit(candidateMax).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *GormRepository) FindKnowledgeDocCandidates(ctx context.Context, tokens []string) ([]suggestionapp.KnowledgeDocCandidate, error) {
	q := applyScopeFilter(r.db.WithContext(ctx).Model(&models.KnowledgeDoc{}), ctx).
		Select("id, title, content, category, tags").
		Order("created_at DESC")

	where, args := suggestionapp.BuildLikeWhereTokens([]string{"title", "content", "tags"}, tokens, 3)
	if where != "" {
		q = q.Where(where, args...)
	}

	var rows []suggestionapp.KnowledgeDocCandidate
	if err := q.Limit(300).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func applyScopeFilter(tx *gorm.DB, ctx context.Context) *gorm.DB {
	if tenantID := platformauth.TenantIDFromContext(ctx); tenantID != "" {
		tx = tx.Where("tenant_id = ?", tenantID)
	}
	if workspaceID := platformauth.WorkspaceIDFromContext(ctx); workspaceID != "" {
		tx = tx.Where("workspace_id = ?", workspaceID)
	}
	return tx
}
