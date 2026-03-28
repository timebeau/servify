// Package errors provides a unified error classification system for Servify.
// It wraps errors with severity, category, and HTTP status mapping to enable
// consistent observability across all modules.
package errors

import (
	"fmt"
)

// Severity indicates who is primarily affected by the error.
type Severity string

const (
	// SeverityUser indicates a user-facing, recoverable error (bad input, not found).
	SeverityUser Severity = "user"
	// SeverityDependency indicates an external service failure (LLM timeout, DB down).
	SeverityDependency Severity = "dependency"
	// SeverityConfig indicates a misconfiguration (bad credentials, missing env).
	SeverityConfig Severity = "config"
	// SeveritySystem indicates an internal invariant violation (nil pointer, panic).
	SeveritySystem Severity = "system"
)

// Category groups errors by subsystem.
type Category string

const (
	CategoryAuth       Category = "auth"
	CategoryDatabase   Category = "database"
	CategoryAI         Category = "ai"
	CategoryRouting    Category = "routing"
	CategoryValidation Category = "validation"
	CategoryRateLimit  Category = "rate_limit"
	CategoryInternal   Category = "internal"
	CategoryNetwork    Category = "network"
)

// AppError is the standard observability-aware error type.
// It wraps an underlying error with severity, category, and HTTP mapping.
type AppError struct {
	Err        error
	Severity   Severity
	Category   Category
	Module     string
	Message    string
	Code       string
	HTTPStatus int
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Module != "" {
		return fmt.Sprintf("[%s:%s] %s: %s", e.Severity, e.Category, e.Module, e.Err.Error())
	}
	return fmt.Sprintf("[%s:%s] %s", e.Severity, e.Category, e.Err.Error())
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Option configures an AppError.
type Option func(*AppError)

// WithModule sets the originating module.
func WithModule(module string) Option {
	return func(e *AppError) { e.Module = module }
}

// WithMessage sets a user-safe message.
func WithMessage(msg string) Option {
	return func(e *AppError) { e.Message = msg }
}

// WithCode sets a machine-readable error code.
func WithCode(code string) Option {
	return func(e *AppError) { e.Code = code }
}

// WithHTTPStatus sets the suggested HTTP response status.
func WithHTTPStatus(status int) Option {
	return func(e *AppError) { e.HTTPStatus = status }
}

// New creates a new AppError wrapping the given error.
func New(err error, severity Severity, category Category, opts ...Option) *AppError {
	if err == nil {
		return nil
	}
	e := &AppError{
		Err:      err,
		Severity: severity,
		Category: category,
		Message:  err.Error(),
	}
	for _, opt := range opts {
		opt(e)
	}
	// Default HTTP status based on severity
	if e.HTTPStatus == 0 {
		e.HTTPStatus = severityToHTTPStatus(severity)
	}
	return e
}

func severityToHTTPStatus(s Severity) int {
	switch s {
	case SeverityUser:
		return 400
	case SeverityDependency:
		return 502
	case SeverityConfig:
		return 500
	case SeveritySystem:
		return 500
	default:
		return 500
	}
}
