package job

import (
	"context"
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	domainjob "github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job"
)

type JobApplicationService interface {
	ProcessJob(ctx context.Context, payload messages.JobPayload) error
}

type ApplicationService struct {
	domainService  *domainjob.JobDomainService
	eventPublisher EventPublisher
	logger         *logger.Logger
}

type EventPublisher interface {
	PublishJobCompleted(ctx context.Context, jobType, assetID, videoID string, success bool, metadata interface{}, errorMessage string) error
}

func NewApplicationService(domainService *domainjob.JobDomainService, eventPublisher EventPublisher) *ApplicationService {
	return &ApplicationService{
		domainService:  domainService,
		eventPublisher: eventPublisher,
		logger:         logger.WithService("job-application-service"),
	}
}

func (s *ApplicationService) ProcessJob(ctx context.Context, payload messages.JobPayload) error {
	s.logger.Info("Processing job", "job_type", payload.JobType, "asset_id", payload.AssetID, "video_id", payload.VideoID, "input", payload.Input)

	var job *domainjob.Job
	var format domainjob.JobFormat
	var outputKey string
	assetIDVO, err := domainjob.NewAssetID(payload.AssetID)
	if err != nil {
		return errors.NewValidationError("invalid asset ID", err)
	}
	videoIDVO, err := domainjob.NewVideoID(payload.VideoID)
	if err != nil {
		return errors.NewValidationError("invalid video ID", err)
	}
	switch payload.JobType {
	case string(domainjob.JobTypeAnalyze):
		job = domainjob.NewAnalyzeJob(*assetIDVO, *videoIDVO, payload.Input)
	case string(domainjob.JobTypeTranscode):
		switch payload.Format {
		case string(domainjob.JobFormatHLS):
			format = domainjob.JobFormatHLS
			outputKey = fmt.Sprintf("%s/hls/%s/playlist.m3u8", payload.AssetID, payload.Quality)
		case string(domainjob.JobFormatDASH):
			format = domainjob.JobFormatDASH
			outputKey = fmt.Sprintf("%s/dash/%s/manifest.mpd", payload.AssetID, payload.Quality)
		default:
			errMsg := fmt.Sprintf("unsupported format: %s", payload.Format)
			s.logger.Error("Unsupported format", "format", payload.Format)
			if pubErr := s.eventPublisher.PublishJobCompleted(ctx, payload.JobType, payload.AssetID, payload.VideoID, false, nil, errMsg); pubErr != nil {
				s.logger.WithError(pubErr).Error("Failed to publish job completed event", "job_type", payload.JobType)
			}
			return errors.NewValidationError(errMsg, nil)
		}
		outputPath := fmt.Sprintf("s3://%s/%s", payload.OutputBucket, outputKey)
		job = domainjob.NewTranscodeJob(*assetIDVO, *videoIDVO, payload.Input, outputPath, payload.Quality, format)
	default:
		errMsg := fmt.Sprintf("unsupported job type: %s", payload.JobType)
		s.logger.Error("Unsupported job type", "job_type", payload.JobType)
		if err := s.eventPublisher.PublishJobCompleted(ctx, payload.JobType, payload.AssetID, payload.VideoID, false, nil, errMsg); err != nil {
			s.logger.WithError(err).Error("Failed to publish job completed event", "job_type", payload.JobType)
		}
		return errors.NewValidationError(errMsg, nil)
	}

	if err := job.Validate(); err != nil {
		s.logger.WithError(err).Error("Job validation failed", "job_id", job.ID().Value())
		if pubErr := s.eventPublisher.PublishJobCompleted(ctx, payload.JobType, job.AssetID().Value(), job.VideoID().Value(), false, nil, err.Error()); pubErr != nil {
			s.logger.WithError(pubErr).Error("Failed to publish job completed event", "job_type", payload.JobType)
		}
		return err
	}

	job.Start()

	metadata, err := s.domainService.ProcessJob(ctx, job)
	if err != nil {
		s.logger.WithError(err).Error("Job processing failed", "job_id", job.ID().Value())
		job.Fail(err.Error())
		if pubErr := s.eventPublisher.PublishJobCompleted(ctx, payload.JobType, job.AssetID().Value(), job.VideoID().Value(), false, nil, err.Error()); pubErr != nil {
			s.logger.WithError(pubErr).Error("Failed to publish job completed event", "job_type", payload.JobType)
		}
		return err
	}

	job.Complete(nil)

	s.logger.Info("Job completed successfully", "job_id", job.ID().Value(), "metadata", metadata)
	if pubErr := s.eventPublisher.PublishJobCompleted(ctx, payload.JobType, job.AssetID().Value(), job.VideoID().Value(), true, metadata, ""); pubErr != nil {
		s.logger.WithError(pubErr).Error("Failed to publish job completed event", "job_type", payload.JobType)
	}
	return nil
}
