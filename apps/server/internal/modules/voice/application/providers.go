package application

import "context"

type StartRecordingCommand struct {
	CallID   string
	Provider string
}

type StopRecordingCommand struct {
	RecordingID string
}

type AppendTranscriptCommand struct {
	CallID    string
	Content   string
	Language  string
	Finalized bool
}

type RecordingProvider interface {
	StartRecording(ctx context.Context, cmd StartRecordingCommand) (string, error)
	StopRecording(ctx context.Context, cmd StopRecordingCommand) error
}

type TranscriptProvider interface {
	AppendTranscript(ctx context.Context, cmd AppendTranscriptCommand) error
}
