package bootstrap

import (
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
