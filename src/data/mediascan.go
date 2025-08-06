package data

import (
	"context"
	"music-streaming/types"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

/*
	TODO: create table for cover art, detect and create cover art entries
	TODO: get song file metadata
*/

func (d *DataLayerPg) MediaScan(musicFolders []string, count chan<- int, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	if len(musicFolders) < 1 {
		log.Print("no top music folders provided")
	}

	ctx := context.Background()

	for _, topDir := range musicFolders {
		entries, err := os.ReadDir(topDir)
		if err != nil {
			log.Printf("Failed reading content of top level directory: %s", err)
			return
		}

		var artists []*types.Artist
		for _, artist := range entries {
			if artist.IsDir() {
				a := &types.Artist{
					Name:     artist.Name(),
					CoverArt: "",
				}
				artists = append(artists, a)
			}
		}

		for _, artist := range artists {
			artistDir := path.Join(topDir, artist.Name)
			entries, err := os.ReadDir(artistDir) //fix
			if err != nil {
				log.Printf("Failed reading content of artist directory: %s", err)
				return
			}

			var albums []*types.Album
			for _, album := range entries {
				if album.IsDir() {
					a := &types.Album{
						Name:     album.Name(),
						CoverArt: "",
						Artist:   artist.Name,
					}
					albums = append(albums, a)
				}
			}

			for _, album := range albums {
				albumDir := path.Join(artistDir, album.Name)
				entries, err = os.ReadDir(albumDir)
				if err != nil {
					log.Printf("Failed reading content of album directory: %s", err) //fix
					return
				}

				var songs []*types.Song
				for _, song := range entries {
					if !song.IsDir() {
						info, err := os.Stat(path.Join(albumDir, song.Name()))
						if err != nil {
							log.Error().Err(err)
							return
						}
						s := &types.Song{
							Title:       song.Name(),
							Album:       album.Name,
							Artist:      artist.Name,
							IsDir:       false,
							CoverArt:    "",
							Duration:    0,
							BitRate:     0,
							Size:        info.Size(),
							Suffix:      "",
							ContentType: "",
							IsVideo:     false,
							Path:        path.Join(albumDir, song.Name()),
						}
						songs = append(songs, s)
					}
				}
				album.Songs = songs
			}

			artist.Albums = albums

			count <- 1 //new artist
			count <- len(artist.Albums)
			for _, a := range artist.Albums {
				count <- len(a.Songs)
			}

			if err := d.exportArtistTree(ctx, artist); err != nil {
				log.Error().Err(err).Msg("Failed to export artist")
			}
		}
	}
}
