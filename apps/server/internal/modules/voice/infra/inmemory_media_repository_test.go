package infra

import (
	"context"
	"testing"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

func TestInMemoryRecordingRepositorySaveAndStop(t *testing.T) {
	repo := NewInMemoryRecordingRepository()
	if err := repo.Save(context.Background(), voiceapp.RecordingDTO{
		ID:     "rec-1",
		CallID: "call-1",
		Status: "recording",
	}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := repo.MarkStopped(context.Background(), "rec-1"); err != nil {
		t.Fatalf("MarkStopped() error = %v", err)
	}
	recording, err := repo.FindByID(context.Background(), "rec-1")
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if recording.Status != "stopped" {
		t.Fatalf("unexpected recording: %+v", recording)
	}
}

func TestInMemoryTranscriptRepositoryAppendAndList(t *testing.T) {
	repo := NewInMemoryTranscriptRepository()
	if err := repo.Append(context.Background(), voiceapp.TranscriptDTO{
		CallID:   "call-1",
		Content:  "hello",
		Language: "en",
	}); err != nil {
		t.Fatalf("Append() error = %v", err)
	}
	items, err := repo.ListByCallID(context.Background(), "call-1")
	if err != nil {
		t.Fatalf("ListByCallID() error = %v", err)
	}
	if len(items) != 1 || items[0].Content != "hello" {
		t.Fatalf("unexpected items: %+v", items)
	}
}
