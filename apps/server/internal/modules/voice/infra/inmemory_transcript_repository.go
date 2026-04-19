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

func (r *InMemoryTranscriptRepository) ListAll(ctx context.Context, page, pageSize int) ([]voiceapp.TranscriptDTO, int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	// Collect all transcripts
	var all []voiceapp.TranscriptDTO
	for _, items := range r.transcripts {
		all = append(all, items...)
	}

	total := int64(len(all))

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(all) {
		return []voiceapp.TranscriptDTO{}, total, nil
	}
	if end > len(all) {
		end = len(all)
	}

	out := make([]voiceapp.TranscriptDTO, end-start)
	copy(out, all[start:end])
	return out, total, nil
}

var _ voiceapp.TranscriptRepository = (*InMemoryTranscriptRepository)(nil)
