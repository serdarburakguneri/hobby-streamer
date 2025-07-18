package job

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
)

func TestAnalyzeRunner_Run_NilPayload(t *testing.T) {
	runner := &AnalyzeRunner{
		logger:   logger.Get().WithService("test"),
		s3Client: nil,
	}

	err := runner.Run(context.Background(), nil)

	if err == nil {
		t.Error("Run() expected error but got none")
	}

	if !errors.IsValidationError(err) {
		t.Errorf("Run() expected validation error, got %v", err)
	}
}

func TestAnalyzeRunner_Run_InvalidPayload(t *testing.T) {
	runner := &AnalyzeRunner{
		logger:   logger.Get().WithService("test"),
		s3Client: nil,
	}

	invalidPayload := []byte(`{"invalid": "json"`)

	err := runner.Run(context.Background(), invalidPayload)

	if err == nil {
		t.Error("Run() expected error but got none")
	}

	if !errors.IsValidationError(err) {
		t.Errorf("Run() expected validation error, got %v", err)
	}
}

func TestAnalyzeRunner_Run_MissingAssetID(t *testing.T) {
	runner := &AnalyzeRunner{
		logger:   logger.Get().WithService("test"),
		s3Client: nil,
	}

	payload := messages.AnalyzePayload{
		VideoID: "video-456",
		Input:   "s3://bucket/video.mp4",
	}

	payloadBytes, _ := json.Marshal(payload)

	err := runner.Run(context.Background(), payloadBytes)

	if err == nil {
		t.Error("Run() expected error but got none")
	}

	if !errors.IsValidationError(err) {
		t.Errorf("Run() expected validation error, got %v", err)
	}
}

func TestAnalyzeRunner_Run_MissingVideoID(t *testing.T) {
	runner := &AnalyzeRunner{
		logger:   logger.Get().WithService("test"),
		s3Client: nil,
	}

	payload := messages.AnalyzePayload{
		AssetID: "asset-123",
		Input:   "s3://bucket/video.mp4",
	}

	payloadBytes, _ := json.Marshal(payload)

	err := runner.Run(context.Background(), payloadBytes)

	if err == nil {
		t.Error("Run() expected error but got none")
	}

	if !errors.IsValidationError(err) {
		t.Errorf("Run() expected validation error, got %v", err)
	}
}

func TestAnalyzeRunner_Run_MissingInput(t *testing.T) {
	runner := &AnalyzeRunner{
		logger:   logger.Get().WithService("test"),
		s3Client: nil,
	}

	payload := messages.AnalyzePayload{
		AssetID: "asset-123",
		VideoID: "video-456",
	}

	payloadBytes, _ := json.Marshal(payload)

	err := runner.Run(context.Background(), payloadBytes)

	if err == nil {
		t.Error("Run() expected error but got none")
	}

	if !errors.IsValidationError(err) {
		t.Errorf("Run() expected validation error, got %v", err)
	}
}

func TestAnalyzeRunner_Run_ValidPayload(t *testing.T) {
	runner := &AnalyzeRunner{
		logger:   logger.Get().WithService("test"),
		s3Client: nil,
	}

	payload := messages.AnalyzePayload{
		AssetID: "asset-123",
		VideoID: "video-456",
		Input:   "/non/existent/file.mp4",
	}

	payloadBytes, _ := json.Marshal(payload)

	err := runner.Run(context.Background(), payloadBytes)

	if err == nil {
		t.Error("Run() expected error (file not found) but got none")
	}

	if !errors.IsNotFoundError(err) {
		t.Errorf("Run() expected not found error, got %v", err)
	}
}

func TestAnalyzeRunner_NewAnalyzeRunner(t *testing.T) {
	runner := NewAnalyzeRunner()

	if runner == nil {
		t.Error("NewAnalyzeRunner() expected runner but got nil")
	}

	if runner.logger == nil {
		t.Error("NewAnalyzeRunner() expected logger to be set")
	}

	if runner.s3Client == nil {
		t.Error("NewAnalyzeRunner() expected s3Client to be set")
	}

	if runner.analyzeProducer != nil {
		t.Error("NewAnalyzeRunner() expected analyzeProducer to be nil")
	}
}

func TestAnalyzeRunner_NewAnalyzeRunnerWithAnalyzeProducer(t *testing.T) {
	runner := NewAnalyzeRunnerWithAnalyzeProducer(nil)

	if runner == nil {
		t.Error("NewAnalyzeRunnerWithAnalyzeProducer() expected runner but got nil")
	}

	if runner.logger == nil {
		t.Error("NewAnalyzeRunnerWithAnalyzeProducer() expected logger to be set")
	}

	if runner.s3Client == nil {
		t.Error("NewAnalyzeRunnerWithAnalyzeProducer() expected s3Client to be set")
	}

	if runner.analyzeProducer != nil {
		t.Error("NewAnalyzeRunnerWithAnalyzeProducer() expected analyzeProducer to be nil when passed nil")
	}
}
