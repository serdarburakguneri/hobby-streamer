package job

import (
	"context"
)

type TranscoderStrategy interface {
	Transcode(ctx context.Context, job *Job, localPath, outputDir string) (string, error)
	ExtractMetadata(ctx context.Context, filePath string) (*TranscodeMetadata, error)
}

type TranscoderRegistry struct {
	strategies map[string]TranscoderStrategy
}

func NewTranscoderRegistry() *TranscoderRegistry {
	return &TranscoderRegistry{
		strategies: map[string]TranscoderStrategy{
			"analyze": &AnalyzeStrategy{},
			"hls":     &HLSTranscoder{},
			"dash":    &DASHTranscoder{},
		},
	}
}

func (r *TranscoderRegistry) Get(format string) TranscoderStrategy {
	return r.strategies[format]
}

func (r *TranscoderRegistry) Register(key string, strategy TranscoderStrategy) {
	r.strategies[key] = strategy
}
