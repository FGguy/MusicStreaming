package domain

type ScanArtist struct {
	Name     string
	CoverArt string
	Albums   []*ScanAlbum
}

func (a *ScanArtist) GetAlbumCount() int {
	return len(a.Albums)
}

type ScanAlbum struct {
	Name     string
	CoverArt string
	Created  string
	Artist   string
	Songs    []*ScanSong
}

func (a *ScanAlbum) GetSongCount() int {
	return len(a.Songs)
}

func (a *ScanAlbum) GetDuration() int {
	duration := 0
	for _, s := range a.Songs {
		duration += s.Duration
	}
	return duration
}

type ScanSong struct {
	Title       string
	Album       string
	Artist      string
	Genre       string
	IsDir       bool
	CoverArt    string
	Created     string
	Duration    int
	BitRate     int
	Size        int64
	Suffix      string
	ContentType string
	IsVideo     bool
	Path        string
}
