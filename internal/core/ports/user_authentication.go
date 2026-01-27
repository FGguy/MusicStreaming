package ports

import (
	"context"
	"fmt"
	"music-streaming/internal/core/domain"
)

// UserAuthenticationPort defines the interface for user authentication operations.
// Implementations of this port handle user credential validation and authentication.
type UserAuthenticationPort interface {
	// AuthenticateUser validates user credentials using MD5 hash with salt.
	// Returns the authenticated user or an error if authentication fails.
	AuthenticateUser(ctx context.Context, username, password, salt string) (domain.User, error)
}

// FailedAuthenticationError indicates that user authentication failed due to invalid credentials.
// It maps to Subsonic API error code 40.
type FailedAuthenticationError struct {
	Username string
}

// Error implements the error interface for FailedAuthenticationError.
func (e *FailedAuthenticationError) Error() string {
	return fmt.Sprintf("Authentication failed for user: %s", e.Username)
}
