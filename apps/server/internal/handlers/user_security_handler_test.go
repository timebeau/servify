//go:build integration
// +build integration

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/platform/usersecurity"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDBForUserSecurity(t *testing.T) *gorm.DB {
	t.Helper()

	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	dsn := "file:user_security_handler_" + name + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("automigrate user: %v", err)
	}

	return db
}

func TestUserSecurityHandler_GetAndRevokeTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := newTestDBForUserSecurity(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	lastLogin := time.Now().UTC().Add(-15 * time.Minute).Round(time.Second)
	if err := db.Create(&models.User{
		ID:           21,
		Username:     "security-user",
		Email:        "security-user@example.com",
		Name:         "Security User",
		Role:         "admin",
		Status:       "active",
		LastLogin:    &lastLogin,
		TokenVersion: 0,
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	h := NewUserSecurityHandler(usersecurity.NewService(db, logger), logger)
	r := gin.New()
	r.GET("/api/security/users/:id", h.GetUserSecurity)
	r.POST("/api/security/users/:id/revoke-tokens", h.RevokeTokens)

	t.Run("get security state", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/security/users/21", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
		}

		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal response: %v body=%s", err, w.Body.String())
		}
		if got := int(body["user_id"].(float64)); got != 21 {
			t.Fatalf("user_id = %d want 21", got)
		}
		if got := body["role"].(string); got != "admin" {
			t.Fatalf("role = %q want admin", got)
		}
		if got := int(body["token_version"].(float64)); got != 0 {
			t.Fatalf("token_version = %d want 0", got)
		}
	})

	t.Run("revoke tokens increments version", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/security/users/21/revoke-tokens", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
		}

		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal response: %v body=%s", err, w.Body.String())
		}
		if got := int(body["token_version"].(float64)); got != 1 {
			t.Fatalf("token_version = %d want 1", got)
		}

		var user models.User
		if err := db.First(&user, 21).Error; err != nil {
			t.Fatalf("reload user: %v", err)
		}
		if user.TokenVersion != 1 {
			t.Fatalf("stored token_version = %d want 1", user.TokenVersion)
		}
		if user.TokenValidAfter == nil || user.TokenValidAfter.IsZero() {
			t.Fatalf("expected token_valid_after to be set")
		}
	})
}

func TestUserSecurityHandler_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	h := NewUserSecurityHandler(nil, logger)
	r := gin.New()
	r.GET("/api/security/users/:id", h.GetUserSecurity)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/security/users/not-a-number", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}
