package job

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type HLSPayload struct {
	Input          string `json:"input"`
	AssetID        string `json:"assetId"`
	VideoType      string `json:"videoType"`
	Format         string `json:"format"`
	OutputBucket   string `json:"outputBucket"`
	OutputKey      string `json:"outputKey"`
	OutputFileName string `json:"outputFileName"`
}

type TranscodeHLSRunner struct {
	logger          *logger.Logger
	analyzeProducer *sqs.Producer
	outputBucket    string
	outputKey       string
	outputFileName  string
	s3Client        *s3.Client
}

func NewTranscodeHLSRunner() *TranscodeHLSRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &TranscodeHLSRunner{
		logger:   logger.WithService("hls-runner"),
		s3Client: s3Client,
	}
}

func NewTranscodeHLSRunnerWithAnalyzeProducer(analyzeProducer *sqs.Producer) *TranscodeHLSRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &TranscodeHLSRunner{
		logger:          logger.WithService("hls-runner"),
		analyzeProducer: analyzeProducer,
		s3Client:        s3Client,
	}
}

func (h *TranscodeHLSRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := h.logger.WithContext(ctx)

	var p HLSPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal HLS payload")
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	h.outputBucket = p.OutputBucket
	h.outputKey = p.OutputKey
	h.outputFileName = p.OutputFileName

	var localInputPath string
	var localOutputPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localInputPath, err = h.s3Client.Download(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download input from S3", "input", p.Input)
			if h.analyzeProducer != nil {
				h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
			}
			return err
		}
		defer os.Remove(localInputPath)
	} else {
		localInputPath = p.Input
	}

	tempDir := os.TempDir()
	localOutputPath = filepath.Join(tempDir, h.outputFileName)

	log.Info("Running FFmpeg HLS transcoding", "input", localInputPath, "output", localOutputPath)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localInputPath, "-c:v", "libx264", "-c:a", "aac", "-f", "hls", localOutputPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg HLS transcoding failed", "input", localInputPath, "output", localOutputPath, "ffmpeg_output", string(out))
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	log.Info("HLS transcoding completed successfully", "input", localInputPath, "output", localOutputPath, "output_length", len(out))
	log.Debug("FFmpeg HLS output", "output", string(out))

	err = h.s3Client.Upload(ctx, localOutputPath, h.outputBucket, h.outputKey)
	if err != nil {
		log.WithError(err).Error("Failed to upload output to S3", "bucket", h.outputBucket, "key", h.outputKey)
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}
	defer os.Remove(localOutputPath)

	if h.analyzeProducer != nil {
		h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, true, "")
	}

	return nil
}

func (h *TranscodeHLSRunner) sendTranscodeCompleted(ctx context.Context, assetID, videoType, format string, success bool, errorMessage string) {
	log := h.logger.WithContext(ctx)

	payload := map[string]interface{}{
		"assetId":   assetID,
		"videoType": videoType,
		"format":    format,
		"success":   success,
	}

	if success {
		payload["bucket"] = h.outputBucket
		payload["key"] = h.outputKey
		payload["fileName"] = h.outputFileName
		payload["url"] = "s3://" + h.outputBucket + "/" + h.outputKey
	}

	if !success && errorMessage != "" {
		payload["error"] = errorMessage
	}

	messageType := "transcode-" + format + "-completed"
	err := h.analyzeProducer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode completed message", "asset_id", assetID, "format", format, "success", success)
	} else {
		log.Info("Transcode completed message sent successfully", "asset_id", assetID, "format", format, "success", success)
	}
}
