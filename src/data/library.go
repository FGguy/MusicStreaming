package data

import (
	"context"
	sqlc "music-streaming/sql/sqlc"
	"music-streaming/types"

	"github.com/jackc/pgx/v5/pgtype"
)

func (d *DataLayerPg) exportArtistTree(ctx context.Context, artist *types.Artist) error {
	conn, err := d.Pg_pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)
	q := sqlc.New(conn)

	ar := sqlc.CreateArtistParams{
		Name:       artist.Name,
		CoverArt:   pgtype.Text{String: artist.CoverArt, Valid: true},
		AlbumCount: pgtype.Int4{Int32: int32(artist.GetAlbumCount()), Valid: true},
	}

	sqlArtist, err := q.CreateArtist(ctx, ar)
	if err != nil {
		return err
	}

	//export albums
	for _, album := range artist.Albums {
		al := sqlc.CreateAlbumParams{
			ArtistID:  pgtype.Int4{Int32: sqlArtist.ArtistID, Valid: true},
			Name:      album.Name,
			CoverArt:  pgtype.Text{String: album.CoverArt, Valid: true},
			SongCount: pgtype.Int4{Int32: int32(album.GetSongCount()), Valid: true},
			Duration:  pgtype.Int4{Int32: int32(album.GetDuration()), Valid: true},
			Artist:    pgtype.Text{String: sqlArtist.Name, Valid: true},
		}
		sqlAlbum, err := q.CreateAlbum(ctx, al)
		if err != nil {
			return err
		}

		if err = exportSongs(ctx, q, sqlAlbum.AlbumID, album.Songs); err != nil {
			return err
		}
	}

	conn.Commit(ctx)
	return nil
}

func exportSongs(ctx context.Context, q *sqlc.Queries, albumId int32, songs []*types.Song) error {
	for _, song := range songs {
		s := sqlc.CreateSongParams{
			AlbumID:     pgtype.Int4{Int32: albumId, Valid: true},
			Title:       song.Title,
			Album:       pgtype.Text{String: song.Album, Valid: true},
			Artist:      pgtype.Text{String: song.Artist, Valid: true},
			IsDir:       pgtype.Bool{Bool: song.IsDir, Valid: true},
			CoverArt:    pgtype.Text{String: song.CoverArt, Valid: true},
			Duration:    pgtype.Int4{Int32: int32(song.Duration), Valid: true},
			BitRate:     pgtype.Int4{Int32: int32(song.BitRate), Valid: true},
			Size:        pgtype.Int4{Int32: int32(song.Size), Valid: true},
			Suffix:      pgtype.Text{String: song.Suffix, Valid: true},
			ContentType: pgtype.Text{String: song.ContentType, Valid: true},
			IsVideo:     pgtype.Bool{Bool: song.IsVideo, Valid: true},
			Path:        song.Path,
		}
		if _, err := q.CreateSong(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
