package application

import (
	"context"
	"errors"
	"fmt"
)

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

type ProviderErrorCode string

const (
	ProviderErrorUnavailable ProviderErrorCode = "provider_unavailable"
	ProviderErrorTimeout     ProviderErrorCode = "provider_timeout"
	ProviderErrorRateLimited ProviderErrorCode = "provider_rate_limited"
	ProviderErrorInvalid     ProviderErrorCode = "provider_invalid_request"
)

type ProviderError struct {
	Code      ProviderErrorCode `json:"code"`
	Message   string            `json:"message"`
	Retryable bool              `json:"retryable"`
	Cause     error             `json:"-"`
}

func (e *ProviderError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return string(e.Code)
}

func (e *ProviderError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

type RetryPolicy struct {
	MaxAttempts int
}

type AsyncCallbackSink interface {
	NotifyRecording(ctx context.Context, recording RecordingDTO) error
	NotifyTranscript(ctx context.Context, transcript TranscriptDTO) error
}

func shouldRetryProviderError(err error) bool {
	var providerErr *ProviderError
	if errors.As(err, &providerErr) {
		return providerErr.Retryable
	}
	return false
}

func applyRetry(attempts int, fn func() error) error {
	if attempts <= 0 {
		attempts = 1
	}
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := fn(); err != nil {
			lastErr = err
			if !shouldRetryProviderError(err) || i == attempts-1 {
				return err
			}
			continue
		}
		return nil
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("provider retry failed without error")
}
