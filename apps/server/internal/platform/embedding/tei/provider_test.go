package tei

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProvider_Embed_Single(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/embed" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		// 返回模拟的 512 维向量
		response := map[string][][]float32{
			"embeddings": {
				make([]float32, 512),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	provider := NewProvider(Config{BaseURL: server.URL})
	ctx := context.Background()
	vectors, err := provider.Embed(ctx, []string{"test text"})
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(vectors) != 1 {
		t.Fatalf("expected 1 vector, got %d", len(vectors))
	}

	if len(vectors[0]) != 512 {
		t.Fatalf("expected dimension 512, got %d", len(vectors[0]))
	}
}

func TestProvider_Embed_Multiple(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string][][]float32{
			"embeddings": {
				make([]float32, 512),
				make([]float32, 512),
				make([]float32, 512),
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	provider := NewProvider(Config{BaseURL: server.URL})
	ctx := context.Background()
	vectors, err := provider.Embed(ctx, []string{"one", "two", "three"})
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

func TestProvider_Embed_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	provider := NewProvider(Config{BaseURL: server.URL})
	ctx := context.Background()
	_, err := provider.Embed(ctx, []string{"test"})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestProvider_Dimension(t *testing.T) {
	tests := []struct {
		model     string
		dimension int
	}{
		{"bge-small-zh-v1.5", 512},
		{"bge-base-zh-v1.5", 768},
		{"bge-large-zh-v1.5", 1024},
		{"unknown", 512}, // 默认
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	provider := NewProvider(Config{BaseURL: server.URL})
	ctx := context.Background()
	if err := provider.HealthCheck(ctx); err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}
}

func TestProvider_HealthCheck_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	provider := NewProvider(Config{BaseURL: server.URL})
	ctx := context.Background()
	if err := provider.HealthCheck(ctx); err == nil {
		t.Fatal("expected error for failed health check")
	}
}
