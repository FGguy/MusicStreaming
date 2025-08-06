package data

import (
	"context"
	sqlc "music-streaming/sql/sqlc"
	"music-streaming/types"
)

type SQLBrowsing interface {
	GetArtist(ctx context.Context, id int32) (*types.Artist, error)
	GetAlbum(ctx context.Context, id int32) (*types.Album, error)
	GetSong(ctx context.Context, id int32) (*types.Song, error)
}

func (d *DataLayerPg) GetArtist(ctx context.Context, id int32) (*types.Artist, error) {
	conn, err := d.Pg_pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	q := sqlc.New(conn)

	sqlArtist, err := q.GetArtist(ctx, id)
	if err != nil {
		return nil, err
	}

	artist := &types.Artist{
		Id:         int(sqlArtist.ArtistID),
		Name:       sqlArtist.Name,
		CoverArt:   sqlArtist.CoverArt.String,
		AlbumCount: int(sqlArtist.AlbumCount.Int32),
	}
	return artist, nil
}

func (d *DataLayerPg) GetAlbum(ctx context.Context, id int32) (*types.Album, error) {
	conn, err := d.Pg_pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	q := sqlc.New(conn)

	sqlAlbum, err := q.GetAlbum(ctx, id)
	if err != nil {
		return nil, err
	}

	album := &types.Album{
		Id:        sqlAlbum.AlbumID,
		ArtistId:  sqlAlbum.ArtistID.Int32,
		Name:      sqlAlbum.Name,
		CoverArt:  sqlAlbum.CoverArt.String,
		SongCount: int(sqlAlbum.SongCount.Int32),
		Created:   sqlAlbum.Created.Time.String(),
		Duration:  int(sqlAlbum.Duration.Int32),
		Artist:    sqlAlbum.Artist.String,
	}
	return album, nil
}

func (d *DataLayerPg) GetSong(ctx context.Context, id int32) (*types.Song, error) {
	conn, err := d.Pg_pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	q := sqlc.New(conn)

	sqlSong, err := q.GetSong(ctx, id)
	if err != nil {
		return nil, err
	}

	song := &types.Song{
		Id:          sqlSong.SongID,
		AlbumId:     int(sqlSong.AlbumID.Int32),
		Title:       sqlSong.Title,
		Album:       sqlSong.Album.String,
		Artist:      sqlSong.Artist.String,
		IsDir:       sqlSong.IsDir.Bool,
		CoverArt:    sqlSong.CoverArt.String,
		Created:     sqlSong.Created.Time.String(),
		Duration:    int(sqlSong.Duration.Int32),
		BitRate:     int(sqlSong.BitRate.Int32),
		Size:        int(sqlSong.Size.Int32),
		Suffix:      sqlSong.Suffix.String,
		ContentType: sqlSong.ContentType.String,
		IsVideo:     sqlSong.IsVideo.Valid,
		Path:        sqlSong.Path,
	}
	return song, nil
}
