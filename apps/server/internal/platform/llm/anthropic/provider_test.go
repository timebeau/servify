package anthropic

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"servify/apps/server/internal/platform/llm"
)

func TestProviderChat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("x-api-key"); got != "test-key" {
			t.Fatalf("unexpected api key header: %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"model":"claude-test",
			"stop_reason":"end_turn",
			"content":[
				{"type":"text","text":"hello from anthropic"},
				{"type":"tool_use","id":"tool-1","name":"lookup","input":{"query":"hello"}}
			],
			"usage":{"input_tokens":10,"output_tokens":5}
		}`))
	}))
	defer srv.Close()

	provider := NewProvider("test-key", srv.URL)
	resp, err := provider.Chat(context.Background(), llm.ChatRequest{
		Messages: []llm.ChatMessage{
			{Role: "system", Content: "You are helpful"},
			{Role: "user", Content: "hello"},
		},
		Tools: []llm.ToolDefinition{
			{Name: "lookup", Description: "lookup docs"},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Provider != "anthropic" || resp.Content != "hello from anthropic" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(resp.ToolCalls) != 1 || resp.ToolCalls[0].Name != "lookup" {
		t.Fatalf("expected tool call, got %+v", resp.ToolCalls)
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
