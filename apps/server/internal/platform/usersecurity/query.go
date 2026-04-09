package usersecurity

import (
	"context"
	"time"

	"servify/apps/server/internal/models"

	"gorm.io/gorm"
)

type RevokedTokenListQuery struct {
	JTI        string
	UserID     *uint
	SessionID  string
	TokenUse   string
	ActiveOnly bool
	Page       int
	PageSize   int
}

type RevokedTokenQueryService interface {
	ListRevokedTokens(ctx context.Context, query RevokedTokenListQuery) ([]models.RevokedToken, int64, error)
}

func (s *Service) ListRevokedTokens(ctx context.Context, query RevokedTokenListQuery) ([]models.RevokedToken, int64, error) {
	if s == nil || s.db == nil {
		return nil, 0, gorm.ErrInvalidDB
	}

	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}

	tx := s.db.WithContext(ctx).Model(&models.RevokedToken{})
	if hasRequestScope(ctx) {
		tx = tx.Where(
			"user_id IN (?) OR user_id IN (?)",
			scopedUserIDQuery(ctx, s.db, &models.Agent{}),
			scopedUserIDQuery(ctx, s.db, &models.Customer{}),
		)
	}
	if query.JTI != "" {
		tx = tx.Where("jti = ?", query.JTI)
	}
	if query.UserID != nil {
		tx = tx.Where("user_id = ?", *query.UserID)
	}
	if query.SessionID != "" {
		tx = tx.Where("session_id = ?", query.SessionID)
	}
	if query.TokenUse != "" {
		tx = tx.Where("token_use = ?", query.TokenUse)
	}
	if query.ActiveOnly {
		tx = tx.Where("expires_at IS NULL OR expires_at >= ?", time.Now().UTC())
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.RevokedToken
	if err := tx.Order("revoked_at DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
