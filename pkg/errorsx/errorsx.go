// Package errorsx provides a structured error handling system with HTTP status codes,
// gRPC integration, and metadata support for enhanced error tracking and debugging.
package errorsx

import (
	"errors"
	"fmt"
	"net/http"

	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// ErrorX defines the error type used in the project to describe detailed error information.
type ErrorX struct {
	// Code represents the HTTP status code of the error, used to identify the error type when interacting with clients.
	Code int `json:"code,omitempty"`

	// Reason represents the cause of the error, usually a business error code, used for precise problem location.
	Reason string `json:"reason,omitempty"`

	// Message represents a brief error message, usually can be directly exposed to users.
	Message string `json:"message,omitempty"`

	// Metadata stores additional meta-information related to the error, which can contain context or debugging information.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// New creates a new error.
func New(code int, reason string, format string, args ...any) *ErrorX {
	return &ErrorX{
		Code:    code,
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
	}
}

// Error implements the `Error` method in the error interface.
func (err *ErrorX) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", err.Code, err.Reason, err.Message, err.Metadata)
}

// WithMessage sets the Message field of the error.
func (err *ErrorX) WithMessage(format string, args ...any) *ErrorX {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

// WithMetadata sets the metadata.
func (err *ErrorX) WithMetadata(md map[string]string) *ErrorX {
	err.Metadata = md
	return err
}

// KV sets metadata using key-value pairs.
func (err *ErrorX) KV(kvs ...string) *ErrorX {
	if err.Metadata == nil {
		err.Metadata = make(map[string]string) // Initialize metadata map
	}

	for i := 0; i < len(kvs); i += 2 {
		// kvs must be in pairs
		if i+1 < len(kvs) {
			err.Metadata[kvs[i]] = kvs[i+1]
		}
	}
	return err
}

// GRPCStatus returns the gRPC status representation.
func (err *ErrorX) GRPCStatus() *status.Status {
	details := errdetails.ErrorInfo{Reason: err.Reason, Metadata: err.Metadata}
	s, _ := status.New(httpstatus.ToGRPCCode(err.Code), err.Message).WithDetails(&details)
	return s
}

// WithRequestID sets the request ID.
func (err *ErrorX) WithRequestID(requestID string) *ErrorX {
	return err.KV("X-Request-ID", requestID) // Set request ID
}

// Is determines whether the current error matches the target error.
// It recursively traverses the error chain and compares the Code and Reason fields of ErrorX instances.
// Returns true if both Code and Reason are equal; otherwise returns false.
func (err *ErrorX) Is(target error) bool {
	if errx := new(ErrorX); errors.As(target, &errx) {
		return errx.Code == err.Code && errx.Reason == err.Reason
	}
	return false
}

// Code returns the HTTP code of the error.
func Code(err error) int {
	if err == nil {
		return http.StatusOK //nolint:mnd
	}
	return FromError(err).Code
}

// Reason returns the reason of the specific error.
func Reason(err error) string {
	if err == nil {
		return ErrInternal.Reason
	}
	return FromError(err).Reason
}

// FromError attempts to convert a generic error to a custom *ErrorX type.
func FromError(err error) *ErrorX {
	// If the passed error is nil, return nil directly, indicating no error needs to be handled.
	if err == nil {
		return nil
	}

	// Check if the passed error is already an instance of ErrorX type.
	// If the error can be converted to *ErrorX type through errors.As, return that instance directly.
	if errx := new(ErrorX); errors.As(err, &errx) {
		return errx
	}

	// gRPC's status.FromError method attempts to convert error to a gRPC error status object.
	// If err cannot be converted to a gRPC error (i.e., not a gRPC status error),
	// return an ErrorX with default values, indicating it's an unknown type of error.
	gs, ok := status.FromError(err)
	if !ok {
		return New(ErrInternal.Code, ErrInternal.Reason, "%s", err.Error())
	}

	// If err is a gRPC error type, it will successfully return a gRPC status object (gs).
	// Create an ErrorX using the error code and message from the gRPC status.
	ret := New(httpstatus.FromGRPCCode(gs.Code()), ErrInternal.Reason, "%s", gs.Message())

	// Traverse all additional information (Details) in the gRPC error details.
	for _, detail := range gs.Details() {
		if typed, ok := detail.(*errdetails.ErrorInfo); ok {
			ret.Reason = typed.Reason
			return ret.WithMetadata(typed.Metadata)
		}
	}

	return ret
}
