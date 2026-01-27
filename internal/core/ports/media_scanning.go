package ports

import (
	"context"
	"music-streaming/internal/core/domain"
)

// MediaScanningPort defines the interface for media library scanning operations.
// It provides methods to scan media directories, extract metadata, and track scan status.
type MediaScanningPort interface {
	// FFProbeProcessFile extracts media metadata from a file using ffprobe.
	// Returns structured metadata information about the media file.
	FFProbeProcessFile(path string) (*domain.FFProbeInfo, error)

	// StartScan initiates a media library scan.
	// Requires admin role permission. Returns the current scan status.
	StartScan(ctx context.Context) (domain.ScanStatus, error)

	// GetScanStatus retrieves the current status of the media library scan.
	// Returns information about whether a scan is in progress and file count.
	GetScanStatus(ctx context.Context) (domain.ScanStatus, error)

	// Scan performs the actual media library scanning operation.
	// This is typically called as a background goroutine by StartScan.
	Scan()
}
