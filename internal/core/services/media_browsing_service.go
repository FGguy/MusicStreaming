package services

import (
	"context"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type MediaBrowsingService struct {
	mediaBrowsingRepo ports.MediaBrowsingRepository
	logger            *slog.Logger
}

func NewMediaBrowsingService(repo ports.MediaBrowsingRepository, logger *slog.Logger) *MediaBrowsingService {
	return &MediaBrowsingService{
		mediaBrowsingRepo: repo,
		logger:            logger,
	}
}

func (s *MediaBrowsingService) GetArtist(ctx context.Context, id int) (domain.Artist, error) {
	s.logger.Info("Getting artist", slog.Int("id", id))
	artist, err := s.mediaBrowsingRepo.GetArtistByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get artist", slog.Int("id", id), slog.String("error", err.Error()))
		return artist, err
	}
	s.logger.Info("Successfully retrieved artist", slog.Int("id", id), slog.String("name", artist.Name))
	return artist, err
}

func (s *MediaBrowsingService) GetAlbum(ctx context.Context, id int) (domain.Album, error) {
	s.logger.Info("Getting album", slog.Int("id", id))
	album, err := s.mediaBrowsingRepo.GetAlbumByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get album", slog.Int("id", id), slog.String("error", err.Error()))
		return album, err
	}
	s.logger.Info("Successfully retrieved album", slog.Int("id", id), slog.String("name", album.Name))
	return album, err
}

func (s *MediaBrowsingService) GetSong(ctx context.Context, id int) (domain.Song, error) {
	s.logger.Info("Getting song", slog.Int("id", id))
	song, err := s.mediaBrowsingRepo.GetSongByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get song", slog.Int("id", id), slog.String("error", err.Error()))
		return song, err
	}
	s.logger.Info("Successfully retrieved song", slog.Int("id", id), slog.String("title", song.Title))
	return song, err
}

func (s *MediaBrowsingService) GetCover(ctx context.Context, id string) (domain.Cover, error) {
	s.logger.Info("Getting cover", slog.String("id", id))
	cover, err := s.mediaBrowsingRepo.GetCoverByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get cover", slog.String("id", id), slog.String("error", err.Error()))
		return cover, err
	}
	s.logger.Info("Successfully retrieved cover", slog.String("id", id))
	return cover, err
}
