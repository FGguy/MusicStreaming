package repositories

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"sync"
)

/*
* artists, albums and songs have their id created by the repository (auto-increment)
 */

type InMemoryMediaBrowsingRepository struct {
	artists map[int]domain.Artist
	albums  map[int]domain.Album
	songs   map[int]domain.Song
	cover   map[string]domain.Cover
	mu      sync.RWMutex
}

func NewInMemoryMediaBrowsingRepository() *InMemoryMediaBrowsingRepository {
	return &InMemoryMediaBrowsingRepository{
		artists: make(map[int]domain.Artist),
		albums:  make(map[int]domain.Album),
		songs:   make(map[int]domain.Song),
		cover:   make(map[string]domain.Cover),
	}
}

func (r *InMemoryMediaBrowsingRepository) GetArtistByID(ctx context.Context, id int) (domain.Artist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artist, exists := r.artists[id]
	if !exists {
		return domain.Artist{}, &ports.NotFoundError{Message: "artist not found"}
	}
	return artist, nil
}

func (r *InMemoryMediaBrowsingRepository) GetAlbumByID(ctx context.Context, id int) (domain.Album, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	album, exists := r.albums[id]
	if !exists {
		return domain.Album{}, &ports.NotFoundError{Message: "album not found"}
	}
	return album, nil
}

func (r *InMemoryMediaBrowsingRepository) GetSongByID(ctx context.Context, id int) (domain.Song, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	song, exists := r.songs[id]
	if !exists {
		return domain.Song{}, &ports.NotFoundError{Message: "song not found"}
	}
	return song, nil
}

func (r *InMemoryMediaBrowsingRepository) GetCoverByID(ctx context.Context, id string) (domain.Cover, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cover, exists := r.cover[id]
	if !exists {
		return domain.Cover{}, &ports.NotFoundError{Message: "cover not found"}
	}
	return cover, nil
}

func (r *InMemoryMediaBrowsingRepository) CreateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// assign a new ID
	artist.Id = len(r.artists) + 1
	r.artists[artist.Id] = artist
	return artist, nil
}

func (r *InMemoryMediaBrowsingRepository) CreateAlbum(ctx context.Context, album domain.Album) (domain.Album, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// assign a new ID
	album.Id = len(r.albums) + 1
	r.albums[album.Id] = album
	return album, nil
}

func (r *InMemoryMediaBrowsingRepository) CreateSong(ctx context.Context, song domain.Song) (domain.Song, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// assign a new ID
	song.Id = len(r.songs) + 1
	r.songs[song.Id] = song
	return song, nil
}

func (r *InMemoryMediaBrowsingRepository) CreateCover(ctx context.Context, cover domain.Cover) (domain.Cover, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cover[cover.Id] = cover
	return cover, nil
}
