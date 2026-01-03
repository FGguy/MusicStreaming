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

func TestMediaRetrievalService_DownloadSong(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		user          *domain.User
		setupMock     func(*mocks.MockMediaBrowsingRepository)
		expectedSong  domain.Song
		expectedError error
	}{
		{
			name: "successful download with admin role",
			id:   1,
			user: &domain.User{
				Username:     "admin",
				AdminRole:    true,
				DownloadRole: false,
			},
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetSongByID(mock.Anything, 1).Return(domain.Song{
					Id:      1,
					Title:   "Test Song",
					Path:    "/music/test.mp3",
					BitRate: 320,
				}, nil)
			},
			expectedSong: domain.Song{
				Id:      1,
				Title:   "Test Song",
				Path:    "/music/test.mp3",
				BitRate: 320,
			},
			expectedError: nil,
		},
		{
			name: "successful download with download role",
			id:   1,
			user: &domain.User{
				Username:     "user",
				AdminRole:    false,
				DownloadRole: true,
			},
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetSongByID(mock.Anything, 1).Return(domain.Song{
					Id:      1,
					Title:   "Test Song",
					Path:    "/music/test.mp3",
					BitRate: 320,
				}, nil)
			},
			expectedSong: domain.Song{
				Id:      1,
				Title:   "Test Song",
				Path:    "/music/test.mp3",
				BitRate: 320,
			},
			expectedError: nil,
		},
		{
			name:          "unauthorized - no roles",
			id:            1,
			user:          &domain.User{Username: "user", AdminRole: false, DownloadRole: false},
			setupMock:     func(m *mocks.MockMediaBrowsingRepository) {},
			expectedSong:  domain.Song{},
			expectedError: &ports.NotAuthorizedError{Username: "user", Action: "download song"},
		},
		{
			// NOTE: This test verifies that the service handles nil users gracefully
			// by safely extracting the username before checking authorization.
			name:          "unauthorized - nil user",
			id:            1,
			user:          nil,
			setupMock:     func(m *mocks.MockMediaBrowsingRepository) {},
			expectedSong:  domain.Song{},
			expectedError: &ports.NotAuthorizedError{Username: "", Action: "download song"},
		},
		{
			name: "song not found",
			id:   999,
			user: &domain.User{
				Username:     "admin",
				AdminRole:    true,
				DownloadRole: false,
			},
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetSongByID(mock.Anything, 999).Return(domain.Song{}, &ports.NotFoundError{Message: "song not found"})
			},
			expectedSong:  domain.Song{},
			expectedError: &ports.NotFoundError{Message: "song not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockMediaBrowsingRepository(t)
			tt.setupMock(repo)
			service := NewMediaRetrievalService(repo, slog.Default())
			ctx := context.Background()
			if tt.user != nil {
				ctx = context.WithValue(ctx, ports.KeyRequestingUserID, tt.user)
			}

			result, err := service.DownloadSong(ctx, tt.id)

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
				if result.Id != tt.expectedSong.Id {
					t.Errorf("expected song ID %d, got %d", tt.expectedSong.Id, result.Id)
				}
				if result.Title != tt.expectedSong.Title {
					t.Errorf("expected song title %s, got %s", tt.expectedSong.Title, result.Title)
				}
			}
		})
	}
}

func TestMediaRetrievalService_StreamSong(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		user          *domain.User
		setupMock     func(*mocks.MockMediaBrowsingRepository)
		expectedSong  domain.Song
		expectedError error
	}{
		{
			name: "successful stream with admin role",
			id:   1,
			user: &domain.User{
				Username:   "admin",
				AdminRole:  true,
				StreamRole: false,
			},
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetSongByID(mock.Anything, 1).Return(domain.Song{
					Id:      1,
					Title:   "Test Song",
					Path:    "/music/test.mp3",
					BitRate: 320,
				}, nil)
			},
			expectedSong: domain.Song{
				Id:      1,
				Title:   "Test Song",
				Path:    "/music/test.mp3",
				BitRate: 320,
			},
			expectedError: nil,
		},
		{
			name: "successful stream with stream role",
			id:   1,
			user: &domain.User{
				Username:   "user",
				AdminRole:  false,
				StreamRole: true,
			},
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetSongByID(mock.Anything, 1).Return(domain.Song{
					Id:      1,
					Title:   "Test Song",
					Path:    "/music/test.mp3",
					BitRate: 320,
				}, nil)
			},
			expectedSong: domain.Song{
				Id:      1,
				Title:   "Test Song",
				Path:    "/music/test.mp3",
				BitRate: 320,
			},
			expectedError: nil,
		},
		{
			name:          "unauthorized - no roles",
			id:            1,
			user:          &domain.User{Username: "user", AdminRole: false, StreamRole: false},
			setupMock:     func(m *mocks.MockMediaBrowsingRepository) {},
			expectedSong:  domain.Song{},
			expectedError: &ports.NotAuthorizedError{Username: "user", Action: "download song"},
		},
		{
			// NOTE: This test verifies that the service handles nil users gracefully
			// by safely extracting the username before checking authorization.
			name:          "unauthorized - nil user",
			id:            1,
			user:          nil,
			setupMock:     func(m *mocks.MockMediaBrowsingRepository) {},
			expectedSong:  domain.Song{},
			expectedError: &ports.NotAuthorizedError{Username: "", Action: "download song"},
		},
		{
			name: "song not found",
			id:   999,
			user: &domain.User{
				Username:   "admin",
				AdminRole:  true,
				StreamRole: false,
			},
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetSongByID(mock.Anything, 999).Return(domain.Song{}, &ports.NotFoundError{Message: "song not found"})
			},
			expectedSong:  domain.Song{},
			expectedError: &ports.NotFoundError{Message: "song not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockMediaBrowsingRepository(t)
			tt.setupMock(repo)
			service := NewMediaRetrievalService(repo, slog.Default())
			ctx := context.Background()
			if tt.user != nil {
				ctx = context.WithValue(ctx, ports.KeyRequestingUserID, tt.user)
			}

			result, err := service.StreamSong(ctx, tt.id)

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
				if result.Id != tt.expectedSong.Id {
					t.Errorf("expected song ID %d, got %d", tt.expectedSong.Id, result.Id)
				}
				if result.Title != tt.expectedSong.Title {
					t.Errorf("expected song title %s, got %s", tt.expectedSong.Title, result.Title)
				}
			}
		})
	}
}

func TestMediaRetrievalService_GetCover(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		setupMock     func(*mocks.MockMediaBrowsingRepository)
		expectedCover domain.Cover
		expectedError error
	}{
		{
			name: "successful retrieval",
			id:   "1",
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetCoverByID(mock.Anything, "1").Return(domain.Cover{
					Id:   "1",
					Path: "/covers/cover1.jpg",
				}, nil)
			},
			expectedCover: domain.Cover{
				Id:   "1",
				Path: "/covers/cover1.jpg",
			},
			expectedError: nil,
		},
		{
			name: "not found error",
			id:   "999",
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetCoverByID(mock.Anything, "999").Return(domain.Cover{}, &ports.NotFoundError{Message: "cover not found"})
			},
			expectedCover: domain.Cover{},
			expectedError: &ports.NotFoundError{Message: "cover not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockMediaBrowsingRepository(t)
			tt.setupMock(repo)
			service := NewMediaRetrievalService(repo, slog.Default())
			ctx := context.Background()

			result, err := service.GetCover(ctx, tt.id)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result.Id != tt.expectedCover.Id {
					t.Errorf("expected cover ID %s, got %s", tt.expectedCover.Id, result.Id)
				}
				if result.Path != tt.expectedCover.Path {
					t.Errorf("expected cover path %s, got %s", tt.expectedCover.Path, result.Path)
				}
			}
		})
	}
}
