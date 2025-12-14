package domain

import "encoding/xml"

type User struct {
	XMLName             xml.Name `xml:"user" json:"-"`
	Username            string   `xml:"username,attr" json:"username"`
	Email               string   `xml:"email,attr" json:"email"`
	Password            string   `xml:"-" json:"password,omitempty"`
	ScrobblingEnabled   bool     `xml:"scrobblingEnabled,attr" json:"scrobblingEnabled"`
	LdapAuthenticated   bool     `xml:"ldapAuthenticated,attr" json:"ldapAuthenticated"`
	AdminRole           bool     `xml:"adminRole,attr" json:"adminRole"`
	SettingsRole        bool     `xml:"settingsRole,attr" json:"settingsRole"`
	StreamRole          bool     `xml:"streamRole,attr" json:"streamRole"`
	JukeboxRole         bool     `xml:"jukeboxRole,attr" json:"jukeboxRole"`
	DownloadRole        bool     `xml:"downloadRole,attr" json:"downloadRole"`
	UploadRole          bool     `xml:"uploadRole,attr" json:"uploadRole"`
	PlaylistRole        bool     `xml:"playlistRole,attr" json:"playlistRole"`
	CoverArtRole        bool     `xml:"coverArtRole,attr" json:"coverArtRole"`
	CommentRole         bool     `xml:"commentRole,attr" json:"commentRole"`
	PodcastRole         bool     `xml:"podcastRole,attr" json:"podcastRole"`
	ShareRole           bool     `xml:"shareRole,attr" json:"shareRole"`
	VideoConversionRole bool     `xml:"videoConversionRole,attr" json:"videoConversionRole"`
	MusicfolderId       []string `xml:"folder,omitempty" json:"folder,omitempty"`
	MaxBitRate          int32    `xml:"maxBitRate,attr" json:"maxBitRate"`
}
