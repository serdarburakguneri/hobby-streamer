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

func TestHLSTranscoder_ExtractMetadata_InvalidFile(t *testing.T) {
	h := &HLSTranscoder{}
	_, err := h.ExtractMetadata(context.Background(), "/nonexistent/file.m3u8", nil)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestHLSTranscoder_ValidateInputOutput(t *testing.T) {
	h := &HLSTranscoder{}
	if err := h.ValidateInput(context.Background(), nil); err != nil {
		t.Errorf("expected nil error for ValidateInput, got %v", err)
	}
	// Valid S3 output
	job := &Job{output: "s3://bucket/key"}
	if err := h.ValidateOutput(job); err != nil {
		t.Errorf("expected nil error for valid S3 output, got %v", err)
	}
	// Invalid S3 output
	job = &Job{output: "not-a-s3-path"}
	if err := h.ValidateOutput(job); err == nil {
		t.Error("expected error for invalid S3 output path")
	}
}
