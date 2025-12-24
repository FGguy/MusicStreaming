package services

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
)

type MediaBrowsingService struct {
	mediaBrowsingRepo ports.MediaBrowsingRepository
}

func NewMediaBrowsingService(repo ports.MediaBrowsingRepository) *MediaBrowsingService {
	return &MediaBrowsingService{
		mediaBrowsingRepo: repo,
	}
}

func (s *MediaBrowsingService) GetArtist(ctx context.Context, id int) (domain.Artist, error) {
	return s.mediaBrowsingRepo.GetArtistByID(ctx, id)
}

func (s *MediaBrowsingService) GetAlbum(ctx context.Context, id int) (domain.Album, error) {
	return s.mediaBrowsingRepo.GetAlbumByID(ctx, id)
}

func (s *MediaBrowsingService) GetSong(ctx context.Context, id int) (domain.Song, error) {
	return s.mediaBrowsingRepo.GetSongByID(ctx, id)
}

func (s *MediaBrowsingService) GetCover(ctx context.Context, id string) (domain.Cover, error) {
	return s.mediaBrowsingRepo.GetCoverByID(ctx, id)
}
