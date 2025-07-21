package job

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
)

type JobDomainService struct {
	logger   *logger.Logger
	s3Client *s3.Client
}

func NewJobDomainService() *JobDomainService {
	s3Client, _ := s3.NewClient(context.Background())
	return &JobDomainService{
		logger:   logger.WithService("job-domain-service"),
		s3Client: s3Client,
	}
}

func (s *JobDomainService) ValidateJob(job *Job) error {
	if job.AssetID() == "" {
		return errors.NewValidationError("asset ID is required", nil)
	}

	if job.VideoID() == "" {
		return errors.NewValidationError("video ID is required", nil)
	}

	if job.Input() == "" {
		return errors.NewValidationError("input is required", nil)
	}

	if job.Type() == JobTypeTranscode && job.Output() == "" {
		return errors.NewValidationError("output is required for transcode jobs", nil)
	}

	if job.Type() == JobTypeTranscode && job.Format() == "" {
		return errors.NewValidationError("format is required for transcode jobs", nil)
	}

	return nil
}

func (s *JobDomainService) AnalyzeVideo(ctx context.Context, job *Job) (*VideoMetadata, error) {
	s.logger.Info("Starting video analysis", "job_id", job.ID(), "input", job.Input())

	localPath, err := s.downloadFromS3(ctx, job.Input())
	if err != nil {
		return nil, errors.NewExternalError("failed to download input file", err)
	}
	defer os.Remove(localPath)

	metadata, err := s.extractVideoMetadata(ctx, localPath)
	if err != nil {
		return nil, errors.NewInternalError("failed to extract video metadata", err)
	}

	if err := s.validateVideo(ctx, localPath); err != nil {
		return nil, errors.NewValidationError("video validation failed", err)
	}

	s.logger.Info("Video analysis completed", "job_id", job.ID(), "metadata", metadata)
	return metadata, nil
}

func (s *JobDomainService) TranscodeVideo(ctx context.Context, job *Job) (*TranscodeMetadata, error) {
	s.logger.Info("Starting video transcoding", "job_id", job.ID(), "input", job.Input(), "output", job.Output(), "format", job.Format())

	localPath, err := s.downloadFromS3(ctx, job.Input())
	if err != nil {
		return nil, errors.NewExternalError("failed to download input file", err)
	}
	defer os.Remove(localPath)

	outputDir := "/tmp/transcode/" + job.ID()
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		return nil, errors.NewInternalError("failed to create output directory", err)
	}
	defer os.RemoveAll(outputDir)

	var outputPath string
	var transcodeErr error

	switch job.Format() {
	case JobFormatHLS:
		outputPath, transcodeErr = s.transcodeToHLS(ctx, localPath, outputDir)
	case JobFormatDASH:
		outputPath, transcodeErr = s.transcodeToDASH(ctx, localPath, outputDir)
	default:
		return nil, errors.NewValidationError(fmt.Sprintf("unsupported format: %s", job.Format()), nil)
	}

	if transcodeErr != nil {
		return nil, errors.NewInternalError("transcoding failed", transcodeErr)
	}

	if !strings.HasPrefix(job.Output(), "s3://") {
		return nil, errors.NewValidationError("output must be an S3 path", nil)
	}
	parts := strings.SplitN(job.Output()[5:], "/", 2)
	if len(parts) != 2 {
		return nil, errors.NewValidationError(fmt.Sprintf("invalid S3 path: %s", job.Output()), nil)
	}
	bucket := parts[0]
	manifestKey := parts[1]
	dirKey := filepath.Dir(manifestKey)

	uploadErr := s.s3Client.UploadDirectory(ctx, outputDir, bucket, dirKey)
	if uploadErr != nil {
		return nil, errors.NewExternalError("failed to upload directory to S3", uploadErr)
	}

	outputURL := "s3://" + bucket + "/" + manifestKey

	metadata, metadataErr := s.extractTranscodeMetadata(ctx, outputPath)
	if metadataErr != nil {
		return nil, errors.NewInternalError("failed to extract transcode metadata", metadataErr)
	}

	metadata.OutputURL = outputURL

	if strings.HasPrefix(outputURL, "s3://") {
		parts := strings.SplitN(outputURL[5:], "/", 2)
		if len(parts) == 2 {
			metadata.Bucket = parts[0]
			metadata.Key = parts[1]
		}
	}

	s.logger.Info("Video transcoding completed", "job_id", job.ID(), "output_url", outputURL)
	return metadata, nil
}

func (s *JobDomainService) downloadFromS3(ctx context.Context, input string) (string, error) {
	if strings.HasPrefix(input, "s3://") {
		var localPath string
		var err error

		retryFunc := func(ctx context.Context) error {
			localPath, err = s.s3Client.Download(ctx, input)
			return err
		}

		retryErr := errors.RetryWithBackoff(ctx, retryFunc, 3)
		if retryErr != nil {
			s.logger.WithError(retryErr).Error("Failed to download from S3 after retries", "input", input)
			return "", errors.NewExternalError("failed to download from S3", retryErr)
		}

		return localPath, nil
	}
	return input, nil
}

func (s *JobDomainService) uploadToS3(ctx context.Context, localPath, s3Path string) (string, error) {
	if !strings.HasPrefix(s3Path, "s3://") {
		return s3Path, nil
	}

	parts := strings.SplitN(s3Path[5:], "/", 2)
	if len(parts) != 2 {
		return "", errors.NewValidationError(fmt.Sprintf("invalid S3 path: %s", s3Path), nil)
	}

	bucket := parts[0]
	key := parts[1]

	retryFunc := func(ctx context.Context) error {
		return s.s3Client.Upload(ctx, localPath, bucket, key)
	}

	retryErr := errors.RetryWithBackoff(ctx, retryFunc, 3)
	if retryErr != nil {
		s.logger.WithError(retryErr).Error("Failed to upload to S3 after retries", "local_path", localPath, "bucket", bucket, "key", key)
		return "", errors.NewExternalError("failed to upload to S3", retryErr)
	}

	s.logger.Info("Successfully uploaded to S3", "local_path", localPath, "bucket", bucket, "key", key)
	return s3Path, nil
}

func (s *JobDomainService) extractVideoMetadata(ctx context.Context, filePath string) (*VideoMetadata, error) {
	var out []byte
	var err error

	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffprobe", // nolint:gosec // ffprobe is a trusted binary with controlled arguments
			"-v", "quiet",
			"-print_format", "json",
			"-show_format",
			"-show_streams",
			filePath)

		out, err = cmd.Output()
		return err
	}

	retryErr := errors.RetryWithBackoff(ctx, retryFunc, 2)
	if retryErr != nil {
		return nil, errors.NewInternalError("ffprobe command failed", retryErr)
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
		return nil, errors.NewInternalError("failed to parse ffprobe output", err)
	}

	metadata := &VideoMetadata{}

	for _, stream := range probeResult.Streams {
		if stream.CodecType == "video" {
			metadata.Width = stream.Width
			metadata.Height = stream.Height
			metadata.Codec = stream.CodecName
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

func (s *JobDomainService) validateVideo(ctx context.Context, filePath string) error {
	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffmpeg", "-i", filePath, "-f", "null", "-") // nolint:gosec // ffmpeg is a trusted binary with controlled arguments
		_, err := cmd.CombinedOutput()
		return err
	}

	retryErr := errors.RetryWithBackoff(ctx, retryFunc, 2)
	if retryErr != nil {
		return errors.NewValidationError("video validation failed", retryErr)
	}
	return nil
}

func (s *JobDomainService) transcodeToHLS(ctx context.Context, inputPath, outputDir string) (string, error) {
	outputPath := filepath.Join(outputDir, "playlist.m3u8")

	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffmpeg", // nolint:gosec // ffmpeg is a trusted binary with controlled arguments
			"-i", inputPath,
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

func (s *JobDomainService) transcodeToDASH(ctx context.Context, inputPath, outputDir string) (string, error) {
	outputPath := filepath.Join(outputDir, "playlist.mpd")

	retryFunc := func(ctx context.Context) error {
		cmd := exec.CommandContext(ctx, "ffmpeg", // nolint:gosec // ffmpeg is a trusted binary with controlled arguments
			"-i", inputPath,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-f", "dash",
			"-seg_duration", "10",
			"-use_template", "1",
			"-use_timeline", "1",
			"-init_seg_name", "init-$RepresentationID$.m4s",
			"-media_seg_name", "chunk-$RepresentationID$-$Number%05d$.m4s",
			outputPath)

		return cmd.Run()
	}

	retryErr := errors.RetryWithBackoff(ctx, retryFunc, 2)
	if retryErr != nil {
		return "", errors.NewInternalError("DASH transcoding failed", retryErr)
	}
	return outputPath, nil
}

func (s *JobDomainService) extractTranscodeMetadata(ctx context.Context, filePath string) (*TranscodeMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.NewInternalError("failed to get file info", err)
	}

	metadata := &TranscodeMetadata{
		Size:        fileInfo.Size(),
		ContentType: "application/x-mpegURL",
	}

	if strings.Contains(filePath, "playlist.m3u8") {
		metadata.ContentType = "application/x-mpegURL"
	} else if strings.Contains(filePath, "playlist.mpd") {
		metadata.ContentType = "application/dash+xml"
	}

	return metadata, nil
}
