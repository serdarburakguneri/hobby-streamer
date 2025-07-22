package job

import (
	"context"
	"os"
	"testing"
)

func TestAnalyzeStrategy_Transcode_Failure(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "input.mp4")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	a := &AnalyzeStrategy{}
	_, err = a.Transcode(context.Background(), nil, tmpFile.Name(), "")
	if err == nil {
		t.Error("expected error when input is invalid or ffprobe/ffmpeg is not available")
	}
}
