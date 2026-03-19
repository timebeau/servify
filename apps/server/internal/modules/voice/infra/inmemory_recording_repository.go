package infra

import (
	"context"
	"fmt"
	"sync"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

type InMemoryRecordingRepository struct {
	mu         sync.Mutex
	recordings map[string]voiceapp.RecordingDTO
}

func NewInMemoryRecordingRepository() *InMemoryRecordingRepository {
	return &InMemoryRecordingRepository{recordings: make(map[string]voiceapp.RecordingDTO)}
}

func (r *InMemoryRecordingRepository) Save(ctx context.Context, recording voiceapp.RecordingDTO) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.recordings[recording.ID] = recording
	return nil
}

func (r *InMemoryRecordingRepository) MarkStopped(ctx context.Context, recordingID string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	recording, ok := r.recordings[recordingID]
	if !ok {
		return fmt.Errorf("recording not found")
	}
	recording.Status = "stopped"
	r.recordings[recordingID] = recording
	return nil
}

func (r *InMemoryRecordingRepository) FindByID(ctx context.Context, recordingID string) (*voiceapp.RecordingDTO, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	recording, ok := r.recordings[recordingID]
	if !ok {
		return nil, fmt.Errorf("recording not found")
	}
	copy := recording
	return &copy, nil
}

var _ voiceapp.RecordingRepository = (*InMemoryRecordingRepository)(nil)
