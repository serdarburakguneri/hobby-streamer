package bucket

import (
	"context"

	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
)

type ApplicationService struct {
	repo           domainbucket.Repository
	domainService  *domainbucket.DomainService
	eventPublisher EventPublisher
}

func NewApplicationService(repo domainbucket.Repository, domainService *domainbucket.DomainService, eventPublisher EventPublisher) *ApplicationService {
	return &ApplicationService{
		repo:           repo,
		domainService:  domainService,
		eventPublisher: eventPublisher,
	}
}

func (s *ApplicationService) CreateBucket(ctx context.Context, cmd CreateBucketCommand) (*domainbucket.Bucket, error) {
	if err := s.domainService.ValidateBucketName(cmd.Name); err != nil {
		return nil, err
	}

	if err := s.domainService.ValidateBucketKey(cmd.Key); err != nil {
		return nil, err
	}

	available, err := s.domainService.CheckKeyAvailability(ctx, cmd.Key)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, domainbucket.ErrKeyAlreadyExists
	}

	bucket, err := domainbucket.NewBucket(cmd.Name, cmd.Key)
	if err != nil {
		return nil, err
	}
	if cmd.Description != nil {
		bucket.UpdateDescription(cmd.Description)
	}
	if cmd.OwnerID != nil {
		bucket.SetOwnerID(cmd.OwnerID)
	}
	if cmd.Status != nil {
		bucket.SetStatus(cmd.Status)
	}

	if err := s.repo.Create(ctx, bucket); err != nil {
		return nil, err
	}

	if err := s.eventPublisher.PublishBucketCreated(ctx, bucket); err != nil {
		return nil, err
	}

	return bucket, nil
}

func (s *ApplicationService) GetBucket(ctx context.Context, cmd GetBucketCommand) (*domainbucket.Bucket, error) {
	return s.repo.GetByID(ctx, cmd.ID)
}

func (s *ApplicationService) GetBucketByKey(ctx context.Context, cmd GetBucketByKeyCommand) (*domainbucket.Bucket, error) {
	return s.repo.GetByKey(ctx, cmd.Key)
}

func (s *ApplicationService) UpdateBucket(ctx context.Context, cmd UpdateBucketCommand) (*domainbucket.Bucket, error) {
	bucket, err := s.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if cmd.Name != nil {
		if err := s.domainService.ValidateBucketName(*cmd.Name); err != nil {
			return nil, err
		}
		if err := bucket.UpdateName(*cmd.Name); err != nil {
			return nil, err
		}
	}

	if cmd.Description != nil {
		bucket.UpdateDescription(cmd.Description)
	}

	if cmd.OwnerID != nil {
		bucket.SetOwnerID(cmd.OwnerID)
	}

	if cmd.Metadata != nil {
		bucket.SetMetadata(cmd.Metadata)
	}

	if cmd.Status != nil {
		bucket.SetStatus(cmd.Status)
	}

	if err := s.repo.Update(ctx, bucket); err != nil {
		return nil, err
	}

	// if err := s.eventPublisher.PublishBucketUpdated(ctx, bucket); err != nil {
	// 	return nil, err
	// }

	return bucket, nil
}

func (s *ApplicationService) DeleteBucket(ctx context.Context, cmd DeleteBucketCommand) error {
	if err := s.domainService.ValidateBucketOwnership(ctx, cmd.ID, cmd.OwnerID); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, cmd.ID); err != nil {
		return err
	}

	if err := s.eventPublisher.PublishBucketDeleted(ctx, cmd.ID); err != nil {
		return err
	}

	return nil
}

func (s *ApplicationService) ListBuckets(ctx context.Context, cmd ListBucketsCommand) (*domainbucket.BucketPage, error) {
	return s.repo.List(ctx, cmd.Limit, cmd.LastKey)
}

func (s *ApplicationService) SearchBuckets(ctx context.Context, cmd SearchBucketsCommand) (*domainbucket.BucketPage, error) {
	return s.repo.Search(ctx, cmd.Query, cmd.Limit, cmd.LastKey)
}

func (s *ApplicationService) GetBucketsByOwner(ctx context.Context, cmd GetBucketsByOwnerCommand) (*domainbucket.BucketPage, error) {
	return s.repo.GetByOwnerID(ctx, cmd.OwnerID, cmd.Limit, cmd.LastKey)
}

func (s *ApplicationService) AddAssetToBucket(ctx context.Context, cmd AddAssetToBucketCommand) error {
	if err := s.domainService.ValidateBucketOwnership(ctx, cmd.BucketID, cmd.OwnerID); err != nil {
		return err
	}

	if err := s.domainService.CanAddAssetToBucket(ctx, cmd.BucketID, cmd.AssetID); err != nil {
		return err
	}

	if err := s.repo.AddAsset(ctx, cmd.BucketID, cmd.AssetID); err != nil {
		return err
	}

	if err := s.eventPublisher.PublishAssetAddedToBucket(ctx, cmd.BucketID, cmd.AssetID); err != nil {
		return err
	}

	return nil
}

func (s *ApplicationService) RemoveAssetFromBucket(ctx context.Context, cmd RemoveAssetFromBucketCommand) error {
	if err := s.domainService.ValidateBucketOwnership(ctx, cmd.BucketID, cmd.OwnerID); err != nil {
		return err
	}

	if err := s.repo.RemoveAsset(ctx, cmd.BucketID, cmd.AssetID); err != nil {
		return err
	}

	if err := s.eventPublisher.PublishAssetRemovedFromBucket(ctx, cmd.BucketID, cmd.AssetID); err != nil {
		return err
	}

	return nil
}

func (s *ApplicationService) GetBucketAssets(ctx context.Context, cmd GetBucketAssetsCommand) ([]string, error) {
	return s.repo.GetAssetIDs(ctx, cmd.BucketID, cmd.Limit, cmd.LastKey)
}
