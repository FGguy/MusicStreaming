package domain

// User represents a user in the system with role-based permissions
type User struct {
	Username            string
	Email               string
	Password            string
	ScrobblingEnabled   bool
	LdapAuthenticated   bool
	AdminRole           bool
	SettingsRole        bool
	StreamRole          bool
	JukeboxRole         bool
	DownloadRole        bool
	UploadRole          bool
	PlaylistRole        bool
	CoverArtRole        bool
	CommentRole         bool
	PodcastRole         bool
	ShareRole           bool
	VideoConversionRole bool
	MusicfolderId       []string
	MaxBitRate          int32
}
