package job

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type HLSTranscoder struct{}

func (h *HLSTranscoder) Transcode(ctx context.Context, job *Job, localPath, outputDir string) (string, error) {
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
	retryErr := errors.RetryWithBackoff(ctx, retryFunc, 2)
	if retryErr != nil {
		return "", errors.NewInternalError("HLS transcoding failed", retryErr)
	}
	return outputPath, nil
}

func (h *HLSTranscoder) ValidateOutput(job *Job) error {
	if !strings.HasPrefix(job.Output(), "s3://") {
		return errors.NewValidationError("output must be an S3 path", nil)
	}
	parts := strings.SplitN(job.Output()[5:], "/", 2)
	if len(parts) != 2 {
		return errors.NewValidationError("invalid S3 path: "+job.Output(), nil)
	}
	return nil
}

func (h *HLSTranscoder) ValidateInput(ctx context.Context, job *Job) error {
	return nil
}

func (h *HLSTranscoder) ExtractMetadata(ctx context.Context, filePath string, job *Job) (*TranscodeMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.NewInternalError("failed to get file info", err)
	}
	metadata := &TranscodeMetadata{
		Size:        fileInfo.Size(),
		Format:      string(JobFormatHLS),
		ContentType: "application/x-mpegURL",
	}
	playlist, err := ioutil.ReadFile(filePath)
	if err == nil {
		lines := strings.Split(string(playlist), "\n")
		var segments []string
		var totalDuration float64
		var segmentCount int
		for i, line := range lines {
			if strings.HasPrefix(line, "#EXTINF:") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					durStr := strings.TrimSuffix(parts[1], ",")
					if dur, err := strconv.ParseFloat(durStr, 64); err == nil {
						totalDuration += dur
					}
				}
				if i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "#") && lines[i+1] != "" {
					segments = append(segments, lines[i+1])
					segmentCount++
				}
			}
			if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
				if strings.Contains(line, "CODECS=") {
					codecStr := line[strings.Index(line, "CODECS=")+7:]
					codecStr = strings.Trim(codecStr, "\"")
					codecs := strings.Split(codecStr, ",")
					if len(codecs) > 0 {
						metadata.VideoCodec = codecs[0]
					}
					if len(codecs) > 1 {
						metadata.AudioCodec = codecs[1]
					}
				}
			}
		}
		metadata.SegmentCount = segmentCount
		metadata.Segments = segments
		if segmentCount > 0 {
			metadata.AvgSegmentDuration = totalDuration / float64(segmentCount)
		}
		metadata.Duration = totalDuration
	}

	if job != nil && strings.HasPrefix(job.Output(), "s3://") {
		parts := strings.SplitN(job.Output()[5:], "/", 2)
		if len(parts) == 2 {
			metadata.OutputURL = "s3://" + parts[0] + "/" + parts[1]
			metadata.Bucket = parts[0]
			metadata.Key = parts[1]
			metadata.Format = string(JobFormatHLS)
		}
	}
	return metadata, nil
}
