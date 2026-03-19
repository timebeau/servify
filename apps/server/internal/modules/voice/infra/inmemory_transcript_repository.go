package infra

import (
	"context"
	"sync"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

type InMemoryTranscriptRepository struct {
	mu          sync.Mutex
	transcripts map[string][]voiceapp.TranscriptDTO
}

func NewInMemoryTranscriptRepository() *InMemoryTranscriptRepository {
	return &InMemoryTranscriptRepository{transcripts: make(map[string][]voiceapp.TranscriptDTO)}
}

func (r *InMemoryTranscriptRepository) Append(ctx context.Context, transcript voiceapp.TranscriptDTO) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transcripts[transcript.CallID] = append(r.transcripts[transcript.CallID], transcript)
	return nil
}

func (r *InMemoryTranscriptRepository) ListByCallID(ctx context.Context, callID string) ([]voiceapp.TranscriptDTO, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	items := r.transcripts[callID]
	out := make([]voiceapp.TranscriptDTO, len(items))
	copy(out, items)
	return out, nil
}

var _ voiceapp.TranscriptRepository = (*InMemoryTranscriptRepository)(nil)
