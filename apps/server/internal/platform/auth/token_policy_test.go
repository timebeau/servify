package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRejectIssuedBefore(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id": 1,
		"iat":     now.Add(-2 * time.Hour).Unix(),
		"exp":     now.Add(10 * time.Minute).Unix(),
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
		Policy: RejectIssuedBefore(now.Add(-1 * time.Hour).Unix()),
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

func TestRequireMinimumTokenVersion(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id":       1,
		"iat":           now.Unix(),
		"exp":           now.Add(10 * time.Minute).Unix(),
		"token_version": 1,
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
		Policy: RequireMinimumTokenVersion(2),
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

func TestComposeTokenPoliciesAllowsCurrentToken(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id":       1,
		"iat":           now.Unix(),
		"exp":           now.Add(10 * time.Minute).Unix(),
		"token_version": 3,
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
		Policy: ComposeTokenPolicies(
			RejectIssuedBefore(now.Add(-1*time.Minute).Unix()),
			RequireMinimumTokenVersion(2),
		),
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
