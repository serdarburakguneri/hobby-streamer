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

type DASHPayload struct {
	Input          string `json:"input"`
	AssetID        string `json:"assetId"`
	VideoType      string `json:"videoType"`
	Format         string `json:"format"`
	OutputBucket   string `json:"outputBucket"`
	OutputKey      string `json:"outputKey"`
	OutputFileName string `json:"outputFileName"`
}

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

	var p DASHPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.WithError(err).Error("Failed to unmarshal DASH payload")
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}
	d.outputBucket = p.OutputBucket
	d.outputKey = p.OutputKey
	d.outputFileName = p.OutputFileName

	var localInputPath string
	var err error

	if strings.HasPrefix(p.Input, "s3://") {
		localInputPath, err = d.s3Client.Download(ctx, p.Input)
		if err != nil {
			log.WithError(err).Error("Failed to download input from S3", "input", p.Input)
			if d.analyzeProducer != nil {
				d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
			}
			return err
		}
		defer os.Remove(localInputPath)
	} else {
		localInputPath = p.Input
	}

	tempDir := os.TempDir()
	outputDir := filepath.Join(tempDir, "dash_output")
	os.MkdirAll(outputDir, 0755)

	manifestPath := filepath.Join(outputDir, d.outputFileName)
	segmentPattern := filepath.Join(outputDir, "segment_%03d.m4s")

	log.Info("Running FFmpeg DASH transcoding", "input", localInputPath, "manifest", manifestPath, "segment_pattern", segmentPattern)
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", localInputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "dash",
		"-init_seg_name", "init_$RepresentationID$.m4s",
		"-media_seg_name", "segment_$RepresentationID$_$Number%03d$.m4s",
		manifestPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("FFmpeg DASH transcoding failed", "input", localInputPath, "manifest", manifestPath, "ffmpeg_output", string(out))
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	log.Info("DASH transcoding completed successfully", "input", localInputPath, "manifest", manifestPath, "output_length", len(out))
	log.Debug("FFmpeg DASH output", "output", string(out))

	err = d.s3Client.Upload(ctx, manifestPath, d.outputBucket, d.outputKey)
	if err != nil {
		log.WithError(err).Error("Failed to upload DASH manifest to S3", "bucket", d.outputBucket, "key", d.outputKey)
		if d.analyzeProducer != nil {
			d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, false, err.Error())
		}
		return err
	}

	segmentFiles, err := filepath.Glob(filepath.Join(outputDir, "*.m4s"))
	if err != nil {
		log.WithError(err).Error("Failed to find DASH segment files")
	} else {
		for _, segmentFile := range segmentFiles {
			segmentName := filepath.Base(segmentFile)
			segmentKey := strings.TrimSuffix(d.outputKey, ".mpd") + "/" + segmentName
			err = d.s3Client.Upload(ctx, segmentFile, d.outputBucket, segmentKey)
			if err != nil {
				log.WithError(err).Error("Failed to upload DASH segment to S3", "bucket", d.outputBucket, "key", segmentKey, "file", segmentFile)
			} else {
				log.Info("Successfully uploaded DASH segment", "bucket", d.outputBucket, "key", segmentKey)
			}
		}
	}
	
	defer os.RemoveAll(outputDir)

	if d.analyzeProducer != nil {
		d.sendTranscodeCompleted(ctx, p.AssetID, p.VideoType, p.Format, true, "")
	}

	return nil
}

func (d *TranscodeDASHRunner) sendTranscodeCompleted(ctx context.Context, assetID, videoType, format string, success bool, errorMessage string) {
	log := d.logger.WithContext(ctx)

	payload := map[string]interface{}{
		"assetId":   assetID,
		"videoType": videoType,
		"format":    format,
		"success":   success,
	}

	if success {
		payload["bucket"] = d.outputBucket
		payload["key"] = d.outputKey
		payload["fileName"] = d.outputFileName
		payload["url"] = "s3://" + d.outputBucket + "/" + d.outputKey
	}

	if !success && errorMessage != "" {
		payload["error"] = errorMessage
	}

	messageType := "transcode-" + format + "-completed"
	err := d.analyzeProducer.SendMessage(ctx, messageType, payload)
	if err != nil {
		log.WithError(err).Error("Failed to send transcode completed message", "asset_id", assetID, "format", format, "success", success)
	} else {
		log.Info("Transcode completed message sent successfully", "asset_id", assetID, "format", format, "success", success)
	}
}
