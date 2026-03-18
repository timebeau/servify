package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectStaticRoot(t *testing.T) {
	tmp := t.TempDir()
	existing := filepath.Join(tmp, "demo-web")
	if err := os.MkdirAll(existing, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	got := detectStaticRoot([]string{filepath.Join(tmp, "missing"), existing})
	if got != existing {
		t.Fatalf("detectStaticRoot() = %q want %q", got, existing)
	}
}
