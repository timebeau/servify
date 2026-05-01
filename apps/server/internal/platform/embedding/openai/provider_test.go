package openai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestProvider_Embed(t *testing.T) {
	// 跳过测试如果没有 API key
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	provider := NewProvider(Config{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  "text-embedding-3-small",
	})

	ctx := context.Background()
	vectors, err := provider.Embed(ctx, []string{"hello world"})
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(vectors) != 1 {
		t.Fatalf("expected 1 vector, got %d", len(vectors))
	}

	if len(vectors[0]) != 1536 {
		t.Fatalf("expected dimension 1536, got %d", len(vectors[0]))
	}

	// 检查向量值不为零
	hasNonZero := false
	for _, v := range vectors[0] {
		if v != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Fatal("embedding vector is all zeros")
	}
}

func TestProvider_Embed_Multiple(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set")
	}

	provider := NewProvider(Config{
		APIKey: os.Getenv("OPENAI_API_KEY"),
	})

	ctx := context.Background()
	vectors, err := provider.Embed(ctx, []string{"first", "second", "third"})
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(vectors) != 3 {
		t.Fatalf("expected 3 vectors, got %d", len(vectors))
	}
}

func TestProvider_Embed_EmptyInput(t *testing.T) {
	provider := NewProvider(Config{})
	ctx := context.Background()
	_, err := provider.Embed(ctx, []string{})
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestProvider_Dimension(t *testing.T) {
	tests := []struct {
		model     string
		dimension int
	}{
		{"text-embedding-3-small", 1536},
		{"text-embedding-3-large", 3072},
		{"text-embedding-ada-002", 1536},
		{"unknown-model", 1536}, // 默认
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			provider := NewProvider(Config{Model: tt.model})
			if dim := provider.Dimension(); dim != tt.dimension {
				t.Fatalf("expected dimension %d, got %d", tt.dimension, dim)
			}
		})
	}
}

func TestProvider_HealthCheck(t *testing.T) {
	// 使用 mock server 测试
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	provider := NewProvider(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	})

	ctx := context.Background()
	if err := provider.HealthCheck(ctx); err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}
}

func TestProvider_HealthCheck_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	provider := NewProvider(Config{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	if err := provider.HealthCheck(ctx); err == nil {
		t.Fatal("expected error for failed health check")
	}
}
