package bootstrap

import (
	"testing"

	"servify/apps/server/internal/config"
)

func TestBuildApp(t *testing.T) {
	cfg := config.GetDefaultConfig()

	app := BuildApp(cfg)
	if app == nil {
		t.Fatal("expected app")
	}
	if app.Config != cfg {
		t.Fatal("expected config to be assigned")
	}
	if app.Logger == nil {
		t.Fatal("expected logger to be initialized")
	}
	if app.EventBus == nil {
		t.Fatal("expected event bus to be initialized")
	}
	if app.Workers == nil {
		t.Fatal("expected workers slice")
	}
	if app.ShutdownHooks == nil {
		t.Fatal("expected shutdown hooks slice")
	}
}
