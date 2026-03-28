package errors

import (
	"fmt"

	llm "servify/apps/server/internal/platform/llm"
)

// Classify wraps any error into an AppError using classification rules.
// If the error is already an *AppError, it returns it unchanged.
// If the error is an *llm.ProviderError, it maps from ProviderErrorCode.
// Otherwise it falls back to SeveritySystem + CategoryInternal.
func Classify(err error) *AppError {
	if err == nil {
		return nil
	}

	// Already classified
	var appErr *AppError
	if As(err, &appErr) {
		return appErr
	}

	// LLM provider error
	var providerErr *llm.ProviderError
	if As(err, &providerErr) {
		return classifyProviderError(providerErr)
	}

	// Unwrap and try again for wrapped errors
	if unwrapped := Unwrap(err); unwrapped != nil {
		if inner := Classify(unwrapped); inner != nil {
			// Preserve the outer error message context
			return &AppError{
				Err:        err,
				Severity:   inner.Severity,
				Category:   inner.Category,
				Module:     inner.Module,
				Message:    inner.Message,
				Code:       inner.Code,
				HTTPStatus: inner.HTTPStatus,
			}
		}
	}

	// Default: system error
	return New(err, SeveritySystem, CategoryInternal)
}

func classifyProviderError(err *llm.ProviderError) *AppError {
	mapping := map[llm.ProviderErrorCode]struct {
		severity Severity
		category Category
		status   int
	}{
		llm.ProviderErrorTimeout:      {SeverityDependency, CategoryAI, 504},
		llm.ProviderErrorRateLimited:  {SeverityDependency, CategoryAI, 429},
		llm.ProviderErrorAuthFailed:   {SeverityConfig, CategoryAuth, 401},
		llm.ProviderErrorUnavailable:  {SeverityDependency, CategoryAI, 502},
		llm.ProviderErrorInvalid:      {SeverityUser, CategoryValidation, 400},
		llm.ProviderErrorUpstream:     {SeverityDependency, CategoryNetwork, 502},
		llm.ProviderErrorNotSupported: {SeverityUser, CategoryValidation, 400},
	}

	m, ok := mapping[err.Code]
	if !ok {
		return New(err, SeverityDependency, CategoryAI,
			WithModule("ai"),
			WithHTTPStatus(502),
		)
	}

	return New(err, m.severity, m.category,
		WithModule("ai"),
		WithCode(string(err.Code)),
		WithHTTPStatus(m.status),
	)
}

// As is a wrapper around errors.As for convenience.
var As = func(err error, target interface{}) bool {
	// Use standard errors.As
	type errorsAs interface{ As(interface{}) bool }
	return fmt.Sprintf("%T", err) != "" && wrapAs(err, target)
}

// Unwrap is a wrapper around errors.Unwrap for convenience.
var Unwrap = func(err error) error {
	type unwrapper interface{ Unwrap() error }
	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}
	return nil
}

func wrapAs(err error, target interface{}) bool {
	// Simple type assertion for *AppError
	if t, ok := target.(**AppError); ok {
		if e, ok := err.(*AppError); ok {
			*t = e
			return true
		}
	}
	// Simple type assertion for *llm.ProviderError
	if t, ok := target.(**llm.ProviderError); ok {
		if e, ok := err.(*llm.ProviderError); ok {
			*t = e
			return true
		}
	}
	// Try unwrapping
	if unwrapped := Unwrap(err); unwrapped != nil {
		return wrapAs(unwrapped, target)
	}
	return false
}
