package application

import (
	"context"
	"time"
)

const (
	RecordingStartedEventName   = "recording.started"
	RecordingStoppedEventName   = "recording.stopped"
	TranscriptAppendedEventName = "transcript.appended"
)

type RecordingService struct {
	provider     RecordingProvider
	repo         RecordingRepository
	publisher    Publisher
	retryPolicy  RetryPolicy
	callbackSink AsyncCallbackSink
}

func NewRecordingService(provider RecordingProvider, repo RecordingRepository, publisher Publisher) *RecordingService {
	return &RecordingService{
		provider:    provider,
		repo:        repo,
		publisher:   publisher,
		retryPolicy: RetryPolicy{MaxAttempts: 2},
	}
}

func (s *RecordingService) StartRecording(ctx context.Context, cmd StartRecordingCommand) (*RecordingDTO, error) {
	var recordingID string
	if err := applyRetry(s.retryPolicy.MaxAttempts, func() error {
		var err error
		recordingID, err = s.provider.StartRecording(ctx, cmd)
		return err
	}); err != nil {
		return nil, err
	}
	recording := &RecordingDTO{
		ID:        recordingID,
		CallID:    cmd.CallID,
		Provider:  cmd.Provider,
		Status:    "recording",
		StartedAt: time.Now(),
	}
	if s.repo != nil {
		if err := s.repo.Save(ctx, *recording); err != nil {
			return nil, err
		}
	}
	s.publish(ctx, RecordingStartedEventName, recording.ID, recording)
	if s.callbackSink != nil {
		_ = s.callbackSink.NotifyRecording(ctx, *recording)
	}
	return recording, nil
}

func (s *RecordingService) StopRecording(ctx context.Context, cmd StopRecordingCommand) error {
	if err := applyRetry(s.retryPolicy.MaxAttempts, func() error {
		return s.provider.StopRecording(ctx, cmd)
	}); err != nil {
		return err
	}
	if s.repo != nil {
		if err := s.repo.MarkStopped(ctx, cmd.RecordingID); err != nil {
			return err
		}
	}
	s.publish(ctx, RecordingStoppedEventName, cmd.RecordingID, map[string]string{"recording_id": cmd.RecordingID})
	return nil
}

func (s *RecordingService) GetRecording(ctx context.Context, recordingID string) (*RecordingDTO, error) {
	if s.repo == nil {
		return nil, nil
	}
	return s.repo.FindByID(ctx, recordingID)
}

func (s *RecordingService) publish(ctx context.Context, name, aggregateID string, payload interface{}) {
	if s.publisher == nil {
		return
	}
	_ = s.publisher.Publish(ctx, NewVoiceEvent(name, aggregateID, payload))
}

func (s *RecordingService) SetRetryPolicy(policy RetryPolicy) {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 1
	}
	s.retryPolicy = policy
}

func (s *RecordingService) SetCallbackSink(sink AsyncCallbackSink) {
	s.callbackSink = sink
}
