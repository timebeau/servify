package xinference

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestProvider_Embed(t *testing.T) {
	if os.Getenv("XINFERENCE_BASE_URL") == "" {
		t.Skip("XINFERENCE_BASE_URL not set")
	}

	provider := NewProvider(Config{
		BaseURL:  os.Getenv("XINFERENCE_BASE_URL"),
		ModelUID: os.Getenv("XINFERENCE_MODEL_UID"),
	})

	ctx := context.Background()
	vectors, err := provider.Embed(ctx, []string{"hello world"})
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(vectors) != 1 {
		t.Fatalf("expected 1 vector, got %d", len(vectors))
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
	if os.Getenv("XINFERENCE_BASE_URL") == "" {
		t.Skip("XINFERENCE_BASE_URL not set")
	}

	provider := NewProvider(Config{
		BaseURL:  os.Getenv("XINFERENCE_BASE_URL"),
		ModelUID: os.Getenv("XINFERENCE_MODEL_UID"),
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
	provider := NewProvider(Config{})
	if dim := provider.Dimension(); dim != 768 {
		t.Fatalf("expected dimension 768, got %d", dim)
	}
}

func TestProvider_HealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json.RawMessage(`[]`))
	}))
	defer server.Close()

	provider := NewProvider(Config{
		BaseURL: server.URL,
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
