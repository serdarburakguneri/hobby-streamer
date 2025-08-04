package transcoding

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	resilience "github.com/serdarburakguneri/hobby-streamer/backend/pkg/resilience"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/valueobjects"
)

type HLSTranscoder struct {
	storage job.Storage
}

func NewHLSTranscoder(storage job.Storage) *HLSTranscoder {
	return &HLSTranscoder{storage: storage}
}

func (h *HLSTranscoder) ValidateInput(ctx context.Context, job *entity.Job) error {
	return nil
}

func (h *HLSTranscoder) Transcode(ctx context.Context, job *entity.Job, localPath, outputDir string) (string, error) {
	outputPath := filepath.Join(outputDir, "playlist.m3u8")
	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-i", localPath,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-hls_time", "10",
			"-hls_list_size", "0",
			"-hls_segment_filename", filepath.Join(outputDir, "segment_%03d.ts"),
			outputPath)
		return cmd.Run()
	}
	if err := resilience.RetryWithBackoff(ctx, retryFunc, 2); err != nil {
		return "", pkgerrors.NewInternalError("HLS transcoding failed", err)
	}
	if job.Type().IsTranscode() && strings.HasPrefix(job.Output(), "s3://") {
		if err := h.storage.Upload(ctx, outputDir, job.Output()); err != nil {
			return "", pkgerrors.NewExternalError("failed to upload HLS output to S3", err)
		}
	}
	return outputPath, nil
}

func (h *HLSTranscoder) ValidateOutput(job *entity.Job) error {
	if !strings.HasPrefix(job.Output(), "s3://") {
		return pkgerrors.NewValidationError("output must be an S3 path", nil)
	}
	parts := strings.SplitN(job.Output()[5:], "/", 2)
	if len(parts) != 2 {
		return pkgerrors.NewValidationError("invalid S3 path: "+job.Output(), nil)
	}
	return nil
}

func (h *HLSTranscoder) ExtractMetadata(ctx context.Context, filePath string, job *entity.Job) (*valueobjects.TranscodeMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to stat HLS playlist", err)
	}
	outputURL := job.Output()
	var bucket, key string
	if strings.HasPrefix(outputURL, "s3://") {
		parts := strings.SplitN(outputURL[5:], "/", 2)
		if len(parts) == 2 {
			bucket = parts[0]
			key = parts[1]
		}
	}
	metadata := &valueobjects.TranscodeMetadata{
		OutputURL:   outputURL,
		Bucket:      bucket,
		Key:         key,
		Size:        fileInfo.Size(),
		Format:      valueobjects.JobFormatHLS.String(),
		ContentType: "application/x-mpegURL",
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return metadata, nil
	}
	lines := strings.Split(string(data), "\n")
	var totalDur float64
	var segments []string
	for _, line := range lines {
		if strings.HasPrefix(line, "#EXTINF:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				durStr := strings.TrimSuffix(parts[1], ",")
				if d, err := strconv.ParseFloat(durStr, 64); err == nil {
					totalDur += d
				}
			}
		} else if strings.HasSuffix(line, ".ts") {
			segments = append(segments, line)
		}
	}
	count := len(segments)
	metadata.Segments = segments
	metadata.SegmentCount = count
	metadata.Duration = totalDur
	if count > 0 {
		metadata.AvgSegmentDuration = totalDur / float64(count)
	}
	return metadata, nil
}
