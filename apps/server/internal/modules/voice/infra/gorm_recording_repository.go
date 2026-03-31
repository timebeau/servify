package infra

import (
	"context"
	"fmt"

	"servify/apps/server/internal/models"
	voiceapp "servify/apps/server/internal/modules/voice/application"

	"gorm.io/gorm"
)

// compile-time interface check
var _ voiceapp.RecordingRepository = (*GormRecordingRepository)(nil)

type GormRecordingRepository struct {
	db *gorm.DB
}

func NewGormRecordingRepository(db *gorm.DB) *GormRecordingRepository {
	return &GormRecordingRepository{db: db}
}

func (r *GormRecordingRepository) Save(ctx context.Context, recording voiceapp.RecordingDTO) error {
	m := models.VoiceRecording{
		ID:        recording.ID,
		CallID:    recording.CallID,
		Provider:  recording.Provider,
		Status:    recording.Status,
		StartedAt: recording.StartedAt,
	}
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("save recording: %w", err)
	}
	return nil
}

func (r *GormRecordingRepository) MarkStopped(ctx context.Context, recordingID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.VoiceRecording{}).
		Where("id = ?", recordingID).
		Update("status", "stopped")
	if result.Error != nil {
		return fmt.Errorf("mark recording stopped: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("recording not found")
	}
	return nil
}

func (r *GormRecordingRepository) FindByID(ctx context.Context, recordingID string) (*voiceapp.RecordingDTO, error) {
	var m models.VoiceRecording
	if err := r.db.WithContext(ctx).First(&m, "id = ?", recordingID).Error; err != nil {
		return nil, fmt.Errorf("recording not found: %w", err)
	}
	dto := voiceapp.RecordingDTO{
		ID:        m.ID,
		CallID:    m.CallID,
		Provider:  m.Provider,
		Status:    m.Status,
		StartedAt: m.StartedAt,
	}
	return &dto, nil
}
