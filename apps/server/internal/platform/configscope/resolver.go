package configscope

import (
	"context"
	"strings"

	"servify/apps/server/internal/config"
)

type PortalConfigProvider interface {
	LoadPortalConfig(ctx context.Context) (config.PortalConfig, bool, error)
}

type Resolver struct {
	system          *config.Config
	tenantPortal    PortalConfigProvider
	workspacePortal PortalConfigProvider
}

func NewResolver(cfg *config.Config, opts ...Option) *Resolver {
	r := &Resolver{system: cfg}
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

type Option func(*Resolver)

func WithTenantPortalProvider(provider PortalConfigProvider) Option {
	return func(r *Resolver) {
		r.tenantPortal = provider
	}
}

func WithWorkspacePortalProvider(provider PortalConfigProvider) Option {
	return func(r *Resolver) {
		r.workspacePortal = provider
	}
}

func (r *Resolver) ResolvePortal(ctx context.Context, runtime *config.PortalConfig) config.PortalConfig {
	var resolved config.PortalConfig
	if r != nil && r.system != nil {
		resolved = r.system.Portal
	}

	if r != nil && r.tenantPortal != nil {
		if next, ok, err := r.tenantPortal.LoadPortalConfig(ctx); err == nil && ok {
			resolved = mergePortalConfig(resolved, next)
		}
	}
	if r != nil && r.workspacePortal != nil {
		if next, ok, err := r.workspacePortal.LoadPortalConfig(ctx); err == nil && ok {
			resolved = mergePortalConfig(resolved, next)
		}
	}
	if runtime != nil {
		resolved = mergePortalConfig(resolved, *runtime)
	}

	return applyPortalDefaults(resolved)
}

func mergePortalConfig(base config.PortalConfig, overlay config.PortalConfig) config.PortalConfig {
	if strings.TrimSpace(overlay.BrandName) != "" {
		base.BrandName = overlay.BrandName
	}
	if strings.TrimSpace(overlay.LogoURL) != "" {
		base.LogoURL = overlay.LogoURL
	}
	if strings.TrimSpace(overlay.PrimaryColor) != "" {
		base.PrimaryColor = overlay.PrimaryColor
	}
	if strings.TrimSpace(overlay.SecondaryColor) != "" {
		base.SecondaryColor = overlay.SecondaryColor
	}
	if strings.TrimSpace(overlay.DefaultLocale) != "" {
		base.DefaultLocale = overlay.DefaultLocale
	}
	if len(overlay.Locales) > 0 {
		base.Locales = append([]string(nil), overlay.Locales...)
	}
	if strings.TrimSpace(overlay.SupportEmail) != "" {
		base.SupportEmail = overlay.SupportEmail
	}
	return base
}

func applyPortalDefaults(p config.PortalConfig) config.PortalConfig {
	if strings.TrimSpace(p.BrandName) == "" {
		p.BrandName = "Servify"
	}
	if strings.TrimSpace(p.PrimaryColor) == "" {
		p.PrimaryColor = "#4299e1"
	}
	if strings.TrimSpace(p.SecondaryColor) == "" {
		p.SecondaryColor = "#764ba2"
	}
	if strings.TrimSpace(p.DefaultLocale) == "" {
		p.DefaultLocale = "zh-CN"
	}
	if len(p.Locales) == 0 {
		p.Locales = []string{"zh-CN", "en-US"}
	}
	return p
}

func (r *Resolver) ResolveOpenAI(runtime *config.OpenAIConfig) config.OpenAIConfig {
	var resolved config.OpenAIConfig
	if r != nil && r.system != nil {
		resolved = r.system.AI.OpenAI
	}
	if runtime != nil {
		resolved = mergeOpenAIConfig(resolved, *runtime)
	}
	return resolved
}

func (r *Resolver) ResolveWeKnora(runtime *config.WeKnoraConfig) config.WeKnoraConfig {
	var resolved config.WeKnoraConfig
	if r != nil && r.system != nil {
		resolved = r.system.WeKnora
	}
	if runtime != nil {
		resolved = mergeWeKnoraConfig(resolved, *runtime)
	}
	return resolved
}

func mergeOpenAIConfig(base config.OpenAIConfig, overlay config.OpenAIConfig) config.OpenAIConfig {
	if strings.TrimSpace(overlay.APIKey) != "" {
		base.APIKey = overlay.APIKey
	}
	if strings.TrimSpace(overlay.BaseURL) != "" {
		base.BaseURL = overlay.BaseURL
	}
	if strings.TrimSpace(overlay.Model) != "" {
		base.Model = overlay.Model
	}
	if overlay.Temperature != 0 {
		base.Temperature = overlay.Temperature
	}
	if overlay.MaxTokens != 0 {
		base.MaxTokens = overlay.MaxTokens
	}
	if overlay.Timeout != 0 {
		base.Timeout = overlay.Timeout
	}
	return base
}

func mergeWeKnoraConfig(base config.WeKnoraConfig, overlay config.WeKnoraConfig) config.WeKnoraConfig {
	if overlay.Enabled {
		base.Enabled = true
	}
	if strings.TrimSpace(overlay.BaseURL) != "" {
		base.BaseURL = overlay.BaseURL
	}
	if strings.TrimSpace(overlay.APIKey) != "" {
		base.APIKey = overlay.APIKey
	}
	if strings.TrimSpace(overlay.TenantID) != "" {
		base.TenantID = overlay.TenantID
	}
	if strings.TrimSpace(overlay.KnowledgeBaseID) != "" {
		base.KnowledgeBaseID = overlay.KnowledgeBaseID
	}
	if overlay.Timeout != 0 {
		base.Timeout = overlay.Timeout
	}
	if overlay.MaxRetries != 0 {
		base.MaxRetries = overlay.MaxRetries
	}
	if overlay.Search.DefaultLimit != 0 {
		base.Search.DefaultLimit = overlay.Search.DefaultLimit
	}
	if overlay.Search.ScoreThreshold != 0 {
		base.Search.ScoreThreshold = overlay.Search.ScoreThreshold
	}
	if strings.TrimSpace(overlay.Search.Strategy) != "" {
		base.Search.Strategy = overlay.Search.Strategy
	}
	if overlay.HealthCheck.Interval != 0 {
		base.HealthCheck.Interval = overlay.HealthCheck.Interval
	}
	if overlay.HealthCheck.Timeout != 0 {
		base.HealthCheck.Timeout = overlay.HealthCheck.Timeout
	}
	return base
}
