package data

import (
	"context"
	"music-streaming/internal/types"
	"music-streaming/internal/util"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

/*
	TODO: create table for cover art, detect and create cover art entries
*/

func (d *DataLayerPg) MediaScan(musicFolders []string, count chan<- int, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	if len(musicFolders) < 1 {
		log.Debug().Msg("no top music folders provided")
	}

	ctx := context.Background()

	for _, topDir := range musicFolders {
		log.Trace().Msgf("Scanning music folder: %s", topDir)
		entries, err := os.ReadDir(topDir)
		if err != nil {
			log.Debug().Msgf("Failed reading content of top level directory: %s", err)
			continue
		}

		var artists []*types.ScanArtist
		for _, artist := range entries {
			if artist.IsDir() {
				a := &types.ScanArtist{
					Name:     artist.Name(),
					CoverArt: "",
				}
				artists = append(artists, a)
			}
		}

		for _, artist := range artists {
			log.Trace().Msgf("Scanning artist: %s", artist.Name)
			artistDir := path.Join(topDir, artist.Name)
			entries, err := os.ReadDir(artistDir) //fix
			if err != nil {
				log.Debug().Msgf("Failed reading content of artist directory: %s", err)
				continue
			}

			var albums []*types.ScanAlbum
			for _, album := range entries {
				if album.IsDir() {
					a := &types.ScanAlbum{
						Name:     album.Name(),
						CoverArt: "",
						Artist:   artist.Name,
					}
					albums = append(albums, a)
				}
			}

			for _, album := range albums {
				log.Trace().Msgf("Scanning album: %s", album.Name)
				albumDir := path.Join(artistDir, album.Name)
				entries, err = os.ReadDir(albumDir)
				if err != nil {
					log.Debug().Msgf("Failed reading content of album directory: %s", err) //fix
					continue
				}

				var songs []*types.ScanSong
				for _, song := range entries {
					log.Trace().Msgf("Scanning song: %s", song.Name())
					if !song.IsDir() {
						info, err := util.FFProbeProcessFile(path.Join(albumDir, song.Name()))
						if err != nil {
							log.Error().Err(err).Msg("Failed processing file with ffprobe")
							continue
						}

						duration, err := strconv.ParseFloat(info.Format.Duration, 32)
						if err != nil {
							log.Warn().Err(err).Msgf("Failed to convert song duration to float")
							continue
						}

						bitrate, err := strconv.ParseInt(info.Format.BitRate, 10, 64)
						if err != nil {
							log.Warn().Err(err).Msgf("Failed to convert song bitrate to int")
							continue
						}

						size, err := strconv.ParseInt(info.Format.Size, 10, 64)
						if err != nil {
							log.Warn().Err(err).Msgf("Failed to convert song size to int")
							continue
						}

						s := &types.ScanSong{
							Title:       song.Name(),
							Album:       album.Name,
							Artist:      artist.Name,
							IsDir:       false,
							CoverArt:    "",
							Duration:    int(duration * 1000),
							BitRate:     int(bitrate),
							Size:        size,
							Suffix:      info.Format.FormatName,
							ContentType: getContentType(info.Format.FormatName),
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

func getContentType(formatName string) string {
	switch strings.ToLower(formatName) {
	case "mp3":
		return "audio/mpeg"
	case "flac":
		return "audio/flac"
	case "wav":
		return "audio/wav"
	case "mp4", "m4a":
		return "audio/mp4"
	case "aac":
		return "audio/aac"
	case "ogg":
		return "audio/ogg"
	case "asf":
		return "audio/x-ms-wma"
	default:
		return "audio/mpeg" // Default fallback
	}
}
