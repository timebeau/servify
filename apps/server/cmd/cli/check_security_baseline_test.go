package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCheckSecurityBaseline_StrictFailure(t *testing.T) {
	configPath := writeTempConfig(t, `
security:
  rate_limiting:
    enabled: false
`)

	var out bytes.Buffer
	err := runCheckSecurityBaseline(configPath, true, &out)
	if err == nil {
		t.Fatal("expected strict mode error")
	}

	output := out.String()
	for _, want := range []string{
		"Security baseline check found",
		"jwt.secret is empty or using the default value",
		"security.rate_limiting is disabled",
		"ai.openai.api_key is empty",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("missing %q in output %q", want, output)
		}
	}
}

func TestRunCheckSecurityBaseline_Pass(t *testing.T) {
	configPath := writeTempConfig(t, `
jwt:
  secret: "prod-secret"
security:
  cors:
    enabled: true
    allowed_origins:
      - "https://app.example.com"
  rate_limiting:
    enabled: true
    requests_per_minute: 120
    burst: 30
    paths:
      - enabled: true
        prefix: "/public/"
        requests_per_minute: 120
        burst: 30
      - enabled: true
        prefix: "/api/v1/ws"
        requests_per_minute: 40
        burst: 20
ai:
  openai:
    api_key: "openai-key"
weknora:
  enabled: false
`)

	var out bytes.Buffer
	if err := runCheckSecurityBaseline(configPath, true, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != "Security baseline check passed." {
		t.Fatalf("unexpected output %q", got)
	}
}

func writeTempConfig(t *testing.T, body string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(strings.TrimSpace(body)+"\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}
