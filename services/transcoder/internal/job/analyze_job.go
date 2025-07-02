package job

import (
	"context"
	"encoding/json"
	"log"
	"os/exec"
)

type AnalyzePayload struct {
	Input string `json:"input"`
}

type AnalyzeRunner struct{}

func (a *AnalyzeRunner) Run(ctx context.Context, payload json.RawMessage) error {
	var p AnalyzePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", p.Input)
	out, err := cmd.CombinedOutput()
	log.Printf("Analyze output: %s", string(out))
	return err
}