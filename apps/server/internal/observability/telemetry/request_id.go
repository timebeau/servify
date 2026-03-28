package telemetry

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestIDMiddleware generates a UUID request ID if not already present,
// sets it in the response header, and stores it in context for downstream
// consumers (logging, audit, metrics).
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}

		// Set in response header for client-side correlation.
		c.Header(requestIDHeader, id)

		// Store in Gin context for other middleware.
		c.Set(FieldRequestID, id)

		// Store in context.Context for non-Gin code paths.
		ctx := WithRequestID(c.Request.Context(), id)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
