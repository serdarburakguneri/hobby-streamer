package job

import (
	"context"
	"os"
	"testing"
)

func TestDASHTranscoder_Transcode_Failure(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "input.mp4")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	d := &DASHTranscoder{}
	_, err = d.Transcode(context.Background(), nil, tmpFile.Name(), os.TempDir())
	if err == nil {
		t.Error("expected error when ffmpeg is not available or input is invalid")
	}
}

func TestDASHTranscoder_ExtractMetadata_InvalidFile(t *testing.T) {
	d := &DASHTranscoder{}
	_, err := d.ExtractMetadata(context.Background(), "/nonexistent/file.mpd", nil)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestDASHTranscoder_ValidateInputOutput(t *testing.T) {
	d := &DASHTranscoder{}
	if err := d.ValidateInput(context.Background(), nil); err != nil {
		t.Errorf("expected nil error for ValidateInput, got %v", err)
	}
	// Valid S3 output
	job := &Job{output: "s3://bucket/key"}
	if err := d.ValidateOutput(job); err != nil {
		t.Errorf("expected nil error for valid S3 output, got %v", err)
	}
	// Invalid S3 output
	job = &Job{output: "not-a-s3-path"}
	if err := d.ValidateOutput(job); err == nil {
		t.Error("expected error for invalid S3 output path")
	}
}
