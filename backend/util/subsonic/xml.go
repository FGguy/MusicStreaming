package subsonic

import (
	"encoding/xml"
)

type SubsonicResponse struct {
	XMLName xml.Name        `xml:"subsonic-response"`
	Xmlns   string          `xml:"xmlns,attr"`
	Status  string          `xml:"status,attr"`
	Version string          `xml:"version,attr"`
	Error   *SubsonicError  `xml:"error,omitempty"`
	User    *SubsonicUser   `xml:"user,omitempty"`
	Users   []*SubsonicUser `xml:"users>user,omitempty"`
}

type SubsonicError struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:"message,attr"`
}

type SubsonicUser struct {
	XMLName             xml.Name `xml:"user"`
	Username            string   `xml:"username,attr"`
	Email               string   `xml:"email,attr"`
	ScrobblingEnabled   bool     `xml:"scrobblingEnabled,attr"`
	LdapAuthenticated   bool     `xml:"ldapAuthenticated,attr"`
	AdminRole           bool     `xml:"adminRole,attr"`
	SettingsRole        bool     `xml:"settingsRole,attr"`
	StreamRole          bool     `xml:"streamRole,attr"`
	JukeboxRole         bool     `xml:"jukeboxRole,attr"`
	DownloadRole        bool     `xml:"downloadRole,attr"`
	UploadRole          bool     `xml:"uploadRole,attr"`
	PlaylistRole        bool     `xml:"playlistRole,attr"`
	CoverArtRole        bool     `xml:"coverArtRole,attr"`
	CommentRole         bool     `xml:"commentRole,attr"`
	PodcastRole         bool     `xml:"podcastRole,attr"`
	ShareRole           bool     `xml:"shareRole,attr"`
	VideoConversionRole bool     `xml:"videoConversionRole,attr"`
	MusicfolderId       []string `xml:"folder,omitempty"`
	MaxBitRate          int32    `xml:"maxBitRate,attr"`
}
