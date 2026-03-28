package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewAppError(t *testing.T) {
	cause := errors.New("something went wrong")
	appErr := New(cause, SeverityUser, CategoryValidation,
		WithModule("ticket"),
		WithMessage("invalid ticket ID"),
		WithCode("TICKET_INVALID_ID"),
	)

	if appErr.Severity != SeverityUser {
		t.Fatalf("expected user severity, got %s", appErr.Severity)
	}
	if appErr.Category != CategoryValidation {
		t.Fatalf("expected validation category, got %s", appErr.Category)
	}
	if appErr.Module != "ticket" {
		t.Fatalf("expected ticket module, got %s", appErr.Module)
	}
	if appErr.Code != "TICKET_INVALID_ID" {
		t.Fatalf("expected TICKET_INVALID_ID code, got %s", appErr.Code)
	}
	if appErr.HTTPStatus != 400 {
		t.Fatalf("expected 400 status, got %d", appErr.HTTPStatus)
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	appErr := New(cause, SeveritySystem, CategoryInternal)

	if !errors.Is(appErr, cause) {
		t.Fatal("expected errors.Is to match root cause")
	}
}

func TestAppError_NilInput(t *testing.T) {
	appErr := New(nil, SeverityUser, CategoryValidation)
	if appErr != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestAppError_DefaultHTTPStatus(t *testing.T) {
	tests := []struct {
		severity     Severity
		expectedCode int
	}{
		{SeverityUser, 400},
		{SeverityDependency, 502},
		{SeverityConfig, 500},
		{SeveritySystem, 500},
	}
	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			appErr := New(errors.New("test"), tt.severity, CategoryInternal)
			if appErr.HTTPStatus != tt.expectedCode {
				t.Fatalf("expected %d for %s, got %d", tt.expectedCode, tt.severity, appErr.HTTPStatus)
			}
		})
	}
}

func TestAppError_WithHTTPStatus(t *testing.T) {
	appErr := New(errors.New("test"), SeverityUser, CategoryAuth, WithHTTPStatus(403))
	if appErr.HTTPStatus != 403 {
		t.Fatalf("expected 403, got %d", appErr.HTTPStatus)
	}
}

func TestAppError_ErrorString(t *testing.T) {
	cause := errors.New("db connection failed")
	appErr := New(cause, SeverityDependency, CategoryDatabase, WithModule("infra"))
	expected := "[dependency:database] infra: db connection failed"
	if appErr.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, appErr.Error())
	}
}

func TestClassify_PlainError(t *testing.T) {
	err := errors.New("unknown error")
	appErr := Classify(err)
	if appErr.Severity != SeveritySystem {
		t.Fatalf("expected system severity, got %s", appErr.Severity)
	}
	if appErr.Category != CategoryInternal {
		t.Fatalf("expected internal category, got %s", appErr.Category)
	}
}

func TestClassify_NilError(t *testing.T) {
	appErr := Classify(nil)
	if appErr != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestClassify_AlreadyAppError(t *testing.T) {
	original := New(errors.New("test"), SeverityUser, CategoryValidation)
	classified := Classify(original)
	if classified != original {
		t.Fatal("expected same AppError instance")
	}
}

func TestClassify_WrappedError(t *testing.T) {
	inner := errors.New("base error")
	wrapped := fmt.Errorf("outer: %w", inner)
	appErr := Classify(wrapped)
	if appErr == nil {
		t.Fatal("expected non-nil AppError")
	}
	if appErr.Err != wrapped {
		t.Fatal("expected outer error to be preserved")
	}
}
