package types

type Artist struct {
	Id         int    `json:"id" xml:"id"`
	Name       string `json:"name" xml:"name"`
	CoverArt   string `json:"coverArt" xml:"coverArt"`
	AlbumCount int    `json:"albumCount" xml:"albumCount"`
}

type Album struct {
	Id        int32  `json:"id" xml:"id"`
	ArtistId  int32  `json:"artistId" xml:"artistId"`
	Name      string `json:"name" xml:"name"`
	CoverArt  string `json:"coverArt" xml:"coverArt"`
	SongCount int    `json:"songCount" xml:"songCount"`
	Created   string `json:"created" xml:"created"`
	Duration  int    `json:"duration" xml:"duration"`
	Artist    string `json:"artist" xml:"artist"`
}

type Song struct {
	Id          int32  `json:"id" xml:"id"`
	AlbumId     int    `json:"albumId" xml:"albumId"`
	Title       string `json:"title" xml:"title"`
	Album       string `json:"album" xml:"album"`
	Artist      string `json:"artist" xml:"artist"`
	IsDir       bool   `json:"isDir" xml:"isDir"`
	CoverArt    string `json:"coverArt" xml:"coverArt"`
	Created     string `json:"created" xml:"created"`
	Duration    int    `json:"duration" xml:"duration"`
	BitRate     int    `json:"bitRate" xml:"bitRate"`
	Size        int    `json:"size" xml:"size"`
	Suffix      string `json:"suffix" xml:"suffix"`
	ContentType string `json:"contentType" xml:"contentType"`
	IsVideo     bool   `json:"isVideo" xml:"isVideo"`
	Path        string `json:"path" xml:"path"`
}
