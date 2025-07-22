package job

import (
	"context"
	"os"
	"testing"
)

func TestHLSTranscoder_Transcode_Failure(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "input.mp4")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	h := &HLSTranscoder{}
	_, err = h.Transcode(context.Background(), nil, tmpFile.Name(), os.TempDir())
	if err == nil {
		t.Error("expected error when ffmpeg is not available or input is invalid")
	}
}
