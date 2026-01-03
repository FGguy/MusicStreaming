package services

import (
	"context"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"
	"sync"

	"music-streaming/internal/core/config"
)

type MediaScanningService struct {
	//repo       ports.MediaBrowsingRepository
	config     *config.Config
	scanStatus *domain.ScanStatus
	mu         sync.Mutex
	logger     *slog.Logger
}

func NewMediaScanningService(config *config.Config, logger *slog.Logger) *MediaScanningService {
	return &MediaScanningService{
		config: config,
		scanStatus: &domain.ScanStatus{
			Scanning: false,
			Count:    0,
		},
		logger: logger,
	}
}

func (s *MediaScanningService) StartScan(ctx context.Context) (domain.ScanStatus, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var username string
	if requestingUser != nil {
		username = requestingUser.Username
	}
	s.logger.Info("Start scan request", slog.String("username", username))

	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		s.logger.Warn("Unauthorized start scan attempt", slog.String("username", username))
		return domain.ScanStatus{}, &ports.NotAuthorizedError{Username: username, Action: "start media scan"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.scanStatus.Scanning {
		s.logger.Warn("Scan already in progress", slog.String("username", username))
		return *s.scanStatus, nil
	}

	s.scanStatus.Scanning = true
	s.logger.Info("Media scan started", slog.String("username", username))
	go s.Scan()

	return *s.scanStatus, nil
}

func (s *MediaScanningService) GetScanStatus(ctx context.Context) (domain.ScanStatus, error) {
	requestingUser, ok := ctx.Value(ports.KeyRequestingUserID).(*domain.User)
	var username string
	if requestingUser != nil {
		username = requestingUser.Username
	}
	s.logger.Debug("Get scan status request", slog.String("username", username))

	if !ok || requestingUser == nil || !requestingUser.AdminRole {
		s.logger.Warn("Unauthorized get scan status attempt", slog.String("username", username))
		return domain.ScanStatus{}, &ports.NotAuthorizedError{Username: username, Action: "get media scan status"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return *s.scanStatus, nil
}

func (s *MediaScanningService) Scan() {
	// Implementation of the scanning logic goes here
	s.logger.Info("Media scan process started")
	// TODO: Implement actual scanning logic
	s.logger.Info("Media scan process completed")
}
