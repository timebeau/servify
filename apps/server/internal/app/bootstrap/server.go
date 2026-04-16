package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"servify/apps/server/internal/config"

	"github.com/sirupsen/logrus"
)

type HTTPServerOptions struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// ListenAddress resolves host/port overrides against config defaults.
func ListenAddress(cfg *config.Config, opts HTTPServerOptions) string {
	host := opts.Host
	port := opts.Port
	if cfg != nil {
		if host == "" {
			host = cfg.Server.Host
		}
		if port == 0 {
			port = cfg.Server.Port
		}
	}
	return fmt.Sprintf("%s:%d", host, port)
}

// NewHTTPServer creates an http.Server with optional timeout overrides.
func NewHTTPServer(cfg *config.Config, handler http.Handler, opts HTTPServerOptions) *http.Server {
	return &http.Server{
		Addr:         ListenAddress(cfg, opts),
		Handler:      handler,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
	}
}

// NewHTTPServerForApp creates the HTTP server from App dependencies and records the router.
func NewHTTPServerForApp(app *App, handler http.Handler, opts HTTPServerOptions) *http.Server {
	if app != nil {
		app.Router = handler
		app.Server = NewHTTPServer(app.Config, handler, opts)
		return app.Server
	}
	return NewHTTPServer(nil, handler, opts)
}

// BuildHTTPServer builds an HTTP server from the app's current router.
func BuildHTTPServer(app *App, opts HTTPServerOptions) *http.Server {
	if app == nil {
		return NewHTTPServer(nil, nil, opts)
	}
	return NewHTTPServerForApp(app, app.Router, opts)
}

// StartHTTPServer launches ListenAndServe in a goroutine.
func StartHTTPServer(server *http.Server, logger *logrus.Logger, message string) {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	if message == "" {
		message = fmt.Sprintf("Starting server on %s", server.Addr)
	}
	go func() {
		logger.Info(message)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()
}

// WaitForShutdownSignal blocks until SIGINT or SIGTERM is received.
func WaitForShutdownSignal() os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return <-quit
}

// ShutdownContext creates a bounded shutdown context.
func ShutdownContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return context.WithTimeout(context.Background(), timeout)
}
