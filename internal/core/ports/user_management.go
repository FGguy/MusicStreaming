package ports

import (
	"context"
	"music-streaming/internal/core/domain"
)

// ContextKey is the type used for context keys in this application.
type ContextKey int

const (
	// KeyRequestingUserID is the context key for storing the requesting user.
	KeyRequestingUserID ContextKey = iota
)

// UserManagementPort defines the interface for user management operations.
// It provides methods for CRUD operations on users and password management.
type UserManagementPort interface {
	// CreateUser creates a new user in the system. Requires admin role.
	CreateUser(ctx context.Context, user domain.User) error

	// UpdateUser updates an existing user's information. Requires admin role or self-update with settings role.
	UpdateUser(ctx context.Context, username string, user domain.User) error

	// DeleteUser removes a user from the system. Requires admin role.
	DeleteUser(ctx context.Context, username string) error

	// GetUser retrieves a specific user by username. Requires admin role or self-query.
	GetUser(ctx context.Context, username string) (domain.User, error)

	// GetUsers retrieves all users in the system. Requires admin role.
	GetUsers(ctx context.Context) ([]domain.User, error)

	// ChangePassword changes a user's password. Requires admin role or self-update.
	ChangePassword(ctx context.Context, username string, newPassword string) error
}

// UserManagementRepository defines the interface for user data persistence.
// Implementations of this interface provide data access operations for users.
type UserManagementRepository interface {
	// CreateUser persists a new user to the data store.
	CreateUser(ctx context.Context, user domain.User) error

	// UpdateUser persists changes to an existing user.
	UpdateUser(ctx context.Context, username string, user domain.User) error

	// DeleteUser removes a user from the data store.
	DeleteUser(ctx context.Context, username string) error

	// GetUser retrieves a specific user by username from the data store.
	GetUser(ctx context.Context, username string) (domain.User, error)

	// GetUsers retrieves all users from the data store.
	GetUsers(ctx context.Context) ([]domain.User, error)
}
