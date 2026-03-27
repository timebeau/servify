package usersecurity

import (
	"context"
	"time"

	"servify/apps/server/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewService(db *gorm.DB, logger *logrus.Logger) *Service {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	return &Service{db: db, logger: logger}
}

func (s *Service) RevokeTokens(ctx context.Context, userID uint) (int, error) {
	return RevokeUserTokens(ctx, s.db, userID, time.Now().UTC())
}

func (s *Service) GetUser(ctx context.Context, userID uint) (*models.User, error) {
	if s == nil || s.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
