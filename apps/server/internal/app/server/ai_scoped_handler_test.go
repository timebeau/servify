package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/models"
	aidelivery "servify/apps/server/internal/modules/ai/delivery"
	platformauth "servify/apps/server/internal/platform/auth"
	"servify/apps/server/internal/services"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type stubFallbackAIHandler struct{}

func (stubFallbackAIHandler) ProcessQuery(context.Context, string, string) (interface{}, error) {
	return nil, nil
}
func (stubFallbackAIHandler) GetStatus(context.Context) map[string]interface{} {
	return map[string]interface{}{"type": "fallback"}
}
func (stubFallbackAIHandler) GetMetrics() (*services.AIMetrics, bool) {
	return &services.AIMetrics{}, true
}
func (stubFallbackAIHandler) UploadDocumentToWeKnora(context.Context, string, string, []string) error {
	return nil
}
func (stubFallbackAIHandler) SyncKnowledgeBase(context.Context) error { return nil }
func (stubFallbackAIHandler) SetWeKnoraEnabled(bool) bool             { return true }
func (stubFallbackAIHandler) ResetCircuitBreaker() bool               { return true }

func openScopedAITestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&models.TenantConfig{}, &models.WorkspaceConfig{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestScopedAIHandlerServiceProcessQueryUsesWorkspaceOpenAIOverride(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"finish_reason": "stop",
				"message":       map[string]interface{}{"content": "scoped-response"},
			}},
			"usage": map[string]interface{}{"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
		})
	}))
	defer srv.Close()

	db := openScopedAITestDB(t)
	if err := db.Create(&models.WorkspaceConfig{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-1",
		OpenAIJSON:  "api_key: scoped-key\nbase_url: " + srv.URL + "\n",
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	handler := NewScopedAIHandlerService(config.GetDefaultConfig(), logrus.New(), db, stubFallbackAIHandler{})
	ctx := platformauth.ContextWithScope(context.Background(), "tenant-a", "workspace-1")
	resp, err := handler.ProcessQuery(ctx, "hello", "session-1")
	if err != nil {
		t.Fatalf("ProcessQuery() error = %v", err)
	}
	enhanced, ok := resp.(*services.EnhancedAIResponse)
	if !ok {
		t.Fatalf("expected enhanced response, got %T", resp)
	}
	if enhanced.Content != "scoped-response" {
		t.Fatalf("content = %q want scoped-response", enhanced.Content)
	}
}

func TestScopedAIHandlerServiceGetStatusUsesWorkspaceWeKnoraOverride(t *testing.T) {
	db := openScopedAITestDB(t)
	if err := db.Create(&models.WorkspaceConfig{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-1",
		WeKnoraJSON: "enabled: true\nbase_url: http://127.0.0.1:1\napi_key: scoped-key\ntenant_id: tenant-a\nknowledge_base_id: kb-scoped\ntimeout: 1s\n",
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	handler := NewScopedAIHandlerService(config.GetDefaultConfig(), logrus.New(), db, stubFallbackAIHandler{})
	ctx := platformauth.ContextWithScope(context.Background(), "tenant-a", "workspace-1")
	status := handler.GetStatus(ctx)
	if enabled, _ := status["weknora_enabled"].(bool); !enabled {
		t.Fatalf("expected scoped weknora enabled status, got %+v", status)
	}
}

func TestScopedAIHandlerServiceGetMetricsFallsBackToBaseHandler(t *testing.T) {
	handler := NewScopedAIHandlerService(config.GetDefaultConfig(), logrus.New(), openScopedAITestDB(t), stubFallbackAIHandler{})
	metrics, ok := handler.GetMetrics()
	if !ok || metrics == nil {
		t.Fatal("expected fallback metrics")
	}
}

func TestRuntimeServiceFromResolvedConfigWithoutWeKnora(t *testing.T) {
	service := runtimeServiceFromResolvedConfig(config.OpenAIConfig{APIKey: "", BaseURL: ""}, config.DifyConfig{}, config.WeKnoraConfig{}, logrus.New())
	status := service.GetStatus(context.Background())
	if typ, _ := status["type"].(string); typ == "" {
		t.Fatalf("expected service status type, got %+v", status)
	}
}

func TestRuntimeServiceFromResolvedConfigWithWeKnoraScopedKnowledgeBase(t *testing.T) {
	service := runtimeServiceFromResolvedConfig(config.OpenAIConfig{}, config.DifyConfig{}, config.WeKnoraConfig{
		Enabled:         true,
		BaseURL:         "http://127.0.0.1:1",
		APIKey:          "key",
		TenantID:        "tenant-a",
		KnowledgeBaseID: "kb-scoped",
		Timeout:         time.Second,
	}, logrus.New())
	status := service.GetStatus(context.Background())
	if enabled, _ := status["weknora_enabled"].(bool); !enabled {
		t.Fatalf("expected weknora enabled, got %+v", status)
	}
}

var _ aidelivery.HandlerService = (*scopedAIHandlerService)(nil)
