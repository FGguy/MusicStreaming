package domain

import "encoding/xml"

type User struct {
	XMLName             xml.Name `xml:"user" json:"-"`
	Username            string   `xml:"username,attr" json:"username" form:"username" binding:"required"`
	Email               string   `xml:"email,attr" json:"email" form:"email" binding:"required"`
	Password            string   `xml:"-" json:"password,omitempty" form:"password" binding:"required"`
	ScrobblingEnabled   bool     `xml:"scrobblingEnabled,attr" json:"scrobblingEnabled" form:"scrobblingEnabled"`
	LdapAuthenticated   bool     `xml:"ldapAuthenticated,attr" json:"ldapAuthenticated" form:"ldapAuthenticated"`
	AdminRole           bool     `xml:"adminRole,attr" json:"adminRole" form:"adminRole"`
	SettingsRole        bool     `xml:"settingsRole,attr" json:"settingsRole" form:"settingsRole"`
	StreamRole          bool     `xml:"streamRole,attr" json:"streamRole" form:"streamRole"`
	JukeboxRole         bool     `xml:"jukeboxRole,attr" json:"jukeboxRole" form:"jukeboxRole"`
	DownloadRole        bool     `xml:"downloadRole,attr" json:"downloadRole" form:"downloadRole"`
	UploadRole          bool     `xml:"uploadRole,attr" json:"uploadRole" form:"uploadRole"`
	PlaylistRole        bool     `xml:"playlistRole,attr" json:"playlistRole" form:"playlistRole"`
	CoverArtRole        bool     `xml:"coverArtRole,attr" json:"coverArtRole" form:"coverArtRole"`
	CommentRole         bool     `xml:"commentRole,attr" json:"commentRole" form:"commentRole"`
	PodcastRole         bool     `xml:"podcastRole,attr" json:"podcastRole" form:"podcastRole"`
	ShareRole           bool     `xml:"shareRole,attr" json:"shareRole" form:"shareRole"`
	VideoConversionRole bool     `xml:"videoConversionRole,attr" json:"videoConversionRole" form:"videoConversionRole"`
	MusicfolderId       []string `xml:"folder,omitempty" json:"folder,omitempty" form:"folder"`
	MaxBitRate          int32    `xml:"maxBitRate,attr" json:"maxBitRate" form:"maxBitRate"`
}
