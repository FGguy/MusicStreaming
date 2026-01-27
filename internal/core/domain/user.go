package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Validate checks if the User has valid field values
func (u *User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New("username is required")
	}
	if len(u.Username) < 3 || len(u.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters, got %d", len(u.Username))
	}

	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("email is not a valid email address: %s", u.Email)
	}

	if u.MaxBitRate < 0 {
		return fmt.Errorf("maxBitRate must be non-negative, got %d", u.MaxBitRate)
	}

	return nil
}
