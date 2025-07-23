package job

import (
	"context"
	"os"

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
		transcoderRegistry: NewTranscoderRegistry(s3Client),
	}
}

func (s *JobDomainService) ProcessJob(ctx context.Context, jobObj *Job) (interface{}, error) {
	localPath, err := downloadFromS3(ctx, s.s3Client, jobObj.Input())
	if err != nil {
		return nil, pkgerrors.NewExternalError("failed to download input file", err)
	}
	defer os.Remove(localPath)

	outputDir := "/tmp/transcode/" + jobObj.ID().Value()
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		return nil, pkgerrors.NewInternalError("failed to create output directory", err)
	}
	defer os.RemoveAll(outputDir)

	strategyKey := string(jobObj.Format())
	if jobObj.Type() == JobTypeAnalyze {
		strategyKey = "analyze"
	}
	strategy := s.transcoderRegistry.Get(strategyKey)
	if strategy == nil {
		return nil, pkgerrors.NewValidationError("unsupported job type/format: "+strategyKey, nil)
	}

	if err := jobObj.Validate(); err != nil {
		return nil, err
	}

	outputPath, err := strategy.Transcode(ctx, jobObj, localPath, outputDir)
	if err != nil {
		return nil, err
	}

	if err := strategy.ValidateOutput(jobObj); err != nil {
		return nil, err
	}

	metadata, metaErr := strategy.ExtractMetadata(ctx, outputPath, jobObj)
	if metaErr != nil {
		return nil, metaErr
	}
	return metadata, nil
}
