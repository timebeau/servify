package bootstrap

import (
	"net/http"
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

func TestNewHTTPServerForAppAssignsRouter(t *testing.T) {
	app := &App{Config: config.GetDefaultConfig()}
	handler := http.NewServeMux()

	srv := NewHTTPServerForApp(app, handler, HTTPServerOptions{Host: "127.0.0.1", Port: 9090})
	if srv.Addr != "127.0.0.1:9090" {
		t.Fatalf("server addr = %q", srv.Addr)
	}
	if app.Router != handler {
		t.Fatal("expected app router to be assigned")
	}
	if app.Server != srv {
		t.Fatal("expected app server to be assigned")
	}
}

func TestBuildHTTPServerUsesAppRouter(t *testing.T) {
	app := &App{
		Config: config.GetDefaultConfig(),
		Router: http.NewServeMux(),
	}
	srv := BuildHTTPServer(app, HTTPServerOptions{Host: "127.0.0.1", Port: 18080})
	if srv.Addr != "127.0.0.1:18080" {
		t.Fatalf("server addr = %q", srv.Addr)
	}
	if app.Server != srv {
		t.Fatal("expected app server to be assigned")
	}
}
