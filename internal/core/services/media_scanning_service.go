package services

import (
	"context"
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
}

func NewMediaScanningService(config *config.Config) *MediaScanningService {
	return &MediaScanningService{
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

	s.scanStatus.Scanning = true
	go s.Scan()

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

func (s *MediaScanningService) Scan() {
	// Implementation of the scanning logic goes here
}
