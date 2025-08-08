package job

import (
	"context"
	"fmt"
	"strings"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/events"
)

type DomainServiceImpl struct {
	storage            Storage
	transcoderRegistry TranscoderRegistry
	eventPublisher     EventPublisher
}

func NewDomainService(storage Storage, transcoderRegistry TranscoderRegistry, eventPublisher EventPublisher) *DomainServiceImpl {
	return &DomainServiceImpl{storage: storage, transcoderRegistry: transcoderRegistry, eventPublisher: eventPublisher}
}

func buildOutputDir(job *entity.Job) string {
	return fmt.Sprintf("/tmp/%s/%s/%s", job.AssetID().Value(), job.Format(), job.Quality())
}

func (s *DomainServiceImpl) ProcessJob(ctx context.Context, jobObj *entity.Job) (interface{}, error) {
	inputIsS3 := strings.HasPrefix(jobObj.Input(), "s3://")
	localPath, err := s.storage.Download(ctx, jobObj.Input())
	if err != nil {
		s.publishJobCompletion(ctx, jobObj, false, nil, err.Error())
		return nil, pkgerrors.NewExternalError("failed to download input file", err)
	}

	outputDir := buildOutputDir(jobObj)
	if err := s.storage.CreateDir(outputDir); err != nil {
		s.publishJobCompletion(ctx, jobObj, false, nil, err.Error())
		return nil, pkgerrors.NewInternalError("failed to create output directory", err)
	}
	defer s.storage.RemoveAll(outputDir)

	strategyKey := string(jobObj.Format())
	if jobObj.Type().IsAnalyze() {
		strategyKey = "analyze"
	}
	strategy := s.transcoderRegistry.Get(strategyKey)
	if strategy == nil {
		if inputIsS3 {
			s.storage.Remove(localPath)
		}
		s.publishJobCompletion(ctx, jobObj, false, nil, "unsupported job type/format: "+strategyKey)
		return nil, pkgerrors.NewValidationError("unsupported job type/format: "+strategyKey, nil)
	}

	if err := jobObj.Validate(); err != nil {
		if inputIsS3 {
			s.storage.Remove(localPath)
		}
		s.publishJobCompletion(ctx, jobObj, false, nil, err.Error())
		return nil, err
	}

	outputPath, err := strategy.Transcode(ctx, jobObj, localPath, outputDir)
	if err != nil {
		if inputIsS3 {
			s.storage.Remove(localPath)
		}
		s.publishJobCompletion(ctx, jobObj, false, nil, err.Error())
		return nil, pkgerrors.NewInternalError("failed to transcode", err)
	}

	if err := strategy.ValidateOutput(jobObj); err != nil {
		if inputIsS3 {
			s.storage.Remove(localPath)
		}
		s.publishJobCompletion(ctx, jobObj, false, nil, err.Error())
		return nil, err
	}

	metadata, metaErr := strategy.ExtractMetadata(ctx, outputPath, jobObj)
	if metaErr != nil {
		if inputIsS3 {
			s.storage.Remove(localPath)
		}
		s.publishJobCompletion(ctx, jobObj, false, nil, metaErr.Error())
		return nil, pkgerrors.NewInternalError("failed to extract metadata", metaErr)
	}

	if inputIsS3 {
		s.storage.Remove(localPath)
	}

	s.publishJobCompletion(ctx, jobObj, true, metadata, "")

	return metadata, nil
}

func (s *DomainServiceImpl) publishJobCompletion(ctx context.Context, jobObj *entity.Job, success bool, metadata interface{}, errorMessage string) {
	completionEvent := events.BuildCompletedEvent(jobObj, success, metadata, errorMessage)
	s.eventPublisher.PublishJobCompleted(ctx, completionEvent)
}
