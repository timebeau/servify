package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/modules/routing/domain"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateAssignment(ctx context.Context, assignment *domain.Assignment) error {
	if assignment == nil {
		return fmt.Errorf("assignment required")
	}
	model := mapTransferRecordModel(*assignment)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	assignment.AssignedAt = model.TransferredAt
	return nil
}

func (r *GormRepository) CreateQueueEntry(ctx context.Context, entry *domain.QueueEntry) error {
	if entry == nil {
		return fmt.Errorf("queue entry required")
	}
	model := mapWaitingRecordModel(*entry)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	*entry = mapQueueEntry(model)
	return nil
}

func (r *GormRepository) GetQueueEntry(ctx context.Context, sessionID string) (*domain.QueueEntry, error) {
	var model models.WaitingRecord
	if err := r.db.WithContext(ctx).
		Order("queued_at DESC").
		First(&model, "session_id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	item := mapQueueEntry(model)
	return &item, nil
}

func (r *GormRepository) ListQueueEntries(ctx context.Context, status string, limit int) ([]domain.QueueEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var items []models.WaitingRecord
	if err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("priority DESC, queued_at ASC").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, err
	}
	out := make([]domain.QueueEntry, 0, len(items))
	for _, item := range items {
		out = append(out, mapQueueEntry(item))
	}
	return out, nil
}

func (r *GormRepository) UpdateQueueEntry(ctx context.Context, entry *domain.QueueEntry) error {
	if entry == nil {
		return fmt.Errorf("queue entry required")
	}
	updates := map[string]interface{}{
		"reason":        entry.Reason,
		"target_skills": marshalSkills(entry.TargetSkills),
		"priority":      entry.Priority,
		"notes":         entry.Notes,
		"status":        string(entry.Status),
		"queued_at":     entry.QueuedAt,
		"assigned_at":   entry.AssignedAt,
		"assigned_to":   entry.AssignedTo,
	}
	return r.db.WithContext(ctx).
		Model(&models.WaitingRecord{}).
		Where("session_id = ?", entry.SessionID).
		Updates(updates).Error
}

func mapTransferRecordModel(item domain.Assignment) models.TransferRecord {
	return models.TransferRecord{
		SessionID:     item.SessionID,
		FromAgentID:   item.FromAgentID,
		ToAgentID:     uintPtr(item.ToAgentID),
		Reason:        item.Reason,
		Notes:         item.Notes,
		TransferredAt: item.AssignedAt,
		CreatedAt:     item.AssignedAt,
	}
}

func mapQueueEntry(model models.WaitingRecord) domain.QueueEntry {
	return domain.QueueEntry{
		SessionID:    model.SessionID,
		Reason:       model.Reason,
		TargetSkills: unmarshalSkills(model.TargetSkills),
		Priority:     model.Priority,
		Notes:        model.Notes,
		Status:       mapQueueStatus(model.Status),
		QueuedAt:     model.QueuedAt,
		AssignedAt:   model.AssignedAt,
		AssignedTo:   model.AssignedTo,
	}
}

func mapWaitingRecordModel(item domain.QueueEntry) models.WaitingRecord {
	return models.WaitingRecord{
		SessionID:    item.SessionID,
		Reason:       item.Reason,
		TargetSkills: marshalSkills(item.TargetSkills),
		Priority:     item.Priority,
		Notes:        item.Notes,
		Status:       string(item.Status),
		QueuedAt:     item.QueuedAt,
		AssignedAt:   item.AssignedAt,
		AssignedTo:   item.AssignedTo,
		CreatedAt:    item.QueuedAt,
	}
}

func mapQueueStatus(status string) domain.QueueStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "transferred":
		return domain.QueueStatusTransferred
	case "cancelled":
		return domain.QueueStatusCancelled
	default:
		return domain.QueueStatusWaiting
	}
}

func marshalSkills(skills []string) string {
	if len(skills) == 0 {
		return ""
	}
	data, err := json.Marshal(skills)
	if err != nil {
		return ""
	}
	return string(data)
}

func unmarshalSkills(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var skills []string
	if err := json.Unmarshal([]byte(raw), &skills); err == nil && len(skills) > 0 {
		return skills
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

func uintPtr(v uint) *uint {
	return &v
}
