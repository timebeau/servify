package telemetry

import (
	"context"
	"testing"
)

func TestRequestIDRoundTrip(t *testing.T) {
	ctx := context.Background()

	// Empty when not set
	if got := RequestIDFromContext(ctx); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}

	// Set and retrieve
	ctx = WithRequestID(ctx, "req-123")
	if got := RequestIDFromContext(ctx); got != "req-123" {
		t.Fatalf("expected req-123, got %q", got)
	}
}

func TestSessionIDRoundTrip(t *testing.T) {
	ctx := context.Background()

	if got := SessionIDFromContext(ctx); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}

	ctx = WithSessionID(ctx, "sess-456")
	if got := SessionIDFromContext(ctx); got != "sess-456" {
		t.Fatalf("expected sess-456, got %q", got)
	}
}

func TestTenantIDRoundTrip(t *testing.T) {
	ctx := context.Background()

	if got := TenantIDFromContext(ctx); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}

	ctx = WithTenantID(ctx, "tenant-789")
	if got := TenantIDFromContext(ctx); got != "tenant-789" {
		t.Fatalf("expected tenant-789, got %q", got)
	}
}

func TestNilContext(t *testing.T) {
	if got := RequestIDFromContext(nil); got != "" {
		t.Fatalf("expected empty for nil context, got %q", got)
	}
	if got := SessionIDFromContext(nil); got != "" {
		t.Fatalf("expected empty for nil context, got %q", got)
	}
	if got := TenantIDFromContext(nil); got != "" {
		t.Fatalf("expected empty for nil context, got %q", got)
	}
}

func TestTraceIDFromContext_NoSpan(t *testing.T) {
	ctx := context.Background()
	if got := TraceIDFromContext(ctx); got != "" {
		t.Fatalf("expected empty with no span, got %q", got)
	}
}
