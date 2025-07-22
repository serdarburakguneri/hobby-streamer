package job

import (
	"context"
	"encoding/json"
	"os/exec"
	"strconv"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type AnalyzeStrategy struct{}

func (a *AnalyzeStrategy) Transcode(ctx context.Context, job *Job, localPath, outputDir string) (string, error) {
	metadata, err := a.ExtractMetadata(ctx, localPath)
	if err != nil {
		return "", pkgerrors.NewInternalError("failed to extract video metadata", err)
	}
	if err := a.validateVideo(ctx, localPath); err != nil {
		return "", pkgerrors.NewValidationError("video validation failed", err)
	}
	_ = metadata
	return "", nil
}

func (a *AnalyzeStrategy) ExtractMetadata(ctx context.Context, filePath string) (*TranscodeMetadata, error) {
	var out []byte
	var err error

	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffprobe",
			"-v", "quiet",
			"-print_format", "json",
			"-show_format",
			"-show_streams",
			filePath)

		out, err = cmd.Output()
		return err
	}

	retryErr := pkgerrors.RetryWithBackoff(ctx, retryFunc, 2)
	if retryErr != nil {
		return nil, pkgerrors.NewInternalError("ffprobe command failed", retryErr)
	}

	var probeResult struct {
		Format struct {
			Duration string `json:"duration"`
			BitRate  string `json:"bit_rate"`
			Size     string `json:"size"`
		} `json:"format"`
		Streams []struct {
			CodecType string `json:"codec_type"`
			CodecName string `json:"codec_name"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(out, &probeResult); err != nil {
		return nil, pkgerrors.NewInternalError("failed to parse ffprobe output", err)
	}

	metadata := &TranscodeMetadata{}

	for _, stream := range probeResult.Streams {
		if stream.CodecType == "video" {
			metadata.VideoCodec = stream.CodecName
			break
		}
	}

	if probeResult.Format.Duration != "" {
		if duration, err := strconv.ParseFloat(probeResult.Format.Duration, 64); err == nil {
			metadata.Duration = duration
		}
	}

	if probeResult.Format.BitRate != "" {
		if bitrate, err := strconv.Atoi(probeResult.Format.BitRate); err == nil {
			metadata.Bitrate = bitrate
		}
	}

	if probeResult.Format.Size != "" {
		if size, err := strconv.ParseInt(probeResult.Format.Size, 10, 64); err == nil {
			metadata.Size = size
		}
	}

	return metadata, nil
}

func (a *AnalyzeStrategy) validateVideo(ctx context.Context, filePath string) error {
	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffmpeg", "-i", filePath, "-f", "null", "-")
		_, err := cmd.CombinedOutput()
		return err
	}

	retryErr := pkgerrors.RetryWithBackoff(ctx, retryFunc, 2)
	if retryErr != nil {
		return pkgerrors.NewValidationError("video validation failed", retryErr)
	}
	return nil
}
