package repositories

import (
	"context"
	"fmt"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"sync"
)

type InMemoryUserManagementRepository struct {
	users map[string]domain.User
	mu    sync.RWMutex
}

func NewInMemoryUserManagementRepository() *InMemoryUserManagementRepository {
	return &InMemoryUserManagementRepository{
		users: make(map[string]domain.User),
	}
}

func (r *InMemoryUserManagementRepository) CreateUser(ctx context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Username]; exists {
		return &ports.FailedOperationError{Description: "User already exists"}
	}
	r.users[user.Username] = user
	return nil
}

func (r *InMemoryUserManagementRepository) UpdateUser(ctx context.Context, username string, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[username]; !exists {
		return &ports.NotFoundError{Message: fmt.Sprintf("User %s not found", username)}
	}
	r.users[username] = user
	return nil
}

func (r *InMemoryUserManagementRepository) DeleteUser(ctx context.Context, username string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[username]; !exists {
		return &ports.NotFoundError{Message: fmt.Sprintf("User %s not found", username)}
	}
	delete(r.users, username)
	return nil
}

func (r *InMemoryUserManagementRepository) GetUser(ctx context.Context, username string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[username]
	if !exists {
		return domain.User{}, &ports.NotFoundError{Message: fmt.Sprintf("User %s not found", username)}
	}
	return user, nil
}

func (r *InMemoryUserManagementRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}
