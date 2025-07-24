package bucket

import (
	"context"

	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type ApplicationService struct {
	repo           domainbucket.Repository
	domainService  *domainbucket.DomainService
	eventPublisher EventPublisher
	logger         *logger.Logger
}

func NewApplicationService(repo domainbucket.Repository, domainService *domainbucket.DomainService, eventPublisher EventPublisher) *ApplicationService {
	return &ApplicationService{
		repo:           repo,
		domainService:  domainService,
		eventPublisher: eventPublisher,
		logger:         logger.WithService("bucket-application-service"),
	}
}

func (s *ApplicationService) CreateBucket(ctx context.Context, cmd CreateBucketCommand) (*domainbucket.Bucket, error) {
	log := s.logger.WithContext(ctx)

	bucket, err := domainbucket.NewBucket(cmd.Name, cmd.Key)
	if err != nil {
		log.WithError(err).Error("Failed to create bucket", "name", cmd.Name, "key", cmd.Key)
		return nil, err
	}

	if cmd.Type != nil {
		err := bucket.SetType(*cmd.Type)
		if err != nil {
			log.WithError(err).Error("Failed to set bucket type", "type", *cmd.Type)
			return nil, err
		}
	}

	available, err := s.domainService.IsKeyAvailable(ctx, cmd.Key)
	if err != nil {
		log.WithError(err).Error("Failed to check key availability", "key", cmd.Key)
		return nil, pkgerrors.NewInternalError("failed to check key availability", err)
	}
	if !available {
		log.Error("Bucket key already exists", "key", cmd.Key)
		return nil, domainbucket.ErrKeyAlreadyExists
	}

	if cmd.Description != nil {
		bucket.UpdateDescription(cmd.Description)
	}
	if cmd.OwnerID != nil {
		bucket.SetOwnerID(cmd.OwnerID)
	}
	if cmd.Status != nil {
		if err := bucket.SetStatus(cmd.Status); err != nil {
			log.WithError(err).Error("Failed to set bucket status", "status", *cmd.Status)
			return nil, err
		}
	}

	if err := s.repo.Create(ctx, bucket); err != nil {
		log.WithError(err).Error("Failed to save bucket", "bucket_id", bucket.ID())
		return nil, pkgerrors.NewInternalError("failed to save bucket", err)
	}

	if err := s.eventPublisher.PublishBucketCreated(ctx, bucket); err != nil {
		log.WithError(err).Error("Failed to publish bucket created event", "bucket_id", bucket.ID())
	}

	log.Info("Bucket created successfully", "bucket_id", bucket.ID(), "key", bucket.Key())
	return bucket, nil
}

func (s *ApplicationService) GetBucket(ctx context.Context, cmd GetBucketCommand) (*domainbucket.Bucket, error) {
	idVO, err := cmd.ToDomainBucketID()
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, *idVO)
}

func (s *ApplicationService) GetBucketByKey(ctx context.Context, cmd GetBucketByKeyCommand) (*domainbucket.Bucket, error) {
	return s.repo.GetByKey(ctx, cmd.Key)
}

func (s *ApplicationService) UpdateBucket(ctx context.Context, cmd UpdateBucketCommand) (*domainbucket.Bucket, error) {
	log := s.logger.WithContext(ctx)
	idVO, err := cmd.ToDomainBucketID()
	if err != nil {
		return nil, err
	}
	bucket, err := s.repo.GetByID(ctx, *idVO)
	if err != nil {
		log.WithError(err).Error("Failed to find bucket for update", "bucket_id", cmd.ID)
		return nil, err
	}

	if cmd.Name != nil {
		if err := bucket.UpdateName(*cmd.Name); err != nil {
			log.WithError(err).Error("Failed to update bucket name", "bucket_id", cmd.ID, "name", *cmd.Name)
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
		if err := bucket.SetStatus(cmd.Status); err != nil {
			log.WithError(err).Error("Failed to set bucket status", "status", *cmd.Status)
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, bucket); err != nil {
		log.WithError(err).Error("Failed to update bucket", "bucket_id", cmd.ID)
		return nil, pkgerrors.NewInternalError("failed to update bucket", err)
	}

	log.Info("Bucket updated successfully", "bucket_id", cmd.ID)
	return bucket, nil
}

func (s *ApplicationService) DeleteBucket(ctx context.Context, cmd DeleteBucketCommand) error {
	idVO, err := cmd.ToDomainBucketID()
	if err != nil {
		return err
	}
	if err := s.domainService.ValidateBucketOwnership(ctx, cmd.ID, cmd.OwnerID); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, *idVO); err != nil {
		return err
	}
	if err := s.eventPublisher.PublishBucketDeleted(ctx, idVO.Value()); err != nil {
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
	idVO, err := cmd.ToDomainBucketID()
	if err != nil {
		return err
	}
	if err := s.domainService.ValidateBucketOwnership(ctx, cmd.BucketID, cmd.OwnerID); err != nil {
		return err
	}
	bucket, err := s.repo.GetByID(ctx, *idVO)
	if err != nil {
		return err
	}
	if err := bucket.CanAddAsset(cmd.AssetID, func(bucketID, assetID string) (bool, error) {
		return s.repo.HasAsset(ctx, *idVO, assetID)
	}); err != nil {
		return err
	}
	if err := s.repo.AddAsset(ctx, *idVO, cmd.AssetID); err != nil {
		return err
	}
	if err := s.eventPublisher.PublishAssetAddedToBucket(ctx, idVO.Value(), cmd.AssetID); err != nil {
		return err
	}
	return nil
}

func (s *ApplicationService) RemoveAssetFromBucket(ctx context.Context, cmd RemoveAssetFromBucketCommand) error {
	idVO, err := cmd.ToDomainBucketID()
	if err != nil {
		return err
	}
	if err := s.domainService.ValidateBucketOwnership(ctx, cmd.BucketID, cmd.OwnerID); err != nil {
		return err
	}
	if err := s.repo.RemoveAsset(ctx, *idVO, cmd.AssetID); err != nil {
		return err
	}
	if err := s.eventPublisher.PublishAssetRemovedFromBucket(ctx, idVO.Value(), cmd.AssetID); err != nil {
		return err
	}
	return nil
}

func (s *ApplicationService) GetBucketAssets(ctx context.Context, cmd GetBucketAssetsCommand) ([]string, error) {
	idVO, err := cmd.ToDomainBucketID()
	if err != nil {
		return nil, err
	}
	return s.repo.GetAssetIDs(ctx, *idVO, cmd.Limit, cmd.LastKey)
}
