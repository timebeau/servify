package scripts

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDifyIntegrationScriptRealModeRejectsLocalHost(t *testing.T) {
	evidenceDir := t.TempDir()

	cmd := exec.Command("bash", "test-dify-integration.sh")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(),
		"DIFY_ACCEPTANCE_MODE=real",
		"SERVIFY_URL=http://127.0.0.1:18080",
		"DIFY_URL=http://127.0.0.1:15001/v1",
		"DIFY_DATASET_ID=dataset-1",
		"EVIDENCE_DIR="+evidenceDir,
	)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected real mode to reject local host, output=%s", string(output))
	}
	if !strings.Contains(string(output), "real 模式拒绝使用本地或私网 Dify 地址") {
		t.Fatalf("expected local-host guard message, output=%s", string(output))
	}

	summary, readErr := os.ReadFile(filepath.Join(evidenceDir, "summary.txt"))
	if readErr != nil {
		t.Fatalf("read summary: %v", readErr)
	}
	if !strings.Contains(string(summary), "real_mode_guard=blocked_private_or_local_host") {
		t.Fatalf("expected real mode guard summary, got %s", string(summary))
	}
}

func TestDifyIntegrationScriptMockModeWritesEvidence(t *testing.T) {
	servify := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/health":
			_, _ = w.Write([]byte(`{"status":"healthy"}`))
		case "/api/v1/ai/status":
			_, _ = w.Write([]byte(`{"success":true,"data":{"type":"orchestrated_enhanced","knowledge_provider_enabled":true,"knowledge_provider":"dify","knowledge_provider_healthy":true}}`))
		case "/api/v1/ai/query":
			_, _ = w.Write([]byte(`{"success":true,"data":{"content":"dify-ok","strategy":"dify"}}`))
		case "/api/v1/ai/knowledge/upload":
			_, _ = w.Write([]byte(`{"success":true}`))
		case "/api/v1/ai/knowledge/sync":
			_, _ = w.Write([]byte(`{"success":true}`))
		case "/api/v1/ai/metrics":
			_, _ = w.Write([]byte(`{"success":true,"data":{"dify_usage_count":2}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer servify.Close()

	difyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/v1/datasets/dataset-1" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"id":"dataset-1","name":"Mock Dify Dataset"}`))
	}))
	defer difyServer.Close()

	evidenceDir := t.TempDir()

	cmd := exec.Command("bash", "test-dify-integration.sh")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(),
		"DIFY_ACCEPTANCE_MODE=mock",
		"SERVIFY_URL="+servify.URL,
		"DIFY_URL="+difyServer.URL+"/v1",
		"DIFY_DATASET_ID=dataset-1",
		"EVIDENCE_DIR="+evidenceDir,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected mock mode success, err=%v output=%s", err, string(output))
	}

	for _, name := range []string{
		"summary.txt",
		"servify-health.json",
		"dify-dataset.json",
		"ai-status.json",
		"ai-query.json",
		"knowledge-upload.json",
		"knowledge-sync.json",
		"ai-metrics.json",
	} {
		if _, statErr := os.Stat(filepath.Join(evidenceDir, name)); statErr != nil {
			t.Fatalf("expected evidence file %s: %v\noutput=%s", name, statErr, string(output))
		}
	}

	summary, err := os.ReadFile(filepath.Join(evidenceDir, "summary.txt"))
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	summaryText := string(summary)
	for _, want := range []string{
		"mode=mock",
		"knowledge_provider=dify",
		"knowledge_provider_enabled=true",
		"knowledge_provider_healthy=true",
		"query_strategy=dify",
		"dify_available=true",
		"dify_usage_count=2",
		"knowledge_upload_ok=true",
		"knowledge_sync_ok=true",
	} {
		if !strings.Contains(summaryText, want) {
			t.Fatalf("expected %q in summary, got %s", want, summaryText)
		}
	}
}
