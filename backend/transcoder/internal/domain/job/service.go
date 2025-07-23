package job

import (
	"context"
	"os"
	"strings"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
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
	if job.AssetID().Value() == "" {
		return pkgerrors.NewValidationError("asset ID is required", nil)
	}

	if job.VideoID().Value() == "" {
		return pkgerrors.NewValidationError("video ID is required", nil)
	}

	if job.Input() == "" {
		return pkgerrors.NewValidationError("input is required", nil)
	}

	if job.Type() == JobTypeTranscode && job.Output() == "" {
		return pkgerrors.NewValidationError("output is required for transcode jobs", nil)
	}

	if job.Type() == JobTypeTranscode && job.Format() == "" {
		return pkgerrors.NewValidationError("format is required for transcode jobs", nil)
	}

	return nil
}

func (s *JobDomainService) ProcessJob(ctx context.Context, job *Job) (interface{}, error) {
	localPath, err := s.downloadFromS3(ctx, job.Input())
	if err != nil {
		return nil, pkgerrors.NewExternalError("failed to download input file", err)
	}
	defer os.Remove(localPath)

	var outputDir string
	if job.Type() == JobTypeTranscode {
		outputDir = "/tmp/transcode/" + job.ID().Value()
		if err := os.MkdirAll(outputDir, 0750); err != nil {
			return nil, pkgerrors.NewInternalError("failed to create output directory", err)
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
		return nil, pkgerrors.NewValidationError("unsupported job type/format: "+strategyKey, nil)
	}

	if err := strategy.ValidateInput(ctx, job); err != nil {
		return nil, err
	}

	outputPath := localPath
	if job.Type() == JobTypeTranscode {
		outputPath, err = strategy.Transcode(ctx, job, localPath, outputDir)
		if err != nil {
			return nil, err
		}
	}
	metadata, metaErr := strategy.ExtractMetadata(ctx, outputPath, job)
	if metaErr != nil {
		return nil, metaErr
	}
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

		retryErr := pkgerrors.RetryWithBackoff(ctx, retryFunc, 3)
		if retryErr != nil {
			s.logger.WithError(retryErr).Error("Failed to download from S3 after retries", "input", input)
			return "", pkgerrors.NewExternalError("failed to download from S3", retryErr)
		}

		return localPath, nil
	}
	return input, nil
}
