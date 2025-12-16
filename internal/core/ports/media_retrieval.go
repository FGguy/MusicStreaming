package ports

import (
	"context"
	"music-streaming/internal/core/domain"
)

// For now only return a file path as string
// Will later handle file format conversion

type MediaRetrievalPort interface {
	DownloadSong(ctx context.Context, id int) (domain.Song, error)
	StreamSong(ctx context.Context, id int) (domain.Song, error)
	GetCover(ctx context.Context, id int) (domain.Cover, error)
}
