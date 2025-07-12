package job

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AnalyzePayload struct {
	Input string `json:"input"`
}

type AnalyzeRunner struct {
	logger *logger.Logger
}

func NewAnalyzeRunner() *AnalyzeRunner {
	return &AnalyzeRunner{
		logger: logger.WithService("analyze-runner"),
	}
}

func (a *AnalyzeRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := a.logger.WithContext(ctx)

	var p AnalyzePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal analyze payload")
		return err
	}

	log.Info("Starting video analysis", "input", p.Input)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", p.Input)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg analysis failed", "input", p.Input, "output", string(out))
		return err
	}

	log.Info("Video analysis completed successfully", "input", p.Input, "output_length", len(out))
	log.Debug("FFmpeg analysis output", "output", string(out))
	return nil
}
