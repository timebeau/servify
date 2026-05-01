package embedding

import (
	"testing"

	"servify/apps/server/internal/platform/embedding/openai"
	"servify/apps/server/internal/platform/embedding/tei"
	"servify/apps/server/internal/platform/embedding/xinference"
)

func TestNewProvider_OpenAI(t *testing.T) {
	cfg := FactoryConfig{
		Provider: "openai",
		OpenAI: OpenAIProviderConfig{
			APIKey:  "test-key",
			BaseURL: "https://api.openai.com/v1",
			Model:   "text-embedding-3-small",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}

	if provider == nil {
		t.Fatal("expected provider to be non-nil")
	}

	// Verify it's an openai.Provider
	_, ok := provider.(*openai.Provider)
	if !ok {
		t.Error("expected provider to be *openai.Provider")
	}
}

func TestNewProvider_OpenAI_MissingAPIKey(t *testing.T) {
	cfg := FactoryConfig{
		Provider: "openai",
		OpenAI: OpenAIProviderConfig{
			BaseURL: "https://api.openai.com/v1",
			Model:   "text-embedding-3-small",
		},
	}

	provider, err := NewProvider(cfg)
	if err == nil {
		t.Fatal("expected error for missing api_key, got nil")
	}
	if provider != nil {
		t.Error("expected provider to be nil when error occurs")
	}
}

func TestNewProvider_TEI(t *testing.T) {
	cfg := FactoryConfig{
		Provider: "tei",
		TEI: TEIProviderConfig{
			BaseURL: "http://localhost:8080",
			Model:   "bge-small-zh-v1.5",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}

	if provider == nil {
		t.Fatal("expected provider to be non-nil")
	}

	// Verify it's a tei.Provider
	_, ok := provider.(*tei.Provider)
	if !ok {
		t.Error("expected provider to be *tei.Provider")
	}
}

func TestNewProvider_TEI_MissingBaseURL(t *testing.T) {
	cfg := FactoryConfig{
		Provider: "tei",
		TEI: TEIProviderConfig{
			Model: "bge-small-zh-v1.5",
		},
	}

	provider, err := NewProvider(cfg)
	if err == nil {
		t.Fatal("expected error for missing base_url, got nil")
	}
	if provider != nil {
		t.Error("expected provider to be nil when error occurs")
	}
}

func TestNewProvider_Xinference(t *testing.T) {
	cfg := FactoryConfig{
		Provider: "xinference",
		Xinference: XinferenceProviderConfig{
			BaseURL:  "http://localhost:9997",
			ModelUID: "test-model",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}

	if provider == nil {
		t.Fatal("expected provider to be non-nil")
	}

	// Verify it's a xinference.Provider
	_, ok := provider.(*xinference.Provider)
	if !ok {
		t.Error("expected provider to be *xinference.Provider")
	}
}

func TestNewProvider_Xinference_MissingBaseURL(t *testing.T) {
	cfg := FactoryConfig{
		Provider: "xinference",
		Xinference: XinferenceProviderConfig{
			ModelUID: "test-model",
		},
	}

	provider, err := NewProvider(cfg)
	if err == nil {
		t.Fatal("expected error for missing base_url, got nil")
	}
	if provider != nil {
		t.Error("expected provider to be nil when error occurs")
	}
}

func TestNewProvider_Xinference_MissingModelUID(t *testing.T) {
	cfg := FactoryConfig{
		Provider: "xinference",
		Xinference: XinferenceProviderConfig{
			BaseURL: "http://localhost:9997",
		},
	}

	provider, err := NewProvider(cfg)
	if err == nil {
		t.Fatal("expected error for missing model_uid, got nil")
	}
	if provider != nil {
		t.Error("expected provider to be nil when error occurs")
	}
}

func TestNewProvider_UnknownProvider(t *testing.T) {
	cfg := FactoryConfig{
		Provider: "unknown",
	}

	provider, err := NewProvider(cfg)
	if err == nil {
		t.Fatal("expected error for unknown provider, got nil")
	}
	if provider != nil {
		t.Error("expected provider to be nil when error occurs")
	}
}
