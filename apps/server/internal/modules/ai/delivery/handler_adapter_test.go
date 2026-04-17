package delivery

import (
	"context"
	"errors"
	"testing"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/services"
)

type stubLegacyAIHandlerService struct {
	processQuery            func(ctx context.Context, query string, sessionID string) (*services.AIResponse, error)
	getStatus               func(ctx context.Context) map[string]interface{}
}

type stubEnhancedAIHandlerService struct {
	stubLegacyAIHandlerService
	processQueryEnhanced         func(ctx context.Context, query string, sessionID string) (*services.EnhancedAIResponse, error)
	uploadKnowledgeDocument      func(ctx context.Context, title, content string, tags []string) error
	getMetrics                   func() *services.AIMetrics
	setKnowledgeProviderEnabled  func(enabled bool)
	resetCircuitBreaker          func()
	syncKnowledgeBase            func(ctx context.Context) error
}

func (s stubLegacyAIHandlerService) ProcessQuery(ctx context.Context, query string, sessionID string) (*services.AIResponse, error) {
	return s.processQuery(ctx, query, sessionID)
}

func (s stubLegacyAIHandlerService) ShouldTransferToHuman(query string, sessionHistory []models.Message) bool {
	return false
}

func (s stubLegacyAIHandlerService) GetSessionSummary(messages []models.Message) (string, error) {
	return "", nil
}

func (s stubLegacyAIHandlerService) InitializeKnowledgeBase() {}

func (s stubLegacyAIHandlerService) GetStatus(ctx context.Context) map[string]interface{} {
	if s.getStatus == nil {
		return nil
	}
	return s.getStatus(ctx)
}

func (s stubEnhancedAIHandlerService) ProcessQueryEnhanced(ctx context.Context, query string, sessionID string) (*services.EnhancedAIResponse, error) {
	return s.processQueryEnhanced(ctx, query, sessionID)
}

func (s stubEnhancedAIHandlerService) UploadKnowledgeDocument(ctx context.Context, title, content string, tags []string) error {
	return s.uploadKnowledgeDocument(ctx, title, content, tags)
}

func (s stubEnhancedAIHandlerService) GetMetrics() *services.AIMetrics {
	return s.getMetrics()
}

func (s stubEnhancedAIHandlerService) SetKnowledgeProviderEnabled(enabled bool) {
	s.setKnowledgeProviderEnabled(enabled)
}

func (s stubEnhancedAIHandlerService) SetFallbackEnabled(enabled bool) {}

func (s stubEnhancedAIHandlerService) ResetCircuitBreaker() {
	s.resetCircuitBreaker()
}

func (s stubEnhancedAIHandlerService) SyncKnowledgeBase(ctx context.Context) error {
	return s.syncKnowledgeBase(ctx)
}

func TestHandlerServiceAdapter_UsesEnhancedSurfaceWhenAvailable(t *testing.T) {
	var enabled bool
	var reset bool
	adapter := NewHandlerServiceAdapter(stubEnhancedAIHandlerService{
		stubLegacyAIHandlerService: stubLegacyAIHandlerService{
			processQuery: func(ctx context.Context, query string, sessionID string) (*services.AIResponse, error) {
				t.Fatal("expected enhanced path")
				return nil, nil
			},
			getStatus: func(ctx context.Context) map[string]interface{} {
				return map[string]interface{}{"type": "enhanced"}
			},
		},
		processQueryEnhanced: func(ctx context.Context, query string, sessionID string) (*services.EnhancedAIResponse, error) {
			return &services.EnhancedAIResponse{
				AIResponse: &services.AIResponse{Content: "enhanced", Source: "ai", Confidence: 0.9},
				Strategy:   "weknora",
			}, nil
		},
		uploadKnowledgeDocument: func(ctx context.Context, title, content string, tags []string) error { return nil },
		getMetrics: func() *services.AIMetrics { return &services.AIMetrics{SuccessCount: 3} },
		setKnowledgeProviderEnabled: func(v bool) { enabled = v },
		resetCircuitBreaker: func() { reset = true },
		syncKnowledgeBase: func(ctx context.Context) error { return nil },
	})

	got, err := adapter.ProcessQuery(context.Background(), "hi", "s1")
	if err != nil {
		t.Fatalf("ProcessQuery() err=%v", err)
	}
	resp, ok := got.(*services.EnhancedAIResponse)
	if !ok || resp.AIResponse.Content != "enhanced" {
		t.Fatalf("ProcessQuery() got=%T %+v", got, got)
	}

	if status := adapter.GetStatus(context.Background()); status["type"] != "enhanced" {
		t.Fatalf("GetStatus() = %+v", status)
	}

	if metrics, ok := adapter.GetMetrics(); !ok || metrics.SuccessCount != 3 {
		t.Fatalf("GetMetrics() = %+v, %v", metrics, ok)
	}

	if err := adapter.UploadKnowledgeDocument(context.Background(), "t", "c", []string{"tag"}); err != nil {
		t.Fatalf("UploadKnowledgeDocument() err=%v", err)
	}
	if err := adapter.SyncKnowledgeBase(context.Background()); err != nil {
		t.Fatalf("SyncKnowledgeBase() err=%v", err)
	}
	if !adapter.SetKnowledgeProviderEnabled(true) || !enabled {
		t.Fatal("SetKnowledgeProviderEnabled() did not delegate")
	}
	if !adapter.ResetCircuitBreaker() || !reset {
		t.Fatal("ResetCircuitBreaker() did not delegate")
	}
}

func TestHandlerServiceAdapter_UsesBaseSurfaceForStandardService(t *testing.T) {
	adapter := NewHandlerServiceAdapter(stubLegacyAIHandlerService{
		processQuery: func(ctx context.Context, query string, sessionID string) (*services.AIResponse, error) {
			return &services.AIResponse{Content: "base", Source: "ai", Confidence: 0.5}, nil
		},
		getStatus: func(ctx context.Context) map[string]interface{} {
			return map[string]interface{}{"type": "base"}
		},
	})

	got, err := adapter.ProcessQuery(context.Background(), "hi", "s1")
	if err != nil {
		t.Fatalf("ProcessQuery() err=%v", err)
	}
	resp, ok := got.(*services.AIResponse)
	if !ok || resp.Content != "base" {
		t.Fatalf("ProcessQuery() got=%T %+v", got, got)
	}

	if metrics, ok := adapter.GetMetrics(); ok || metrics != nil {
		t.Fatalf("GetMetrics() = %+v, %v", metrics, ok)
	}
	if err := adapter.UploadKnowledgeDocument(context.Background(), "t", "c", nil); err == nil {
		t.Fatal("UploadKnowledgeDocument() expected error")
	}
	if err := adapter.SyncKnowledgeBase(context.Background()); err == nil {
		t.Fatal("SyncKnowledgeBase() expected error")
	}
	if adapter.SetKnowledgeProviderEnabled(true) {
		t.Fatal("SetKnowledgeProviderEnabled() = true, want false")
	}
	if adapter.ResetCircuitBreaker() {
		t.Fatal("ResetCircuitBreaker() = true, want false")
	}
}

func TestHandlerServiceAdapter_ReturnsConfiguredErrorForNilService(t *testing.T) {
	adapter := NewHandlerServiceAdapter(nil)
	if _, err := adapter.ProcessQuery(context.Background(), "hi", "s1"); err == nil {
		t.Fatal("ProcessQuery() expected error")
	}
	if err := adapter.UploadKnowledgeDocument(context.Background(), "t", "c", nil); err == nil {
		t.Fatal("UploadKnowledgeDocument() expected error")
	}
	if err := adapter.SyncKnowledgeBase(context.Background()); err == nil {
		t.Fatal("SyncKnowledgeBase() expected error")
	}
}

func TestHandlerServiceAdapter_PropagatesEnhancedErrors(t *testing.T) {
	expectedErr := errors.New("boom")
	adapter := NewHandlerServiceAdapter(stubEnhancedAIHandlerService{
		stubLegacyAIHandlerService: stubLegacyAIHandlerService{
			processQuery: func(ctx context.Context, query string, sessionID string) (*services.AIResponse, error) {
				return nil, expectedErr
			},
		},
		processQueryEnhanced: func(ctx context.Context, query string, sessionID string) (*services.EnhancedAIResponse, error) {
			return nil, expectedErr
		},
		getMetrics: func() *services.AIMetrics { return nil },
		uploadKnowledgeDocument: func(ctx context.Context, title, content string, tags []string) error { return expectedErr },
		setKnowledgeProviderEnabled: func(enabled bool) {},
		resetCircuitBreaker: func() {},
		syncKnowledgeBase: func(ctx context.Context) error { return expectedErr },
	})

	if _, err := adapter.ProcessQuery(context.Background(), "hi", "s1"); !errors.Is(err, expectedErr) {
		t.Fatalf("ProcessQuery() err=%v", err)
	}
	if err := adapter.UploadKnowledgeDocument(context.Background(), "t", "c", nil); !errors.Is(err, expectedErr) {
		t.Fatalf("UploadKnowledgeDocument() err=%v", err)
	}
	if err := adapter.SyncKnowledgeBase(context.Background()); !errors.Is(err, expectedErr) {
		t.Fatalf("SyncKnowledgeBase() err=%v", err)
	}
}
