package ports

import (
	"context"
	"music-streaming/internal/core/domain"
)

/*
TODO:
- GetMusicFolders
- GetIndexes
- GetMusicDirectory
- GetGenres
- GetArtists
- GetRandomSong
- GetStarred
- GetFiles

LastFM Integration
- GetAlbumInfo
- GetArtistInfo
*/

type MediaBrowsingPort interface {
	GetArtist(ctx context.Context, id int) (domain.Artist, error)
	GetAlbum(ctx context.Context, id int) (domain.Album, error)
	GetSong(ctx context.Context, id int) (domain.Song, error)
	GetCover(ctx context.Context, id string) (domain.Cover, error)
}

type MediaBrowsingRepository interface {
	GetArtistByID(ctx context.Context, id int) (domain.Artist, error)
	GetAlbumByID(ctx context.Context, id int) (domain.Album, error)
	GetSongByID(ctx context.Context, id int) (domain.Song, error)
	GetCoverByID(ctx context.Context, id string) (domain.Cover, error)

	// Used by media scanning service, not directly accessed by application services
	CreateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error)
	CreateAlbum(ctx context.Context, album domain.Album) (domain.Album, error)
	CreateSong(ctx context.Context, song domain.Song) (domain.Song, error)
	CreateCover(ctx context.Context, cover domain.Cover) (domain.Cover, error)
}
