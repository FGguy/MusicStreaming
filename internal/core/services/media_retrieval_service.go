package services

import (
	"context"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type MediaRetrievalService struct {
	MediaBrowsingRepository ports.MediaBrowsingRepository
	logger                  *slog.Logger
}

func NewMediaRetrievalService(mediaBrowsingRepository ports.MediaBrowsingRepository, logger *slog.Logger) *MediaRetrievalService {
	return &MediaRetrievalService{
		MediaBrowsingRepository: mediaBrowsingRepository,
		logger:                  logger,
	}
}

func (s *MediaRetrievalService) DownloadSong(ctx context.Context, id int) (domain.Song, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var username string
	if requestingUser != nil {
		username = requestingUser.Username
	}
	s.logger.Info("Download song request", slog.Int("id", id), slog.String("username", username))

	if !ok || requestingUser == nil || (!requestingUser.AdminRole && !requestingUser.DownloadRole) {
		s.logger.Warn("Unauthorized download song attempt", slog.Int("id", id), slog.String("username", username))
		return domain.Song{}, &ports.NotAuthorizedError{Username: username, Action: "download song"}
	}

	song, err := s.MediaBrowsingRepository.GetSongByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get song for download", slog.Int("id", id), slog.String("username", username), slog.String("error", err.Error()))
		return domain.Song{}, err
	}
	s.logger.Info("Song download successful", slog.Int("id", id), slog.String("title", song.Title), slog.String("username", username))
	return song, nil
}

func (s *MediaRetrievalService) StreamSong(ctx context.Context, id int) (domain.Song, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var username string
	if requestingUser != nil {
		username = requestingUser.Username
	}
	s.logger.Info("Stream song request", slog.Int("id", id), slog.String("username", username))

	if !ok || requestingUser == nil || (!requestingUser.AdminRole && !requestingUser.StreamRole) {
		s.logger.Warn("Unauthorized stream song attempt", slog.Int("id", id), slog.String("username", username))
		return domain.Song{}, &ports.NotAuthorizedError{Username: username, Action: "download song"}
	}

	song, err := s.MediaBrowsingRepository.GetSongByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get song for streaming", slog.Int("id", id), slog.String("username", username), slog.String("error", err.Error()))
		return domain.Song{}, err
	}
	s.logger.Info("Song stream started", slog.Int("id", id), slog.String("title", song.Title), slog.String("username", username))
	return song, nil
}

func (s *MediaRetrievalService) GetCover(ctx context.Context, id int) (domain.Cover, error) {
	s.logger.Info("Getting cover", slog.Int("id", id))
	cover, err := s.MediaBrowsingRepository.GetCoverByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get cover", slog.Int("id", id), slog.String("error", err.Error()))
		return domain.Cover{}, err
	}
	s.logger.Info("Successfully retrieved cover", slog.Int("id", id))
	return cover, err
}
