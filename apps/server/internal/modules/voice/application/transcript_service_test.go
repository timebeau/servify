package application

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/eventbus"
)

type stubTranscriptProvider struct{}

func (s *stubTranscriptProvider) AppendTranscript(ctx context.Context, cmd AppendTranscriptCommand) error {
	return nil
}

type stubTranscriptBus struct{ events []eventbus.Event }

func (s *stubTranscriptBus) Publish(ctx context.Context, event eventbus.Event) error {
	s.events = append(s.events, event)
	return nil
}

func TestTranscriptServiceAppend(t *testing.T) {
	bus := &stubTranscriptBus{}
	svc := NewTranscriptService(&stubTranscriptProvider{}, nil, bus)

	transcript, err := svc.Append(context.Background(), AppendTranscriptCommand{
		CallID:    "call-1",
		Content:   "hello world",
		Language:  "en",
		Finalized: true,
	})
	if err != nil {
		t.Fatalf("Append() error = %v", err)
	}
	if transcript.CallID != "call-1" || !transcript.Finalized {
		t.Fatalf("unexpected transcript: %+v", transcript)
	}
	if len(bus.events) != 1 || bus.events[0].Name() != TranscriptAppendedEventName {
		t.Fatalf("unexpected events: %+v", bus.events)
	}
}
