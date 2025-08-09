package pipeline

import (
	"context"

	domain "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/pipeline/entity"
)

type Repository interface {
	Upsert(ctx context.Context, p *domain.Pipeline) error
	Get(ctx context.Context, assetID, videoID string) (*domain.Pipeline, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) MarkRequested(ctx context.Context, assetID, videoID, step, jobID, correlationID string) error {
	p, _ := s.repo.Get(ctx, assetID, videoID)
	if p == nil {
		p = domain.NewPipeline(assetID, videoID)
	}
	p.SetRequested(step, jobID, correlationID)
	return s.repo.Upsert(ctx, p)
}

func (s *Service) MarkCompleted(ctx context.Context, assetID, videoID, step string) error {
	p, _ := s.repo.Get(ctx, assetID, videoID)
	if p == nil {
		p = domain.NewPipeline(assetID, videoID)
	}
	p.SetCompleted(step)
	return s.repo.Upsert(ctx, p)
}

func (s *Service) MarkFailed(ctx context.Context, assetID, videoID, step, errMsg string) error {
	p, _ := s.repo.Get(ctx, assetID, videoID)
	if p == nil {
		p = domain.NewPipeline(assetID, videoID)
	}
	p.SetFailed(step, errMsg)
	return s.repo.Upsert(ctx, p)
}

func (s *Service) Get(ctx context.Context, assetID, videoID string) (*domain.Pipeline, error) {
	return s.repo.Get(ctx, assetID, videoID)
}
