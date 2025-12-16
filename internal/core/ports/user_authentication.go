package ports

import (
	"context"
	"fmt"
	"music-streaming/internal/core/domain"
)

type UserAuthenticationPort interface {
	AuthenticateUser(ctx context.Context, username, password, salt string) (domain.User, error)
}

type FailedAuthenticationError struct {
	Username string
}

func (e *FailedAuthenticationError) Error() string {
	return fmt.Sprintf("Authentication failed for user: %s", e.Username)
}
