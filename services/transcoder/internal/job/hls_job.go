package job

import (
	"context"
	"encoding/json"
	"log"
	"os/exec"
)

type HLSPayload struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type TranscodeHLSRunner struct{}

func (t *TranscodeHLSRunner) Run(ctx context.Context, payload json.RawMessage) error {
	var p HLSPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", p.Input, "-f", "hls", p.Output)
	out, err := cmd.CombinedOutput()
	log.Printf("HLS output: %s", string(out))
	return err
}