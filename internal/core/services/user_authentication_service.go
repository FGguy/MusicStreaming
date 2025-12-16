package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type UserAuthenticationService struct {
	userRepo ports.UserManagementRepository
}

func NewUserAuthenticationService(userRepo ports.UserManagementRepository) *UserAuthenticationService {
	return &UserAuthenticationService{
		userRepo: userRepo,
	}
}

func (s *UserAuthenticationService) AuthenticateUser(ctx context.Context, username, password, salt string) (domain.User, error) {
	user, err := s.userRepo.GetUser(ctx, username)
	if err != nil {
		return domain.User{}, &ports.FailedAuthenticationError{Username: username}
	}

	if !validatePassword(user.Password, salt, password) {
		return domain.User{}, &ports.FailedAuthenticationError{Username: username}
	}

	return user, nil
}

func validatePassword(hashedPass string, salt string, pass string) bool {
	h := md5.Sum([]byte(pass + salt))
	return hashedPass == hex.EncodeToString(h[:])
}
