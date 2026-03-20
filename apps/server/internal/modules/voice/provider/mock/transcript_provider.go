package mock

import (
	"context"
	"sync"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

type TranscriptProvider struct {
	mu          sync.Mutex
	Transcripts []voiceapp.AppendTranscriptCommand
	AppendErr   error
}

func NewTranscriptProvider() *TranscriptProvider {
	return &TranscriptProvider{}
}

func (p *TranscriptProvider) AppendTranscript(ctx context.Context, cmd voiceapp.AppendTranscriptCommand) error {
	_ = ctx
	if p.AppendErr != nil {
		return p.AppendErr
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Transcripts = append(p.Transcripts, cmd)
	return nil
}

var _ voiceapp.TranscriptProvider = (*TranscriptProvider)(nil)
