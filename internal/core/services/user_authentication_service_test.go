package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"music-streaming/internal/core/services/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestUserAuthenticationService_AuthenticateUser(t *testing.T) {
	// Helper to create a hashed password
	createHashedPassword := func(password, salt string) string {
		h := md5.Sum([]byte(password + salt))
		return hex.EncodeToString(h[:])
	}

	tests := []struct {
		name          string
		username      string
		password      string
		salt          string
		setupMock     func(*mocks.MockUserManagementRepository, string, string)
		expectedUser  domain.User
		expectedError error
	}{
		{
			name:     "successful authentication",
			username: "testuser",
			password: "password123",
			salt:     "somesalt",
			setupMock: func(m *mocks.MockUserManagementRepository, username, hashedPass string) {
				m.EXPECT().GetUser(mock.Anything, username).Return(domain.User{
					Username: username,
					Password: hashedPass,
					Email:    "test@example.com",
				}, nil)
			},
			expectedUser: domain.User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedError: nil,
		},
		{
			name:     "user not found",
			username: "nonexistent",
			password: "password123",
			salt:     "somesalt",
			setupMock: func(m *mocks.MockUserManagementRepository, username, hashedPass string) {
				m.EXPECT().GetUser(mock.Anything, username).Return(domain.User{}, &ports.NotFoundError{Message: "user not found"})
			},
			expectedUser:  domain.User{},
			expectedError: &ports.FailedAuthenticationError{Username: "nonexistent"},
		},
		{
			name:     "wrong password",
			username: "testuser",
			password: "wrongpassword",
			salt:     "somesalt",
			setupMock: func(m *mocks.MockUserManagementRepository, username, hashedPass string) {
				// Return a user with the correct password hash (different from wrongpassword hash)
				correctPasswordHash := createHashedPassword("correctpassword", "somesalt")
				m.EXPECT().GetUser(mock.Anything, username).Return(domain.User{
					Username: username,
					Password: correctPasswordHash,
					Email:    "test@example.com",
				}, nil)
			},
			expectedUser:  domain.User{},
			expectedError: &ports.FailedAuthenticationError{Username: "testuser"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserManagementRepository(t)
			hashedPass := createHashedPassword(tt.password, tt.salt)
			tt.setupMock(repo, tt.username, hashedPass)
			service := NewUserAuthenticationService(repo, slog.Default())
			ctx := context.Background()

			result, err := service.AuthenticateUser(ctx, tt.username, tt.password, tt.salt)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					if authErr, ok := err.(*ports.FailedAuthenticationError); ok {
						expectedErr := tt.expectedError.(*ports.FailedAuthenticationError)
						if authErr.Username != expectedErr.Username {
							t.Errorf("expected error %v, got %v", expectedErr, authErr)
						}
					} else {
						t.Errorf("expected FailedAuthenticationError, got %T", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result.Username != tt.expectedUser.Username {
					t.Errorf("expected username %s, got %s", tt.expectedUser.Username, result.Username)
				}
				if result.Email != tt.expectedUser.Email {
					t.Errorf("expected email %s, got %s", tt.expectedUser.Email, result.Email)
				}
			}
		})
	}
}
