package services

import (
	"context"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type UserManagementService struct {
	repo   ports.UserManagementRepository
	logger *slog.Logger
}

func NewUserManagementService(repo ports.UserManagementRepository, logger *slog.Logger) *UserManagementService {
	return &UserManagementService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserManagementService) CreateUser(ctx context.Context, user domain.User) error {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var username string
	if requestingUser != nil {
		username = requestingUser.Username
	}

	s.logger.Info("Create user request", slog.String("requesting_user", username), slog.String("target_username", user.Username))

	// Permission checking
	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		s.logger.Warn("Unauthorized create user attempt", slog.String("requesting_user", username), slog.String("target_username", user.Username))
		return &ports.NotAuthorizedError{Username: username, Action: "create user"}
	}

	// Parameter validation
	if user.Username == "" || user.Email == "" || user.Password == "" {
		s.logger.Warn("Invalid parameters for create user", slog.String("requesting_user", username))
		return &ports.MissingOrInvalidParameterError{ParameterName: "username, email or password"}
	}

	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("Failed to create user", slog.String("requesting_user", username), slog.String("target_username", user.Username), slog.String("error", err.Error()))
		return err
	}
	s.logger.Info("User created successfully", slog.String("requesting_user", username), slog.String("target_username", user.Username))
	return err
}

func (s *UserManagementService) UpdateUser(ctx context.Context, username string, user domain.User) error {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var requestingUsername string
	if requestingUser != nil {
		requestingUsername = requestingUser.Username
	}
	s.logger.Info("Update user request", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))

	if !ok || requestingUser == nil || (!requestingUser.AdminRole && (requestingUser.Username != username || !requestingUser.SettingsRole)) {
		s.logger.Warn("Unauthorized update user attempt", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
		return &ports.NotAuthorizedError{Username: requestingUsername, Action: "update user"}
	}

	if user.Username == "" || user.Email == "" || user.Password == "" || username == "" {
		s.logger.Warn("Invalid parameters for update user", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
		return &ports.MissingOrInvalidParameterError{ParameterName: "username, email or password"}
	}

	err := s.repo.UpdateUser(ctx, username, user)
	if err != nil {
		s.logger.Error("Failed to update user", slog.String("requesting_user", requestingUsername), slog.String("target_username", username), slog.String("error", err.Error()))
		return err
	}
	s.logger.Info("User updated successfully", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
	return err
}

func (s *UserManagementService) DeleteUser(ctx context.Context, username string) error {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var requestingUsername string
	if requestingUser != nil {
		requestingUsername = requestingUser.Username
	}
	s.logger.Info("Delete user request", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))

	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		s.logger.Warn("Unauthorized delete user attempt", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
		return &ports.NotAuthorizedError{Username: requestingUsername, Action: "delete user"}
	}

	if username == "" {
		s.logger.Warn("Invalid parameters for delete user", slog.String("requesting_user", requestingUsername))
		return &ports.MissingOrInvalidParameterError{ParameterName: "username"}
	}

	err := s.repo.DeleteUser(ctx, username)
	if err != nil {
		s.logger.Error("Failed to delete user", slog.String("requesting_user", requestingUsername), slog.String("target_username", username), slog.String("error", err.Error()))
		return err
	}
	s.logger.Info("User deleted successfully", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
	return err
}

func (s *UserManagementService) GetUser(ctx context.Context, username string) (domain.User, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var requestingUsername string
	if requestingUser != nil {
		requestingUsername = requestingUser.Username
	}
	s.logger.Info("Get user request", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))

	if !ok || requestingUser == nil || (!requestingUser.AdminRole && requestingUser.Username != username) {
		s.logger.Warn("Unauthorized get user attempt", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
		return domain.User{}, &ports.NotAuthorizedError{Username: requestingUsername, Action: "get user"}
	}

	if username == "" {
		s.logger.Warn("Invalid parameters for get user", slog.String("requesting_user", requestingUsername))
		return domain.User{}, &ports.MissingOrInvalidParameterError{ParameterName: "username"}
	}

	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		s.logger.Error("Failed to get user", slog.String("requesting_user", requestingUsername), slog.String("target_username", username), slog.String("error", err.Error()))
		return user, err
	}
	s.logger.Info("User retrieved successfully", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
	return user, err
}

func (s *UserManagementService) GetUsers(ctx context.Context) ([]domain.User, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var requestingUsername string
	if requestingUser != nil {
		requestingUsername = requestingUser.Username
	}
	s.logger.Info("Get users request", slog.String("requesting_user", requestingUsername))

	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		s.logger.Warn("Unauthorized get users attempt", slog.String("requesting_user", requestingUsername))
		return make([]domain.User, 0), &ports.NotAuthorizedError{Username: requestingUsername, Action: "get users"}
	}

	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		s.logger.Error("Failed to get users", slog.String("requesting_user", requestingUsername), slog.String("error", err.Error()))
		return users, err
	}
	s.logger.Info("Users retrieved successfully", slog.String("requesting_user", requestingUsername), slog.Int("count", len(users)))
	return users, err
}

func (s *UserManagementService) ChangePassword(ctx context.Context, username string, newPassword string) error {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var requestingUsername string
	if requestingUser != nil {
		requestingUsername = requestingUser.Username
	}
	s.logger.Info("Change password request", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))

	if !ok || requestingUser == nil || (!requestingUser.AdminRole && requestingUser.Username != username) {
		s.logger.Warn("Unauthorized change password attempt", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
		return &ports.NotAuthorizedError{Username: requestingUsername, Action: "change password"}
	}

	if username == "" || newPassword == "" {
		s.logger.Warn("Invalid parameters for change password", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
		return &ports.MissingOrInvalidParameterError{ParameterName: "username or new password"}
	}

	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		s.logger.Error("Failed to get user for password change", slog.String("requesting_user", requestingUsername), slog.String("target_username", username), slog.String("error", err.Error()))
		return err
	}
	user.Password = newPassword

	err = s.repo.UpdateUser(ctx, username, user)
	if err != nil {
		s.logger.Error("Failed to update password", slog.String("requesting_user", requestingUsername), slog.String("target_username", username), slog.String("error", err.Error()))
		return err
	}
	s.logger.Info("Password changed successfully", slog.String("requesting_user", requestingUsername), slog.String("target_username", username))
	return err
}
