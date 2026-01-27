package ports

import (
	"context"
	"music-streaming/internal/core/domain"
)

/*
TODO: Future Subsonic API endpoints to implement:
- GetMusicFolders
- GetIndexes
- GetMusicDirectory
- GetGenres
- GetArtists
- GetRandomSong
- GetStarred
- GetFiles

LastFM Integration:
- GetAlbumInfo
- GetArtistInfo
*/

// MediaBrowsingPort defines the interface for browsing media catalog operations.
// It provides methods to retrieve artists, albums, songs, and cover art.
type MediaBrowsingPort interface {
	// GetArtist retrieves an artist by their unique ID.
	GetArtist(ctx context.Context, id int) (domain.Artist, error)

	// GetAlbum retrieves an album by its unique ID.
	GetAlbum(ctx context.Context, id int) (domain.Album, error)

	// GetSong retrieves a song by its unique ID.
	GetSong(ctx context.Context, id int) (domain.Song, error)

	// GetCover retrieves cover art metadata by its unique ID.
	GetCover(ctx context.Context, id string) (domain.Cover, error)
}

// MediaBrowsingRepository defines the interface for media catalog data persistence.
// Implementations provide data access operations for artists, albums, songs, and covers.
type MediaBrowsingRepository interface {
	// GetArtistByID retrieves an artist from the data store by ID.
	GetArtistByID(ctx context.Context, id int) (domain.Artist, error)

	// GetAlbumByID retrieves an album from the data store by ID.
	GetAlbumByID(ctx context.Context, id int) (domain.Album, error)

	// GetSongByID retrieves a song from the data store by ID.
	GetSongByID(ctx context.Context, id int) (domain.Song, error)

	// GetCoverByID retrieves cover art metadata from the data store by ID.
	GetCoverByID(ctx context.Context, id string) (domain.Cover, error)

	// CreateArtist persists a new artist to the data store.
	// Used by media scanning service during library indexing.
	CreateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error)

	// CreateAlbum persists a new album to the data store.
	// Used by media scanning service during library indexing.
	CreateAlbum(ctx context.Context, album domain.Album) (domain.Album, error)

	// CreateSong persists a new song to the data store.
	// Used by media scanning service during library indexing.
	CreateSong(ctx context.Context, song domain.Song) (domain.Song, error)

	// CreateCover persists new cover art metadata to the data store.
	// Used by media scanning service during library indexing.
	CreateCover(ctx context.Context, cover domain.Cover) (domain.Cover, error)
}
