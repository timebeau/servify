package application

import "context"

type Repository interface {
	StartCall(ctx context.Context, cmd StartCallCommand) (*CallDTO, error)
	AnswerCall(ctx context.Context, cmd AnswerCallCommand) (*CallDTO, error)
	HoldCall(ctx context.Context, cmd HoldCallCommand) (*CallDTO, error)
	ResumeCall(ctx context.Context, cmd ResumeCallCommand) (*CallDTO, error)
	EndCall(ctx context.Context, cmd EndCallCommand) (*CallDTO, error)
	TransferCall(ctx context.Context, cmd TransferCallCommand) (*CallDTO, error)
}

type RecordingRepository interface {
	Save(ctx context.Context, recording RecordingDTO) error
	MarkStopped(ctx context.Context, recordingID string) error
	FindByID(ctx context.Context, recordingID string) (*RecordingDTO, error)
}

type TranscriptRepository interface {
	Append(ctx context.Context, transcript TranscriptDTO) error
	ListByCallID(ctx context.Context, callID string) ([]TranscriptDTO, error)
	ListAll(ctx context.Context, page, pageSize int) ([]TranscriptDTO, int64, error)
}
