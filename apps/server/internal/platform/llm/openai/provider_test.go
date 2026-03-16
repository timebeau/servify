package openai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"servify/apps/server/internal/platform/llm"
)

func TestProviderChat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected auth header: %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"choices":[{"finish_reason":"stop","message":{"content":"hello from model"}}],
			"usage":{"prompt_tokens":11,"completion_tokens":7,"total_tokens":18}
		}`))
	}))
	defer srv.Close()

	provider := NewProvider("test-key", srv.URL)
	resp, err := provider.Chat(context.Background(), llm.ChatRequest{
		Model: "gpt-test",
		Messages: []llm.ChatMessage{
			{Role: "user", Content: "hello"},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Content != "hello from model" {
		t.Fatalf("unexpected content: %s", resp.Content)
	}
	if resp.TokenUsage == nil || resp.TokenUsage.TotalTokens != 18 {
		t.Fatalf("unexpected token usage: %+v", resp.TokenUsage)
	}
}

func TestProviderHealthCheck(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	provider := NewProvider("test-key", srv.URL)
	if err := provider.HealthCheck(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
