package delivery

import (
	"context"
	"testing"

	voiceapp "servify/apps/server/internal/modules/voice/application"
	voiceinfra "servify/apps/server/internal/modules/voice/infra"
	"servify/apps/server/internal/platform/eventbus"
)

func TestCoordinatorStartRecordingAndAppendTranscript(t *testing.T) {
	bus := &stubBus{}
	callService := voiceapp.NewService(voiceinfra.NewInMemoryRepository(), bus)
	mediaProvider := voiceinfra.NewInMemoryMediaProvider()
	recordingService := voiceapp.NewRecordingService(mediaProvider, voiceinfra.NewInMemoryRecordingRepository(), bus)
	transcriptService := voiceapp.NewTranscriptService(mediaProvider, voiceinfra.NewInMemoryTranscriptRepository(), bus)
	coord := NewCoordinator(callService, recordingService, transcriptService)

	rec, err := coord.StartRecording(context.Background(), voiceapp.StartRecordingCommand{CallID: "call-1", Provider: "mock"})
	if err != nil {
		t.Fatalf("StartRecording() error = %v", err)
	}
	if rec == nil || rec.ID == "" {
		t.Fatalf("expected recording id")
	}
	tr, err := coord.AppendTranscript(context.Background(), voiceapp.AppendTranscriptCommand{CallID: "call-1", Content: "hi", Language: "en"})
	if err != nil {
		t.Fatalf("AppendTranscript() error = %v", err)
	}
	if tr == nil || tr.Content != "hi" {
		t.Fatalf("unexpected transcript: %+v", tr)
	}
}

type stubBus struct{}

func (s *stubBus) Publish(ctx context.Context, event eventbus.Event) error {
	return nil
}
