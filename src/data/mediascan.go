package data

import (
	"io"
	"log"
	"music-streaming/types"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/dhowden/tag"
	"github.com/gabriel-vasile/mimetype"
	"github.com/tcolgate/mp3"
)

/*
	TODO: create table for cover art, detect and create cover art entries
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
						s, err := readInfoFromFile(path.Join(albumDir, song.Name()))
						if err != nil {
							log.Printf("Failed getting info from song file Error:%s", err)
							continue
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
					log.Println(s)
				}
			}
		}
	}
}

func readInfoFromFile(path string) (*types.Song, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		return nil, err
	}

	m, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var created time.Time
	stat, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		log.Print("Failed to get file information.")
		created = time.Unix(0, 0)
	} else {
		created = time.Unix(0, stat.CreationTime.Nanoseconds())
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var (
		duration time.Duration
	)
	d := mp3.NewDecoder(file)
	var frame mp3.Frame
	skipped := 0

	for {
		if err := d.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		duration += frame.Duration()
	}

	s := &types.Song{
		Title:       m.Title(),
		Album:       m.Album(),
		Artist:      m.Artist(),
		Genre:       m.Genre(),
		IsDir:       false,
		CoverArt:    "", //TODO
		Created:     created.String(),
		Duration:    int(duration.Seconds()), //TODO
		BitRate:     0,                       //TODO
		Size:        info.Size(),
		Suffix:      string(m.FileType()),
		ContentType: mtype.String(),
		IsVideo:     false,
		Path:        path,
	}

	return s, nil
}
