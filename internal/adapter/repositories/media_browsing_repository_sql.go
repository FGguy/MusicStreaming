package repositories

import (
	"context"
	"fmt"
	sqlc "music-streaming/internal/adapter/sql/sqlc"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SQLMediaBrowsingRepository struct {
	queries *sqlc.Queries
	db      *pgx.Conn
}

func NewSQLMediaBrowsingRepository(db *pgx.Conn) *SQLMediaBrowsingRepository {
	return &SQLMediaBrowsingRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

func (r *SQLMediaBrowsingRepository) GetArtistByID(ctx context.Context, id int) (domain.Artist, error) {
	sqlArtist, err := r.queries.GetArtist(ctx, int32(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Artist{}, &ports.NotFoundError{Message: "artist not found"}
		}
		return domain.Artist{}, fmt.Errorf("failed to get artist: %w", err)
	}

	return toDomainArtist(sqlArtist), nil
}

func (r *SQLMediaBrowsingRepository) GetAlbumByID(ctx context.Context, id int) (domain.Album, error) {
	sqlAlbum, err := r.queries.GetAlbum(ctx, int32(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Album{}, &ports.NotFoundError{Message: "album not found"}
		}
		return domain.Album{}, fmt.Errorf("failed to get album: %w", err)
	}

	return toDomainAlbum(sqlAlbum), nil
}

func (r *SQLMediaBrowsingRepository) GetSongByID(ctx context.Context, id int) (domain.Song, error) {
	sqlSong, err := r.queries.GetSong(ctx, int32(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Song{}, &ports.NotFoundError{Message: "song not found"}
		}
		return domain.Song{}, fmt.Errorf("failed to get song: %w", err)
	}

	return toDomainSong(sqlSong), nil
}

func (r *SQLMediaBrowsingRepository) GetCoverByID(ctx context.Context, id string) (domain.Cover, error) {
	sqlCover, err := r.queries.GetCover(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Cover{}, &ports.NotFoundError{Message: "cover not found"}
		}
		return domain.Cover{}, fmt.Errorf("failed to get cover: %w", err)
	}

	return domain.Cover{
		Id:   sqlCover.CoverID,
		Path: sqlCover.Path,
	}, nil
}

func (r *SQLMediaBrowsingRepository) CreateArtist(ctx context.Context, artist domain.Artist) (domain.Artist, error) {
	var coverArt pgtype.Text
	if artist.CoverArt != "" {
		coverArt = pgtype.Text{String: artist.CoverArt, Valid: true}
	}
	var albumCount pgtype.Int4
	if artist.AlbumCount > 0 {
		albumCount = pgtype.Int4{Int32: int32(artist.AlbumCount), Valid: true}
	}

	sqlArtist, err := r.queries.CreateArtist(ctx, sqlc.CreateArtistParams{
		Name:       artist.Name,
		CoverArt:   coverArt,
		AlbumCount: albumCount,
	})
	if err != nil {
		return domain.Artist{}, fmt.Errorf("failed to create artist: %w", err)
	}

	return toDomainArtist(sqlArtist), nil
}

func (r *SQLMediaBrowsingRepository) CreateAlbum(ctx context.Context, album domain.Album) (domain.Album, error) {
	var artistID pgtype.Int4
	if album.ArtistId > 0 {
		artistID = pgtype.Int4{Int32: int32(album.ArtistId), Valid: true}
	}
	var coverArt pgtype.Text
	if album.CoverArt != "" {
		coverArt = pgtype.Text{String: album.CoverArt, Valid: true}
	}
	var songCount pgtype.Int4
	if album.SongCount > 0 {
		songCount = pgtype.Int4{Int32: int32(album.SongCount), Valid: true}
	}
	var duration pgtype.Int4
	if album.Duration > 0 {
		duration = pgtype.Int4{Int32: int32(album.Duration), Valid: true}
	}
	var artist pgtype.Text
	if album.Artist != "" {
		artist = pgtype.Text{String: album.Artist, Valid: true}
	}

	var created pgtype.Timestamp
	if album.Created != "" {
		t, err := time.Parse(time.RFC3339, album.Created)
		if err == nil {
			created = pgtype.Timestamp{Time: t, Valid: true}
		}
	}

	sqlAlbum, err := r.queries.CreateAlbum(ctx, sqlc.CreateAlbumParams{
		ArtistID:  artistID,
		Name:      album.Name,
		CoverArt:  coverArt,
		SongCount: songCount,
		Created:   created,
		Duration:  duration,
		Artist:    artist,
	})
	if err != nil {
		return domain.Album{}, fmt.Errorf("failed to create album: %w", err)
	}

	return toDomainAlbum(sqlAlbum), nil
}

func (r *SQLMediaBrowsingRepository) CreateSong(ctx context.Context, song domain.Song) (domain.Song, error) {
	var albumID pgtype.Int4
	if song.AlbumId > 0 {
		albumID = pgtype.Int4{Int32: int32(song.AlbumId), Valid: true}
	}
	var album pgtype.Text
	if song.Album != "" {
		album = pgtype.Text{String: song.Album, Valid: true}
	}
	var artist pgtype.Text
	if song.Artist != "" {
		artist = pgtype.Text{String: song.Artist, Valid: true}
	}
	isDir := pgtype.Bool{Bool: song.IsDir, Valid: true}
	var coverArt pgtype.Text
	if song.CoverArt != "" {
		coverArt = pgtype.Text{String: song.CoverArt, Valid: true}
	}
	var created pgtype.Timestamp
	if song.Created != "" {
		t, err := time.Parse(time.RFC3339, song.Created)
		if err == nil {
			created = pgtype.Timestamp{Time: t, Valid: true}
		}
	}
	var duration pgtype.Int4
	if song.Duration > 0 {
		duration = pgtype.Int4{Int32: int32(song.Duration), Valid: true}
	}
	var bitRate pgtype.Int4
	if song.BitRate > 0 {
		bitRate = pgtype.Int4{Int32: int32(song.BitRate), Valid: true}
	}
	var size pgtype.Int4
	if song.Size > 0 {
		// Note: pgtype.Int4 can only handle up to int32, but domain uses int64
		// This might cause issues with very large files
		size = pgtype.Int4{Int32: int32(song.Size), Valid: true}
	}
	var suffix pgtype.Text
	if song.Suffix != "" {
		suffix = pgtype.Text{String: song.Suffix, Valid: true}
	}
	var contentType pgtype.Text
	if song.ContentType != "" {
		contentType = pgtype.Text{String: song.ContentType, Valid: true}
	}
	isVideo := pgtype.Bool{Bool: song.IsVideo, Valid: true}

	sqlSong, err := r.queries.CreateSong(ctx, sqlc.CreateSongParams{
		AlbumID:     albumID,
		Title:       song.Title,
		Album:       album,
		Artist:      artist,
		IsDir:       isDir,
		CoverArt:    coverArt,
		Created:     created,
		Duration:    duration,
		BitRate:     bitRate,
		Size:        size,
		Suffix:      suffix,
		ContentType: contentType,
		IsVideo:     isVideo,
		Path:        song.Path,
	})
	if err != nil {
		return domain.Song{}, fmt.Errorf("failed to create song: %w", err)
	}

	return toDomainSong(sqlSong), nil
}

func (r *SQLMediaBrowsingRepository) CreateCover(ctx context.Context, cover domain.Cover) (domain.Cover, error) {
	sqlCover, err := r.queries.CreateCover(ctx, sqlc.CreateCoverParams{
		CoverID: cover.Id,
		Path:    cover.Path,
	})
	if err != nil {
		return domain.Cover{}, fmt.Errorf("failed to create cover: %w", err)
	}

	return domain.Cover{
		Id:   sqlCover.CoverID,
		Path: sqlCover.Path,
	}, nil
}

// Helper functions to convert between SQL and domain models

func toDomainArtist(sqlArtist sqlc.Artist) domain.Artist {
	artist := domain.Artist{
		Id:   int(sqlArtist.ArtistID),
		Name: sqlArtist.Name,
	}
	if sqlArtist.CoverArt.Valid {
		artist.CoverArt = sqlArtist.CoverArt.String
	}
	if sqlArtist.AlbumCount.Valid {
		artist.AlbumCount = int(sqlArtist.AlbumCount.Int32)
	}
	return artist
}

func toDomainAlbum(sqlAlbum sqlc.Album) domain.Album {
	album := domain.Album{
		Id:   int(sqlAlbum.AlbumID),
		Name: sqlAlbum.Name,
	}
	if sqlAlbum.ArtistID.Valid {
		album.ArtistId = int(sqlAlbum.ArtistID.Int32)
	}
	if sqlAlbum.CoverArt.Valid {
		album.CoverArt = sqlAlbum.CoverArt.String
	}
	if sqlAlbum.SongCount.Valid {
		album.SongCount = int(sqlAlbum.SongCount.Int32)
	}
	if sqlAlbum.Created.Valid {
		album.Created = sqlAlbum.Created.Time.Format(time.RFC3339)
	}
	if sqlAlbum.Duration.Valid {
		album.Duration = int(sqlAlbum.Duration.Int32)
	}
	if sqlAlbum.Artist.Valid {
		album.Artist = sqlAlbum.Artist.String
	}
	return album
}

func toDomainSong(sqlSong sqlc.Song) domain.Song {
	song := domain.Song{
		Id:    int(sqlSong.SongID),
		Title: sqlSong.Title,
		Path:  sqlSong.Path,
	}
	if sqlSong.AlbumID.Valid {
		song.AlbumId = int(sqlSong.AlbumID.Int32)
	}
	if sqlSong.Album.Valid {
		song.Album = sqlSong.Album.String
	}
	if sqlSong.Artist.Valid {
		song.Artist = sqlSong.Artist.String
	}
	if sqlSong.IsDir.Valid {
		song.IsDir = sqlSong.IsDir.Bool
	}
	if sqlSong.CoverArt.Valid {
		song.CoverArt = sqlSong.CoverArt.String
	}
	if sqlSong.Created.Valid {
		song.Created = sqlSong.Created.Time.Format(time.RFC3339)
	}
	if sqlSong.Duration.Valid {
		song.Duration = int(sqlSong.Duration.Int32)
	}
	if sqlSong.BitRate.Valid {
		song.BitRate = int(sqlSong.BitRate.Int32)
	}
	if sqlSong.Size.Valid {
		song.Size = int64(sqlSong.Size.Int32)
	}
	if sqlSong.Suffix.Valid {
		song.Suffix = sqlSong.Suffix.String
	}
	if sqlSong.ContentType.Valid {
		song.ContentType = sqlSong.ContentType.String
	}
	if sqlSong.IsVideo.Valid {
		song.IsVideo = sqlSong.IsVideo.Bool
	}
	return song
}
