package data

import (
	"fmt"
	"log"
	"music-streaming/types"
	"os"
)

/*
	TODO: detect cover art
	TODO: join paths more elegantly
	TODO: get song file metadata
	TODO: export to sql
*/

func (d *DataLayerPg) MediaScan(musicFolders []string, count chan<- int, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	if len(musicFolders) < 1 {
		log.Print("no top music folders provided")
	}

	for _, topDir := range musicFolders {
		log.Printf("Processing folder: %s", topDir)
		//check content of folder
		entries, err := os.ReadDir(topDir)
		if err != nil {
			log.Printf("Failed reading content of top level directory: %s", err)
			return
		}

		//find artists
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

		//for each artist find albums
		for _, artist := range artists {
			log.Printf("Processing artist: %s", artist.Name)

			artistDir := fmt.Sprintf("%s/%s", topDir, artist.Name)
			entries, err := os.ReadDir(artistDir) //fix
			if err != nil {
				log.Printf("Failed reading content of artist directory: %s", err)
				return
			}

			var albums []*types.Album
			for _, album := range entries {
				if album.IsDir() {
					a := &types.Album{
						Name:   album.Name(),
						Artist: artist.Name,
					}
					albums = append(albums, a)
				}
			}

			//scan for songs
			for _, album := range albums {
				log.Printf("Processing album: %s", album.Name)

				albumDir := fmt.Sprintf("%s/%s", artistDir, album.Name)
				entries, err = os.ReadDir(albumDir)
				if err != nil {
					log.Printf("Failed reading content of album directory: %s", err) //fix
					return
				}

				var songs []*types.Song
				for _, song := range entries {
					if !song.IsDir() {
						s := &types.Song{
							Title:  song.Name(),
							Album:  album.Name,
							Artist: artist.Name,
						}
						songs = append(songs, s)
					}
				}
				album.Songs = songs
			}

			artist.Albums = albums

			//export artist to sql
			log.Println(artist.Name)
			count <- 1 //new artist
			count <- len(artist.Albums)
			for _, a := range artist.Albums {
				count <- len(a.Songs)
				log.Println(a.Name)
				for _, s := range a.Songs {
					log.Println(s.Title)
				}
			}
		}
	}
}
