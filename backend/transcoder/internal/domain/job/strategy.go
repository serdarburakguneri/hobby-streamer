package job

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/valueobjects"
)

type TranscodeStrategy interface {
	ValidateInput(ctx context.Context, job *entity.Job) error
	Transcode(ctx context.Context, job *entity.Job, localPath, outputDir string) (string, error)
	ExtractMetadata(ctx context.Context, filePath string, job *entity.Job) (*valueobjects.TranscodeMetadata, error)
	ValidateOutput(job *entity.Job) error
}

type TranscoderRegistry interface {
	Get(format string) TranscodeStrategy
}
