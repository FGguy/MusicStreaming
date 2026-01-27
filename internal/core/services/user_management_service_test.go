package services

import (
	"context"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"music-streaming/internal/core/services/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestUserManagementService_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		user           domain.User
		requestingUser *domain.User
		setupMock      func(*mocks.MockUserManagementRepository)
		expectedError  error
	}{
		{
			name: "successful creation with admin role",
			user: domain.User{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().CreateUser(mock.Anything, domain.User{
					Username: "newuser",
					Email:    "newuser@example.com",
					Password: "password123",
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "unauthorized - no admin role",
			user: domain.User{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			requestingUser: &domain.User{
				Username:  "user",
				AdminRole: false,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.NotAuthorizedError{Username: "user", Action: "create user"},
		},
		{
			name: "missing username",
			user: domain.User{
				Username: "",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.MissingOrInvalidParameterError{ParameterName: "username, email or password"},
		},
		{
			name: "missing email",
			user: domain.User{
				Username: "newuser",
				Email:    "",
				Password: "password123",
			},
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.MissingOrInvalidParameterError{ParameterName: "username, email or password"},
		},
		{
			name: "missing password",
			user: domain.User{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "",
			},
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.MissingOrInvalidParameterError{ParameterName: "username, email or password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserManagementRepository(t)
			tt.setupMock(repo)
			service := NewUserManagementService(repo, slog.Default())
			ctx := context.Background()
			if tt.requestingUser != nil {
				ctx = context.WithValue(ctx, ports.KeyRequestingUserID, tt.requestingUser)
			}

			err := service.CreateUser(ctx, tt.user)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					switch expectedErr := tt.expectedError.(type) {
					case *ports.NotAuthorizedError:
						if authErr, ok := err.(*ports.NotAuthorizedError); ok {
							if authErr.Username != expectedErr.Username || authErr.Action != expectedErr.Action {
								t.Errorf("expected error %v, got %v", expectedErr, authErr)
							}
						} else {
							t.Errorf("expected NotAuthorizedError, got %T", err)
						}
					case *ports.MissingOrInvalidParameterError:
						if paramErr, ok := err.(*ports.MissingOrInvalidParameterError); ok {
							if paramErr.ParameterName != expectedErr.ParameterName {
								t.Errorf("expected error %v, got %v", expectedErr, paramErr)
							}
						} else {
							t.Errorf("expected MissingOrInvalidParameterError, got %T", err)
						}
					default:
						if err.Error() != tt.expectedError.Error() {
							t.Errorf("expected error %v, got %v", tt.expectedError, err)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUserManagementService_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		user           domain.User
		requestingUser *domain.User
		setupMock      func(*mocks.MockUserManagementRepository)
		expectedError  error
	}{
		{
			name:     "successful update with admin role",
			username: "testuser",
			user: domain.User{
				Username: "testuser",
				Email:    "updated@example.com",
				Password: "newpassword",
			},
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().UpdateUser(mock.Anything, "testuser", domain.User{
					Username: "testuser",
					Email:    "updated@example.com",
					Password: "newpassword",
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "successful update with settings role on own account",
			username: "testuser",
			user: domain.User{
				Username: "testuser",
				Email:    "updated@example.com",
				Password: "newpassword",
			},
			requestingUser: &domain.User{
				Username:     "testuser",
				AdminRole:    false,
				SettingsRole: true,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().UpdateUser(mock.Anything, "testuser", domain.User{
					Username: "testuser",
					Email:    "updated@example.com",
					Password: "newpassword",
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "unauthorized - no admin role and not own account",
			username: "otheruser",
			user: domain.User{
				Username: "otheruser",
				Email:    "updated@example.com",
				Password: "newpassword",
			},
			requestingUser: &domain.User{
				Username:     "testuser",
				AdminRole:    false,
				SettingsRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.NotAuthorizedError{Username: "testuser", Action: "update user"},
		},
		{
			name:     "missing username",
			username: "",
			user: domain.User{
				Username: "testuser",
				Email:    "updated@example.com",
				Password: "newpassword",
			},
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.MissingOrInvalidParameterError{ParameterName: "username, email or password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserManagementRepository(t)
			tt.setupMock(repo)
			service := NewUserManagementService(repo, slog.Default())
			ctx := context.Background()
			if tt.requestingUser != nil {
				ctx = context.WithValue(ctx, ports.KeyRequestingUserID, tt.requestingUser)
			}

			err := service.UpdateUser(ctx, tt.username, tt.user)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					switch expectedErr := tt.expectedError.(type) {
					case *ports.NotAuthorizedError:
						if authErr, ok := err.(*ports.NotAuthorizedError); ok {
							if authErr.Username != expectedErr.Username || authErr.Action != expectedErr.Action {
								t.Errorf("expected error %v, got %v", expectedErr, authErr)
							}
						} else {
							t.Errorf("expected NotAuthorizedError, got %T", err)
						}
					case *ports.MissingOrInvalidParameterError:
						if paramErr, ok := err.(*ports.MissingOrInvalidParameterError); ok {
							if paramErr.ParameterName != expectedErr.ParameterName {
								t.Errorf("expected error %v, got %v", expectedErr, paramErr)
							}
						} else {
							t.Errorf("expected MissingOrInvalidParameterError, got %T", err)
						}
					default:
						if err.Error() != tt.expectedError.Error() {
							t.Errorf("expected error %v, got %v", tt.expectedError, err)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUserManagementService_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		requestingUser *domain.User
		setupMock      func(*mocks.MockUserManagementRepository)
		expectedError  error
	}{
		{
			name:     "successful deletion with admin role",
			username: "testuser",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().DeleteUser(mock.Anything, "testuser").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "unauthorized - no admin role",
			username: "testuser",
			requestingUser: &domain.User{
				Username:  "user",
				AdminRole: false,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.NotAuthorizedError{Username: "user", Action: "delete user"},
		},
		{
			name:     "missing username",
			username: "",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.MissingOrInvalidParameterError{ParameterName: "username"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserManagementRepository(t)
			tt.setupMock(repo)
			service := NewUserManagementService(repo, slog.Default())
			ctx := context.Background()
			if tt.requestingUser != nil {
				ctx = context.WithValue(ctx, ports.KeyRequestingUserID, tt.requestingUser)
			}

			err := service.DeleteUser(ctx, tt.username)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					switch expectedErr := tt.expectedError.(type) {
					case *ports.NotAuthorizedError:
						if authErr, ok := err.(*ports.NotAuthorizedError); ok {
							if authErr.Username != expectedErr.Username || authErr.Action != expectedErr.Action {
								t.Errorf("expected error %v, got %v", expectedErr, authErr)
							}
						} else {
							t.Errorf("expected NotAuthorizedError, got %T", err)
						}
					case *ports.MissingOrInvalidParameterError:
						if paramErr, ok := err.(*ports.MissingOrInvalidParameterError); ok {
							if paramErr.ParameterName != expectedErr.ParameterName {
								t.Errorf("expected error %v, got %v", expectedErr, paramErr)
							}
						} else {
							t.Errorf("expected MissingOrInvalidParameterError, got %T", err)
						}
					default:
						if err.Error() != tt.expectedError.Error() {
							t.Errorf("expected error %v, got %v", tt.expectedError, err)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUserManagementService_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		requestingUser *domain.User
		setupMock      func(*mocks.MockUserManagementRepository)
		expectedUser   domain.User
		expectedError  error
	}{
		{
			name:     "successful get with admin role",
			username: "testuser",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().GetUser(mock.Anything, "testuser").Return(domain.User{
					Username: "testuser",
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
			name:     "successful get own user",
			username: "testuser",
			requestingUser: &domain.User{
				Username:  "testuser",
				AdminRole: false,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().GetUser(mock.Anything, "testuser").Return(domain.User{
					Username: "testuser",
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
			name:     "unauthorized - get other user without admin role",
			username: "otheruser",
			requestingUser: &domain.User{
				Username:  "testuser",
				AdminRole: false,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedUser:  domain.User{},
			expectedError: &ports.NotAuthorizedError{Username: "testuser", Action: "get user"},
		},
		{
			name:     "missing username",
			username: "",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedUser:  domain.User{},
			expectedError: &ports.MissingOrInvalidParameterError{ParameterName: "username"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserManagementRepository(t)
			tt.setupMock(repo)
			service := NewUserManagementService(repo, slog.Default())
			ctx := context.Background()
			if tt.requestingUser != nil {
				ctx = context.WithValue(ctx, ports.KeyRequestingUserID, tt.requestingUser)
			}

			result, err := service.GetUser(ctx, tt.username)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					switch expectedErr := tt.expectedError.(type) {
					case *ports.NotAuthorizedError:
						if authErr, ok := err.(*ports.NotAuthorizedError); ok {
							if authErr.Username != expectedErr.Username || authErr.Action != expectedErr.Action {
								t.Errorf("expected error %v, got %v", expectedErr, authErr)
							}
						} else {
							t.Errorf("expected NotAuthorizedError, got %T", err)
						}
					case *ports.MissingOrInvalidParameterError:
						if paramErr, ok := err.(*ports.MissingOrInvalidParameterError); ok {
							if paramErr.ParameterName != expectedErr.ParameterName {
								t.Errorf("expected error %v, got %v", expectedErr, paramErr)
							}
						} else {
							t.Errorf("expected MissingOrInvalidParameterError, got %T", err)
						}
					default:
						if err.Error() != tt.expectedError.Error() {
							t.Errorf("expected error %v, got %v", tt.expectedError, err)
						}
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

func TestUserManagementService_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		requestingUser *domain.User
		setupMock      func(*mocks.MockUserManagementRepository)
		expectedUsers  []domain.User
		expectedError  error
	}{
		{
			name: "successful get with admin role",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().GetUsers(mock.Anything).Return([]domain.User{
					{Username: "user1", Email: "user1@example.com"},
					{Username: "user2", Email: "user2@example.com"},
				}, nil)
			},
			expectedUsers: []domain.User{
				{Username: "user1", Email: "user1@example.com"},
				{Username: "user2", Email: "user2@example.com"},
			},
			expectedError: nil,
		},
		{
			name: "unauthorized - no admin role",
			requestingUser: &domain.User{
				Username:  "user",
				AdminRole: false,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedUsers: []domain.User{},
			expectedError: &ports.NotAuthorizedError{Username: "user", Action: "get users"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserManagementRepository(t)
			tt.setupMock(repo)
			service := NewUserManagementService(repo, slog.Default())
			ctx := context.Background()
			if tt.requestingUser != nil {
				ctx = context.WithValue(ctx, ports.KeyRequestingUserID, tt.requestingUser)
			}

			result, err := service.GetUsers(ctx)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					if authErr, ok := err.(*ports.NotAuthorizedError); ok {
						expectedErr := tt.expectedError.(*ports.NotAuthorizedError)
						if authErr.Username != expectedErr.Username || authErr.Action != expectedErr.Action {
							t.Errorf("expected error %v, got %v", expectedErr, authErr)
						}
					} else {
						t.Errorf("expected NotAuthorizedError, got %T", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != len(tt.expectedUsers) {
					t.Errorf("expected %d users, got %d", len(tt.expectedUsers), len(result))
				}
			}
		})
	}
}

func TestUserManagementService_ChangePassword(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		newPassword    string
		requestingUser *domain.User
		setupMock      func(*mocks.MockUserManagementRepository)
		expectedError  error
	}{
		{
			name:        "successful change with admin role",
			username:    "testuser",
			newPassword: "newpassword123",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().GetUser(mock.Anything, "testuser").Return(domain.User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "oldpassword",
				}, nil)
				m.EXPECT().UpdateUser(mock.Anything, "testuser", domain.User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "newpassword123",
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "successful change own password",
			username:    "testuser",
			newPassword: "newpassword123",
			requestingUser: &domain.User{
				Username:  "testuser",
				AdminRole: false,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().GetUser(mock.Anything, "testuser").Return(domain.User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "oldpassword",
				}, nil)
				m.EXPECT().UpdateUser(mock.Anything, "testuser", domain.User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "newpassword123",
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "unauthorized - change other user password without admin role",
			username:    "otheruser",
			newPassword: "newpassword123",
			requestingUser: &domain.User{
				Username:  "testuser",
				AdminRole: false,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.NotAuthorizedError{Username: "testuser", Action: "change password"},
		},
		{
			name:        "missing username",
			username:    "",
			newPassword: "newpassword123",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.MissingOrInvalidParameterError{ParameterName: "username or new password"},
		},
		{
			name:        "missing new password",
			username:    "testuser",
			newPassword: "",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock:     func(m *mocks.MockUserManagementRepository) {},
			expectedError: &ports.MissingOrInvalidParameterError{ParameterName: "username or new password"},
		},
		{
			name:        "user not found",
			username:    "nonexistent",
			newPassword: "newpassword123",
			requestingUser: &domain.User{
				Username:  "admin",
				AdminRole: true,
			},
			setupMock: func(m *mocks.MockUserManagementRepository) {
				m.EXPECT().GetUser(mock.Anything, "nonexistent").Return(domain.User{}, &ports.NotFoundError{Message: "user not found"})
			},
			expectedError: &ports.NotFoundError{Message: "user not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserManagementRepository(t)
			tt.setupMock(repo)
			service := NewUserManagementService(repo, slog.Default())
			ctx := context.Background()
			if tt.requestingUser != nil {
				ctx = context.WithValue(ctx, ports.KeyRequestingUserID, tt.requestingUser)
			}

			err := service.ChangePassword(ctx, tt.username, tt.newPassword)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					switch expectedErr := tt.expectedError.(type) {
					case *ports.NotAuthorizedError:
						if authErr, ok := err.(*ports.NotAuthorizedError); ok {
							if authErr.Username != expectedErr.Username || authErr.Action != expectedErr.Action {
								t.Errorf("expected error %v, got %v", expectedErr, authErr)
							}
						} else {
							t.Errorf("expected NotAuthorizedError, got %T", err)
						}
					case *ports.MissingOrInvalidParameterError:
						if paramErr, ok := err.(*ports.MissingOrInvalidParameterError); ok {
							if paramErr.ParameterName != expectedErr.ParameterName {
								t.Errorf("expected error %v, got %v", expectedErr, paramErr)
							}
						} else {
							t.Errorf("expected MissingOrInvalidParameterError, got %T", err)
						}
					case *ports.NotFoundError:
						if notFoundErr, ok := err.(*ports.NotFoundError); ok {
							if notFoundErr.Message != expectedErr.Message {
								t.Errorf("expected error %v, got %v", expectedErr, notFoundErr)
							}
						} else {
							t.Errorf("expected NotFoundError, got %T", err)
						}
					default:
						if err.Error() != tt.expectedError.Error() {
							t.Errorf("expected error %v, got %v", tt.expectedError, err)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
