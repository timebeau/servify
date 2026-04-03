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
	auditplatform "servify/apps/server/internal/platform/audit"
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

	if err := db.AutoMigrate(&models.User{}, &models.UserAuthSession{}); err != nil {
		t.Fatalf("automigrate user/auth session: %v", err)
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

func TestUserSecurityHandler_RevokeTokensWritesAuditSnapshots(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := newTestDBForUserSecurity(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	if err := db.Create(&models.User{
		ID:           31,
		Username:     "audit-user",
		Email:        "audit-user@example.com",
		Name:         "Audit User",
		Role:         "admin",
		Status:       "active",
		TokenVersion: 2,
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	recorder := &ticketAuditRecorder{}
	h := NewUserSecurityHandler(usersecurity.NewService(db, logger), logger)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(99))
		c.Set("principal_kind", "admin")
		c.Next()
	})
	r.Use(auditplatform.Middleware(recorder))
	r.POST("/api/security/users/:id/revoke-tokens", h.RevokeTokens)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/security/users/31/revoke-tokens", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if len(recorder.entries) != 1 {
		t.Fatalf("expected 1 audit entry got %d", len(recorder.entries))
	}
	entry := recorder.entries[0]
	if entry.BeforeJSON == "" || entry.AfterJSON == "" {
		t.Fatalf("expected before/after snapshots, got before=%q after=%q", entry.BeforeJSON, entry.AfterJSON)
	}
	if !strings.Contains(entry.BeforeJSON, `"token_version":2`) {
		t.Fatalf("unexpected before snapshot: %s", entry.BeforeJSON)
	}
	if !strings.Contains(entry.AfterJSON, `"token_version":3`) {
		t.Fatalf("unexpected after snapshot: %s", entry.AfterJSON)
	}
	if !strings.Contains(entry.AfterJSON, `"user_id":31`) {
		t.Fatalf("unexpected after snapshot user id: %s", entry.AfterJSON)
	}
}

func TestUserSecurityHandler_BatchRevokeTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := newTestDBForUserSecurity(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	users := []models.User{
		{ID: 41, Username: "batch-1", Email: "batch-1@example.com", Name: "Batch 1", Role: "admin", Status: "active", TokenVersion: 1},
		{ID: 42, Username: "batch-2", Email: "batch-2@example.com", Name: "Batch 2", Role: "agent", Status: "active", TokenVersion: 3},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("seed users: %v", err)
	}

	recorder := &ticketAuditRecorder{}
	h := NewUserSecurityHandler(usersecurity.NewService(db, logger), logger)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(99))
		c.Set("principal_kind", "admin")
		c.Next()
	})
	r.Use(auditplatform.Middleware(recorder))
	r.POST("/api/security/users/revoke-tokens", h.BatchRevokeTokens)

	body := strings.NewReader(`{"user_ids":[41,42]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/security/users/revoke-tokens", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v body=%s", err, w.Body.String())
	}
	if got := int(resp["count"].(float64)); got != 2 {
		t.Fatalf("count = %d want 2", got)
	}

	var updated []models.User
	if err := db.Order("id asc").Find(&updated, []uint{41, 42}).Error; err != nil {
		t.Fatalf("reload users: %v", err)
	}
	if updated[0].TokenVersion != 2 || updated[1].TokenVersion != 4 {
		t.Fatalf("unexpected token versions: %+v", updated)
	}

	if len(recorder.entries) != 1 {
		t.Fatalf("expected 1 audit entry got %d", len(recorder.entries))
	}
	entry := recorder.entries[0]
	if !strings.Contains(entry.BeforeJSON, `"user_id":41`) || !strings.Contains(entry.BeforeJSON, `"token_version":1`) {
		t.Fatalf("unexpected before snapshot: %s", entry.BeforeJSON)
	}
	if !strings.Contains(entry.AfterJSON, `"user_id":42`) || !strings.Contains(entry.AfterJSON, `"token_version":4`) {
		t.Fatalf("unexpected after snapshot: %s", entry.AfterJSON)
	}
}

func TestUserSecurityHandler_QueryUsersSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := newTestDBForUserSecurity(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	validAfter := time.Now().UTC().Add(-2 * time.Hour).Round(time.Second)
	lastLogin := time.Now().UTC().Add(-30 * time.Minute).Round(time.Second)
	users := []models.User{
		{ID: 51, Username: "query-1", Email: "query-1@example.com", Name: "Query 1", Role: "admin", Status: "active", TokenVersion: 2, LastLogin: &lastLogin, TokenValidAfter: &validAfter},
		{ID: 52, Username: "query-2", Email: "query-2@example.com", Name: "Query 2", Role: "agent", Status: "inactive", TokenVersion: 0},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("seed users: %v", err)
	}

	h := NewUserSecurityHandler(usersecurity.NewService(db, logger), logger)
	r := gin.New()
	r.POST("/api/security/users/query", h.QueryUsersSecurity)

	t.Run("returns ordered preview items", func(t *testing.T) {
		body := strings.NewReader(`{"user_ids":[52,51,52]}`)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/security/users/query", body)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
		}

		var resp struct {
			Count int `json:"count"`
			Items []struct {
				UserID           int    `json:"user_id"`
				Username         string `json:"username"`
				Status           string `json:"status"`
				TokenVersion     int    `json:"token_version"`
				NextTokenVersion int    `json:"next_token_version"`
			} `json:"items"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v body=%s", err, w.Body.String())
		}
		if resp.Count != 3 {
			t.Fatalf("count = %d want 3", resp.Count)
		}
		if len(resp.Items) != 3 {
			t.Fatalf("len(items) = %d want 3", len(resp.Items))
		}
		if resp.Items[0].UserID != 52 || resp.Items[1].UserID != 51 || resp.Items[2].UserID != 52 {
			t.Fatalf("unexpected item order: %+v", resp.Items)
		}
		if resp.Items[0].TokenVersion != 0 || resp.Items[0].NextTokenVersion != 1 {
			t.Fatalf("unexpected preview for first item: %+v", resp.Items[0])
		}
		if resp.Items[1].Username != "query-1" || resp.Items[1].Status != "active" {
			t.Fatalf("unexpected second item: %+v", resp.Items[1])
		}
	})

	t.Run("returns not found when any user missing", func(t *testing.T) {
		body := strings.NewReader(`{"user_ids":[51,999]}`)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/security/users/query", body)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404 got %d body=%s", w.Code, w.Body.String())
		}
	})
}

func TestUserSecurityHandler_ListAndRevokeSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := newTestDBForUserSecurity(t)
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	if err := db.Create(&models.User{
		ID:       61,
		Username: "session-user",
		Email:    "session-user@example.com",
		Name:     "Session User",
		Role:     "admin",
		Status:   "active",
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Create(&models.UserAuthSession{
		ID:           "auth-session-61",
		UserID:       61,
		Status:       "active",
		TokenVersion: 1,
	}).Error; err != nil {
		t.Fatalf("seed auth session: %v", err)
	}

	recorder := &ticketAuditRecorder{}
	h := NewUserSecurityHandler(usersecurity.NewService(db, logger), logger)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(99))
		c.Set("principal_kind", "admin")
		c.Next()
	})
	r.Use(auditplatform.Middleware(recorder))
	r.GET("/api/security/users/:id/sessions", h.ListUserSessions)
	r.POST("/api/security/users/:id/sessions/revoke", h.RevokeSession)

	t.Run("list sessions", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/security/users/61/sessions", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), `"session_id":"auth-session-61"`) {
			t.Fatalf("unexpected body: %s", w.Body.String())
		}
	})

	t.Run("revoke session writes audit snapshots", func(t *testing.T) {
		body := strings.NewReader(`{"session_id":"auth-session-61"}`)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/security/users/61/sessions/revoke", body)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), `"status":"revoked"`) {
			t.Fatalf("unexpected body: %s", w.Body.String())
		}

		var session models.UserAuthSession
		if err := db.First(&session, "id = ?", "auth-session-61").Error; err != nil {
			t.Fatalf("reload session: %v", err)
		}
		if session.Status != "revoked" || session.TokenVersion != 2 {
			t.Fatalf("unexpected session state: %+v", session)
		}

		if len(recorder.entries) == 0 {
			t.Fatalf("expected audit entries")
		}
		entry := recorder.entries[len(recorder.entries)-1]
		if !strings.Contains(entry.BeforeJSON, `"session_id":"auth-session-61"`) {
			t.Fatalf("unexpected before snapshot: %s", entry.BeforeJSON)
		}
		if !strings.Contains(entry.AfterJSON, `"status":"revoked"`) {
			t.Fatalf("unexpected after snapshot: %s", entry.AfterJSON)
		}
	})
}
