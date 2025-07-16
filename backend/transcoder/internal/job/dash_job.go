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

type TranscodeDASHRunner struct {
	logger          *logger.Logger
	analyzeProducer *sqs.Producer
	outputBucket    string
	outputKey       string
	outputFileName  string
	s3Client        *s3.Client
}

func NewTranscodeDASHRunner() *TranscodeDASHRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &TranscodeDASHRunner{
		logger:   logger.WithService("dash-runner"),
		s3Client: s3Client,
	}
}

func NewTranscodeDASHRunnerWithAnalyzeProducer(analyzeProducer *sqs.Producer) *TranscodeDASHRunner {
	s3Client, _ := s3.NewClient(context.Background())
	return &TranscodeDASHRunner{
		logger:          logger.WithService("dash-runner"),
		analyzeProducer: analyzeProducer,
		s3Client:        s3Client,
	}
}

func (d *TranscodeDASHRunner) Run(ctx context.Context, payload json.RawMessage) error {
	log := d.logger.WithContext(ctx)

	var p messages.TranscodePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal DASH payload")
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
		}
		return err
	}

	d.outputBucket = p.OutputBucket
	d.outputKey = p.OutputKey
	d.outputFileName = p.OutputFileName

	var localInputPath string
	var localOutputPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localInputPath, err = d.s3Client.Download(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download input from S3", "input", p.Input)
			if d.analyzeProducer != nil {
				d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
			}
			return err
		}
		defer os.Remove(localInputPath)
	} else {
		localInputPath = p.Input
	}

	tempDir := os.TempDir()
	dashDir := filepath.Join(tempDir, "dash_output")
	if err := os.MkdirAll(dashDir, 0755); err != nil {
		log.WithError(err).Error("Failed to create DASH output directory", "dir", dashDir)
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
		}
		return err
	}
	defer os.RemoveAll(dashDir)

	localOutputPath = filepath.Join(dashDir, d.outputFileName)

	log.Info("Running FFmpeg DASH transcoding", "input", localInputPath, "output", localOutputPath)
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", localInputPath, "-c:v", "libx264", "-c:a", "aac", "-f", "dash", localOutputPath)
	cmd.Dir = dashDir
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg DASH transcoding failed", "input", localInputPath, "output", localOutputPath, "ffmpeg_output", string(out))
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
		}
		return err
	}

	log.Info("DASH transcoding completed successfully", "input", localInputPath, "output", localOutputPath, "output_length", len(out))
	log.Debug("FFmpeg DASH output", "output", string(out))

	// Upload all .mpd and .m4s files from dashDir to S3
	entries, err := os.ReadDir(dashDir)
	if err != nil {
		log.WithError(err).Error("Failed to read DASH output directory", "dir", dashDir)
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, err.Error())
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".mpd") || strings.HasSuffix(name, ".m4s") {
			localPath := filepath.Join(dashDir, name)
			s3Key := filepath.Join(p.AssetID, p.VideoID, name)
			if upErr := d.s3Client.Upload(ctx, localPath, d.outputBucket, s3Key); upErr != nil {
				log.WithError(upErr).Error("Failed to upload file to S3", "bucket", d.outputBucket, "key", s3Key)
				if d.analyzeProducer != nil {
					d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, false, upErr.Error())
				}
				return upErr
			}
		}
	}

	if d.analyzeProducer != nil {
		d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoID, p.Format, true, "")
	}

	return nil
}

func (d *TranscodeDASHRunner) sendTranscodeCompleted(ctx context.Context, assetID, videoID, format string, success bool, errorMessage string) {
	log := d.logger.WithContext(ctx)

	payload := messages.TranscodeCompletionPayload{
		AssetID: assetID,
		VideoID: videoID,
		Format:  format,
		Success: success,
	}

	if success {
		payload.Bucket = d.outputBucket
		payload.Key = d.outputKey
		payload.FileName = d.outputFileName
		payload.URL = "s3://" + d.outputBucket + "/" + d.outputKey
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

	err := d.analyzeProducer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode completed message", "asset_id", assetID, "format", format, "success", success)
	} else {
		log.Info("Transcode completed message sent successfully", "asset_id", assetID, "format", format, "success", success)
	}
}
