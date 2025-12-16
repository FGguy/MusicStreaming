package ports

import (
	"context"
	"fmt"
	"music-streaming/internal/core/domain"
)

type ContextKey int

const (
	KeyRequestingUserID ContextKey = iota
)

type UserManagementPort interface {
	CreateUser(ctx context.Context, user domain.User) error
	UpdateUser(ctx context.Context, username string, user domain.User) error
	DeleteUser(ctx context.Context, username string) error
	GetUser(ctx context.Context, username string) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	ChangePassword(ctx context.Context, username string, newPassword string) error
}

type UserManagementRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	UpdateUser(ctx context.Context, username string, user domain.User) error
	DeleteUser(ctx context.Context, username string) error
	GetUser(ctx context.Context, username string) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
}

type UserNotFoundError struct {
	Username string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("No user with username: %s", e.Username)
}

type MissingOrInvalidParameterError struct {
	ParameterName string
}

func (e *MissingOrInvalidParameterError) Error() string {
	return fmt.Sprintf("Missing or invalid parameter: %s", e.ParameterName)
}

type FailedUserOperationError struct {
	Description string
}

func (e *FailedUserOperationError) Error() string {
	return e.Description
}

type UserNotAuthorizedError struct {
	Username string
	Action   string
}

func (e *UserNotAuthorizedError) Error() string {
	return fmt.Sprintf("User %s is not authorized to perform action: %s", e.Username, e.Action)
}
