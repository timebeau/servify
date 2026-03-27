package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"servify/apps/server/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testAuthDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:auth_user_state_policy?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migrate user: %v", err)
	}
	return db
}

func TestUserStateTokenPolicyRejectsInactiveUser(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	db := testAuthDB(t)
	if err := db.Create(&models.User{
		ID:       7,
		Username: "u7",
		Email:    "u7@example.com",
		Status:   "inactive",
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id": 7,
		"iat":     now.Unix(),
		"exp":     now.Add(10 * time.Minute).Unix(),
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
		Policy: NewUserStateTokenPolicy(db),
	}))
	r.GET("/claims", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/claims", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUserStateTokenPolicyRejectsTokenIssuedBeforeUserCutoff(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cutoff := now.Add(-30 * time.Minute)
	db := testAuthDB(t)
	if err := db.Create(&models.User{
		ID:              8,
		Username:        "u8",
		Email:           "u8@example.com",
		Status:          "active",
		TokenValidAfter: &cutoff,
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id": 8,
		"iat":     now.Add(-2 * time.Hour).Unix(),
		"exp":     now.Add(10 * time.Minute).Unix(),
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
		Policy: NewUserStateTokenPolicy(db),
	}))
	r.GET("/claims", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/claims", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUserStateTokenPolicyRejectsStaleTokenVersion(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	db := testAuthDB(t)
	if err := db.Create(&models.User{
		ID:           9,
		Username:     "u9",
		Email:        "u9@example.com",
		Status:       "active",
		TokenVersion: 2,
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id":       9,
		"iat":           now.Unix(),
		"exp":           now.Add(10 * time.Minute).Unix(),
		"token_version": 1,
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
		Policy: NewUserStateTokenPolicy(db),
	}))
	r.GET("/claims", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/claims", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUserStateTokenPolicyAllowsCurrentActiveUser(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cutoff := now.Add(-30 * time.Minute)
	db := testAuthDB(t)
	if err := db.Create(&models.User{
		ID:              10,
		Username:        "u10",
		Email:           "u10@example.com",
		Status:          "active",
		TokenValidAfter: &cutoff,
		TokenVersion:    2,
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id":       10,
		"iat":           now.Unix(),
		"exp":           now.Add(10 * time.Minute).Unix(),
		"token_version": 2,
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
		Policy: NewUserStateTokenPolicy(db),
	}))
	r.GET("/claims", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/claims", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
}
