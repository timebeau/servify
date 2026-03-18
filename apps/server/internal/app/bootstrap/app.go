package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/platform/eventbus"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Worker is the minimal lifecycle contract for background jobs managed by App.
type Worker interface {
	Name() string
	Start() error
	Stop(context.Context) error
}

// App is the bootstrap root for server runtime wiring.
// The initial skeleton only collects shared runtime dependencies.
type App struct {
	Config        *config.Config
	Logger        *logrus.Logger
	DB            *gorm.DB
	Router        http.Handler
	EventBus      eventbus.Bus
	Workers       []Worker
	ShutdownHooks []func() error
}

// BuildApp creates the application runtime skeleton.
// Later tasks will move config, logging, db, router, and worker wiring here.
func BuildApp(cfg *config.Config) *App {
	return &App{
		Config:        cfg,
		Logger:        logrus.StandardLogger(),
		EventBus:      eventbus.NewInMemoryBus(),
		Workers:       make([]Worker, 0),
		ShutdownHooks: make([]func() error, 0),
	}
}

// RegisterWorker appends a managed background worker.
func (a *App) RegisterWorker(w Worker) {
	if w == nil {
		return
	}
	a.Workers = append(a.Workers, w)
}

// StartWorkers starts all registered workers in order.
func (a *App) StartWorkers() error {
	for _, w := range a.Workers {
		if err := w.Start(); err != nil {
			return fmt.Errorf("start worker %s: %w", w.Name(), err)
		}
	}
	return nil
}

// AddShutdownHook appends a shutdown hook executed during app termination.
func (a *App) AddShutdownHook(hook func() error) {
	if hook == nil {
		return
	}
	a.ShutdownHooks = append(a.ShutdownHooks, hook)
}

// RunShutdownHooks runs hooks in reverse order.
func (a *App) RunShutdownHooks() error {
	for i := len(a.ShutdownHooks) - 1; i >= 0; i-- {
		if err := a.ShutdownHooks[i](); err != nil {
			return err
		}
	}
	return nil
}

// StopWorkers stops all registered workers in reverse order.
func (a *App) StopWorkers(ctx context.Context) error {
	for i := len(a.Workers) - 1; i >= 0; i-- {
		if err := a.Workers[i].Stop(ctx); err != nil {
			return fmt.Errorf("stop worker %s: %w", a.Workers[i].Name(), err)
		}
	}
	return nil
}
