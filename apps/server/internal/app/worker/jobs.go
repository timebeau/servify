package worker

import (
	"context"
	"sync"
	"time"

	"servify/apps/server/internal/app/bootstrap"
)

type statisticsService interface {
	StartDailyStatsWorkerContext(context.Context, time.Duration)
}

// StatisticsWorker runs periodic daily-stats aggregation.
type StatisticsWorker struct {
	service  statisticsService
	interval time.Duration

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func NewStatisticsWorker(service statisticsService, interval time.Duration) bootstrap.Worker {
	if interval <= 0 {
		interval = time.Hour
	}
	return &StatisticsWorker{service: service, interval: interval}
}

func (w *StatisticsWorker) Name() string { return "statistics-daily-stats" }

func (w *StatisticsWorker) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil || w.service == nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	w.cancel = cancel
	w.done = done
	go func() {
		defer close(done)
		w.service.StartDailyStatsWorkerContext(ctx, w.interval)
	}()
	return nil
}

func (w *StatisticsWorker) Stop(ctx context.Context) error {
	w.mu.Lock()
	cancel := w.cancel
	done := w.done
	w.cancel = nil
	w.done = nil
	w.mu.Unlock()

	if cancel == nil {
		return nil
	}
	cancel()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type slaMonitorService interface {
	StartSLAMonitor(context.Context, time.Duration)
}

// SLAMonitorWorker runs periodic SLA violation scanning.
type SLAMonitorWorker struct {
	service  slaMonitorService
	interval time.Duration

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func NewSLAMonitorWorker(service slaMonitorService, interval time.Duration) bootstrap.Worker {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	return &SLAMonitorWorker{service: service, interval: interval}
}

func (w *SLAMonitorWorker) Name() string { return "sla-monitor" }

func (w *SLAMonitorWorker) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil || w.service == nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	w.cancel = cancel
	w.done = done
	go func() {
		defer close(done)
		w.service.StartSLAMonitor(ctx, w.interval)
	}()
	return nil
}

func (w *SLAMonitorWorker) Stop(ctx context.Context) error {
	w.mu.Lock()
	cancel := w.cancel
	done := w.done
	w.cancel = nil
	w.done = nil
	w.mu.Unlock()

	if cancel == nil {
		return nil
	}
	cancel()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
