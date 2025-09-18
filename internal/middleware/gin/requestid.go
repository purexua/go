package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/purexua/go/internal/known"
	"github.com/purexua/go/pkg/contextx"
)

// RequestIDMiddleware is a Gin middleware that injects `x-request-id` key-value pairs
// into the context and response of each HTTP request.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get `x-request-id` from request headers, generate new UUID if not present
		requestID := c.Request.Header.Get(known.XRequestID)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Save RequestID to context.Context for subsequent program use
		ctx := contextx.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// Save RequestID to HTTP response headers with key `x-request-id`
		c.Writer.Header().Set(known.XRequestID, requestID)

		// Continue processing the request
		c.Next()
	}
}
