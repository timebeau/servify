package llm

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestHTTPErrorClassification(t *testing.T) {
	err := HTTPError("openai", 429, "rate limited")
	if err.Code != ProviderErrorRateLimited || !err.Retryable {
		t.Fatalf("unexpected rate limit classification: %+v", err)
	}

	err = HTTPError("openai", 401, "unauthorized")
	if err.Code != ProviderErrorAuthFailed || err.Retryable {
		t.Fatalf("unexpected auth classification: %+v", err)
	}
}

func TestRetryStopsAfterSuccess(t *testing.T) {
	attempts := 0
	err := Retry(context.Background(), RetryPolicy{
		MaxAttempts:   3,
		BaseDelayMs:   1,
		BackoffFactor: 1,
	}, func(ctx context.Context) error {
		attempts++
		if attempts < 2 {
			return &ProviderError{Provider: "openai", Code: ProviderErrorUnavailable, Retryable: true}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestWithRequestTimeoutUsesRequestOverride(t *testing.T) {
	ctx, cancel := WithRequestTimeout(context.Background(), RequestOptions{TimeoutMs: 5}, time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		t.Fatal("timeout fired too early")
	case <-time.After(1 * time.Millisecond):
	}
}

func TestRetryDoesNotRetryNonProviderErrors(t *testing.T) {
	attempts := 0
	err := Retry(context.Background(), RetryPolicy{MaxAttempts: 3}, func(ctx context.Context) error {
		attempts++
		return errors.New("boom")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}
