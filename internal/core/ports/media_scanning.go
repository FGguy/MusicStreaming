package ports

import (
	"context"
	"music-streaming/internal/core/domain"
)

type MediaScanningPort interface {
	StartScan(ctx context.Context) (domain.ScanStatus, error)
	GetScanStatus(ctx context.Context) (domain.ScanStatus, error)
	Scan()
}
