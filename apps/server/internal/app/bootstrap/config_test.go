package bootstrap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFindsRepoRootConfigFromNestedDir(t *testing.T) {
	root := t.TempDir()
	serverDir := filepath.Join(root, "apps", "server")
	if err := os.MkdirAll(serverDir, 0o755); err != nil {
		t.Fatalf("mkdir nested dir: %v", err)
	}

	configPath := filepath.Join(root, "config.yml")
	if err := os.WriteFile(configPath, []byte("server:\n  port: 19090\nweknora:\n  enabled: true\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		if chdirErr := os.Chdir(cwd); chdirErr != nil {
			t.Fatalf("restore cwd: %v", chdirErr)
		}
	}()

	if err := os.Chdir(serverDir); err != nil {
		t.Fatalf("chdir nested dir: %v", err)
	}

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.Server.Port != 19090 {
		t.Fatalf("server port = %d, want 19090", cfg.Server.Port)
	}
	if !cfg.WeKnora.Enabled {
		t.Fatal("expected weknora.enabled to be loaded from repo-root config")
	}
}
