package infra

import (
	"context"
	"fmt"
	"sync"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

type InMemoryMediaProvider struct {
	mu           sync.Mutex
	recordings   map[string]voiceapp.StartRecordingCommand
	transcripts  []voiceapp.AppendTranscriptCommand
	nextSequence int
}

func NewInMemoryMediaProvider() *InMemoryMediaProvider {
	return &InMemoryMediaProvider{
		recordings: make(map[string]voiceapp.StartRecordingCommand),
	}
}

func (p *InMemoryMediaProvider) StartRecording(ctx context.Context, cmd voiceapp.StartRecordingCommand) (string, error) {
	_ = ctx
	p.mu.Lock()
	defer p.mu.Unlock()
	p.nextSequence++
	id := fmt.Sprintf("rec-%d", p.nextSequence)
	p.recordings[id] = cmd
	return id, nil
}

func (p *InMemoryMediaProvider) StopRecording(ctx context.Context, cmd voiceapp.StopRecordingCommand) error {
	_ = ctx
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.recordings, cmd.RecordingID)
	return nil
}

func (p *InMemoryMediaProvider) AppendTranscript(ctx context.Context, cmd voiceapp.AppendTranscriptCommand) error {
	_ = ctx
	p.mu.Lock()
	defer p.mu.Unlock()
	p.transcripts = append(p.transcripts, cmd)
	return nil
}

var _ voiceapp.RecordingProvider = (*InMemoryMediaProvider)(nil)
var _ voiceapp.TranscriptProvider = (*InMemoryMediaProvider)(nil)
