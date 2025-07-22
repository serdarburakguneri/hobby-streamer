package bucket

import (
	"context"
	"errors"
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
		return errors.New("bucket name cannot be empty")
	}
	if len(name) > 100 {
		return errors.New("bucket name too long")
	}
	return nil
}

func (s *DomainService) ValidateBucketKey(key string) error {
	if !isValidKey(key) {
		return errors.New("invalid bucket key format")
	}
	return nil
}

func (s *DomainService) CheckKeyAvailability(ctx context.Context, key string) (bool, error) {
	_, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		if errors.Is(err, ErrBucketNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (s *DomainService) ValidateBucketOwnership(ctx context.Context, bucketID string, ownerID string) error {
	_, err := s.repo.GetByID(ctx, bucketID)
	if err != nil {
		return err
	}

	//if bucket.OwnerID() == nil || *bucket.OwnerID() != ownerID {
	//	return errors.New("unauthorized access to bucket")
	//}

	return nil
}

func (s *DomainService) CanAddAssetToBucket(ctx context.Context, bucketID string, assetID string) error {
	exists, err := s.repo.HasAsset(ctx, bucketID, assetID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("asset already exists in bucket")
	}
	return nil
}

func (s *DomainService) ValidateBucketNotEmpty(ctx context.Context, bucketID string) error {
	count, err := s.repo.AssetCount(ctx, bucketID)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("cannot perform operation on empty bucket")
	}
	return nil
}
