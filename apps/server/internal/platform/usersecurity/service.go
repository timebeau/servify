package usersecurity

import (
	"context"
	"fmt"
	"sort"
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

func (s *Service) BatchRevokeTokens(ctx context.Context, userIDs []uint) (map[uint]int, error) {
	if s == nil || s.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if len(userIDs) == 0 {
		return map[uint]int{}, nil
	}

	results := make(map[uint]int, len(userIDs))
	revokeAt := time.Now().UTC()
	for _, userID := range userIDs {
		if userID == 0 {
			return nil, fmt.Errorf("user_id required")
		}
		version, err := RevokeUserTokens(ctx, s.db, userID, revokeAt)
		if err != nil {
			return nil, err
		}
		results[userID] = version
	}
	return results, nil
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

func (s *Service) GetUsers(ctx context.Context, userIDs []uint) ([]models.User, error) {
	if s == nil || s.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if len(userIDs) == 0 {
		return []models.User{}, nil
	}

	orderedIDs := make([]uint, 0, len(userIDs))
	seen := make(map[uint]struct{}, len(userIDs))
	for _, userID := range userIDs {
		if userID == 0 {
			return nil, fmt.Errorf("user_id required")
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}
		orderedIDs = append(orderedIDs, userID)
	}

	var users []models.User
	if err := s.db.WithContext(ctx).Where("id IN ?", orderedIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) != len(orderedIDs) {
		found := make(map[uint]struct{}, len(users))
		for _, user := range users {
			found[user.ID] = struct{}{}
		}
		missing := make([]uint, 0, len(orderedIDs))
		for _, userID := range orderedIDs {
			if _, ok := found[userID]; !ok {
				missing = append(missing, userID)
			}
		}
		sort.Slice(missing, func(i, j int) bool { return missing[i] < missing[j] })
		return nil, fmt.Errorf("user not found: %v", missing)
	}

	byID := make(map[uint]models.User, len(users))
	for _, user := range users {
		byID[user.ID] = user
	}

	orderedUsers := make([]models.User, 0, len(userIDs))
	for _, userID := range userIDs {
		user, ok := byID[userID]
		if !ok {
			return nil, fmt.Errorf("user not found: %d", userID)
		}
		orderedUsers = append(orderedUsers, user)
	}
	return orderedUsers, nil
}

func (s *Service) ListUserSessions(ctx context.Context, userID uint) ([]models.UserAuthSession, error) {
	if s == nil || s.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if userID == 0 {
		return nil, fmt.Errorf("user_id required")
	}

	var sessions []models.UserAuthSession
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at desc, created_at desc").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *Service) RevokeSession(ctx context.Context, userID uint, sessionID string) (*models.UserAuthSession, error) {
	if s == nil || s.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if userID == 0 {
		return nil, fmt.Errorf("user_id required")
	}
	if sessionID == "" {
		return nil, fmt.Errorf("session_id required")
	}

	now := time.Now().UTC()
	result := s.db.WithContext(ctx).Model(&models.UserAuthSession{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Updates(map[string]interface{}{
			"status":        "revoked",
			"revoked_at":    now,
			"token_version": gorm.Expr("COALESCE(token_version, 0) + 1"),
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("session not found")
	}

	var session models.UserAuthSession
	if err := s.db.WithContext(ctx).First(&session, "id = ? AND user_id = ?", sessionID, userID).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
