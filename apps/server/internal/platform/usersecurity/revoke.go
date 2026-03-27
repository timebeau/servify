package usersecurity

import (
	"context"
	"fmt"
	"time"

	"servify/apps/server/internal/models"

	"gorm.io/gorm"
)

func RevokeUserTokens(ctx context.Context, db *gorm.DB, userID uint, revokeAt time.Time) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("db is required")
	}
	if userID == 0 {
		return 0, fmt.Errorf("user_id required")
	}
	if revokeAt.IsZero() {
		revokeAt = time.Now().UTC()
	}

	result := db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"token_valid_after": revokeAt,
		"token_version":     gorm.Expr("COALESCE(token_version, 0) + 1"),
	})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to revoke user tokens: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("user not found")
	}

	var user models.User
	if err := db.WithContext(ctx).Select("token_version").First(&user, userID).Error; err != nil {
		return 0, fmt.Errorf("failed to load updated token version: %w", err)
	}
	return user.TokenVersion, nil
}
