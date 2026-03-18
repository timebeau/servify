package bootstrap

import (
	"testing"
	"time"

	"servify/apps/server/internal/config"
)

func TestListenAddress(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Server.Host = "0.0.0.0"
	cfg.Server.Port = 8080

	if got := ListenAddress(cfg, HTTPServerOptions{}); got != "0.0.0.0:8080" {
		t.Fatalf("ListenAddress() = %q want %q", got, "0.0.0.0:8080")
	}
	if got := ListenAddress(cfg, HTTPServerOptions{Host: "127.0.0.1", Port: 9090}); got != "127.0.0.1:9090" {
		t.Fatalf("ListenAddress() override = %q want %q", got, "127.0.0.1:9090")
	}
}

func TestShutdownContextDefault(t *testing.T) {
	ctx, cancel := ShutdownContext(0)
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline")
	}
	if time.Until(deadline) <= 0 {
		t.Fatal("expected future deadline")
	}
}
