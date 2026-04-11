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

type stubOpenAIProvider struct {
	value config.OpenAIConfig
	ok    bool
	err   error
}

func (s stubOpenAIProvider) LoadOpenAIConfig(ctx context.Context) (config.OpenAIConfig, bool, error) {
	return s.value, s.ok, s.err
}

type stubWeKnoraProvider struct {
	value config.WeKnoraConfig
	ok    bool
	err   error
}

func (s stubWeKnoraProvider) LoadWeKnoraConfig(ctx context.Context) (config.WeKnoraConfig, bool, error) {
	return s.value, s.ok, s.err
}

type stubSessionRiskProvider struct {
	value config.SessionRiskPolicyConfig
	ok    bool
	err   error
}

func (s stubSessionRiskProvider) LoadSessionRiskConfig(ctx context.Context) (config.SessionRiskPolicyConfig, bool, error) {
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

	got := resolver.ResolveOpenAI(context.Background(), &config.OpenAIConfig{
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

	got := resolver.ResolveWeKnora(context.Background(), &config.WeKnoraConfig{
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

func TestResolverResolveOpenAIPrecedence(t *testing.T) {
	resolver := NewResolver(
		&config.Config{
			AI: config.AIConfig{
				OpenAI: config.OpenAIConfig{
					APIKey:      "system-key",
					BaseURL:     "https://system.example",
					Model:       "gpt-system",
					Temperature: 0.2,
				},
			},
		},
		WithTenantOpenAIProvider(stubOpenAIProvider{
			ok: true,
			value: config.OpenAIConfig{
				Model:       "gpt-tenant",
				Temperature: 0.3,
			},
		}),
		WithWorkspaceOpenAIProvider(stubOpenAIProvider{
			ok: true,
			value: config.OpenAIConfig{
				Model: "gpt-workspace",
			},
		}),
	)

	got := resolver.ResolveOpenAI(context.Background(), &config.OpenAIConfig{
		Temperature: 0.8,
	})

	if got.APIKey != "system-key" {
		t.Fatalf("api key = %q want system-key", got.APIKey)
	}
	if got.Model != "gpt-workspace" {
		t.Fatalf("model = %q want workspace override", got.Model)
	}
	if got.Temperature != 0.8 {
		t.Fatalf("temperature = %v want 0.8", got.Temperature)
	}
}

func TestResolverResolveWeKnoraPrecedence(t *testing.T) {
	resolver := NewResolver(
		&config.Config{
			WeKnora: config.WeKnoraConfig{
				Enabled:         true,
				BaseURL:         "https://wk-system.example",
				APIKey:          "system-key",
				TenantID:        "tenant-system",
				KnowledgeBaseID: "kb-system",
				MaxRetries:      2,
			},
		},
		WithTenantWeKnoraProvider(stubWeKnoraProvider{
			ok: true,
			value: config.WeKnoraConfig{
				TenantID:        "tenant-tenant",
				KnowledgeBaseID: "kb-tenant",
			},
		}),
		WithWorkspaceWeKnoraProvider(stubWeKnoraProvider{
			ok: true,
			value: config.WeKnoraConfig{
				KnowledgeBaseID: "kb-workspace",
				MaxRetries:      4,
			},
		}),
	)

	got := resolver.ResolveWeKnora(context.Background(), &config.WeKnoraConfig{
		MaxRetries: 6,
	})

	if got.BaseURL != "https://wk-system.example" {
		t.Fatalf("base url = %q want system", got.BaseURL)
	}
	if got.TenantID != "tenant-tenant" {
		t.Fatalf("tenant id = %q want tenant override", got.TenantID)
	}
	if got.KnowledgeBaseID != "kb-workspace" {
		t.Fatalf("kb = %q want workspace override", got.KnowledgeBaseID)
	}
	if got.MaxRetries != 6 {
		t.Fatalf("max retries = %d want 6", got.MaxRetries)
	}
}

func TestResolverResolveSessionRiskPrecedence(t *testing.T) {
	resolver := NewResolver(
		&config.Config{
			Server: config.ServerConfig{
				Environment: "production",
			},
			Security: config.SecurityConfig{
				SessionRisk: config.SessionRiskPolicyConfig{
					HotRefreshWindowMinutes: 15,
					ManySessionsThreshold:   3,
					MediumRiskScore:         2,
					HighRiskScore:           4,
				},
				SessionRiskProfiles: map[string]config.SessionRiskPolicyConfig{
					"production": {
						ManySessionsThreshold: 4,
						HighRiskScore:         6,
					},
				},
			},
		},
		WithTenantSessionRiskProvider(stubSessionRiskProvider{
			ok: true,
			value: config.SessionRiskPolicyConfig{
				ManySessionsThreshold: 5,
			},
		}),
		WithWorkspaceSessionRiskProvider(stubSessionRiskProvider{
			ok: true,
			value: config.SessionRiskPolicyConfig{
				HighRiskScore: 8,
			},
		}),
	)

	got := resolver.ResolveSessionRisk(context.Background(), &config.SessionRiskPolicyConfig{
		HotRefreshWindowMinutes: 30,
	})

	if got.HotRefreshWindowMinutes != 30 {
		t.Fatalf("hot refresh window = %d want runtime override", got.HotRefreshWindowMinutes)
	}
	if got.ManySessionsThreshold != 5 {
		t.Fatalf("many sessions threshold = %d want tenant override", got.ManySessionsThreshold)
	}
	if got.HighRiskScore != 8 {
		t.Fatalf("high risk score = %d want workspace override", got.HighRiskScore)
	}
	if got.MediumRiskScore != 2 {
		t.Fatalf("medium risk score = %d want system default", got.MediumRiskScore)
	}
}
