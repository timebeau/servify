package bootstrap

import (
	"fmt"
	"strings"

	"servify/apps/server/internal/config"

	"github.com/sirupsen/logrus"
)

const defaultJWTSecret = "default-secret-key"

func SecurityWarnings(cfg *config.Config) []string {
	if cfg == nil {
		return []string{"config is nil; security defaults may be unsafe"}
	}

	var warnings []string

	if strings.TrimSpace(cfg.JWT.Secret) == "" || strings.TrimSpace(cfg.JWT.Secret) == defaultJWTSecret {
		warnings = append(warnings, "jwt.secret is empty or using the default value")
	}
	if cfg.Security.CORS.Enabled && len(cfg.Security.CORS.AllowedOrigins) == 1 && strings.TrimSpace(cfg.Security.CORS.AllowedOrigins[0]) == "*" {
		warnings = append(warnings, "security.cors.allowed_origins allows all origins")
	}
	if !cfg.Security.RateLimiting.Enabled {
		warnings = append(warnings, "security.rate_limiting is disabled")
	} else {
		if !hasRateLimitPrefix(cfg.Security.RateLimiting.Paths, "/public/") {
			warnings = append(warnings, "security.rate_limiting.paths has no dedicated limit for /public/")
		}
		if !hasRateLimitPrefix(cfg.Security.RateLimiting.Paths, "/api/v1/ws") {
			warnings = append(warnings, "security.rate_limiting.paths has no dedicated limit for /api/v1/ws")
		}
	}
	if strings.TrimSpace(cfg.AI.OpenAI.APIKey) == "" {
		warnings = append(warnings, "ai.openai.api_key is empty")
	}
	if cfg.WeKnora.Enabled && strings.TrimSpace(cfg.WeKnora.APIKey) == "" {
		warnings = append(warnings, "weknora is enabled but weknora.api_key is empty")
	}

	return warnings
}

func hasRateLimitPrefix(paths []config.PathRateLimitConfig, prefix string) bool {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return false
	}
	for _, path := range paths {
		if !path.Enabled {
			continue
		}
		if strings.TrimSpace(path.Prefix) == prefix {
			return true
		}
	}
	return false
}

func LogSecurityWarnings(logger *logrus.Logger, cfg *config.Config) {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	for _, warning := range SecurityWarnings(cfg) {
		logger.Warn(fmt.Sprintf("security baseline warning: %s", warning))
	}
}
