package ports

import (
	"context"
	"music-streaming/internal/core/domain"
)

// MediaRetrievalPort defines the interface for media streaming and download operations.
// It provides methods to retrieve songs for streaming or downloading, and cover art.
//
// Note: Current implementation returns file paths as strings.
// Future enhancements will include format conversion and transcoding.
type MediaRetrievalPort interface {
	// DownloadSong retrieves a song for download by its ID.
	// Requires download role permission.
	DownloadSong(ctx context.Context, id int) (domain.Song, error)

	// StreamSong retrieves a song for streaming by its ID.
	// Requires stream role permission. May apply bitrate transcoding based on user preferences.
	StreamSong(ctx context.Context, id int) (domain.Song, error)

	// GetCover retrieves cover art metadata for display.
	// Requires cover art role permission.
	GetCover(ctx context.Context, id string) (domain.Cover, error)
}
