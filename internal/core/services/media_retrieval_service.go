package services

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type MediaRetrievalService struct {
	MediaBrowsingRepository ports.MediaBrowsingRepository
}

func NewMediaRetrievalService(mediaBrowsingRepository ports.MediaBrowsingRepository) *MediaRetrievalService {
	return &MediaRetrievalService{
		MediaBrowsingRepository: mediaBrowsingRepository,
	}
}

func (s *MediaRetrievalService) DownloadSong(ctx context.Context, id int) (domain.Song, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || (!requestingUser.AdminRole && !requestingUser.DownloadRole) {
		return domain.Song{}, &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "download song"}
	}

	song, err := s.MediaBrowsingRepository.GetSongByID(ctx, id)
	if err != nil {
		return domain.Song{}, err
	}
	return song, nil
}

func (s *MediaRetrievalService) StreamSong(ctx context.Context, id int) (domain.Song, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || (!requestingUser.AdminRole && !requestingUser.StreamRole) {
		return domain.Song{}, &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "download song"}
	}

	song, err := s.MediaBrowsingRepository.GetSongByID(ctx, id)
	if err != nil {
		return domain.Song{}, err
	}
	return song, nil
}

func (s *MediaRetrievalService) GetCover(ctx context.Context, id string) (domain.Cover, error) {
	cover, err := s.MediaBrowsingRepository.GetCoverByID(ctx, id)
	if err != nil {
		return domain.Cover{}, err
	}
	return cover, nil
}
