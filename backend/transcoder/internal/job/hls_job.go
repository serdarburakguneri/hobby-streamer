package job

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

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

	var p messages.TranscodePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal HLS payload")
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
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
				h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
			}
			return err
		}
		defer os.Remove(localInputPath)
	} else {
		localInputPath = p.Input
	}

	tempDir := os.TempDir()
	hlsDir := filepath.Join(tempDir, "hls_output")
	if err := os.MkdirAll(hlsDir, 0755); err != nil {
		log.WithError(err).Error("Failed to create HLS output directory", "dir", hlsDir)
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
		}
		return err
	}
	defer os.RemoveAll(hlsDir)

	baseFileName := strings.TrimSuffix(h.outputFileName, ".m3u8")
	segmentDir := filepath.Join(hlsDir, baseFileName)
	if err := os.MkdirAll(segmentDir, 0755); err != nil {
		log.WithError(err).Error("Failed to create HLS segment directory", "dir", segmentDir)
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
		}
		return err
	}

	// Write both manifest and segments to the same directory
	segmentPattern := filepath.Join(hlsDir, baseFileName+"%d.ts")
	localOutputPath = filepath.Join(hlsDir, h.outputFileName)

	log.Info("Running FFmpeg HLS transcoding", "input", localInputPath, "output", localOutputPath, "segment_pattern", segmentPattern)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localInputPath, "-c:v", "libx264", "-c:a", "aac", "-f", "hls", "-hls_segment_filename", segmentPattern, localOutputPath)
	cmd.Dir = hlsDir
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg HLS transcoding failed", "input", localInputPath, "output", localOutputPath, "ffmpeg_output", string(out))
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
		}
		return err
	}

	log.Info("HLS transcoding completed successfully", "input", localInputPath, "output", localOutputPath, "output_length", len(out))
	log.Debug("FFmpeg HLS output", "output", string(out))

	// Upload all .m3u8 and .ts files from hlsDir to S3
	entries, err := os.ReadDir(hlsDir)
	if err != nil {
		log.WithError(err).Error("Failed to read HLS output directory", "dir", hlsDir)
		if h.analyzeProducer != nil {
			h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".m3u8") || strings.HasSuffix(name, ".ts") {
			localPath := filepath.Join(hlsDir, name)
			s3Key := filepath.Join(p.AssetID, p.VideoID, name)
			if upErr := h.s3Client.Upload(ctx, localPath, h.outputBucket, s3Key); upErr != nil {
				log.WithError(upErr).Error("Failed to upload file to S3", "bucket", h.outputBucket, "key", s3Key)
				if h.analyzeProducer != nil {
					h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, upErr.Error())
				}
				return upErr
			}
		}
	}

	if h.analyzeProducer != nil {
		h.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, true, "")
	}

	return nil
}

func (h *TranscodeHLSRunner) sendTranscodeCompleted(ctx context.Context, assetID, videoID, format string, success bool, errorMessage string) {
	log := h.logger.WithContext(ctx)

	payload := messages.TranscodeCompletionPayload{
		AssetID: assetID,
		VideoID: videoID,
		Format:  format,
		Success: success,
	}

	if success {
		payload.Bucket = h.outputBucket
		payload.Key = h.outputKey
		payload.FileName = h.outputFileName
		payload.URL = "s3://" + h.outputBucket + "/" + h.outputKey
	}

	if !success && errorMessage != "" {
		payload.Error = errorMessage
	}

	var messageType string
	switch format {
	case "hls":
		messageType = messages.MessageTypeTranscodeHLSCompleted
	case "dash":
		messageType = messages.MessageTypeTranscodeDASHCompleted
	default:
		log.Error("Unknown format for transcode completion", "format", format)
		return
	}

	err := h.analyzeProducer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode completed message", "asset_id", assetID, "format", format, "success", success)
	} else {
		log.Info("Transcode completed message sent successfully", "asset_id", assetID, "format", format, "success", success)
	}
}
