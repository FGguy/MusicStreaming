package subsonic

import (
	"encoding/xml"
)

type SubsonicXmlResponse struct {
	XMLName xml.Name           `xml:"subsonic-response"`
	Xmlns   string             `xml:"xmlns,attr"`
	Status  string             `xml:"status,attr"`
	Version string             `xml:"version,attr"`
	Error   *SubsonicXmlError  `xml:"error,omitempty"`
	User    *SubsonicXmlUser   `xml:"user,omitempty"`
	Users   []*SubsonicXmlUser `xml:"users>user,omitempty"`
}

type SubsonicXmlError struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:"message,attr"`
}

type SubsonicXmlUser struct {
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
