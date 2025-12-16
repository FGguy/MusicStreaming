package repositories

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type InMemoryUserManagementRepository struct {
	users map[string]domain.User
}

func NewInMemoryUserManagementRepository() *InMemoryUserManagementRepository {
	return &InMemoryUserManagementRepository{
		users: make(map[string]domain.User),
	}
}

func (r *InMemoryUserManagementRepository) CreateUser(ctx context.Context, user domain.User) error {
	if _, exists := r.users[user.Username]; exists {
		return &ports.FailedUserOperationError{Description: "User already exists"}
	}
	r.users[user.Username] = user
	return nil
}

func (r *InMemoryUserManagementRepository) UpdateUser(ctx context.Context, username string, user domain.User) error {
	if _, exists := r.users[username]; !exists {
		return &ports.UserNotFoundError{Username: username}
	}
	r.users[username] = user
	return nil
}

func (r *InMemoryUserManagementRepository) DeleteUser(ctx context.Context, username string) error {
	if _, exists := r.users[username]; !exists {
		return &ports.UserNotFoundError{Username: username}
	}
	delete(r.users, username)
	return nil
}

func (r *InMemoryUserManagementRepository) GetUser(ctx context.Context, username string) (domain.User, error) {
	user, exists := r.users[username]
	if !exists {
		return domain.User{}, &ports.UserNotFoundError{Username: username}
	}
	return user, nil
}

func (r *InMemoryUserManagementRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	users := make([]domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}
