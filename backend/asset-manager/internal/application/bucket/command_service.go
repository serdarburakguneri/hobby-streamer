package bucket

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket/commands"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type CommandService struct {
	saver    bucket.Saver
	finder   bucket.Finder
	relation bucket.Relation
	logger   *logger.Logger
}

func NewCommandService(
	saver bucket.Saver,
	finder bucket.Finder,
	relation bucket.Relation,
	logger *logger.Logger,
) *CommandService {
	return &CommandService{
		saver:    saver,
		finder:   finder,
		relation: relation,
		logger:   logger,
	}
}

func (s *CommandService) CreateBucket(ctx context.Context, cmd commands.CreateBucketCommand) (*entity.Bucket, error) {
	bucket, err := entity.NewBucket(cmd.Name, cmd.Key)
	if err != nil {
		return nil, errors.NewValidationError("failed to create new bucket", err)
	}

	if cmd.OwnerID != nil {
		bucket.UpdateOwnerID(cmd.OwnerID)
	}

	if err := s.saver.Save(ctx, bucket); err != nil {
		return nil, errors.NewInternalError("failed to save bucket", err)
	}

	return bucket, nil
}

func (s *CommandService) UpdateBucket(ctx context.Context, cmd commands.UpdateBucketCommand) error {
	bucket, err := s.finder.FindByID(ctx, cmd.ID)
	if err != nil {
		return errors.NewInternalError("failed to find bucket", err)
	}
	if bucket == nil {
		return errors.NewNotFoundError("bucket not found", nil)
	}

	if cmd.Name != nil {
		bucket.UpdateName(*cmd.Name)
	}
	if cmd.Description != nil {
		bucket.UpdateDescription(cmd.Description)
	}
	if cmd.Status != nil {
		bucket.UpdateStatus(cmd.Status)
	}
	if cmd.Type != nil {
		bucket.UpdateType(cmd.Type)
	}
	if cmd.OwnerID != nil {
		bucket.UpdateOwnerID(cmd.OwnerID)
	}
	if cmd.Metadata != nil {
		bucket.UpdateMetadata(cmd.Metadata)
	}

	return s.saver.Update(ctx, bucket)
}

func (s *CommandService) DeleteBucket(ctx context.Context, cmd commands.DeleteBucketCommand) error {
	return s.saver.Delete(ctx, cmd.ID)
}

func (s *CommandService) AddAssetToBucket(ctx context.Context, cmd commands.AddAssetToBucketCommand) error {
	if err := s.relation.AddAsset(ctx, cmd.BucketID, cmd.AssetID); err != nil {
		return errors.NewValidationError("failed to add asset to bucket", err)
	}
	return nil
}

func (s *CommandService) RemoveAssetFromBucket(ctx context.Context, cmd commands.RemoveAssetFromBucketCommand) error {
	if err := s.relation.RemoveAsset(ctx, cmd.BucketID, cmd.AssetID); err != nil {
		return errors.NewValidationError("failed to remove asset from bucket", err)
	}
	return nil
}
