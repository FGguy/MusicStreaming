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

func TestMediaBrowsingService_GetArtist(t *testing.T) {
	tests := []struct {
		name           string
		id             int
		setupMock      func(*mocks.MockMediaBrowsingRepository)
		expectedArtist domain.Artist
		expectedError  error
	}{
		{
			name: "successful retrieval",
			id:   1,
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetArtistByID(mock.Anything, 1).Return(domain.Artist{
					Id:         1,
					Name:       "Test Artist",
					CoverArt:   "cover1",
					AlbumCount: 5,
				}, nil)
			},
			expectedArtist: domain.Artist{
				Id:         1,
				Name:       "Test Artist",
				CoverArt:   "cover1",
				AlbumCount: 5,
			},
			expectedError: nil,
		},
		{
			name: "not found error",
			id:   999,
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetArtistByID(mock.Anything, 999).Return(domain.Artist{}, &ports.NotFoundError{Message: "artist not found"})
			},
			expectedArtist: domain.Artist{},
			expectedError:  &ports.NotFoundError{Message: "artist not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockMediaBrowsingRepository(t)
			tt.setupMock(repo)
			service := NewMediaBrowsingService(repo, slog.Default())
			ctx := context.Background()

			result, err := service.GetArtist(ctx, tt.id)

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
				if result.Id != tt.expectedArtist.Id {
					t.Errorf("expected artist ID %d, got %d", tt.expectedArtist.Id, result.Id)
				}
				if result.Name != tt.expectedArtist.Name {
					t.Errorf("expected artist name %s, got %s", tt.expectedArtist.Name, result.Name)
				}
			}
		})
	}
}

func TestMediaBrowsingService_GetAlbum(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		setupMock     func(*mocks.MockMediaBrowsingRepository)
		expectedAlbum domain.Album
		expectedError error
	}{
		{
			name: "successful retrieval",
			id:   1,
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetAlbumByID(mock.Anything, 1).Return(domain.Album{
					Id:        1,
					ArtistId:  1,
					Name:      "Test Album",
					CoverArt:  "cover1",
					SongCount: 10,
					Created:   "2024-01-01",
					Duration:  3600,
					Artist:    "Test Artist",
				}, nil)
			},
			expectedAlbum: domain.Album{
				Id:        1,
				ArtistId:  1,
				Name:      "Test Album",
				CoverArt:  "cover1",
				SongCount: 10,
				Created:   "2024-01-01",
				Duration:  3600,
				Artist:    "Test Artist",
			},
			expectedError: nil,
		},
		{
			name: "not found error",
			id:   999,
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetAlbumByID(mock.Anything, 999).Return(domain.Album{}, &ports.NotFoundError{Message: "album not found"})
			},
			expectedAlbum: domain.Album{},
			expectedError: &ports.NotFoundError{Message: "album not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockMediaBrowsingRepository(t)
			tt.setupMock(repo)
			service := NewMediaBrowsingService(repo, slog.Default())
			ctx := context.Background()

			result, err := service.GetAlbum(ctx, tt.id)

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
				if result.Id != tt.expectedAlbum.Id {
					t.Errorf("expected album ID %d, got %d", tt.expectedAlbum.Id, result.Id)
				}
				if result.Name != tt.expectedAlbum.Name {
					t.Errorf("expected album name %s, got %s", tt.expectedAlbum.Name, result.Name)
				}
			}
		})
	}
}

func TestMediaBrowsingService_GetSong(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		setupMock     func(*mocks.MockMediaBrowsingRepository)
		expectedSong  domain.Song
		expectedError error
	}{
		{
			name: "successful retrieval",
			id:   1,
			setupMock: func(m *mocks.MockMediaBrowsingRepository) {
				m.EXPECT().GetSongByID(mock.Anything, 1).Return(domain.Song{
					Id:          1,
					AlbumId:     1,
					Title:       "Test Song",
					Album:       "Test Album",
					Artist:      "Test Artist",
					IsDir:       false,
					CoverArt:    "cover1",
					Created:     "2024-01-01",
					Duration:    180,
					BitRate:     320,
					Size:        5000000,
					Suffix:      "mp3",
					ContentType: "audio/mpeg",
					IsVideo:     false,
					Path:        "/music/test.mp3",
				}, nil)
			},
			expectedSong: domain.Song{
				Id:          1,
				AlbumId:     1,
				Title:       "Test Song",
				Album:       "Test Album",
				Artist:      "Test Artist",
				IsDir:       false,
				CoverArt:    "cover1",
				Created:     "2024-01-01",
				Duration:    180,
				BitRate:     320,
				Size:        5000000,
				Suffix:      "mp3",
				ContentType: "audio/mpeg",
				IsVideo:     false,
				Path:        "/music/test.mp3",
			},
			expectedError: nil,
		},
		{
			name: "not found error",
			id:   999,
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
			service := NewMediaBrowsingService(repo, slog.Default())
			ctx := context.Background()

			result, err := service.GetSong(ctx, tt.id)

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

func TestMediaBrowsingService_GetCover(t *testing.T) {
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
			service := NewMediaBrowsingService(repo, slog.Default())
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
