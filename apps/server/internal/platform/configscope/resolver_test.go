package configscope

import (
	"context"
	"testing"

	"servify/apps/server/internal/config"
)

type stubPortalProvider struct {
	value config.PortalConfig
	ok    bool
	err   error
}

func (s stubPortalProvider) LoadPortalConfig(ctx context.Context) (config.PortalConfig, bool, error) {
	return s.value, s.ok, s.err
}

func TestResolverResolvePortalPrecedence(t *testing.T) {
	resolver := NewResolver(
		&config.Config{
			Portal: config.PortalConfig{
				BrandName:      "system",
				PrimaryColor:   "#111111",
				SecondaryColor: "#222222",
				DefaultLocale:  "zh-CN",
				Locales:        []string{"zh-CN"},
			},
		},
		WithTenantPortalProvider(stubPortalProvider{
			ok: true,
			value: config.PortalConfig{
				BrandName:    "tenant",
				PrimaryColor: "#333333",
			},
		}),
		WithWorkspacePortalProvider(stubPortalProvider{
			ok: true,
			value: config.PortalConfig{
				BrandName:     "workspace",
				DefaultLocale: "en-US",
			},
		}),
	)

	got := resolver.ResolvePortal(context.Background(), &config.PortalConfig{
		BrandName:      "runtime",
		SecondaryColor: "#444444",
	})

	if got.BrandName != "runtime" {
		t.Fatalf("brand = %q want runtime", got.BrandName)
	}
	if got.PrimaryColor != "#333333" {
		t.Fatalf("primary color = %q want tenant override", got.PrimaryColor)
	}
	if got.SecondaryColor != "#444444" {
		t.Fatalf("secondary color = %q want runtime override", got.SecondaryColor)
	}
	if got.DefaultLocale != "en-US" {
		t.Fatalf("default locale = %q want workspace override", got.DefaultLocale)
	}
}

func TestResolverResolvePortalDefaults(t *testing.T) {
	resolver := NewResolver(nil)
	got := resolver.ResolvePortal(context.Background(), nil)

	if got.BrandName != "Servify" {
		t.Fatalf("brand = %q want Servify", got.BrandName)
	}
	if got.PrimaryColor != "#4299e1" {
		t.Fatalf("primary color = %q want default", got.PrimaryColor)
	}
	if got.DefaultLocale != "zh-CN" {
		t.Fatalf("default locale = %q want zh-CN", got.DefaultLocale)
	}
	if len(got.Locales) != 2 {
		t.Fatalf("locales = %+v want defaults", got.Locales)
	}
}

func TestResolverResolveOpenAI(t *testing.T) {
	resolver := NewResolver(&config.Config{
		AI: config.AIConfig{
			OpenAI: config.OpenAIConfig{
				APIKey:      "system-key",
				BaseURL:     "https://system.example",
				Model:       "gpt-system",
				Temperature: 0.2,
				MaxTokens:   100,
			},
		},
	})

	got := resolver.ResolveOpenAI(&config.OpenAIConfig{
		Model:       "gpt-runtime",
		Temperature: 0.7,
	})

	if got.APIKey != "system-key" {
		t.Fatalf("api key = %q want system-key", got.APIKey)
	}
	if got.Model != "gpt-runtime" {
		t.Fatalf("model = %q want gpt-runtime", got.Model)
	}
	if got.Temperature != 0.7 {
		t.Fatalf("temperature = %v want 0.7", got.Temperature)
	}
}

func TestResolverResolveWeKnora(t *testing.T) {
	resolver := NewResolver(&config.Config{
		WeKnora: config.WeKnoraConfig{
			Enabled:         true,
			BaseURL:         "https://wk-system.example",
			APIKey:          "system-key",
			TenantID:        "tenant-a",
			KnowledgeBaseID: "kb-a",
			MaxRetries:      2,
		},
	})

	got := resolver.ResolveWeKnora(&config.WeKnoraConfig{
		KnowledgeBaseID: "kb-runtime",
		MaxRetries:      5,
	})

	if got.BaseURL != "https://wk-system.example" {
		t.Fatalf("base url = %q want system", got.BaseURL)
	}
	if got.KnowledgeBaseID != "kb-runtime" {
		t.Fatalf("kb = %q want kb-runtime", got.KnowledgeBaseID)
	}
	if got.MaxRetries != 5 {
		t.Fatalf("max retries = %d want 5", got.MaxRetries)
	}
}
