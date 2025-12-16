package ports

import (
	"context"
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
