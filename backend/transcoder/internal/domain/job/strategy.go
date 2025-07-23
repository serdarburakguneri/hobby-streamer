package job

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
)

type TranscoderStrategy interface {
	ValidateOutput(job *Job) error
	Transcode(ctx context.Context, job *Job, localPath, outputDir string) (string, error)
	ExtractMetadata(ctx context.Context, filePath string, job *Job) (*TranscodeMetadata, error)
	ValidateInput(ctx context.Context, job *Job) error
}

type TranscoderRegistry struct {
	strategies map[string]TranscoderStrategy
}

func NewTranscoderRegistry(s3Client *s3.Client) *TranscoderRegistry {
	return &TranscoderRegistry{
		strategies: map[string]TranscoderStrategy{
			"analyze": &AnalyzeTranscoder{},
			"hls":     NewHLSTranscoder(s3Client),
			"dash":    NewDASHTranscoder(s3Client),
		},
	}
}

func (r *TranscoderRegistry) Get(format string) TranscoderStrategy {
	return r.strategies[format]
}

func (r *TranscoderRegistry) Register(key string, strategy TranscoderStrategy) {
	r.strategies[key] = strategy
}
