// Package core provides utility functions for data copying and type conversion,
// particularly for handling time.Time to string conversions and vice versa.
package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/purexua/go/pkg/errorsx"
)

// ErrorResponse defines the error response structure.
type ErrorResponse struct {
	// The reason for the error, identifying the error type.
	Reason string `json:"reason,omitempty"`
	// The detailed description of the error.
	Message string `json:"message,omitempty"`
	// The accompanying metadata information.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// WriteResponse writes the response to the client.
func WriteResponse(c *gin.Context, data any, err error) {
	if err != nil {
		// If an error occurs, generate an error response
		errx := errorsx.FromError(err) // Extract error details
		c.JSON(errx.Code, ErrorResponse{
			Reason:   errx.Reason,
			Message:  errx.Message,
			Metadata: errx.Metadata,
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
