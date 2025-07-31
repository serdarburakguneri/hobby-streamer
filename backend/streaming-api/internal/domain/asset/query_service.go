package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
)

type QueryService struct {
	repo Repository
}

func NewQueryService(repo Repository) *QueryService {
	return &QueryService{
		repo: repo,
	}
}

func (s *QueryService) GetBySlug(ctx context.Context, slug valueobjects.Slug) (*entity.Asset, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *QueryService) GetAll(ctx context.Context) ([]*entity.Asset, error) {
	return s.repo.GetAll(ctx)
}

func (s *QueryService) GetPublic(ctx context.Context) ([]*entity.Asset, error) {
	return s.repo.GetPublic(ctx)
}

func (s *QueryService) GetByType(ctx context.Context, assetType valueobjects.AssetType) ([]*entity.Asset, error) {
	return s.repo.GetByType(ctx, assetType)
}

func (s *QueryService) GetByGenre(ctx context.Context, genre valueobjects.Genre) ([]*entity.Asset, error) {
	return s.repo.GetByGenre(ctx, genre)
}
