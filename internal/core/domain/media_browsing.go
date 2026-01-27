package domain

import (
	"errors"
	"fmt"
	"strings"
)

// Artist represents a music artist in the domain
type Artist struct {
	Id         int
	Name       string
	CoverArt   string
	AlbumCount int
}

// Validate checks if the Artist has valid field values
func (a *Artist) Validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("artist name is required")
	}
	if a.AlbumCount < 0 {
		return fmt.Errorf("album count must be non-negative, got %d", a.AlbumCount)
	}
	return nil
}

// Album represents a music album in the domain
type Album struct {
	Id        int
	ArtistId  int
	Name      string
	CoverArt  string
	SongCount int
	Created   string
	Duration  int
	Artist    string
}

// Validate checks if the Album has valid field values
func (a *Album) Validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("album name is required")
	}
	if a.SongCount < 0 {
		return fmt.Errorf("song count must be non-negative, got %d", a.SongCount)
	}
	if a.Duration < 0 {
		return fmt.Errorf("duration must be non-negative, got %d", a.Duration)
	}
	return nil
}

// Song represents a music track in the domain
type Song struct {
	Id          int
	AlbumId     int
	Title       string
	Album       string
	Artist      string
	IsDir       bool
	CoverArt    string
	Created     string
	Duration    int
	BitRate     int
	Size        int64
	Suffix      string
	ContentType string
	IsVideo     bool
	Path        string
}

// Validate checks if the Song has valid field values
func (s *Song) Validate() error {
	if strings.TrimSpace(s.Title) == "" {
		return errors.New("song title is required")
	}
	if s.Duration < 0 {
		return fmt.Errorf("duration must be non-negative, got %d", s.Duration)
	}
	if s.BitRate < 0 {
		return fmt.Errorf("bitrate must be non-negative, got %d", s.BitRate)
	}
	if s.Size < 0 {
		return fmt.Errorf("size must be non-negative, got %d", s.Size)
	}
	if !s.IsDir && strings.TrimSpace(s.Path) == "" {
		return errors.New("path is required for non-directory songs")
	}
	return nil
}

// Cover represents album or artist cover art
type Cover struct {
	Id   string
	Path string
}

// Validate checks if the Cover has valid field values
func (c *Cover) Validate() error {
	if strings.TrimSpace(c.Id) == "" {
		return errors.New("cover id is required")
	}
	if strings.TrimSpace(c.Path) == "" {
		return errors.New("cover path is required")
	}
	return nil
}
