package llm

import (
	"fmt"
	"net/http"
)

type ProviderErrorCode string

const (
	ProviderErrorTimeout      ProviderErrorCode = "timeout"
	ProviderErrorUnavailable  ProviderErrorCode = "unavailable"
	ProviderErrorRateLimited  ProviderErrorCode = "rate_limited"
	ProviderErrorAuthFailed   ProviderErrorCode = "auth_failed"
	ProviderErrorInvalid      ProviderErrorCode = "invalid_request"
	ProviderErrorUpstream     ProviderErrorCode = "upstream_error"
	ProviderErrorNotSupported ProviderErrorCode = "not_supported"
)

type ProviderError struct {
	Provider   string
	Code       ProviderErrorCode
	Message    string
	StatusCode int
	Retryable  bool
	Cause      error
}

func (e *ProviderError) Error() string {
	if e == nil {
		return ""
	}
	if e.StatusCode > 0 {
		return fmt.Sprintf("%s provider error (%s:%d): %s", e.Provider, e.Code, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("%s provider error (%s): %s", e.Provider, e.Code, e.Message)
}

func (e *ProviderError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func HTTPError(provider string, statusCode int, message string) *ProviderError {
	code := ProviderErrorUpstream
	retryable := false

	switch {
	case statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden:
		code = ProviderErrorAuthFailed
	case statusCode == http.StatusTooManyRequests:
		code = ProviderErrorRateLimited
		retryable = true
	case statusCode == http.StatusRequestTimeout || statusCode == http.StatusGatewayTimeout:
		code = ProviderErrorTimeout
		retryable = true
	case statusCode >= 500:
		code = ProviderErrorUnavailable
		retryable = true
	case statusCode >= 400:
		code = ProviderErrorInvalid
	}

	return &ProviderError{
		Provider:   provider,
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Retryable:  retryable,
	}
}
