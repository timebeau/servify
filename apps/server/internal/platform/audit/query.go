package audit

import (
	"context"
	"time"

	"servify/apps/server/internal/models"

	"gorm.io/gorm"
)

type ListQuery struct {
	Action        string
	ResourceType  string
	ResourceID    string
	PrincipalKind string
	ActorUserID   *uint
	Success       *bool
	From          *time.Time
	To            *time.Time
	Page          int
	PageSize      int
}

type QueryService interface {
	List(ctx context.Context, query ListQuery) ([]models.AuditLog, int64, error)
}

type GormQueryService struct {
	db *gorm.DB
}

func NewGormQueryService(db *gorm.DB) *GormQueryService {
	if db == nil {
		return nil
	}
	return &GormQueryService{db: db}
}

func (s *GormQueryService) List(ctx context.Context, query ListQuery) ([]models.AuditLog, int64, error) {
	if s == nil || s.db == nil {
		return nil, 0, nil
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

	tx := s.db.WithContext(ctx).Model(&models.AuditLog{})
	if query.Action != "" {
		tx = tx.Where("action = ?", query.Action)
	}
	if query.ResourceType != "" {
		tx = tx.Where("resource_type = ?", query.ResourceType)
	}
	if query.ResourceID != "" {
		tx = tx.Where("resource_id = ?", query.ResourceID)
	}
	if query.PrincipalKind != "" {
		tx = tx.Where("principal_kind = ?", query.PrincipalKind)
	}
	if query.ActorUserID != nil {
		tx = tx.Where("actor_user_id = ?", *query.ActorUserID)
	}
	if query.Success != nil {
		tx = tx.Where("success = ?", *query.Success)
	}
	if query.From != nil {
		tx = tx.Where("created_at >= ?", *query.From)
	}
	if query.To != nil {
		tx = tx.Where("created_at <= ?", *query.To)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []models.AuditLog
	if err := tx.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
