package types

type Artist struct {
	Name     string
	CoverArt string
	Albums   []*Album
}

func (a *Artist) GetAlbumCount() int {
	return len(a.Albums)
}

type Album struct {
	Name     string
	CoverArt string
	Created  string
	Artist   string
	Songs    []*Song
}

func (a *Album) GetSongCount() int {
	return len(a.Songs)
}

func (a *Album) GetDuration() int {
	duration := 0
	for _, s := range a.Songs {
		duration += s.Duration
	}
	return duration
}

type Song struct {
	Title       string
	Album       string
	Artist      string
	IsDir       bool
	CoverArt    string
	Created     string
	Duration    int
	BitRate     int
	Size        int
	Suffix      string
	ContentType string
	IsVideo     bool
	Path        string
}
