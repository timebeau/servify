package application

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/eventbus"
)

type stubRepo struct {
	call *CallDTO
}

func (s *stubRepo) StartCall(ctx context.Context, cmd StartCallCommand) (*CallDTO, error) {
	s.call = &CallDTO{ID: cmd.CallID, SessionID: cmd.SessionID, Status: "started"}
	return s.call, nil
}
func (s *stubRepo) AnswerCall(ctx context.Context, cmd AnswerCallCommand) (*CallDTO, error) {
	s.call.Status = "answered"
	return s.call, nil
}
func (s *stubRepo) EndCall(ctx context.Context, cmd EndCallCommand) (*CallDTO, error) {
	s.call.Status = "ended"
	return s.call, nil
}
func (s *stubRepo) TransferCall(ctx context.Context, cmd TransferCallCommand) (*CallDTO, error) {
	return s.call, nil
}

type stubBus struct{ events []eventbus.Event }

func (s *stubBus) Publish(ctx context.Context, event eventbus.Event) error {
	s.events = append(s.events, event)
	return nil
}

func TestStartAndEndCallPublishesEvents(t *testing.T) {
	repo := &stubRepo{}
	bus := &stubBus{}
	svc := NewService(repo, bus)
	call, err := svc.StartCall(context.Background(), StartCallCommand{CallID: "c1", SessionID: "s1"})
	if err != nil {
		t.Fatalf("StartCall() error = %v", err)
	}
	if call.Status != "started" {
		t.Fatalf("unexpected status: %s", call.Status)
	}
	if _, err := svc.EndCall(context.Background(), EndCallCommand{CallID: "c1"}); err != nil {
		t.Fatalf("EndCall() error = %v", err)
	}
	if len(bus.events) != 2 || bus.events[0].Name() != CallStartedEventName || bus.events[1].Name() != CallEndedEventName {
		t.Fatalf("unexpected events: %+v", bus.events)
	}
}
