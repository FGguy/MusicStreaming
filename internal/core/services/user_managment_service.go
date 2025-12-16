package services

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type UserManagementService struct {
	repo ports.UserManagementRepository
}

func NewUserManagementService(repo ports.UserManagementRepository) *UserManagementService {
	return &UserManagementService{
		repo: repo,
	}
}

func (s *UserManagementService) CreateUser(ctx context.Context, user domain.User) error {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)

	// Permission checking
	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		return &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "create user"}
	}

	// Parameter validation
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return &ports.MissingOrInvalidParameterError{ParameterName: "username, email or password"}
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *UserManagementService) UpdateUser(ctx context.Context, username string, user domain.User) error {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || (!requestingUser.AdminRole && !(requestingUser.Username == username && requestingUser.SettingsRole)) {
		return &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "update user"}
	}

	if user.Username == "" || user.Email == "" || user.Password == "" || username == "" {
		return &ports.MissingOrInvalidParameterError{ParameterName: "username, email or password"}
	}

	return s.repo.UpdateUser(ctx, username, user)
}

func (s *UserManagementService) DeleteUser(ctx context.Context, username string) error {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		return &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "delete user"}
	}

	if username == "" {
		return &ports.MissingOrInvalidParameterError{ParameterName: "username"}
	}

	return s.repo.DeleteUser(ctx, username)
}

func (s *UserManagementService) GetUser(ctx context.Context, username string) (domain.User, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || (!requestingUser.AdminRole && requestingUser.Username != username) {
		return domain.User{}, &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "get user"}
	}

	if username == "" {
		return domain.User{}, &ports.MissingOrInvalidParameterError{ParameterName: "username"}
	}

	return s.repo.GetUser(ctx, username)
}

func (s *UserManagementService) GetUsers(ctx context.Context) ([]domain.User, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		return make([]domain.User, 0), &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "get users"}
	}

	return s.repo.GetUsers(ctx)
}

func (s *UserManagementService) ChangePassword(ctx context.Context, username string, newPassword string) error {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || (!requestingUser.AdminRole && requestingUser.Username != username) {
		return &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "change password"}
	}

	if username == "" || newPassword == "" {
		return &ports.MissingOrInvalidParameterError{ParameterName: "username or new password"}
	}

	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return err
	}
	user.Password = newPassword

	return s.repo.UpdateUser(ctx, username, user)
}
