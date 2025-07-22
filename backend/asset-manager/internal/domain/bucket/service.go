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

func (s *DomainService) ValidateBucketName(name string) error {
	if name == "" {
		return pkgerrors.NewValidationError("bucket name cannot be empty", nil)
	}
	if len(name) > 100 {
		return pkgerrors.NewValidationError("bucket name too long", nil)
	}
	return nil
}

func (s *DomainService) ValidateBucketKey(key string) error {
	if !isValidKey(key) {
		return pkgerrors.NewValidationError("invalid bucket key format", nil)
	}
	return nil
}

func (s *DomainService) CheckKeyAvailability(ctx context.Context, key string) (bool, error) {
	_, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		if pkgerrors.IsNotFoundError(err) {
			return true, nil
		}
		return false, pkgerrors.WithContext(err, map[string]interface{}{"operation": "CheckKeyAvailability", "key": key})
	}
	return false, nil
}

func (s *DomainService) ValidateBucketOwnership(ctx context.Context, bucketID string, ownerID string) error {
	_, err := s.repo.GetByID(ctx, bucketID)
	if err != nil {
		return pkgerrors.WithContext(err, map[string]interface{}{"operation": "ValidateBucketOwnership", "bucketID": bucketID, "ownerID": ownerID})
	}

	//if bucket.OwnerID() == nil || *bucket.OwnerID() != ownerID {
	//	return errors.New("unauthorized access to bucket")
	//}

	return nil
}

func (s *DomainService) CanAddAssetToBucket(ctx context.Context, bucketID string, assetID string) error {
	exists, err := s.repo.HasAsset(ctx, bucketID, assetID)
	if err != nil {
		return pkgerrors.WithContext(err, map[string]interface{}{"operation": "CanAddAssetToBucket", "bucketID": bucketID, "assetID": assetID})
	}
	if exists {
		return pkgerrors.NewValidationError("asset already exists in bucket", nil)
	}
	return nil
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
