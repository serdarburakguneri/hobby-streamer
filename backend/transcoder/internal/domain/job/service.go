package job

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3"
)

type JobDomainService struct {
	logger             *logger.Logger
	s3Client           *s3.Client
	transcoderRegistry *TranscoderRegistry
}

func NewJobDomainService() *JobDomainService {
	s3Client, _ := s3.NewClient(context.Background())
	return &JobDomainService{
		logger:             logger.WithService("job-domain-service"),
		s3Client:           s3Client,
		transcoderRegistry: NewTranscoderRegistry(),
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

func (s *JobDomainService) ProcessJob(ctx context.Context, job *Job) (interface{}, error) {
	localPath, err := s.downloadFromS3(ctx, job.Input())
	if err != nil {
		return nil, errors.NewExternalError("failed to download input file", err)
	}
	defer os.Remove(localPath)

	var outputDir string
	if job.Type() == JobTypeTranscode {
		outputDir = "/tmp/transcode/" + job.ID()
		if err := os.MkdirAll(outputDir, 0750); err != nil {
			return nil, errors.NewInternalError("failed to create output directory", err)
		}
		defer os.RemoveAll(outputDir)
	}

	var strategyKey string
	if job.Type() == JobTypeAnalyze {
		strategyKey = "analyze"
	} else {
		strategyKey = string(job.Format())
	}
	strategy := s.transcoderRegistry.Get(strategyKey)
	if strategy == nil {
		return nil, errors.NewValidationError("unsupported job type/format: "+strategyKey, nil)
	}
	outputPath, stratErr := strategy.Transcode(ctx, job, localPath, outputDir)
	if stratErr != nil {
		return nil, stratErr
	}

	if job.Type() == JobTypeTranscode {
		if !strings.HasPrefix(job.Output(), "s3://") {
			return nil, errors.NewValidationError("output must be an S3 path", nil)
		}
		parts := strings.SplitN(job.Output()[5:], "/", 2)
		if len(parts) != 2 {
			return nil, errors.NewValidationError("invalid S3 path: "+job.Output(), nil)
		}
		bucket := parts[0]
		manifestKey := parts[1]
		dirKey := filepath.Dir(manifestKey)

		uploadErr := s.s3Client.UploadDirectory(ctx, outputDir, bucket, dirKey)
		if uploadErr != nil {
			return nil, errors.NewExternalError("failed to upload directory to S3", uploadErr)
		}

		outputURL := "s3://" + bucket + "/" + manifestKey
		metadata, metadataErr := strategy.ExtractMetadata(ctx, outputPath)
		if metadataErr != nil {
			return nil, errors.NewInternalError("failed to extract transcode metadata", metadataErr)
		}
		if metadata != nil {
			metadata.OutputURL = outputURL
			metadata.Bucket = bucket
			metadata.Key = manifestKey
			metadata.Format = string(job.Format())
		}
		return metadata, nil
	}

	return nil, nil
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
