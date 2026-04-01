package worker

import (
	"context"
	"testing"
	"time"
)

type fakeLoop struct {
	started chan struct{}
	stopped chan struct{}
}

func (f *fakeLoop) run(ctx context.Context) {
	close(f.started)
	<-ctx.Done()
	close(f.stopped)
}

type fakeStatisticsService struct {
	loop *fakeLoop
}

func (f *fakeStatisticsService) StartDailyStatsWorkerContext(ctx context.Context, interval time.Duration) {
	f.loop.run(ctx)
}

type fakeSLAService struct {
	loop *fakeLoop
}

func (f *fakeSLAService) StartSLAMonitor(ctx context.Context, interval time.Duration) {
	f.loop.run(ctx)
}

type fakeAuditRetentionService struct {
	calls chan struct{}
}

func (f *fakeAuditRetentionService) Cleanup(ctx context.Context, now time.Time) (int64, error) {
	select {
	case f.calls <- struct{}{}:
	default:
	}
	<-ctx.Done()
	return 0, ctx.Err()
}

func TestStatisticsWorkerLifecycle(t *testing.T) {
	loop := &fakeLoop{started: make(chan struct{}), stopped: make(chan struct{})}
	w := &StatisticsWorker{
		service:  &fakeStatisticsService{loop: loop},
		interval: 100 * time.Millisecond,
	}
	if err := w.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	select {
	case <-loop.started:
	case <-time.After(time.Second):
		t.Fatal("worker did not start")
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := w.Stop(stopCtx); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	select {
	case <-loop.stopped:
	case <-time.After(time.Second):
		t.Fatal("worker did not stop")
	}
}

func TestSLAMonitorWorkerLifecycle(t *testing.T) {
	loop := &fakeLoop{started: make(chan struct{}), stopped: make(chan struct{})}
	w := &SLAMonitorWorker{
		service:  &fakeSLAService{loop: loop},
		interval: 100 * time.Millisecond,
	}
	if err := w.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	select {
	case <-loop.started:
	case <-time.After(time.Second):
		t.Fatal("worker did not start")
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := w.Stop(stopCtx); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	select {
	case <-loop.stopped:
	case <-time.After(time.Second):
		t.Fatal("worker did not stop")
	}
}

func TestAuditCleanupWorkerLifecycle(t *testing.T) {
	calls := make(chan struct{}, 1)
	w := &AuditCleanupWorker{
		service:  &fakeAuditRetentionService{calls: calls},
		interval: 100 * time.Millisecond,
		now:      time.Now,
	}
	if err := w.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	select {
	case <-calls:
	case <-time.After(time.Second):
		t.Fatal("worker did not execute cleanup")
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := w.Stop(stopCtx); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
}
