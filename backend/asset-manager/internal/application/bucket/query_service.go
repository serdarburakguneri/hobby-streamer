package bucket

import (
	"context"
	"strconv"

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

// ListBucketsPage handles pagination logic for listing buckets.
func (s *QueryService) ListBucketsPage(ctx context.Context, query queries.ListBucketsQuery) (*entity.BucketPage, error) {
	items, err := s.pager.List(ctx, query.Limit, query.Offset)
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

	return &entity.BucketPage{Items: items, HasMore: hasMore, LastKey: lastKey}, nil
}

// SearchBucketsPage handles pagination logic for searching buckets.
func (s *QueryService) SearchBucketsPage(ctx context.Context, query queries.SearchBucketsQuery) (*entity.BucketPage, error) {
	items, err := s.pager.Search(ctx, query.Query, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

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

	hasMore := len(items) >= limitVal
	lastKey := make(map[string]interface{})
	if hasMore {
		lastKey["key"] = strconv.Itoa(offsetVal + len(items))
	}

	return &entity.BucketPage{Items: items, HasMore: hasMore, LastKey: lastKey}, nil
}

// GetBucketsByOwnerPage handles pagination logic for buckets by owner.
func (s *QueryService) GetBucketsByOwnerPage(ctx context.Context, query queries.GetBucketsByOwnerQuery) (*entity.BucketPage, error) {
	items, err := s.pager.FindByOwnerID(ctx, query.OwnerID, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

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

	hasMore := len(items) >= limitVal
	lastKey := make(map[string]interface{})
	if hasMore {
		lastKey["key"] = strconv.Itoa(offsetVal + len(items))
	}

	return &entity.BucketPage{Items: items, HasMore: hasMore, LastKey: lastKey}, nil
}
