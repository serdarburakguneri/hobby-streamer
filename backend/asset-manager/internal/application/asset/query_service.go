package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/queries"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type QueryService struct {
	finder  asset.Finder
	querier asset.Querier
	logger  *logger.Logger
}

func NewQueryService(
	finder asset.Finder,
	querier asset.Querier,
	logger *logger.Logger,
) *QueryService {
	return &QueryService{
		finder:  finder,
		querier: querier,
		logger:  logger,
	}
}

func (s *QueryService) GetAsset(ctx context.Context, query queries.GetAssetQuery) (*entity.Asset, error) {
	assetID, err := valueobjects.NewAssetID(query.ID)
	if err != nil {
		return nil, err
	}
	return s.finder.FindByID(ctx, *assetID)
}

func (s *QueryService) ListAssets(ctx context.Context, query queries.ListAssetsQuery) ([]*entity.Asset, error) {
	return s.querier.List(ctx, query.Limit, query.Offset)
}

func (s *QueryService) SearchAssets(ctx context.Context, query queries.SearchAssetsQuery) ([]*entity.Asset, error) {
	return s.querier.Search(ctx, query.Query, query.Limit, query.Offset)
}
