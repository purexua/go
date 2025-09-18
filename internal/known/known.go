// Package known defines common constants used across the application,
// including HTTP/gRPC headers and other shared values.
package known

// HTTP/gRPC Header definitions.
// gRPC uses HTTP/2 as the underlying transport protocol, and HTTP/2 specification
// requires header keys to be lowercase. Therefore, in gRPC, all header keys are
// forced to be converted to lowercase to comply with HTTP/2 requirements.
// In HTTP/1.x, many implementations preserve user-set case formatting,
// but some HTTP frameworks or utility libraries (such as certain web servers or proxies)
// may automatically convert headers to lowercase to simplify processing logic.
// For compatibility, headers are uniformly set to lowercase here.
// Additionally, header keys starting with x- indicate custom headers.
const (
	// XRequestID defines the context key representing the request ID.
	XRequestID = "x-request-id"

	// XUserID defines the context key representing the request user ID. UserID is unique throughout the user's lifecycle.
	XUserID = "x-user-id"

	// XUsername defines the context key representing the request username.
	XUsername = "x-username"
)

// Other constant definitions.
const (
	// AdminUsername defines the admin username.
	AdminUsername = "root"

	// MaxErrGroupConcurrency defines the maximum number of concurrent tasks for errgroup.
	// Used to limit the number of Goroutines executing simultaneously in errgroup,
	// thereby preventing resource exhaustion and improving program stability.
	// This value can be adjusted based on scenario requirements.
	MaxErrGroupConcurrency = 1000
)
