package async

import (
	"context"
	"errors"
	"testing"
	"time"

	"servify/apps/server/internal/observability/metrics"
	"servify/apps/server/internal/platform/eventbus"
)

type testEvent struct {
	base eventbus.BaseEvent
}

func (e *testEvent) ID() string            { return e.base.EventID }
func (e *testEvent) Name() string          { return e.base.EventName }
func (e *testEvent) OccurredAt() time.Time { return e.base.EventOccurredAt }
func (e *testEvent) TenantID() string      { return e.base.EventTenantID }
func (e *testEvent) AggregateID() string   { return e.base.EventAggregateID }

func TestBusMiddleware_Success(t *testing.T) {
	reg := metrics.NewRegistry()
	busMetrics := NewBusMetrics(reg)

	called := false
	handler := eventbus.HandlerFunc(func(ctx context.Context, event eventbus.Event) error {
		called = true
		return nil
	})

	wrapped := WrapHandler("test.event", handler, busMetrics, nil, nil)

	evt := &testEvent{base: eventbus.BaseEvent{
		EventID:         "evt-1",
		EventName:       "test.event",
		EventOccurredAt: time.Now(),
	}}

	err := wrapped.Handle(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected handler to be called")
	}

	// Verify metrics
	mfs, _ := reg.Gatherer().Gather()
	found := false
	for _, mf := range mfs {
		if mf.GetName() == "eventbus_handled_total" {
			found = true
			if mf.GetMetric()[0].GetCounter().GetValue() != 1 {
				t.Fatal("expected handled counter = 1")
			}
		}
	}
	if !found {
		t.Fatal("expected eventbus_handled_total metric")
	}
}

func TestBusMiddleware_Failure(t *testing.T) {
	reg := metrics.NewRegistry()
	busMetrics := NewBusMetrics(reg)
	dlr := NewInMemoryDeadLetterRecorder(100)

	handler := eventbus.HandlerFunc(func(ctx context.Context, event eventbus.Event) error {
		return errors.New("handler failed")
	})

	wrapped := WrapHandler("test.event", handler, busMetrics, dlr, nil)

	evt := &testEvent{base: eventbus.BaseEvent{
		EventID:         "evt-2",
		EventName:       "test.event",
		EventOccurredAt: time.Now(),
	}}

	err := wrapped.Handle(context.Background(), evt)
	if err == nil {
		t.Fatal("expected error from handler")
	}

	// Verify dead letter recorded
	entries, err := dlr.List(context.Background(), "", 10)
	if err != nil {
		t.Fatalf("unexpected error listing dead letters: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 dead letter entry, got %d", len(entries))
	}
	if entries[0].Error != "handler failed" {
		t.Fatalf("expected 'handler failed', got %q", entries[0].Error)
	}
}
