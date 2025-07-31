package job

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job"
)

type JobApplicationService interface {
	ProcessJob(ctx context.Context, payload messages.JobPayload) error
}

type ApplicationService struct {
	domainService job.DomainService
	jobFactory    *JobFactory
	logger        *logger.Logger
}

func NewApplicationService(domainService job.DomainService, cfg config.ServiceConfig) *ApplicationService {
	return &ApplicationService{
		domainService: domainService,
		jobFactory:    NewJobFactory(cfg),
		logger:        logger.WithService("job-application-service"),
	}
}

func (s *ApplicationService) ProcessJob(ctx context.Context, payload messages.JobPayload) error {
	s.logger.Info("Processing job", "job_type", payload.JobType, "asset_id", payload.AssetID, "video_id", payload.VideoID, "input", payload.Input)

	job, err := s.jobFactory.CreateJob(payload)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create job", "job_type", payload.JobType)
		return err
	}

	if err := job.Validate(); err != nil {
		s.logger.WithError(err).Error("Job validation failed", "job_id", job.ID().Value())
		return err
	}

	job.Start()

	metadata, err := s.domainService.ProcessJob(ctx, job)
	if err != nil {
		s.logger.WithError(err).Error("Job processing failed", "job_id", job.ID().Value())
		job.Fail(err.Error())
		return err
	}

	job.Complete(nil)

	s.logger.Info("Job completed successfully", "job_id", job.ID().Value(), "metadata", metadata)
	return nil
}
