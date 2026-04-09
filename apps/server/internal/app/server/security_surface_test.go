package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"servify/apps/server/internal/config"

	"github.com/gin-gonic/gin"
)

func TestRouteSecurityWarnings_BuildRouterCovered(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.JWT.Secret = "test-secret"

	router := BuildRouter(Dependencies{
		Config:           cfg,
		AIHandlerService: stubAIHandlerService{},
	})

	warnings := RouteSecurityWarnings(router.Routes(), cfg)
	if len(warnings) != 0 {
		t.Fatalf("expected no route security warnings, got %v", warnings)
	}
}

func TestRouteSecurityWarnings_MissingCatalogEntry(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.JWT.Secret = "test-secret"

	warnings := RouteSecurityWarnings(gin.RoutesInfo{
		{Method: http.MethodGet, Path: "/public/raw-export"},
	}, cfg)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning got %d (%v)", len(warnings), warnings)
	}
	if !strings.Contains(warnings[0], "GET /public/raw-export") {
		t.Fatalf("unexpected warning %q", warnings[0])
	}
}

func TestBuildRouter_AuthPublicRoutesRemainAnonymous(t *testing.T) {
	router := BuildRouter(Dependencies{
		Config: testRouterConfig(),
	})

	for _, path := range []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
	} {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusUnauthorized || w.Code == http.StatusForbidden {
			t.Fatalf("expected %s to remain anonymous, got %d body=%s", path, w.Code, w.Body.String())
		}
	}
}
