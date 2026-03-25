package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/services"
)

type stubAIHandlerService struct{}

func (stubAIHandlerService) ProcessQuery(_ context.Context, _ string, _ string) (interface{}, error) {
	return nil, nil
}

func (stubAIHandlerService) GetStatus(_ context.Context) map[string]interface{} {
	return map[string]interface{}{"status": "ok"}
}

func (stubAIHandlerService) GetMetrics() (*services.AIMetrics, bool) {
	return nil, false
}

func (stubAIHandlerService) UploadDocumentToWeKnora(_ context.Context, _, _ string, _ []string) error {
	return nil
}

func (stubAIHandlerService) SyncKnowledgeBase(_ context.Context) error {
	return nil
}

func (stubAIHandlerService) SetWeKnoraEnabled(bool) bool {
	return true
}

func (stubAIHandlerService) ResetCircuitBreaker() bool {
	return true
}

func createTestHS256JWT(t *testing.T, payload map[string]interface{}, secret string) string {
	t.Helper()

	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	enc := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
	h := enc(headerJSON)
	p := enc(payloadJSON)
	signing := h + "." + p

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signing))
	return signing + "." + enc(mac.Sum(nil))
}

func testRouterConfig() *config.Config {
	cfg := config.GetDefaultConfig()
	cfg.JWT.Secret = "test-secret"
	return cfg
}

func TestBuildRouter_AIV1RequiresManagementPrincipal(t *testing.T) {
	router := BuildRouter(Dependencies{
		Config:           testRouterConfig(),
		AIHandlerService: stubAIHandlerService{},
	})

	t.Run("missing token rejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("end user token rejected", func(t *testing.T) {
		token := createTestHS256JWT(t, map[string]interface{}{
			"sub":            "customer-1",
			"principal_kind": "end_user",
		}, "test-secret")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/status", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403 got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("agent token allowed", func(t *testing.T) {
		token := createTestHS256JWT(t, map[string]interface{}{
			"user_id":        7,
			"principal_kind": "agent",
			"roles":          []string{"agent"},
		}, "test-secret")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/status", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), `"success":true`) {
			t.Fatalf("expected success response, body=%s", w.Body.String())
		}
	})
}

func TestBuildRouter_MetricsIngestIsServiceOnly(t *testing.T) {
	router := BuildRouter(Dependencies{
		Config: testRouterConfig(),
	})

	t.Run("agent token rejected", func(t *testing.T) {
		token := createTestHS256JWT(t, map[string]interface{}{
			"user_id":        9,
			"principal_kind": "agent",
			"roles":          []string{"agent"},
		}, "test-secret")

		req := httptest.NewRequest(http.MethodPost, "/api/v1/metrics/ingest", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403 got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("service token allowed through auth layer", func(t *testing.T) {
		token := createTestHS256JWT(t, map[string]interface{}{
			"token_type": "service",
		}, "test-secret")

		req := httptest.NewRequest(http.MethodPost, "/api/v1/metrics/ingest", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 after passing auth, got %d body=%s", w.Code, w.Body.String())
		}
	})
}

func TestBuildRouter_PublicRoutesRemainAnonymous(t *testing.T) {
	router := BuildRouter(Dependencies{
		Config: testRouterConfig(),
	})

	req := httptest.NewRequest(http.MethodGet, "/public/portal/config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
}
