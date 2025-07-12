package job

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type HLSPayload struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type TranscodeHLSRunner struct {
	logger *logger.Logger
}

func NewTranscodeHLSRunner() *TranscodeHLSRunner {
	return &TranscodeHLSRunner{
		logger: logger.WithService("hls-runner"),
	}
}

func (h *TranscodeHLSRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := h.logger.WithContext(ctx)

	var p HLSPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal HLS payload")
		return err
	}

	log.Info("Starting HLS transcoding", "input", p.Input, "output", p.Output)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", p.Input, "-c:v", "libx264", "-c:a", "aac", "-f", "hls", p.Output)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg HLS transcoding failed", "input", p.Input, "output", p.Output, "ffmpeg_output", string(out))
		return err
	}

	log.Info("HLS transcoding completed successfully", "input", p.Input, "output", p.Output, "output_length", len(out))
	log.Debug("FFmpeg HLS output", "output", string(out))
	return nil
}
