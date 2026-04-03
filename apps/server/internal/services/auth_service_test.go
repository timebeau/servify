package services

import (
	"context"
	"testing"
	"time"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/models"
	platformauth "servify/apps/server/internal/platform/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newAuthServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:auth_service?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.UserAuthSession{}); err != nil {
		t.Fatalf("migrate models: %v", err)
	}
	return db
}

func testAuthConfig() *config.Config {
	cfg := config.GetDefaultConfig()
	cfg.JWT.Secret = "test-secret"
	cfg.JWT.ExpiresIn = time.Hour
	return cfg
}

func TestAuthServiceLoginAndRefreshRotateSession(t *testing.T) {
	db := newAuthServiceTestDB(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := &models.User{
		ID:           77,
		Username:     "auth-user",
		Email:        "auth-user@example.com",
		Password:     string(hash),
		Status:       "active",
		Role:         "admin",
		TokenVersion: 2,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	svc := NewAuthService(db, testAuthConfig())

	loginResult, err := svc.Login(context.Background(), LoginInput{
		Username: "auth-user",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if loginResult.SessionID == "" {
		t.Fatalf("expected session id in login result")
	}

	payload1, err := (platformauth.Validator{Secret: testAuthConfig().JWT.Secret}).ValidateToken(loginResult.Token)
	if err != nil {
		t.Fatalf("validate login token: %v", err)
	}
	if got := int(payload1["token_version"].(float64)); got != 2 {
		t.Fatalf("token_version = %d want 2", got)
	}
	if got := payload1["session_id"].(string); got != loginResult.SessionID {
		t.Fatalf("session_id = %q want %q", got, loginResult.SessionID)
	}
	if got := int(payload1["session_token_version"].(float64)); got != 0 {
		t.Fatalf("session_token_version = %d want 0", got)
	}

	refreshResult, err := svc.RefreshToken(context.Background(), user.ID, loginResult.SessionID)
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}
	if refreshResult.SessionID != loginResult.SessionID {
		t.Fatalf("refresh session id = %q want %q", refreshResult.SessionID, loginResult.SessionID)
	}

	payload2, err := (platformauth.Validator{Secret: testAuthConfig().JWT.Secret}).ValidateToken(refreshResult.Token)
	if err != nil {
		t.Fatalf("validate refresh token: %v", err)
	}
	if got := int(payload2["session_token_version"].(float64)); got != 1 {
		t.Fatalf("session_token_version = %d want 1", got)
	}

	var session models.UserAuthSession
	if err := db.First(&session, "id = ?", loginResult.SessionID).Error; err != nil {
		t.Fatalf("load session: %v", err)
	}
	if session.TokenVersion != 1 {
		t.Fatalf("stored session token_version = %d want 1", session.TokenVersion)
	}
	if session.LastRefreshedAt == nil || session.LastRefreshedAt.IsZero() {
		t.Fatalf("expected last_refreshed_at to be set")
	}
}
