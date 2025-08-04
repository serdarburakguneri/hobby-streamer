package graphql

import (
	"context"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	resilience "github.com/serdarburakguneri/hobby-streamer/backend/pkg/resilience"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	queries "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/graphql/queries"
)

type AssetRepository struct {
	client         *Client
	logger         *logger.Logger
	circuitBreaker *resilience.CircuitBreaker
}

func NewAssetRepository(client *Client, circuitBreaker *resilience.CircuitBreaker) asset.Repository {
	return &AssetRepository{
		client:         client,
		logger:         logger.Get().WithService("asset-graphql-repository"),
		circuitBreaker: circuitBreaker,
	}
}

func (r *AssetRepository) GetByID(ctx context.Context, id assetvalueobjects.AssetID) (*entity.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, a := range assets {
		if a.ID().Equals(id) {
			return a, nil
		}
	}

	return nil, pkgerrors.NewNotFoundError("asset not found", nil)
}

func (r *AssetRepository) GetBySlug(ctx context.Context, slug assetvalueobjects.Slug) (*entity.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, a := range assets {
		if a.Slug().Equals(slug) {
			return a, nil
		}
	}

	return nil, pkgerrors.NewNotFoundError("asset not found", nil)
}

func (r *AssetRepository) GetAll(ctx context.Context) ([]*entity.Asset, error) {
	query := queries.GetAssetsQuery

	var response struct {
		Assets struct {
			Items []*GraphQLAsset `json:"items"`
		} `json:"assets"`
	}

	err := r.circuitBreaker.Execute(ctx, func() error {
		return r.client.Query(ctx, query, nil, &response)
	})

	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "get_all_assets",
		})
	}

	assets := make([]*entity.Asset, len(response.Assets.Items))
	for i, graphQLAsset := range response.Assets.Items {
		domainAsset, err := ConvertGraphQLAssetToDomain(graphQLAsset)
		if err != nil {
			r.logger.WithError(err).Error("Failed to convert GraphQL asset to domain", "asset_id", graphQLAsset.ID)
			continue
		}
		assets[i] = domainAsset
	}

	return assets, nil
}

func (r *AssetRepository) GetPublic(ctx context.Context) ([]*entity.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var publicAssets []*entity.Asset
	for _, a := range assets {
		if a.IsPublished() {
			publicAssets = append(publicAssets, a)
		}
	}

	return publicAssets, nil
}

func (r *AssetRepository) GetByType(ctx context.Context, assetType assetvalueobjects.AssetType) ([]*entity.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredAssets []*entity.Asset
	for _, a := range assets {
		if a.Type().Equals(assetType) {
			filteredAssets = append(filteredAssets, a)
		}
	}

	return filteredAssets, nil
}

func (r *AssetRepository) GetByGenre(ctx context.Context, genre assetvalueobjects.Genre) ([]*entity.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredAssets []*entity.Asset
	for _, a := range assets {
		if a.Genre() != nil && a.Genre().Equals(genre) {
			filteredAssets = append(filteredAssets, a)
		}
	}

	return filteredAssets, nil
}

func (r *AssetRepository) GetByOwner(ctx context.Context, ownerID assetvalueobjects.OwnerID) ([]*entity.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredAssets []*entity.Asset
	for _, a := range assets {
		if a.OwnerID() != nil && a.OwnerID().Equals(ownerID) {
			filteredAssets = append(filteredAssets, a)
		}
	}

	return filteredAssets, nil
}

func (r *AssetRepository) GetReady(ctx context.Context) ([]*entity.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var readyAssets []*entity.Asset
	for _, a := range assets {
		if a.IsReady() {
			readyAssets = append(readyAssets, a)
		}
	}

	return readyAssets, nil
}

func (r *AssetRepository) GetPublished(ctx context.Context) ([]*entity.Asset, error) {
	assets, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var publishedAssets []*entity.Asset
	for _, a := range assets {
		if a.IsPublished() {
			publishedAssets = append(publishedAssets, a)
		}
	}

	return publishedAssets, nil
}

type GraphQLAsset struct {
	ID          string              `json:"id"`
	Slug        string              `json:"slug"`
	Title       *string             `json:"title"`
	Description *string             `json:"description"`
	Type        string              `json:"type"`
	Genre       *string             `json:"genre"`
	Genres      []string            `json:"genres"`
	Tags        []string            `json:"tags"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Metadata    *string             `json:"metadata"`
	OwnerID     *string             `json:"ownerId"`
	Videos      []GraphQLVideo      `json:"videos"`
	Images      []GraphQLImage      `json:"images"`
	PublishRule *GraphQLPublishRule `json:"publishRule"`
}

type GraphQLVideo struct {
	ID                 string                  `json:"id"`
	Label              string                  `json:"label"`
	Type               VideoType               `json:"type"`
	Format             VideoFormat             `json:"format"`
	StorageLocation    GraphQLS3Object         `json:"storageLocation"`
	Width              *int                    `json:"width"`
	Height             *int                    `json:"height"`
	Duration           *float64                `json:"duration"`
	Bitrate            *int                    `json:"bitrate"`
	Codec              *string                 `json:"codec"`
	Size               *int                    `json:"size"`
	ContentType        *string                 `json:"contentType"`
	StreamInfo         *GraphQLStreamInfo      `json:"streamInfo"`
	Metadata           []string                `json:"metadata"`
	Status             VideoStatus             `json:"status"`
	Thumbnail          *GraphQLImage           `json:"thumbnail"`
	CreatedAt          time.Time               `json:"createdAt"`
	UpdatedAt          time.Time               `json:"updatedAt"`
	Quality            *string                 `json:"quality"`
	IsReady            bool                    `json:"isReady"`
	IsProcessing       bool                    `json:"isProcessing"`
	IsFailed           bool                    `json:"isFailed"`
	SegmentCount       *int                    `json:"segmentCount"`
	VideoCodec         *string                 `json:"videoCodec"`
	AudioCodec         *string                 `json:"audioCodec"`
	AvgSegmentDuration *float64                `json:"avgSegmentDuration"`
	Segments           []string                `json:"segments"`
	FrameRate          *string                 `json:"frameRate"`
	AudioChannels      *int                    `json:"audioChannels"`
	AudioSampleRate    *int                    `json:"audioSampleRate"`
	TranscodingInfo    *GraphQLTranscodingInfo `json:"transcodingInfo"`
}

type GraphQLImage struct {
	ID              string           `json:"id"`
	FileName        string           `json:"fileName"`
	URL             string           `json:"url"`
	Type            ImageType        `json:"type"`
	StorageLocation *GraphQLS3Object `json:"storageLocation"`
	Width           *int             `json:"width"`
	Height          *int             `json:"height"`
	Size            *int             `json:"size"`
	ContentType     *string          `json:"contentType"`
	Metadata        []string         `json:"metadata"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

type GraphQLS3Object struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	URL    string `json:"url"`
}

type GraphQLStreamInfo struct {
	DownloadURL *string `json:"downloadUrl"`
	CDNPrefix   *string `json:"cdnPrefix"`
	URL         *string `json:"url"`
}

type GraphQLPublishRule struct {
	PublishAt   *time.Time `json:"publishAt"`
	UnpublishAt *time.Time `json:"unpublishAt"`
	Regions     []string   `json:"regions"`
	AgeRating   *string    `json:"ageRating"`
}

type GraphQLTranscodingInfo struct {
	JobID       *string    `json:"jobId"`
	Progress    *float64   `json:"progress"`
	OutputURL   *string    `json:"outputUrl"`
	Error       *string    `json:"error"`
	CompletedAt *time.Time `json:"completedAt"`
}

type VideoType string
type VideoFormat string
type VideoStatus string
type ImageType string
