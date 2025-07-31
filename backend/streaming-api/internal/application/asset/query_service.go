package asset

import (
	"context"
	"strings"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"

	assetrepo "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	bucketrepo "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
	bucketvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/valueobjects"
)

type Service struct {
	repo       assetrepo.Repository
	bucketRepo bucketrepo.Repository
}

func NewService(repo assetrepo.Repository, bucketRepo bucketrepo.Repository) *Service {
	return &Service{repo: repo, bucketRepo: bucketRepo}
}

func (s *Service) GetAsset(ctx context.Context, slug assetvalueobjects.Slug) (*assetentity.Asset, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *Service) GetAssets(ctx context.Context) ([]*assetentity.Asset, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetPublicAssets(ctx context.Context) ([]*assetentity.Asset, error) {
	return s.repo.GetPublic(ctx)
}

func (s *Service) GetAssetsByType(ctx context.Context, typ assetvalueobjects.AssetType) ([]*assetentity.Asset, error) {
	return s.repo.GetByType(ctx, typ)
}

func (s *Service) GetAssetsByGenre(ctx context.Context, genre assetvalueobjects.Genre) ([]*assetentity.Asset, error) {
	return s.repo.GetByGenre(ctx, genre)
}

func (s *Service) GetAssetsInBucket(ctx context.Context, key bucketvalueobjects.BucketKey) ([]*assetentity.Asset, error) {
	bkt, err := s.bucketRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return s.bucketRepo.GetAssets(ctx, bkt)
}

func (s *Service) SearchAssets(ctx context.Context, query string, filters *SearchFilters) ([]*assetentity.Asset, error) {
	list, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []*assetentity.Asset
	for _, a := range list {
		if query != "" && !strings.Contains(strings.ToLower(a.Title().Value()), strings.ToLower(query)) {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

func (s *Service) GetStreamingInfo(ctx context.Context, slug assetvalueobjects.Slug, userID, region string, userAge int) (*StreamingInfo, error) {
	// TODO: implement streaming rules
	return nil, nil
}

func (s *Service) GetRecommendedAssets(ctx context.Context, slug assetvalueobjects.Slug, limit int) ([]*assetentity.Asset, error) {
	// TODO: implement recommendations
	return nil, nil
}

func (s *Service) GetPublishStatus(ctx context.Context, slug assetvalueobjects.Slug) (constants.PublishStatus, error) {
	// TODO: implement publish status
	return constants.PublishStatusInvalid, nil
}
