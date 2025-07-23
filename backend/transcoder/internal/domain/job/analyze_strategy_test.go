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

func TestAnalyzeStrategy_ExtractMetadata_InvalidFile(t *testing.T) {
	a := &AnalyzeStrategy{}
	_, err := a.ExtractMetadata(context.Background(), "/nonexistent/file.mp4", nil)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestAnalyzeStrategy_ValidateInputOutput(t *testing.T) {
	a := &AnalyzeStrategy{}
	if err := a.ValidateInput(context.Background(), nil); err != nil {
		t.Errorf("expected nil error for ValidateInput, got %v", err)
	}
	if err := a.ValidateOutput(nil); err != nil {
		t.Errorf("expected nil error for ValidateOutput, got %v", err)
	}
}

func TestJob_StateTransitions(t *testing.T) {
	assetID, _ := NewAssetID("asset-1")
	videoID, _ := NewVideoID("video-1")
	job := NewAnalyzeJob(*assetID, *videoID, "input.mp4")
	if !job.IsPending() {
		t.Error("expected job to be pending initially")
	}
	job.Start()
	if !job.IsRunning() {
		t.Error("expected job to be running after Start")
	}
	job.UpdateProgress(50)
	if job.Progress() != 50 {
		t.Errorf("expected progress 50, got %v", job.Progress())
	}
	job.Complete(nil)
	if !job.IsCompleted() {
		t.Error("expected job to be completed after Complete")
	}
	job.Fail("fail")
	if !job.IsFailed() {
		t.Error("expected job to be failed after Fail")
	}
}
