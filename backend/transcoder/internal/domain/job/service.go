package job

import (
	"context"
	"fmt"
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

func buildOutputDir(job *Job) string {
	return fmt.Sprintf("/tmp/%s/%s/%s", job.AssetID().Value(), job.Format(), job.Quality())
}

func (s *JobDomainService) ProcessJob(ctx context.Context, jobObj *Job) (interface{}, error) {
	localPath, err := downloadFromS3(ctx, s.s3Client, jobObj.Input())
	if err != nil {
		return nil, pkgerrors.NewExternalError("failed to download input file", err)
	}

	outputDir := buildOutputDir(jobObj)
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		os.Remove(localPath)
		return nil, pkgerrors.NewInternalError("failed to create output directory", err)
	}
	defer os.RemoveAll(outputDir)

	strategyKey := string(jobObj.Format())
	if jobObj.Type() == JobTypeAnalyze {
		strategyKey = "analyze"
	}
	strategy := s.transcoderRegistry.Get(strategyKey)
	if strategy == nil {
		os.Remove(localPath)
		return nil, pkgerrors.NewValidationError("unsupported job type/format: "+strategyKey, nil)
	}

	if err := jobObj.Validate(); err != nil {
		os.Remove(localPath)
		return nil, err
	}

	outputPath, err := strategy.Transcode(ctx, jobObj, localPath, outputDir)
	if err != nil {
		os.Remove(localPath)
		return nil, pkgerrors.NewInternalError("failed to transcode", err)
	}

	if err := strategy.ValidateOutput(jobObj); err != nil {
		os.Remove(localPath)
		return nil, err
	}

	metadata, metaErr := strategy.ExtractMetadata(ctx, outputPath, jobObj)
	if metaErr != nil {
		os.Remove(localPath)
		return nil, pkgerrors.NewInternalError("failed to extract metadata", metaErr)
	}

	os.Remove(localPath)
	return metadata, nil
}
