package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestEnforceRequestScope_AgentCannotOverrideScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id":      1,
		"roles":        []string{"agent"},
		"tenant_id":    "tenant-a",
		"workspace_id": "workspace-1",
		"iat":          now.Unix(),
		"exp":          now.Add(10 * time.Minute).Unix(),
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
	}))
	r.Use(EnforceRequestScope())
	r.GET("/scoped", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scoped?tenant_id=tenant-b", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestEnforceRequestScope_ServiceCanNarrowScopeWhenTokenIsUnscoped(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"token_type": "service",
		"iat":        now.Unix(),
		"exp":        now.Add(10 * time.Minute).Unix(),
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
	}))
	r.Use(EnforceRequestScope())
	r.GET("/scoped", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"tenant_id":    c.MustGet("tenant_id"),
			"workspace_id": c.MustGet("workspace_id"),
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scoped?tenant_id=tenant-a&workspace_id=workspace-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if body == "" || !strings.Contains(body, "tenant-a") || !strings.Contains(body, "workspace-1") {
		t.Fatalf("expected projected request scope, body=%s", body)
	}
}

func TestEnforceRequestScope_RejectsConflictingRequestSelectors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"token_type": "service",
		"iat":        now.Unix(),
		"exp":        now.Add(10 * time.Minute).Unix(),
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		Now:    func() time.Time { return now },
	}))
	r.Use(EnforceRequestScope())
	r.GET("/scoped", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scoped?tenant_id=tenant-a", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Tenant-ID", "tenant-b")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}
