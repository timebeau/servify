package mock

import (
	"context"
	"fmt"
	"sync"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

type RecordingProvider struct {
	mu           sync.Mutex
	recordings   map[string]voiceapp.StartRecordingCommand
	nextSequence int
	StartErr     error
	StopErr      error
}

func NewRecordingProvider() *RecordingProvider {
	return &RecordingProvider{
		recordings: make(map[string]voiceapp.StartRecordingCommand),
	}
}

func (p *RecordingProvider) StartRecording(ctx context.Context, cmd voiceapp.StartRecordingCommand) (string, error) {
	_ = ctx
	if p.StartErr != nil {
		return "", p.StartErr
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.nextSequence++
	id := fmt.Sprintf("rec-%d", p.nextSequence)
	p.recordings[id] = cmd
	return id, nil
}

func (p *RecordingProvider) StopRecording(ctx context.Context, cmd voiceapp.StopRecordingCommand) error {
	_ = ctx
	if p.StopErr != nil {
		return p.StopErr
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.recordings, cmd.RecordingID)
	return nil
}

var _ voiceapp.RecordingProvider = (*RecordingProvider)(nil)
