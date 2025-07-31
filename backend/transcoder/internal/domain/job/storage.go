package job

import "context"

//go:generate mockgen -destination=mock_storage.go -package job github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job Storage

type Storage interface {
	Download(ctx context.Context, input string) (string, error)
	CreateDir(path string) error
	Remove(path string) error
	RemoveAll(path string) error
	Upload(ctx context.Context, localDir, s3Path string) error
}
