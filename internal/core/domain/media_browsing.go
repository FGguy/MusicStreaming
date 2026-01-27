package domain

// Artist represents a music artist in the domain
type Artist struct {
	Id         int
	Name       string
	CoverArt   string
	AlbumCount int
}

// Album represents a music album in the domain
type Album struct {
	Id        int
	ArtistId  int
	Name      string
	CoverArt  string
	SongCount int
	Created   string
	Duration  int
	Artist    string
}

// Song represents a music track in the domain
type Song struct {
	Id          int
	AlbumId     int
	Title       string
	Album       string
	Artist      string
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

// Cover represents album or artist cover art
type Cover struct {
	Id   string
	Path string
}
