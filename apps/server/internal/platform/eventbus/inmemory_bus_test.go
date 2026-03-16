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
