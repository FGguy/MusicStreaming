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

func NewUser(username string, email string, password string, musicFolderId []string, options ...UserOptions) *User {
	user := &User{
		Username:            username,
		Email:               email,
		Password:            password,
		ScrobblingEnabled:   false,
		LdapAuthenticated:   false,
		AdminRole:           false,
		SettingsRole:        false,
		StreamRole:          true,
		JukeboxRole:         false,
		DownloadRole:        false,
		UploadRole:          false,
		PlaylistRole:        false,
		CoverArtRole:        false,
		CommentRole:         false,
		PodcastRole:         false,
		ShareRole:           false,
		VideoConversionRole: false,
		MaxBitRate:          0,
		MusicfolderId:       musicFolderId,
	}

	for _, option := range options {
		option(user)
	}

	return user
}

type UserOptions func(*User)

func WithAdminRole(adminRole bool) UserOptions {
	return func(u *User) {
		u.AdminRole = adminRole
	}
}

func WithSettingsRole(settingsRole bool) UserOptions {
	return func(u *User) {
		u.SettingsRole = settingsRole
	}
}

func WithStreamRole(streamRole bool) UserOptions {
	return func(u *User) {
		u.StreamRole = streamRole
	}
}

func WithJukeboxRole(jukeboxRole bool) UserOptions {
	return func(u *User) {
		u.JukeboxRole = jukeboxRole
	}
}

func WithDownloadRole(downloadRole bool) UserOptions {
	return func(u *User) {
		u.DownloadRole = downloadRole
	}
}

func WithUploadRole(uploadRole bool) UserOptions {
	return func(u *User) {
		u.UploadRole = uploadRole
	}
}

func WithPlaylistRole(playlistRole bool) UserOptions {
	return func(u *User) {
		u.PlaylistRole = playlistRole
	}
}

func WithCoverArtRole(coverArtRole bool) UserOptions {
	return func(u *User) {
		u.CoverArtRole = coverArtRole
	}
}

func WithCommentRole(commentRole bool) UserOptions {
	return func(u *User) {
		u.CommentRole = commentRole
	}
}

func WithPodcastRole(podcastRole bool) UserOptions {
	return func(u *User) {
		u.PodcastRole = podcastRole
	}
}

func WithShareRole(shareRole bool) UserOptions {
	return func(u *User) {
		u.ShareRole = shareRole
	}
}

func WithVideoConversionRole(videoConversionRole bool) UserOptions {
	return func(u *User) {
		u.VideoConversionRole = videoConversionRole
	}
}

func WithScrobblingEnabled(scrobblingEnabled bool) UserOptions {
	return func(u *User) {
		u.ScrobblingEnabled = scrobblingEnabled
	}
}

func WithLdapAuthenticated(ldapAuthenticated bool) UserOptions {
	return func(u *User) {
		u.LdapAuthenticated = ldapAuthenticated
	}
}

func WithMaxBitRate(maxBitRate int32) UserOptions {
	return func(u *User) {
		u.MaxBitRate = maxBitRate
	}
}
