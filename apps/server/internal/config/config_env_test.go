package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestLoad_ExpandsEnvironmentPlaceholders(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	t.Setenv("OPENAI_API_KEY", "resolved-openai-key")
	t.Setenv("SERVIFY_JWT_SECRET", "resolved-jwt-secret")

	viper.Set("ai.openai.api_key", "${OPENAI_API_KEY}")
	viper.Set("jwt.secret", "${SERVIFY_JWT_SECRET}")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AI.OpenAI.APIKey != "resolved-openai-key" {
		t.Fatalf("expected expanded openai key, got %q", cfg.AI.OpenAI.APIKey)
	}
	if cfg.JWT.Secret != "resolved-jwt-secret" {
		t.Fatalf("expected expanded jwt secret, got %q", cfg.JWT.Secret)
	}
}

func TestLoad_ClearsUnsetEnvironmentPlaceholders(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	if err := os.Unsetenv("OPENAI_API_KEY"); err != nil {
		t.Fatalf("Unsetenv() error = %v", err)
	}

	viper.Set("ai.openai.api_key", "${OPENAI_API_KEY}")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AI.OpenAI.APIKey != "" {
		t.Fatalf("expected unresolved placeholder to be cleared, got %q", cfg.AI.OpenAI.APIKey)
	}
}
