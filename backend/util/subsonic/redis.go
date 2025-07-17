package subsonic

type SubsonicRedisUser struct {
	Username            string   `redis:"username"`
	Email               string   `redis:"email"`
	Password            string   `redis:"password"`
	ScrobblingEnabled   bool     `redis:"scrobblingEnabled"`
	LdapAuthenticated   bool     `redis:"ldapAuthenticated"`
	AdminRole           bool     `redis:"adminRole"`
	SettingsRole        bool     `redis:"settingsRole"`
	StreamRole          bool     `redis:"streamRole"`
	JukeboxRole         bool     `redis:"jukeboxRole"`
	DownloadRole        bool     `redis:"downloadRole"`
	UploadRole          bool     `redis:"uploadRole"`
	PlaylistRole        bool     `redis:"playlistRole"`
	CoverArtRole        bool     `redis:"coverArtRole"`
	CommentRole         bool     `redis:"commentRole"`
	PodcastRole         bool     `redis:"podcastRole"`
	ShareRole           bool     `redis:"shareRole"`
	VideoConversionRole bool     `redis:"videoConversionRole"`
	MusicfolderId       []string `redis:"folder"`
	MaxBitRate          int32    `redis:"maxBitRate"`
}
