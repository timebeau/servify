package telemetry

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFieldsFromContext_AllFields(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-1")
	ctx = WithSessionID(ctx, "sess-2")
	ctx = WithTenantID(ctx, "tenant-3")

	fields := FieldsFromContext(ctx)

	if fields[FieldRequestID] != "req-1" {
		t.Fatalf("expected req-1, got %v", fields[FieldRequestID])
	}
	if fields[FieldSessionID] != "sess-2" {
		t.Fatalf("expected sess-2, got %v", fields[FieldSessionID])
	}
	if fields[FieldTenantID] != "tenant-3" {
		t.Fatalf("expected tenant-3, got %v", fields[FieldTenantID])
	}
}

func TestFieldsFromContext_PartialFields(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-1")

	fields := FieldsFromContext(ctx)

	if fields[FieldRequestID] != "req-1" {
		t.Fatalf("expected req-1, got %v", fields[FieldRequestID])
	}
	if _, exists := fields[FieldSessionID]; exists {
		t.Fatal("expected session_id to be absent")
	}
	if _, exists := fields[FieldTenantID]; exists {
		t.Fatal("expected tenant_id to be absent")
	}
}

func TestFieldsFromContext_Empty(t *testing.T) {
	fields := FieldsFromContext(context.Background())
	if len(fields) != 0 {
		t.Fatalf("expected 0 fields, got %d", len(fields))
	}
}

func TestLoggerWithRequestID(t *testing.T) {
	logger := logrus.New()
	entry := LoggerWithRequestID(logger, "req-abc")

	if entry.Data[FieldRequestID] != "req-abc" {
		t.Fatalf("expected req-abc, got %v", entry.Data[FieldRequestID])
	}
}
