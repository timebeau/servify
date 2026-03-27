package usersecurity

import (
	"context"
	"testing"
	"time"

	"servify/apps/server/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRevokeUserTokens(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:usersecurity_revoke?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := db.Create(&models.User{
		ID:       1,
		Username: "u1",
		Email:    "u1@example.com",
		Status:   "active",
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	revokeAt := time.Unix(1_700_000_000, 0).UTC()
	version, err := RevokeUserTokens(context.Background(), db, 1, revokeAt)
	if err != nil {
		t.Fatalf("RevokeUserTokens() error = %v", err)
	}
	if version != 1 {
		t.Fatalf("version = %d want 1", version)
	}

	var user models.User
	if err := db.First(&user, 1).Error; err != nil {
		t.Fatalf("load user: %v", err)
	}
	if user.TokenVersion != 1 {
		t.Fatalf("token_version = %d want 1", user.TokenVersion)
	}
	if user.TokenValidAfter == nil || !user.TokenValidAfter.Equal(revokeAt) {
		t.Fatalf("token_valid_after = %v want %v", user.TokenValidAfter, revokeAt)
	}
}
