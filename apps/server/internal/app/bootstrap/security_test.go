package bootstrap

import (
	"bytes"
	"strings"
	"testing"

	"servify/apps/server/internal/config"

	"github.com/sirupsen/logrus"
)

func TestSecurityWarnings(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Security.RateLimiting.Enabled = false
	cfg.Security.CORS.AllowedOrigins = []string{"*"}
	cfg.AI.OpenAI.APIKey = ""
	cfg.WeKnora.Enabled = true
	cfg.WeKnora.APIKey = ""

	warnings := SecurityWarnings(cfg)
	joined := strings.Join(warnings, "\n")

	for _, want := range []string{
		"jwt.secret is empty or using the default value",
		"security.cors.allowed_origins allows all origins",
		"security.rate_limiting is disabled",
		"ai.openai.api_key is empty",
		"weknora is enabled but weknora.api_key is empty",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing warning %q in %q", want, joined)
		}
	}
}

func TestLogSecurityWarnings(t *testing.T) {
	cfg := config.GetDefaultConfig()

	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})

	LogSecurityWarnings(logger, cfg)

	output := buf.String()
	if !strings.Contains(output, "security baseline warning") {
		t.Fatalf("expected warning output, got %q", output)
	}
	if !strings.Contains(output, "jwt.secret is empty or using the default value") {
		t.Fatalf("expected jwt warning in %q", output)
	}
}
