package application

import "context"

type Repository interface {
	StartCall(ctx context.Context, cmd StartCallCommand) (*CallDTO, error)
	AnswerCall(ctx context.Context, cmd AnswerCallCommand) (*CallDTO, error)
	EndCall(ctx context.Context, cmd EndCallCommand) (*CallDTO, error)
	TransferCall(ctx context.Context, cmd TransferCallCommand) (*CallDTO, error)
}
