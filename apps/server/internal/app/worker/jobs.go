package worker

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"servify/apps/server/internal/app/bootstrap"

	"github.com/sirupsen/logrus"
)

// jitter returns a random duration in [0, fraction*base).
// Returns 0 if base is negative or fraction <= 0.
func jitter(base time.Duration, fraction float64) time.Duration {
	if fraction <= 0 || base <= 0 {
		return 0
	}
	return time.Duration(float64(base) * fraction * rand.Float64())
}

type statisticsService interface {
	StartDailyStatsWorkerContext(context.Context, time.Duration)
}

// StatisticsWorker runs periodic daily-stats aggregation.
type StatisticsWorker struct {
	service  statisticsService
	interval time.Duration
	logger   *logrus.Logger

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func NewStatisticsWorker(service statisticsService, interval time.Duration) bootstrap.Worker {
	if interval <= 0 {
		interval = time.Hour
	}
	return &StatisticsWorker{
		service:  service,
		interval: interval,
		logger:   logrus.StandardLogger(),
	}
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
		// Add initial jitter to stagger workers on restart.
		initialDelay := jitter(w.interval, 0.1)
		if initialDelay > 0 {
			if w.logger != nil {
				w.logger.Debugf("statistics worker: initial jitter delay %v", initialDelay)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(initialDelay):
			}
		}
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
	logger   *logrus.Logger

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func NewSLAMonitorWorker(service slaMonitorService, interval time.Duration) bootstrap.Worker {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	return &SLAMonitorWorker{
		service:  service,
		interval: interval,
		logger:   logrus.StandardLogger(),
	}
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
		initialDelay := jitter(w.interval, 0.1)
		if initialDelay > 0 {
			if w.logger != nil {
				w.logger.Debugf("sla-monitor worker: initial jitter delay %v", initialDelay)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(initialDelay):
			}
		}
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

type auditRetentionService interface {
	Cleanup(context.Context, time.Time) (int64, error)
}

// AuditCleanupWorker periodically deletes expired audit logs.
type AuditCleanupWorker struct {
	service  auditRetentionService
	interval time.Duration
	logger   *logrus.Logger
	now      func() time.Time

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func NewAuditCleanupWorker(service auditRetentionService, interval time.Duration) bootstrap.Worker {
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	return &AuditCleanupWorker{
		service:  service,
		interval: interval,
		logger:   logrus.StandardLogger(),
		now:      time.Now,
	}
}

func (w *AuditCleanupWorker) Name() string { return "audit-retention-cleanup" }

func (w *AuditCleanupWorker) Start() error {
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
		initialDelay := jitter(w.interval, 0.1)
		if initialDelay > 0 {
			if w.logger != nil {
				w.logger.Debugf("audit-cleanup worker: initial jitter delay %v", initialDelay)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(initialDelay):
			}
		}

		run := func() bool {
			deleted, err := w.service.Cleanup(ctx, w.now().UTC())
			if err != nil {
				if w.logger != nil {
					w.logger.WithError(err).Warn("audit cleanup worker: cleanup failed")
				}
				return false
			}
			if deleted > 0 && w.logger != nil {
				w.logger.Infof("audit cleanup worker: deleted %d expired audit logs", deleted)
			}
			return true
		}

		if !run() && ctx.Err() != nil {
			return
		}

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
	return nil
}

func (w *AuditCleanupWorker) Stop(ctx context.Context) error {
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
