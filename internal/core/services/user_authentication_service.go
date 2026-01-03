package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type UserAuthenticationService struct {
	userRepo ports.UserManagementRepository
	logger   *slog.Logger
}

func NewUserAuthenticationService(userRepo ports.UserManagementRepository, logger *slog.Logger) *UserAuthenticationService {
	return &UserAuthenticationService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *UserAuthenticationService) AuthenticateUser(ctx context.Context, username, password, salt string) (domain.User, error) {
	s.logger.Info("Authentication attempt", slog.String("username", username))
	user, err := s.userRepo.GetUser(ctx, username)
	if err != nil {
		s.logger.Warn("Authentication failed - user not found", slog.String("username", username))
		return domain.User{}, &ports.FailedAuthenticationError{Username: username}
	}

	if !validatePassword(user.Password, salt, password) {
		s.logger.Warn("Authentication failed - invalid password", slog.String("username", username))
		return domain.User{}, &ports.FailedAuthenticationError{Username: username}
	}

	s.logger.Info("Authentication successful", slog.String("username", username))
	return user, nil
}

func validatePassword(hashedPass string, salt string, pass string) bool {
	h := md5.Sum([]byte(pass + salt))
	return hashedPass == hex.EncodeToString(h[:])
}
