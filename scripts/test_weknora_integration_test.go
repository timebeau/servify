package scripts

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWeKnoraIntegrationScriptRealModeRejectsLocalHost(t *testing.T) {
	evidenceDir := t.TempDir()

	cmd := exec.Command("bash", "test-weknora-integration.sh")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(),
		"WEKNORA_ACCEPTANCE_MODE=real",
		"SERVIFY_URL=http://127.0.0.1:18080",
		"WEKNORA_URL=http://127.0.0.1:19000",
		"EVIDENCE_DIR="+evidenceDir,
	)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected real mode to reject local host, output=%s", string(output))
	}

	if !strings.Contains(string(output), "real 模式拒绝使用本地或私网 WeKnora 地址") {
		t.Fatalf("expected local-host guard message, output=%s", string(output))
	}

	summaryPath := filepath.Join(evidenceDir, "summary.txt")
	summary, readErr := os.ReadFile(summaryPath)
	if readErr != nil {
		t.Fatalf("read summary: %v", readErr)
	}
	if !strings.Contains(string(summary), "real_mode_guard=blocked_private_or_local_host") {
		t.Fatalf("expected real mode guard summary, got %s", string(summary))
	}
}

func TestWeKnoraIntegrationScriptMockModeWritesEvidence(t *testing.T) {
	servify := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/health":
			_, _ = w.Write([]byte(`{"status":"healthy"}`))
		case "/api/v1/ai/query":
			_, _ = w.Write([]byte(`{"success":true,"data":{"content":"ok"}}`))
		case "/api/v1/ai/status":
			_, _ = w.Write([]byte(`{"success":true,"data":{"type":"enhanced","weknora_enabled":true}}`))
		case "/api/v1/ai/metrics":
			_, _ = w.Write([]byte(`{"success":true,"data":{"query_count":3,"weknora_usage_count":2,"fallback_usage_count":1}}`))
		case "/api/v1/ai/knowledge/upload":
			_, _ = w.Write([]byte(`{"success":true,"data":{"document_id":"doc-1"}}`))
		case "/api/v1/ai/knowledge/sync":
			_, _ = w.Write([]byte(`{"success":true}`))
		case "/api/v1/ws/stats":
			_, _ = w.Write([]byte(`{"success":true,"data":{"client_count":0}}`))
		case "/api/v1/webrtc/connections":
			_, _ = w.Write([]byte(`{"success":true,"data":{"connection_count":0}}`))
		case "/api/customers/stats":
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"missing token"}`))
				return
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"total":1}}`))
		case "/api/statistics/dashboard":
			claims := jwtClaimsFromAuthHeader(r.Header.Get("Authorization"))
			if claims["principal_kind"] == "agent" {
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"error":"forbidden"}`))
				return
			}
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"missing token"}`))
				return
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"tickets":1}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer servify.Close()

	weknora := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/api/v1/health" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"service":"weknora-mock","status":"ok"}`))
	}))
	defer weknora.Close()

	evidenceDir := t.TempDir()

	cmd := exec.Command("bash", "test-weknora-integration.sh")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(),
		"WEKNORA_ACCEPTANCE_MODE=mock",
		"SERVIFY_URL="+servify.URL,
		"WEKNORA_URL="+weknora.URL,
		"EVIDENCE_DIR="+evidenceDir,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected mock mode success, err=%v output=%s", err, string(output))
	}

	for _, name := range []string{
		"summary.txt",
		"servify-health.json",
		"weknora-health.json",
		"ai-status.json",
		"ai-query.json",
		"knowledge-upload.json",
		"knowledge-sync.json",
		"ws-stats.json",
		"webrtc-connections.json",
	} {
		path := filepath.Join(evidenceDir, name)
		if _, statErr := os.Stat(path); statErr != nil {
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
		"overall_status=healthy",
		"service_type=enhanced",
		"weknora_available=true",
		"knowledge_upload_ok=true",
		"knowledge_sync_ok=true",
	} {
		if !strings.Contains(summaryText, want) {
			t.Fatalf("expected %q in summary, got %s", want, summaryText)
		}
	}
}

func jwtClaimsFromAuthHeader(header string) map[string]any {
	header = strings.TrimSpace(header)
	if !strings.HasPrefix(header, "Bearer ") {
		return map[string]any{}
	}
	token := strings.TrimPrefix(header, "Bearer ")
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return map[string]any{}
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return map[string]any{}
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return map[string]any{}
	}
	return claims
}
