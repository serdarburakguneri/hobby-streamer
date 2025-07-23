package job

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type AnalyzeTranscoder struct{}

func (a *AnalyzeTranscoder) Transcode(ctx context.Context, job *Job, localPath, outputDir string) (string, error) {
	if err := a.validateVideo(ctx, localPath); err != nil {
		return "", pkgerrors.NewValidationError("video validation failed", err)
	}

	return localPath, nil
}

func (a *AnalyzeTranscoder) ValidateOutput(job *Job) error {
	return nil
}

func (a *AnalyzeTranscoder) ValidateInput(ctx context.Context, job *Job) error {
	return nil
}

func getContentTypeFromFormat(formatName string) string {
	switch formatName {
	case "mov,mp4,m4a,3gp,3g2,mj2":
		return "video/mp4"
	case "matroska,webm":
		return "video/webm"
	case "avi":
		return "video/x-msvideo"
	case "wmv":
		return "video/x-ms-wmv"
	case "flv":
		return "video/x-flv"
	case "m4v":
		return "video/x-m4v"
	default:
		return "video/mp4"
	}
}

func (a *AnalyzeTranscoder) ExtractMetadata(ctx context.Context, filePath string, job *Job) (*TranscodeMetadata, error) {
	var out []byte
	var err error

	if fileInfo, statErr := os.Stat(filePath); statErr != nil {
		return nil, pkgerrors.NewInternalError(fmt.Sprintf("file does not exist or cannot be accessed: %s", filePath), statErr)
	} else {
		logger.Get().Info("File exists and is accessible", "file_path", filePath, "file_size", fileInfo.Size())
	}

	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffprobe",
			"-v", "quiet",
			"-print_format", "json",
			"-show_format",
			"-show_streams",
			filePath)

		out, err = cmd.CombinedOutput()
		if err != nil {
			logger.Get().WithError(err).Error("ffprobe command failed", "file_path", filePath, "output", string(out))
		}
		return err
	}

	retryErr := pkgerrors.RetryWithBackoff(ctx, retryFunc, 2)
	if retryErr != nil {
		return nil, pkgerrors.NewInternalError("ffprobe command failed", retryErr)
	}

	var probeResult struct {
		Format struct {
			Duration   string `json:"duration"`
			BitRate    string `json:"bit_rate"`
			Size       string `json:"size"`
			FormatName string `json:"format_name"`
		} `json:"format"`
		Streams []struct {
			CodecType string `json:"codec_type"`
			CodecName string `json:"codec_name"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(out, &probeResult); err != nil {
		logger.Get().WithError(err).Error("Failed to parse ffprobe output", "output", string(out))
		return nil, pkgerrors.NewInternalError("failed to parse ffprobe output", err)
	}

	metadata := &TranscodeMetadata{}

	for _, stream := range probeResult.Streams {
		if stream.CodecType == "video" {
			metadata.VideoCodec = stream.CodecName
			metadata.Codec = stream.CodecName
			metadata.Width = stream.Width
			metadata.Height = stream.Height
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

	if probeResult.Format.FormatName != "" {
		metadata.ContentType = getContentTypeFromFormat(probeResult.Format.FormatName)
	}

	return metadata, nil
}

func (a *AnalyzeTranscoder) validateVideo(ctx context.Context, filePath string) error {
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
