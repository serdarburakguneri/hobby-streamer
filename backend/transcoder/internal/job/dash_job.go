package job

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type DASHPayload struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type TranscodeDASHRunner struct {
	logger *logger.Logger
}

func NewTranscodeDASHRunner() *TranscodeDASHRunner {
	return &TranscodeDASHRunner{
		logger: logger.WithService("dash-runner"),
	}
}

func (d *TranscodeDASHRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := d.logger.WithContext(ctx)

	var p DASHPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal DASH payload")
		return err
	}

	log.Info("Starting DASH transcoding", "input", p.Input, "output", p.Output)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", p.Input, "-c:v", "libx264", "-c:a", "aac", "-f", "dash", p.Output)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg DASH transcoding failed", "input", p.Input, "output", p.Output, "ffmpeg_output", string(out))
		return err
	}

	log.Info("DASH transcoding completed successfully", "input", p.Input, "output", p.Output, "output_length", len(out))
	log.Debug("FFmpeg DASH output", "output", string(out))
	return nil
}
