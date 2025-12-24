package ports

import (
	"context"
	"music-streaming/internal/core/domain"
)

type MediaScanningPort interface {
	FFProbeProcessFile(path string) (*domain.FFProbeInfo, error)
	StartScan(ctx context.Context) (domain.ScanStatus, error)
	GetScanStatus(ctx context.Context) (domain.ScanStatus, error)
	Scan()
}
