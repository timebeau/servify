package bootstrap

import (
	"context"
	"errors"
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

type stubWorker struct {
	name     string
	started  bool
	stopped  bool
	startErr error
	stopErr  error
}

func (w *stubWorker) Name() string { return w.name }
func (w *stubWorker) Start() error {
	w.started = true
	return w.startErr
}
func (w *stubWorker) Stop(context.Context) error {
	w.stopped = true
	return w.stopErr
}

func TestAppWorkerLifecycle(t *testing.T) {
	app := BuildApp(config.GetDefaultConfig())
	w1 := &stubWorker{name: "w1"}
	w2 := &stubWorker{name: "w2"}
	app.RegisterWorker(w1)
	app.RegisterWorker(w2)

	if err := app.StartWorkers(); err != nil {
		t.Fatalf("StartWorkers() error = %v", err)
	}
	if !w1.started || !w2.started {
		t.Fatalf("expected all workers started: %+v %+v", w1, w2)
	}

	if err := app.StopWorkers(context.Background()); err != nil {
		t.Fatalf("StopWorkers() error = %v", err)
	}
	if !w1.stopped || !w2.stopped {
		t.Fatalf("expected all workers stopped: %+v %+v", w1, w2)
	}
}

func TestAppStartWorkersError(t *testing.T) {
	app := BuildApp(config.GetDefaultConfig())
	app.RegisterWorker(&stubWorker{name: "broken", startErr: errors.New("boom")})

	if err := app.StartWorkers(); err == nil {
		t.Fatal("expected start error")
	}
}

func TestAppShutdownHooks(t *testing.T) {
	app := BuildApp(config.GetDefaultConfig())
	order := make([]string, 0, 2)
	app.AddShutdownHook(func() error {
		order = append(order, "first")
		return nil
	})
	app.AddShutdownHook(func() error {
		order = append(order, "second")
		return nil
	})

	if err := app.RunShutdownHooks(); err != nil {
		t.Fatalf("RunShutdownHooks() error = %v", err)
	}
	if len(order) != 2 || order[0] != "second" || order[1] != "first" {
		t.Fatalf("unexpected hook order: %v", order)
	}
}
