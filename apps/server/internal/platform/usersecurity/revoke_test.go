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

func TestServiceGetUsers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:usersecurity_get_users?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	users := []models.User{
		{ID: 11, Username: "u11", Email: "u11@example.com", Status: "active", Role: "admin", TokenVersion: 2},
		{ID: 12, Username: "u12", Email: "u12@example.com", Status: "inactive", Role: "agent", TokenVersion: 0},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("seed users: %v", err)
	}

	svc := NewService(db, nil)

	got, err := svc.GetUsers(context.Background(), []uint{12, 11, 12})
	if err != nil {
		t.Fatalf("GetUsers() error = %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len(GetUsers()) = %d want 3", len(got))
	}
	if got[0].ID != 12 || got[1].ID != 11 || got[2].ID != 12 {
		t.Fatalf("unexpected order: %+v", got)
	}

	if _, err := svc.GetUsers(context.Background(), []uint{11, 99}); err == nil {
		t.Fatalf("expected missing user error")
	}
}

func TestServiceListAndRevokeSession(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:usersecurity_sessions?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.UserAuthSession{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := db.Create(&models.User{
		ID:       21,
		Username: "u21",
		Email:    "u21@example.com",
		Status:   "active",
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.UserAuthSession{
		ID:           "auth-session-test",
		UserID:       21,
		Status:       "active",
		TokenVersion: 1,
	}).Error; err != nil {
		t.Fatalf("seed auth session: %v", err)
	}

	svc := NewService(db, nil)

	sessions, err := svc.ListUserSessions(context.Background(), 21)
	if err != nil {
		t.Fatalf("ListUserSessions() error = %v", err)
	}
	if len(sessions) != 1 || sessions[0].ID != "auth-session-test" {
		t.Fatalf("unexpected sessions: %+v", sessions)
	}

	session, err := svc.RevokeSession(context.Background(), 21, "auth-session-test")
	if err != nil {
		t.Fatalf("RevokeSession() error = %v", err)
	}
	if session.Status != "revoked" {
		t.Fatalf("status = %q want revoked", session.Status)
	}
	if session.TokenVersion != 2 {
		t.Fatalf("token_version = %d want 2", session.TokenVersion)
	}
	if session.RevokedAt == nil || session.RevokedAt.IsZero() {
		t.Fatalf("expected revoked_at to be set")
	}
}
