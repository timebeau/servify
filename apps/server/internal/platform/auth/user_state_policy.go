package auth

import (
	"errors"
	"strings"
	"time"

	"servify/apps/server/internal/models"

	"gorm.io/gorm"
)

func NewUserStateTokenPolicy(db *gorm.DB) TokenPolicy {
	if db == nil {
		return nil
	}

	return func(payload map[string]interface{}, claims Claims, now time.Time) error {
		if !claims.HasUserID || claims.UserID == 0 {
			return nil
		}

		var user models.User
		if err := db.WithContext(ContextWithScope(nil, claims.TenantID, claims.WorkspaceID)).
			Select("id", "status", "token_valid_after", "token_version").
			First(&user, claims.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("token user no longer exists")
			}
			return errors.New("failed to evaluate token state")
		}

		status := strings.ToLower(strings.TrimSpace(user.Status))
		if status != "" && status != "active" {
			return errors.New("token user is not active")
		}

		if user.TokenValidAfter != nil && !user.TokenValidAfter.IsZero() {
			if err := RejectIssuedBefore(user.TokenValidAfter.Unix())(payload, claims, now); err != nil {
				return err
			}
		}
		if user.TokenVersion > 0 {
			if err := RequireMinimumTokenVersion(int64(user.TokenVersion))(payload, claims, now); err != nil {
				return err
			}
		}

		return nil
	}
}
