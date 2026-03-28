package telemetry

import (
	"context"

	"github.com/sirupsen/logrus"
)

// FieldsFromContext builds logrus.Fields from observability context values.
// Only non-empty fields are included.
func FieldsFromContext(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}

	if id := RequestIDFromContext(ctx); id != "" {
		fields[FieldRequestID] = id
	}
	if id := SessionIDFromContext(ctx); id != "" {
		fields[FieldSessionID] = id
	}
	if id := TenantIDFromContext(ctx); id != "" {
		fields[FieldTenantID] = id
	}
	if id := TraceIDFromContext(ctx); id != "" {
		fields[FieldTraceID] = id
	}
	if id := SpanIDFromContext(ctx); id != "" {
		fields[FieldSpanID] = id
	}

	return fields
}

// LoggerWithFields returns a logrus.Entry with observability fields from the context.
func LoggerWithFields(logger *logrus.Logger, ctx context.Context) *logrus.Entry {
	return logger.WithFields(FieldsFromContext(ctx))
}

// LoggerWithRequestID returns a logrus.Entry with only the request_id field set.
func LoggerWithRequestID(logger *logrus.Logger, requestID string) *logrus.Entry {
	return logger.WithField(FieldRequestID, requestID)
}
