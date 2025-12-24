package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"

	"music-streaming/internal/core/config"
)

type MediaScanningService struct {
	repo   ports.MediaBrowsingRepository
	logger *slog.Logger
	config *config.Config

	scanStatus *domain.ScanStatus
	mu         sync.Mutex
}

func NewMediaScanningService(repo ports.MediaBrowsingRepository, logger *slog.Logger, config *config.Config) *MediaScanningService {
	return &MediaScanningService{
		repo:   repo,
		logger: logger,
		config: config,
		scanStatus: &domain.ScanStatus{
			Scanning: false,
			Count:    0,
		},
	}
}

func (s *MediaScanningService) StartScan(ctx context.Context) (domain.ScanStatus, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		return domain.ScanStatus{}, &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "start media scan"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// If a scan is already in progress, return the current status
	if s.scanStatus.Scanning {
		return *s.scanStatus, nil
	} else {
		s.scanStatus.Count = 0
		s.scanStatus.Scanning = true
		go s.Scan()
	}

	return *s.scanStatus, nil
}

func (s *MediaScanningService) GetScanStatus(ctx context.Context) (domain.ScanStatus, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		return domain.ScanStatus{}, &ports.NotAuthorizedError{Username: requestingUser.Username, Action: "get media scan status"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return *s.scanStatus, nil
}

// TODO: Refactor into smaller functions
func (s *MediaScanningService) Scan() {

	// Validate music directories
	if s.config.MusicDirectories == nil || len(s.config.MusicDirectories) == 0 {
		s.logger.Error("No music directories configured for scanning.")
		s.mu.Lock()
		s.scanStatus.Scanning = false
		s.mu.Unlock()
		return
	}

	s.logger.Debug("Starting Media Scan with music directories: %q", s.config.MusicDirectories)

	var (
		ctx     = context.Background()
		albums  = make([]domain.Album, 0)
		artists = make([]domain.Artist, 0)
	)

	defer func() {
		s.mu.Lock()
		s.scanStatus.Scanning = false
		s.mu.Unlock()
		s.logger.Debug("Media Scan completed.")
	}()

	for _, dir := range s.config.MusicDirectories {
		s.logger.Debug("Scanning directory: %s", dir)

		// scan all artists in the directory
		entries, err := os.ReadDir(dir)
		if err != nil {
			s.logger.Debug("Failed reading content of top level directory: %s", err)
			continue
		}

		for _, artistEntry := range entries {
			if artistEntry.IsDir() {
				fullpath := path.Join(dir, artistEntry.Name())
				coverId, coverPath := detectCover(fullpath)
				artist := domain.Artist{
					Name:     artistEntry.Name(),
					CoverArt: coverId,
				}
				if coverId != "" {
					cover := domain.Cover{
						Id:   coverId,
						Path: coverPath,
					}

					s.logger.Debug("Exporting %s", coverId)
					if _, err := s.repo.CreateCover(ctx, cover); err != nil {
						s.logger.Error("Failed to create cover art in database: %s", err)
					}
				}

				// return artist with ID
				s.mu.Lock()
				s.scanStatus.Count++
				s.mu.Unlock()

				artist, err = s.repo.CreateArtist(ctx, artist)
				if err != nil {
					s.logger.Error("Failed to create artist in database: %s", err)
					continue
				}

				artists = append(artists, artist)
			}
		}

		// scan all albums for an artist
		for _, artist := range artists {
			artistPath := path.Join(dir, artist.Name)
			albumEntries, err := os.ReadDir(artistPath)
			if err != nil {
				s.logger.Debug("Failed reading content of artist directory: %s", err)
				continue
			}

			for _, albumEntry := range albumEntries {
				if albumEntry.IsDir() {
					fullAlbumPath := path.Join(artistPath, albumEntry.Name())
					coverId, coverPath := detectCover(fullAlbumPath)
					album := domain.Album{
						ArtistId: artist.Id,
						Artist:   artist.Name,
						Name:     albumEntry.Name(),
						CoverArt: coverId,
					}
					if coverId != "" {
						cover := domain.Cover{
							Id:   coverId,
							Path: coverPath,
						}

						s.logger.Debug("Exporting %s", coverId)
						if _, err := s.repo.CreateCover(ctx, cover); err != nil {
							s.logger.Error("Failed to create cover art in database: %s", err)
						}
					}

					s.mu.Lock()
					s.scanStatus.Count++
					s.mu.Unlock()

					album, err = s.repo.CreateAlbum(ctx, album)
					if err != nil {
						s.logger.Error("Failed to create album in database: %s", err)
						continue
					}
					albums = append(albums, album)
				}
			}
		}

		// scan all songs for an album
		for _, album := range albums {
			albumPath := path.Join(dir, album.Artist, album.Name)
			songEntries, err := os.ReadDir(albumPath)
			if err != nil {
				s.logger.Debug("Failed reading content of album directory: %s", err)
				continue
			}

			for _, songEntry := range songEntries {
				if !songEntry.IsDir() {
					songPath := path.Join(albumPath, songEntry.Name())
					info, err := s.FFProbeProcessFile(songPath)
					if err != nil {
						s.logger.Warn("FFProbe failed processing file %s: %s", songEntry.Name(), err)
						continue
					}

					duration, err := strconv.ParseFloat(info.Format.Duration, 32)
					if err != nil {
						s.logger.Warn("Failed to convert song duration to float: %s", err)
						continue
					}

					bitrate, err := strconv.ParseInt(info.Format.BitRate, 10, 64)
					if err != nil {
						s.logger.Warn("Failed to convert song bitrate to int: %s", err)
						continue
					}

					size, err := strconv.ParseInt(info.Format.Size, 10, 64)
					if err != nil {
						s.logger.Warn("Failed to convert song size to int: %s", err)
						continue
					}

					song := domain.Song{
						Title:       songEntry.Name(),
						Album:       album.Name,
						Artist:      album.Artist,
						IsDir:       false,
						CoverArt:    "",
						Duration:    int(duration * 1000),
						BitRate:     int(bitrate),
						Size:        size,
						Suffix:      info.Format.FormatName,
						ContentType: getContentType(info.Format.FormatName),
						IsVideo:     false,
						Path:        songPath,
					}

					s.mu.Lock()
					s.scanStatus.Count++
					s.mu.Unlock()

					_, err = s.repo.CreateSong(ctx, song)
					if err != nil {
						s.logger.Error("Failed to create song in database: %s", err)
						continue
					}
				}
			}
		}
	}
}

var validCoverFormats = []string{
	"cover.jpg", "cover.jpeg", "cover.png",
}

func detectCover(topPath string) (cover_id string, cover_path string) {
	for _, format := range validCoverFormats {
		cover_path := path.Join(topPath, format)
		_, err := os.Stat(cover_path)
		if err == nil {
			h := sha256.New()
			h.Write([]byte(cover_path))
			return hex.EncodeToString(h.Sum(nil))[0:16], cover_path
		}
	}
	return "", ""
}

func getContentType(formatName string) string {
	switch strings.ToLower(formatName) {
	case "mp3":
		return "audio/mpeg"
	case "flac":
		return "audio/flac"
	case "wav":
		return "audio/wav"
	case "mp4", "m4a":
		return "audio/mp4"
	case "aac":
		return "audio/aac"
	case "ogg":
		return "audio/ogg"
	case "asf":
		return "audio/x-ms-wma"
	default:
		return "audio/mpeg" // Default fallback
	}
}

func (s *MediaScanningService) FFProbeProcessFile(filePath string) (*domain.FFProbeInfo, error) {
	//build command
	cmd := exec.Command("ffprobe", "-v", "error", "-hide_banner", "-print_format", "json", "-show_format", filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info domain.FFProbeInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
