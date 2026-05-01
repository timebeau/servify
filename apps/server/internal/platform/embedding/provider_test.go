package embedding

import (
	"context"
	"testing"
)

// MockProvider is a mock implementation for testing
type MockProvider struct {
	embedFunc    func(ctx context.Context, texts []string) ([][]float32, error)
	dimensionVal int
}

func (m *MockProvider) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if m.embedFunc != nil {
		return m.embedFunc(ctx, texts)
	}
	// Return zero vectors with fixed dimension
	result := make([][]float32, len(texts))
	for i := range result {
		result[i] = make([]float32, m.Dimension())
	}
	return result, nil
}

func (m *MockProvider) Dimension() int {
	if m.dimensionVal > 0 {
		return m.dimensionVal
	}
	return 1536
}

func (m *MockProvider) HealthCheck(ctx context.Context) error {
	return nil
}

func TestMockProvider(t *testing.T) {
	ctx := context.Background()
	provider := &MockProvider{dimensionVal: 512}

	// Test Embed
	vectors, err := provider.Embed(ctx, []string{"test"})
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}
	if len(vectors) != 1 {
		t.Fatalf("expected 1 vector, got %d", len(vectors))
	}
	if len(vectors[0]) != 512 {
		t.Fatalf("expected dimension 512, got %d", len(vectors[0]))
	}

	// Test Dimension
	if provider.Dimension() != 512 {
		t.Fatalf("expected dimension 512, got %d", provider.Dimension())
	}

	// Test HealthCheck
	if err := provider.HealthCheck(ctx); err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}
}

func TestMockProviderMultipleTexts(t *testing.T) {
	ctx := context.Background()
	provider := &MockProvider{dimensionVal: 768}

	texts := []string{"first", "second", "third"}
	vectors, err := provider.Embed(ctx, texts)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}
	if len(vectors) != len(texts) {
		t.Fatalf("expected %d vectors, got %d", len(texts), len(vectors))
	}
	for i, v := range vectors {
		if len(v) != 768 {
			t.Fatalf("vector %d: expected dimension 768, got %d", i, len(v))
		}
	}
}

func TestMockProviderDefaultDimension(t *testing.T) {
	provider := &MockProvider{}
	if provider.Dimension() != 1536 {
		t.Fatalf("expected default dimension 1536, got %d", provider.Dimension())
	}
}

func TestMockProviderCustomEmbedFunc(t *testing.T) {
	ctx := context.Background()
	called := false
	provider := &MockProvider{
		dimensionVal: 256,
		embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
			called = true
			result := make([][]float32, len(texts))
			for i := range result {
				result[i] = make([]float32, 256)
				result[i][0] = 1.0
			}
			return result, nil
		},
	}

	vectors, err := provider.Embed(ctx, []string{"test"})
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}
	if !called {
		t.Fatalf("custom embedFunc was not called")
	}
	if vectors[0][0] != 1.0 {
		t.Fatalf("expected first element to be 1.0, got %f", vectors[0][0])
	}
}
