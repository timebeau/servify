package application

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/eventbus"
)

type stubRecordingProvider struct {
	startID string
}

func (s *stubRecordingProvider) StartRecording(ctx context.Context, cmd StartRecordingCommand) (string, error) {
	return s.startID, nil
}

func (s *stubRecordingProvider) StopRecording(ctx context.Context, cmd StopRecordingCommand) error {
	return nil
}

type stubRecordingBus struct{ events []eventbus.Event }

func (s *stubRecordingBus) Publish(ctx context.Context, event eventbus.Event) error {
	s.events = append(s.events, event)
	return nil
}

func TestRecordingServiceStartStop(t *testing.T) {
	bus := &stubRecordingBus{}
	svc := NewRecordingService(&stubRecordingProvider{startID: "rec-1"}, nil, bus)

	recording, err := svc.StartRecording(context.Background(), StartRecordingCommand{CallID: "call-1", Provider: "mock"})
	if err != nil {
		t.Fatalf("StartRecording() error = %v", err)
	}
	if recording.ID != "rec-1" || recording.Status != "recording" {
		t.Fatalf("unexpected recording: %+v", recording)
	}
	if err := svc.StopRecording(context.Background(), StopRecordingCommand{RecordingID: "rec-1"}); err != nil {
		t.Fatalf("StopRecording() error = %v", err)
	}
	if len(bus.events) != 2 || bus.events[0].Name() != RecordingStartedEventName || bus.events[1].Name() != RecordingStoppedEventName {
		t.Fatalf("unexpected events: %+v", bus.events)
	}
}
