package asset

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset/commands"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type CommandService struct {
	saver  asset.Saver
	finder asset.Finder
	logger *logger.Logger
}

func NewCommandService(
	saver asset.Saver,
	finder asset.Finder,
	logger *logger.Logger,
) *CommandService {
	return &CommandService{
		saver:  saver,
		finder: finder,
		logger: logger,
	}
}

func (s *CommandService) CreateAsset(ctx context.Context, cmd commands.CreateAssetCommand) (*entity.Asset, error) {
	asset, err := entity.NewAsset(cmd.Slug, cmd.Title, cmd.AssetType)
	if err != nil {
		return nil, errors.NewValidationError("failed to create new asset", err)
	}

	if cmd.OwnerID != nil {
		asset.SetOwnerID(cmd.OwnerID)
	}
	if cmd.ParentID != nil {
		asset.SetParentID(cmd.ParentID)
	}

	if err := s.saver.Save(ctx, asset); err != nil {
		return nil, errors.NewInternalError("failed to save asset", err)
	}

	return asset, nil
}

func (s *CommandService) DeleteAsset(ctx context.Context, cmd commands.DeleteAssetCommand) error {
	return s.saver.Delete(ctx, cmd.ID)
}

func (s *CommandService) AddVideo(ctx context.Context, cmd commands.AddVideoCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}

	if _, err := asset.AddVideo(
		cmd.Label,
		cmd.Format,
		cmd.StorageLocation,
		cmd.Width,
		cmd.Height,
		cmd.Duration,
		cmd.Bitrate,
		cmd.Codec,
		cmd.Size,
		cmd.ContentType,
		cmd.VideoCodec,
		cmd.AudioCodec,
		cmd.FrameRate,
		cmd.AudioChannels,
		cmd.AudioSampleRate,
		cmd.StreamInfo,
	); err != nil {
		return errors.NewValidationError("failed to add video", err)
	}

	return s.saver.Update(ctx, asset)
}

func (s *CommandService) RemoveVideo(ctx context.Context, cmd commands.RemoveVideoCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}

	if err := asset.RemoveVideo(cmd.VideoID); err != nil {
		return errors.NewValidationError("failed to remove video", err)
	}

	return s.saver.Update(ctx, asset)
}

func (s *CommandService) UpdateVideoStatus(ctx context.Context, cmd commands.UpdateVideoStatusCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}

	if err := asset.UpdateVideoStatus(cmd.VideoID, cmd.Status); err != nil {
		return errors.NewValidationError("failed to update video status", err)
	}

	return s.saver.Update(ctx, asset)
}

func (s *CommandService) UpdateVideoMetadata(ctx context.Context, cmd commands.UpdateVideoMetadataCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}
	contentTypeVO, err := valueobjects.NewContentType(cmd.ContentType)
	if err != nil {
		return errors.NewValidationError("invalid content type", err)
	}
	transcodingInfo := valueobjects.NewTranscodingInfo(cmd.Width, cmd.Height, cmd.Duration, cmd.Bitrate, cmd.Codec, cmd.Size, *contentTypeVO)
	if err := asset.UpdateVideoTranscodingInfo(cmd.VideoID, *transcodingInfo); err != nil {
		return errors.NewValidationError("failed to update video metadata", err)
	}
	return s.saver.Update(ctx, asset)
}

func (s *CommandService) AddImage(ctx context.Context, cmd commands.AddImageCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}

	asset.AddImage(cmd.Image)

	return s.saver.Update(ctx, asset)
}

func (s *CommandService) RemoveImage(ctx context.Context, cmd commands.RemoveImageCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}

	if err := asset.RemoveImage(cmd.ImageID); err != nil {
		return errors.NewValidationError("failed to remove image", err)
	}

	return s.saver.Update(ctx, asset)
}

func (s *CommandService) SetPublishRule(ctx context.Context, cmd commands.SetAssetPublishRuleCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}
	if err := asset.SetPublishRule(&cmd.PublishRule); err != nil {
		return errors.NewValidationError("failed to set publish rule", err)
	}
	return s.saver.Update(ctx, asset)
}

func (s *CommandService) ClearPublishRule(ctx context.Context, cmd commands.ClearAssetPublishRuleCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}
	if err := asset.SetPublishRule(nil); err != nil {
		return errors.NewValidationError("failed to clear publish rule", err)
	}
	return s.saver.Update(ctx, asset)
}

func (s *CommandService) UpdateAssetTitle(ctx context.Context, cmd commands.UpdateAssetTitleCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}
	asset.UpdateTitle(&cmd.Title)
	return s.saver.Update(ctx, asset)
}

func (s *CommandService) UpdateAssetDescription(ctx context.Context, cmd commands.UpdateAssetDescriptionCommand) error {
	asset, err := s.finder.FindByID(ctx, cmd.AssetID)
	if err != nil {
		return errors.NewInternalError("failed to find asset", err)
	}
	if asset == nil {
		return errors.NewNotFoundError("asset not found", nil)
	}
	asset.UpdateDescription(&cmd.Description)
	return s.saver.Update(ctx, asset)
}
