package bucket

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket/queries"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type QueryService struct {
	finder   bucket.Finder
	pager    bucket.Pager
	relation bucket.Relation
	logger   *logger.Logger
}

func NewQueryService(
	finder bucket.Finder,
	pager bucket.Pager,
	relation bucket.Relation,
	logger *logger.Logger,
) *QueryService {
	return &QueryService{
		finder:   finder,
		pager:    pager,
		relation: relation,
		logger:   logger,
	}
}

func (s *QueryService) GetBucket(ctx context.Context, query queries.GetBucketQuery) (*entity.Bucket, error) {
	return s.finder.FindByID(ctx, query.ID)
}

func (s *QueryService) GetBucketByKey(ctx context.Context, query queries.GetBucketByKeyQuery) (*entity.Bucket, error) {
	return s.finder.FindByKey(ctx, query.Key)
}

func (s *QueryService) ListBuckets(ctx context.Context, query queries.ListBucketsQuery) ([]*entity.Bucket, error) {
	return s.pager.List(ctx, query.Limit, query.Offset)
}

func (s *QueryService) SearchBuckets(ctx context.Context, query queries.SearchBucketsQuery) ([]*entity.Bucket, error) {
	return s.pager.Search(ctx, query.Query, query.Limit, query.Offset)
}

func (s *QueryService) GetBucketsByOwner(ctx context.Context, query queries.GetBucketsByOwnerQuery) ([]*entity.Bucket, error) {
	return s.pager.FindByOwnerID(ctx, query.OwnerID, query.Limit, query.Offset)
}

func (s *QueryService) GetBucketAssets(ctx context.Context, query queries.GetBucketAssetsQuery) ([]string, error) {
	return s.relation.GetAssetIDs(ctx, query.BucketID, query.Limit, nil)
}
