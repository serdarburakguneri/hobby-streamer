package job

import (
	"context"
	"encoding/json"
	"log"
	"os/exec"
)

type DASHPayload struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type TranscodeDASHRunner struct{}

func (t *TranscodeDASHRunner) Run(ctx context.Context, payload json.RawMessage) error {
	var p DASHPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", p.Input, "-f", "dash", p.Output)
	out, err := cmd.CombinedOutput()
	log.Printf("DASH output: %s", string(out))
	return err
}