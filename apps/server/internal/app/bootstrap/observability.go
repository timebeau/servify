package bootstrap

import (
	"context"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/observability"
)

// SetupObservability initializes tracing and optionally registers its shutdown hook.
func SetupObservability(ctx context.Context, cfg *config.Config, app *App) error {
	shutdown, err := observability.SetupTracing(ctx, cfg)
	if err != nil {
		return err
	}
	if app != nil {
		app.AddShutdownHook(func() error {
			return shutdown(context.Background())
		})
	}
	return nil
}
