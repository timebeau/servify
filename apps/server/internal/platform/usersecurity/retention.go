package usersecurity

import (
	"context"
	"time"

	"servify/apps/server/internal/models"

	"gorm.io/gorm"
)

type RevokedTokenRetentionService interface {
	Cleanup(ctx context.Context, now time.Time) (int64, error)
}

type GormRevokedTokenRetentionService struct {
	db        *gorm.DB
	batchSize int
}

func NewGormRevokedTokenRetentionService(db *gorm.DB, batchSize int) *GormRevokedTokenRetentionService {
	if db == nil {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 500
	}
	return &GormRevokedTokenRetentionService{
		db:        db,
		batchSize: batchSize,
	}
}

func (s *GormRevokedTokenRetentionService) Cleanup(ctx context.Context, now time.Time) (int64, error) {
	if s == nil || s.db == nil {
		return 0, nil
	}

	var deleted int64
	for {
		var ids []uint
		if err := s.db.WithContext(ctx).
			Model(&models.RevokedToken{}).
			Where("expires_at IS NOT NULL AND expires_at < ?", now).
			Order("expires_at ASC").
			Limit(s.batchSize).
			Pluck("id", &ids).Error; err != nil {
			return deleted, err
		}
		if len(ids) == 0 {
			return deleted, nil
		}

		res := s.db.WithContext(ctx).Delete(&models.RevokedToken{}, ids)
		if res.Error != nil {
			return deleted, res.Error
		}
		deleted += res.RowsAffected
		if len(ids) < s.batchSize {
			return deleted, nil
		}
		if err := ctx.Err(); err != nil {
			return deleted, err
		}
	}
}
