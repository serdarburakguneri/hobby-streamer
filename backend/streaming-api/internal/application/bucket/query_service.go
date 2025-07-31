package bucket

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
	bucketvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/valueobjects"
)

type Service struct {
	repo bucket.Repository
}

func NewService(repo bucket.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBucket(ctx context.Context, key bucketvalueobjects.BucketKey) (*bucketentity.Bucket, error) {
	return s.repo.GetByKey(ctx, key)
}

func (s *Service) GetBuckets(ctx context.Context, limit int, nextKey *string) ([]*bucketentity.Bucket, error) {
	return s.repo.GetAll(ctx, limit, nextKey)
}
