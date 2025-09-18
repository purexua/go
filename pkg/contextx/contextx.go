// Package contextx provides utilities for managing context values including
// user information, access tokens, and request IDs across the application.
package contextx

import (
	"context"
)

// Context key type definitions.
type (
	// usernameKey defines the context key for username.
	usernameKey struct{}
	// userIDKey defines the context key for user ID.
	userIDKey struct{}
	// accessTokenKey defines the context key for access token.
	accessTokenKey struct{}
	// requestIDKey defines the context key for request ID.
	requestIDKey struct{}
)

// WithUserID stores the user ID in the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserID extracts the user ID from the context.
func UserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey{}).(string)
	return userID
}

// WithUsername stores the username in the context.
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey{}, username)
}

// Username extracts the username from the context.
func Username(ctx context.Context) string {
	username, _ := ctx.Value(usernameKey{}).(string)
	return username
}

// WithAccessToken stores the access token in the context.
func WithAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, accessTokenKey{}, accessToken)
}

// AccessToken extracts the access token from the context.
func AccessToken(ctx context.Context) string {
	accessToken, _ := ctx.Value(accessTokenKey{}).(string)
	return accessToken
}

// WithRequestID stores the request ID in the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// RequestID extracts the request ID from the context.
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}
