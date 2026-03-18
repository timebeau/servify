package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"servify/apps/server/internal/config"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_RBACRoleExpansion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Unix(1_700_000_000, 0)
	secret := "test-secret"
	token := createTestHS256JWT(t, map[string]interface{}{
		"user_id": 1,
		"roles":   []string{"viewer"},
		"iat":     now.Unix(),
		"exp":     now.Add(10 * time.Minute).Unix(),
	}, secret)

	r := gin.New()
	r.Use(AuthMiddleware(MiddlewareConfig{
		Secret: secret,
		RBAC: config.RBACConfig{
			Enabled: true,
			Roles: map[string][]string{
				"viewer": {"tickets.read"},
			},
		},
		Now: func() time.Time { return now },
	}))
	r.Use(RequireResourcePermission("tickets"))
	r.GET("/tickets", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.POST("/tickets", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tickets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/tickets", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusForbidden {
		t.Fatalf("POST expected 403 got %d body=%s", w2.Code, w2.Body.String())
	}
}
