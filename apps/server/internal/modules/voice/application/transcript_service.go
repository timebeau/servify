package application

import (
	"context"
	"time"
)

type TranscriptService struct {
	provider  TranscriptProvider
	repo      TranscriptRepository
	publisher Publisher
}

func NewTranscriptService(provider TranscriptProvider, repo TranscriptRepository, publisher Publisher) *TranscriptService {
	return &TranscriptService{provider: provider, repo: repo, publisher: publisher}
}

func (s *TranscriptService) Append(ctx context.Context, cmd AppendTranscriptCommand) (*TranscriptDTO, error) {
	if err := s.provider.AppendTranscript(ctx, cmd); err != nil {
		return nil, err
	}
	transcript := &TranscriptDTO{
		CallID:     cmd.CallID,
		Content:    cmd.Content,
		Language:   cmd.Language,
		Finalized:  cmd.Finalized,
		AppendedAt: time.Now(),
	}
	if s.repo != nil {
		if err := s.repo.Append(ctx, *transcript); err != nil {
			return nil, err
		}
	}
	s.publish(ctx, TranscriptAppendedEventName, cmd.CallID, transcript)
	return transcript, nil
}

func (s *TranscriptService) ListByCallID(ctx context.Context, callID string) ([]TranscriptDTO, error) {
	if s.repo == nil {
		return nil, nil
	}
	return s.repo.ListByCallID(ctx, callID)
}

func (s *TranscriptService) publish(ctx context.Context, name, aggregateID string, payload interface{}) {
	if s.publisher == nil {
		return
	}
	_ = s.publisher.Publish(ctx, NewVoiceEvent(name, aggregateID, payload))
}
