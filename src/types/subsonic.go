package types

import (
	"encoding/xml"
	sqlc "music-streaming/sql/sqlc"
	"strings"
)

type SubsonicResponse struct {
	XMLName xml.Name        `xml:"subsonic-response" json:"-"`
	Xmlns   string          `xml:"xmlns,attr" json:"-"`
	Status  string          `xml:"status,attr" json:"status"`
	Version string          `xml:"version,attr" json:"version"`
	Error   *SubsonicError  `xml:"error,omitempty" json:"error,omitempty"`
	User    *SubsonicUser   `xml:"user,omitempty" json:"user,omitempty"`
	Users   []*SubsonicUser `xml:"users,omitempty" json:"users,omitempty"`
}

type SubsonicError struct {
	XMLName xml.Name `xml:"error,omitempty" json:"-"`
	Code    string   `xml:"code,attr" json:"code"`
	Message string   `xml:"message,attr" json:"message"`
}

type SubsonicUser struct {
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

func MapSqlUserToSubsonicUser(user *sqlc.User, password string) *SubsonicUser {
	return &SubsonicUser{
		Username:            user.Username,
		Email:               user.Email,
		Password:            password,
		ScrobblingEnabled:   user.Scrobblingenabled,
		LdapAuthenticated:   user.Ldapauthenticated,
		AdminRole:           user.Adminrole,
		SettingsRole:        user.Settingsrole,
		StreamRole:          user.Streamrole,
		JukeboxRole:         user.Jukeboxrole,
		DownloadRole:        user.Downloadrole,
		UploadRole:          user.Uploadrole,
		PlaylistRole:        user.Playlistrole,
		CoverArtRole:        user.Coverartrole,
		CommentRole:         user.Commentrole,
		PodcastRole:         user.Podcastrole,
		ShareRole:           user.Sharerole,
		VideoConversionRole: user.Videoconversionrole,
		MusicfolderId:       strings.Split(user.Musicfolderid.String, ";"),
		MaxBitRate:          user.Maxbitrate,
	}
}
