package repositories

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"sync"
)

type MediaBrowsingRepositoryInMemory struct {
	artists map[int]domain.Artist
	albums  map[int]domain.Album
	songs   map[int]domain.Song
	cover   map[int]domain.Cover
	mu      sync.RWMutex
}

func NewMediaBrowsingRepositoryInMemory() *MediaBrowsingRepositoryInMemory {
	return &MediaBrowsingRepositoryInMemory{
		artists: make(map[int]domain.Artist),
		albums:  make(map[int]domain.Album),
		songs:   make(map[int]domain.Song),
		cover:   make(map[int]domain.Cover),
	}
}

func (r *MediaBrowsingRepositoryInMemory) GetArtistByID(ctx context.Context, id int) (domain.Artist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artist, exists := r.artists[id]
	if !exists {
		return domain.Artist{}, &ports.NotFoundError{Message: "artist not found"}
	}
	return artist, nil
}

func (r *MediaBrowsingRepositoryInMemory) GetAlbumByID(ctx context.Context, id int) (domain.Album, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	album, exists := r.albums[id]
	if !exists {
		return domain.Album{}, &ports.NotFoundError{Message: "album not found"}
	}
	return album, nil
}

func (r *MediaBrowsingRepositoryInMemory) GetSongByID(ctx context.Context, id int) (domain.Song, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	song, exists := r.songs[id]
	if !exists {
		return domain.Song{}, &ports.NotFoundError{Message: "song not found"}
	}
	return song, nil
}

func (r *MediaBrowsingRepositoryInMemory) GetCoverByID(ctx context.Context, id int) (domain.Cover, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cover, exists := r.cover[id]
	if !exists {
		return domain.Cover{}, &ports.NotFoundError{Message: "cover not found"}
	}
	return cover, nil
}

func (r *MediaBrowsingRepositoryInMemory) CreateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.artists[artist.Id] = artist
	return artist, nil
}

func (r *MediaBrowsingRepositoryInMemory) CreateAlbum(ctx context.Context, album domain.Album) (domain.Album, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.albums[album.Id] = album
	return album, nil
}

func (r *MediaBrowsingRepositoryInMemory) CreateSong(ctx context.Context, song domain.Song) (domain.Song, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.songs[song.Id] = song
	return song, nil
}

func (r *MediaBrowsingRepositoryInMemory) CreateCover(ctx context.Context, cover domain.Cover) (domain.Cover, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cover[cover.Id] = cover
	return cover, nil
}
