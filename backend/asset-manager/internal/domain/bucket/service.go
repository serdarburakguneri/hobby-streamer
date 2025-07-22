package bucket

import (
	"context"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type DomainService struct {
	repo Repository
}

func NewDomainService(repo Repository) *DomainService {
	return &DomainService{
		repo: repo,
	}
}

func (s *DomainService) IsKeyAvailable(ctx context.Context, key string) (bool, error) {
	_, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		if pkgerrors.IsNotFoundError(err) {
			return true, nil
		}
		return false, pkgerrors.WithContext(err, map[string]interface{}{"operation": "IsKeyAvailable", "key": key})
	}
	return false, nil
}

func (s *DomainService) ValidateBucketOwnership(ctx context.Context, bucketID string, ownerID string) error {
	bucket, err := s.repo.GetByID(ctx, bucketID)
	if err != nil {
		return pkgerrors.WithContext(err, map[string]interface{}{"operation": "ValidateBucketOwnership", "bucketID": bucketID, "ownerID": ownerID})
	}

	return bucket.ValidateOwnership(ownerID)
}

func (s *DomainService) ValidateBucketNotEmpty(ctx context.Context, bucketID string) error {
	count, err := s.repo.AssetCount(ctx, bucketID)
	if err != nil {
		return pkgerrors.WithContext(err, map[string]interface{}{"operation": "ValidateBucketNotEmpty", "bucketID": bucketID})
	}
	if count == 0 {
		return pkgerrors.NewValidationError("cannot perform operation on empty bucket", nil)
	}
	return nil
}
