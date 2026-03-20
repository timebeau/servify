package aiprovider

import (
	"testing"
	"time"
)

func TestMatrixEnabledSortsByPriority(t *testing.T) {
	matrix := Matrix{
		Providers: []ProviderDescriptor{
			{ID: "openai-secondary", Kind: KindLLM, Enabled: true, Fallback: FallbackPolicy{Priority: 2}},
			{ID: "weknora", Kind: KindKnowledge, Enabled: true, Fallback: FallbackPolicy{Priority: 1}},
			{ID: "openai-primary", Kind: KindLLM, Enabled: true, Fallback: FallbackPolicy{Priority: 1}},
			{ID: "disabled", Kind: KindLLM, Enabled: false, Fallback: FallbackPolicy{Priority: 0}},
		},
	}

	enabled := matrix.Enabled(KindLLM)
	if len(enabled) != 2 {
		t.Fatalf("expected 2 enabled llm providers, got %d", len(enabled))
	}
	if enabled[0].ID != "openai-primary" || enabled[1].ID != "openai-secondary" {
		t.Fatalf("unexpected provider order: %#v", enabled)
	}
}

func TestMatrixFindReturnsDescriptor(t *testing.T) {
	matrix := Matrix{
		Providers: []ProviderDescriptor{
			{
				ID:      "weknora",
				Kind:    KindKnowledge,
				Driver:  "weknora",
				Enabled: true,
				Fallback: FallbackPolicy{
					Priority: 1,
					CircuitBreaker: CircuitBreakerPolicy{
						Enabled:      true,
						MaxFailures:  3,
						ResetTimeout: time.Minute,
					},
				},
			},
		},
	}

	provider, ok := matrix.Find("weknora")
	if !ok {
		t.Fatal("expected provider to be found")
	}
	if provider.Fallback.CircuitBreaker.MaxFailures != 3 {
		t.Fatalf("expected circuit breaker policy to be preserved, got %#v", provider.Fallback)
	}
}
