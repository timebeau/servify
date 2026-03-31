package infra

import (
	"context"

	"servify/apps/server/internal/models"
	voiceapp "servify/apps/server/internal/modules/voice/application"

	"gorm.io/gorm"
)

// compile-time interface check
var _ voiceapp.TranscriptRepository = (*GormTranscriptRepository)(nil)

type GormTranscriptRepository struct {
	db *gorm.DB
}

func NewGormTranscriptRepository(db *gorm.DB) *GormTranscriptRepository {
	return &GormTranscriptRepository{db: db}
}

func (r *GormTranscriptRepository) Append(ctx context.Context, transcript voiceapp.TranscriptDTO) error {
	m := models.VoiceTranscript{
		CallID:    transcript.CallID,
		Content:   transcript.Content,
		Language:  transcript.Language,
		Finalized: transcript.Finalized,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormTranscriptRepository) ListByCallID(ctx context.Context, callID string) ([]voiceapp.TranscriptDTO, error) {
	var records []models.VoiceTranscript
	if err := r.db.WithContext(ctx).
		Where("call_id = ?", callID).
		Order("created_at ASC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	out := make([]voiceapp.TranscriptDTO, len(records))
	for i, r := range records {
		out[i] = voiceapp.TranscriptDTO{
			CallID:     r.CallID,
			Content:    r.Content,
			Language:   r.Language,
			Finalized:  r.Finalized,
			AppendedAt: r.CreatedAt,
		}
	}
	return out, nil
}
