package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCheckObservabilityBaseline_StrictFailure(t *testing.T) {
	configPath := writeTempConfig(t, `
monitoring:
  enabled: false
  metrics_path: "metrics"
  tracing:
    enabled: true
    endpoint: ""
    sample_ratio: 0
    service_name: ""
`)

	var out bytes.Buffer
	err := runCheckObservabilityBaseline(configPath, true, t.TempDir(), &out)
	if err == nil {
		t.Fatal("expected strict mode error")
	}

	output := out.String()
	for _, want := range []string{
		"Observability baseline check found",
		"monitoring is disabled",
		"monitoring.metrics_path must start with /",
		"monitoring.tracing.endpoint is empty while tracing is enabled",
		"monitoring.tracing.sample_ratio must be within (0,1]",
		"monitoring.tracing.service_name is empty while tracing is enabled",
		"observability asset missing: deploy/observability/alerts/rules.yaml",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("missing %q in output %q", want, output)
		}
	}
}

func TestRunCheckObservabilityBaseline_Pass(t *testing.T) {
	configPath := writeTempConfig(t, `
monitoring:
  enabled: true
  metrics_path: "/metrics"
  tracing:
    enabled: true
    endpoint: "http://otel-collector:4317"
    sample_ratio: 0.1
    service_name: "servify"
`)

	repoRoot := tempObservabilityRepoRoot(t)

	var out bytes.Buffer
	if err := runCheckObservabilityBaseline(configPath, true, repoRoot, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != "Observability baseline check passed." {
		t.Fatalf("unexpected output %q", got)
	}
}

func tempObservabilityRepoRoot(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	for _, rel := range []string{
		"infra/compose/docker-compose.observability.yml",
		"infra/observability/otel-collector-config.yaml",
		"deploy/observability/alerts/rules.yaml",
		"deploy/observability/dashboards/servify-service.json",
		"deploy/observability/dashboards/servify-business.json",
		"deploy/observability/runbook/operational-runbook.md",
	} {
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte("ok\n"), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	return root
}
