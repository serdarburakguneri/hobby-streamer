package asset

import (
	"context"
	"strconv"

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

// ListAssetsPage handles pagination logic for listing assets.
func (s *QueryService) ListAssetsPage(ctx context.Context, query queries.ListAssetsQuery) (*entity.AssetPage, error) {
	items, err := s.querier.List(ctx, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	// determine limit and offset values
	limitVal := 0
	if query.Limit != nil {
		limitVal = *query.Limit
	} else {
		limitVal = len(items)
	}
	offsetVal := 0
	if query.Offset != nil {
		offsetVal = *query.Offset
	}

	// compute hasMore and lastKey
	hasMore := len(items) >= limitVal
	lastKey := make(map[string]interface{})
	if hasMore {
		lastKey["key"] = strconv.Itoa(offsetVal + len(items))
	}

	return &entity.AssetPage{Items: items, HasMore: hasMore, LastKey: lastKey}, nil
}

// SearchAssetsPage handles pagination logic for searching assets.
func (s *QueryService) SearchAssetsPage(ctx context.Context, query queries.SearchAssetsQuery) (*entity.AssetPage, error) {
	items, err := s.querier.Search(ctx, query.Query, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	// determine limit and offset values
	limitVal := 0
	if query.Limit != nil {
		limitVal = *query.Limit
	} else {
		limitVal = len(items)
	}
	offsetVal := 0
	if query.Offset != nil {
		offsetVal = *query.Offset
	}

	// compute hasMore and lastKey
	hasMore := len(items) >= limitVal
	lastKey := make(map[string]interface{})
	if hasMore {
		lastKey["key"] = strconv.Itoa(offsetVal + len(items))
	}

	return &entity.AssetPage{Items: items, HasMore: hasMore, LastKey: lastKey}, nil
}
