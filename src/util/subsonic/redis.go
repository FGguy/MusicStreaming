package subsonic

type SubsonicRedisUser struct {
	Username            string   `redis:"username" json:"username"`
	Email               string   `redis:"email" json:"email"`
	Password            string   `redis:"password" json:"password"`
	ScrobblingEnabled   bool     `redis:"scrobblingEnabled" json:"scrobblingEnabled"`
	LdapAuthenticated   bool     `redis:"ldapAuthenticated" json:"ldapAuthenticated"`
	AdminRole           bool     `redis:"adminRole" json:"adminRole"`
	SettingsRole        bool     `redis:"settingsRole" json:"settingsRole"`
	StreamRole          bool     `redis:"streamRole" json:"streamRole"`
	JukeboxRole         bool     `redis:"jukeboxRole" json:"jukeboxRole"`
	DownloadRole        bool     `redis:"downloadRole" json:"downloadRole"`
	UploadRole          bool     `redis:"uploadRole" json:"uploadRole"`
	PlaylistRole        bool     `redis:"playlistRole" json:"playlistRole"`
	CoverArtRole        bool     `redis:"coverArtRole" json:"coverArtRole"`
	CommentRole         bool     `redis:"commentRole" json:"commentRole"`
	PodcastRole         bool     `redis:"podcastRole" json:"podcastRole"`
	ShareRole           bool     `redis:"shareRole" json:"shareRole"`
	VideoConversionRole bool     `redis:"videoConversionRole" json:"videoConversionRole"`
	MusicfolderId       []string `redis:"folder" json:"folder"`
	MaxBitRate          int32    `redis:"maxBitRate" json:"maxBitRate"`
}
