package disabled

import (
	"context"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

const (
	recordingDisabledMessage  = "voice recording provider is disabled"
	transcriptDisabledMessage = "voice transcript provider is disabled"
)

type RecordingProvider struct{}

func NewRecordingProvider() *RecordingProvider {
	return &RecordingProvider{}
}

func (p *RecordingProvider) StartRecording(ctx context.Context, cmd voiceapp.StartRecordingCommand) (string, error) {
	_ = ctx
	_ = cmd
	return "", &voiceapp.ProviderError{
		Code:      voiceapp.ProviderErrorUnavailable,
		Message:   recordingDisabledMessage,
		Retryable: false,
	}
}

func (p *RecordingProvider) StopRecording(ctx context.Context, cmd voiceapp.StopRecordingCommand) error {
	_ = ctx
	_ = cmd
	return &voiceapp.ProviderError{
		Code:      voiceapp.ProviderErrorUnavailable,
		Message:   recordingDisabledMessage,
		Retryable: false,
	}
}

type TranscriptProvider struct{}

func NewTranscriptProvider() *TranscriptProvider {
	return &TranscriptProvider{}
}

func (p *TranscriptProvider) AppendTranscript(ctx context.Context, cmd voiceapp.AppendTranscriptCommand) error {
	_ = ctx
	_ = cmd
	return &voiceapp.ProviderError{
		Code:      voiceapp.ProviderErrorUnavailable,
		Message:   transcriptDisabledMessage,
		Retryable: false,
	}
}

var _ voiceapp.RecordingProvider = (*RecordingProvider)(nil)
var _ voiceapp.TranscriptProvider = (*TranscriptProvider)(nil)
