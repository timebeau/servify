package audit

import (
	"context"
	"time"

	"servify/apps/server/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RetentionService deletes expired audit logs according to a configured policy.
type RetentionService interface {
	Cleanup(ctx context.Context, now time.Time) (int64, error)
}

type GormRetentionService struct {
	db        *gorm.DB
	retention time.Duration
	batchSize int
	logger    *logrus.Logger
}

func NewGormRetentionService(db *gorm.DB, retention time.Duration, batchSize int) *GormRetentionService {
	if db == nil || retention <= 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 500
	}
	return &GormRetentionService{
		db:        db,
		retention: retention,
		batchSize: batchSize,
		logger:    logrus.StandardLogger(),
	}
}

func (s *GormRetentionService) Cleanup(ctx context.Context, now time.Time) (int64, error) {
	if s == nil || s.db == nil || s.retention <= 0 {
		return 0, nil
	}

	cutoff := now.Add(-s.retention)
	var deleted int64

	for {
		var ids []uint
		if err := s.db.WithContext(ctx).
			Model(&models.AuditLog{}).
			Where("created_at < ?", cutoff).
			Order("created_at ASC").
			Limit(s.batchSize).
			Pluck("id", &ids).Error; err != nil {
			return deleted, err
		}
		if len(ids) == 0 {
			return deleted, nil
		}

		res := s.db.WithContext(ctx).Delete(&models.AuditLog{}, ids)
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
