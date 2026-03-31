package eventbus

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestInMemoryBusPublishToSubscribers(t *testing.T) {
	bus := NewInMemoryBus()
	event := BaseEvent{
		EventID:          "evt-1",
		EventName:        "ticket.created",
		EventOccurredAt:  time.Now(),
		EventTenantID:    "tenant-1",
		EventAggregateID: "ticket-1",
	}

	called := 0
	bus.Subscribe("ticket.created", HandlerFunc(func(ctx context.Context, got Event) error {
		called++
		if got.Name() != "ticket.created" {
			t.Fatalf("unexpected event name: %s", got.Name())
		}
		return nil
	}))
	bus.Subscribe("ticket.created", HandlerFunc(func(ctx context.Context, got Event) error {
		called++
		return nil
	}))

	if err := bus.Publish(context.Background(), event); err != nil {
		t.Fatalf("publish failed: %v", err)
	}
	if called != 2 {
		t.Fatalf("expected 2 handlers to be called, got %d", called)
	}
}

func TestInMemoryBusPublishReturnsHandlerError(t *testing.T) {
	bus := NewInMemoryBus()
	want := errors.New("boom")
	event := BaseEvent{
		EventID:          "evt-2",
		EventName:        "routing.agent_assigned",
		EventOccurredAt:  time.Now(),
		EventTenantID:    "tenant-1",
		EventAggregateID: "assignment-1",
	}

	bus.Subscribe("routing.agent_assigned", HandlerFunc(func(ctx context.Context, got Event) error {
		return want
	}))

	err := bus.Publish(context.Background(), event)
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestInMemoryBusContinuesAfterHandlerError(t *testing.T) {
	bus := NewInMemoryBus()
	event := BaseEvent{
		EventID:         "evt-3",
		EventName:       "test.error_propagation",
		EventOccurredAt: time.Now(),
	}

	handler1Called := false
	handler2Called := false

	bus.Subscribe("test.error_propagation", HandlerFunc(func(ctx context.Context, got Event) error {
		handler1Called = true
		return errors.New("handler1 failed")
	}))
	bus.Subscribe("test.error_propagation", HandlerFunc(func(ctx context.Context, got Event) error {
		handler2Called = true
		return nil
	}))

	_ = bus.Publish(context.Background(), event)

	if !handler1Called {
		t.Error("handler1 should have been called")
	}
	if !handler2Called {
		t.Error("handler2 should have been called even after handler1 failed")
	}
}

func TestInMemoryBusDeadLetterTracking(t *testing.T) {
	bus := NewInMemoryBus()
	event := BaseEvent{
		EventID:         "evt-dl",
		EventName:       "test.dead_letter",
		EventOccurredAt: time.Now(),
	}

	bus.Subscribe("test.dead_letter", HandlerFunc(func(ctx context.Context, got Event) error {
		return errors.New("always fails")
	}))

	_ = bus.Publish(context.Background(), event)

	health := bus.Health()
	if health.PublishedCount != 1 {
		t.Errorf("PublishedCount = %d, want 1", health.PublishedCount)
	}
	if health.FailedCount != 1 {
		t.Errorf("FailedCount = %d, want 1", health.FailedCount)
	}
	if health.DeadLetterCount != 1 {
		t.Errorf("DeadLetterCount = %d, want 1", health.DeadLetterCount)
	}
	if len(health.DeadLetters) != 1 {
		t.Fatalf("len(DeadLetters) = %d, want 1", len(health.DeadLetters))
	}
	if health.DeadLetters[0].EventID != "evt-dl" {
		t.Errorf("DeadLetter EventID = %q, want %q", health.DeadLetters[0].EventID, "evt-dl")
	}
}

func TestInMemoryBusDeadLetterRingBuffer(t *testing.T) {
	bus := NewInMemoryBus()
	bus.maxDeadLetters = 3

	event := BaseEvent{EventID: "overflow", EventName: "test.overflow", EventOccurredAt: time.Now()}
	bus.Subscribe("test.overflow", HandlerFunc(func(ctx context.Context, got Event) error {
		return errors.New("fail")
	}))

	for i := 0; i < 5; i++ {
		_ = bus.Publish(context.Background(), event)
	}

	health := bus.Health()
	if len(health.DeadLetters) != 3 {
		t.Errorf("len(DeadLetters) = %d, want 3 (ring buffer cap)", len(health.DeadLetters))
	}
	if health.DeadLetterCount != 5 {
		t.Errorf("DeadLetterCount = %d, want 5 (total count)", health.DeadLetterCount)
	}
}
