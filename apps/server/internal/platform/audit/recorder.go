package audit

import (
	"context"

	"servify/apps/server/internal/models"

	"gorm.io/gorm"
)

type Entry struct {
	ActorUserID   *uint
	PrincipalKind string
	Action        string
	ResourceType  string
	ResourceID    string
	Route         string
	Method        string
	StatusCode    int
	Success       bool
	RequestID     string
	ClientIP      string
	UserAgent     string
	TenantID      string
	WorkspaceID   string
	RequestJSON   string
	BeforeJSON    string
	AfterJSON     string
}

type Recorder interface {
	Record(ctx context.Context, entry Entry) error
}

type GormRecorder struct {
	db *gorm.DB
}

func NewGormRecorder(db *gorm.DB) *GormRecorder {
	if db == nil {
		return nil
	}
	return &GormRecorder{db: db}
}

func (r *GormRecorder) Record(ctx context.Context, entry Entry) error {
	if r == nil || r.db == nil {
		return nil
	}
	record := models.AuditLog{
		ActorUserID:   entry.ActorUserID,
		PrincipalKind: entry.PrincipalKind,
		Action:        entry.Action,
		ResourceType:  entry.ResourceType,
		ResourceID:    entry.ResourceID,
		Route:         entry.Route,
		Method:        entry.Method,
		StatusCode:    entry.StatusCode,
		Success:       entry.Success,
		RequestID:     entry.RequestID,
		ClientIP:      entry.ClientIP,
		UserAgent:     entry.UserAgent,
		TenantID:      entry.TenantID,
		WorkspaceID:   entry.WorkspaceID,
		RequestJSON:   entry.RequestJSON,
		BeforeJSON:    entry.BeforeJSON,
		AfterJSON:     entry.AfterJSON,
	}
	return r.db.WithContext(ctx).Create(&record).Error
}
