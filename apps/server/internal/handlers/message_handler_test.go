package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubMessageRouterRuntime struct {
	stats map[string]interface{}
}

func (s *stubMessageRouterRuntime) Start() error {
	return nil
}

func (s *stubMessageRouterRuntime) Stop() error {
	return nil
}

func (s *stubMessageRouterRuntime) GetPlatformStats() map[string]interface{} {
	return s.stats
}

func TestMessageHandlerGetPlatformStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewMessageHandler(&stubMessageRouterRuntime{
		stats: map[string]interface{}{
			"total_platforms":  2,
			"active_platforms": []string{"telegram", "wechat"},
		},
	})

	router := gin.New()
	router.GET("/api/v1/messages/platforms", handler.GetPlatformStats)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/messages/platforms", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var got struct {
		Success bool `json:"success"`
		Data    struct {
			TotalPlatforms int      `json:"total_platforms"`
			ActivePlatforms []string `json:"active_platforms"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if !got.Success {
		t.Fatalf("expected success=true, body=%s", w.Body.String())
	}
	if got.Data.TotalPlatforms != 2 {
		t.Fatalf("expected total_platforms=2, got %d", got.Data.TotalPlatforms)
	}
	if len(got.Data.ActivePlatforms) != 2 {
		t.Fatalf("expected 2 active platforms, got %+v", got.Data.ActivePlatforms)
	}
}
