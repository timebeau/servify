package application

import "context"

type Service struct {
	repo      Repository
	publisher Publisher
}

func NewService(repo Repository, publisher Publisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

func (s *Service) StartCall(ctx context.Context, cmd StartCallCommand) (*CallDTO, error) {
	call, err := s.repo.StartCall(ctx, cmd)
	if err != nil {
		return nil, err
	}
	s.publish(ctx, CallStartedEventName, call.ID, call)
	return call, nil
}

func (s *Service) AnswerCall(ctx context.Context, cmd AnswerCallCommand) (*CallDTO, error) {
	return s.repo.AnswerCall(ctx, cmd)
}

func (s *Service) EndCall(ctx context.Context, cmd EndCallCommand) (*CallDTO, error) {
	call, err := s.repo.EndCall(ctx, cmd)
	if err != nil {
		return nil, err
	}
	s.publish(ctx, CallEndedEventName, call.ID, call)
	return call, nil
}

func (s *Service) TransferCall(ctx context.Context, cmd TransferCallCommand) (*CallDTO, error) {
	return s.repo.TransferCall(ctx, cmd)
}

func (s *Service) publish(ctx context.Context, name, aggregateID string, payload interface{}) {
	if s.publisher == nil {
		return
	}
	_ = s.publisher.Publish(ctx, NewVoiceEvent(name, aggregateID, payload))
}
