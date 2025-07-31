package transcoding

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job"
)

type Registry struct {
	strategies map[string]job.TranscodeStrategy
}

func NewRegistry(storage job.Storage) *Registry {
	return &Registry{
		strategies: map[string]job.TranscodeStrategy{
			"analyze": NewAnalyzeTranscoder(),
			"hls":     NewHLSTranscoder(storage),
			"dash":    NewDASHTranscoder(storage),
		},
	}
}

func (r *Registry) Get(format string) job.TranscodeStrategy {
	return r.strategies[format]
}
