package entity

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/valueobjects"
)

type Bucket struct {
	id          valueobjects.BucketID
	key         valueobjects.BucketKey
	name        valueobjects.BucketName
	description *valueobjects.BucketDescription
	bucketType  valueobjects.BucketType
	status      *valueobjects.BucketStatus
	assetIDs    *valueobjects.AssetIDs
	createdAt   valueobjects.CreatedAt
	updatedAt   valueobjects.UpdatedAt
	assets      []*entity.Asset
}

func NewBucket(
	id valueobjects.BucketID,
	key valueobjects.BucketKey,
	name valueobjects.BucketName,
	description *valueobjects.BucketDescription,
	bucketType valueobjects.BucketType,
	status *valueobjects.BucketStatus,
	assetIDs *valueobjects.AssetIDs,
	createdAt valueobjects.CreatedAt,
	updatedAt valueobjects.UpdatedAt,
	assets []*entity.Asset,
) *Bucket {
	return &Bucket{
		id:          id,
		key:         key,
		name:        name,
		description: description,
		bucketType:  bucketType,
		status:      status,
		assetIDs:    assetIDs,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		assets:      assets,
	}
}

func (b *Bucket) ID() valueobjects.BucketID {
	return b.id
}

func (b *Bucket) Key() valueobjects.BucketKey {
	return b.key
}

func (b *Bucket) Name() valueobjects.BucketName {
	return b.name
}

func (b *Bucket) Description() *valueobjects.BucketDescription {
	return b.description
}

func (b *Bucket) Type() valueobjects.BucketType {
	return b.bucketType
}

func (b *Bucket) Status() *valueobjects.BucketStatus {
	return b.status
}

func (b *Bucket) AssetIDs() *valueobjects.AssetIDs {
	return b.assetIDs
}

func (b *Bucket) CreatedAt() valueobjects.CreatedAt {
	return b.createdAt
}

func (b *Bucket) UpdatedAt() valueobjects.UpdatedAt {
	return b.updatedAt
}

func (b *Bucket) Assets() []*entity.Asset {
	return b.assets
}

func (b *Bucket) AssetCount() int {
	return len(b.assets)
}

func (b *Bucket) IsActive() bool {
	if b.status == nil {
		return true
	}
	return b.status.Value() == "active"
}

func (b *Bucket) IsCollection() bool {
	return b.bucketType.Value() == "collection"
}

func (b *Bucket) IsPlaylist() bool {
	return b.bucketType.Value() == "playlist"
}

func (b *Bucket) IsCategory() bool {
	return b.bucketType.Value() == "category"
}

func (b *Bucket) IsFeatured() bool {
	return b.bucketType.Value() == "featured"
}

func (b *Bucket) IsTrending() bool {
	return b.bucketType.Value() == "trending"
}

func (b *Bucket) GetPublicAssets() []*entity.Asset {
	var publicAssets []*entity.Asset
	for _, asset := range b.assets {
		if asset.IsPublished() {
			publicAssets = append(publicAssets, asset)
		}
	}
	return publicAssets
}

func (b *Bucket) GetAssetsByType(assetType assetvalueobjects.AssetType) []*entity.Asset {
	var filteredAssets []*entity.Asset
	for _, asset := range b.assets {
		if asset.Type().Equals(assetType) {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	return filteredAssets
}

func (b *Bucket) GetAssetsByGenre(genre assetvalueobjects.Genre) []*entity.Asset {
	var filteredAssets []*entity.Asset
	for _, asset := range b.assets {
		if asset.Genre() != nil && asset.Genre().Equals(genre) {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	return filteredAssets
}

func (b *Bucket) ContainsAsset(assetID string) bool {
	if b.assetIDs == nil {
		return false
	}
	return b.assetIDs.Contains(assetID)
}

func (b *Bucket) GetReadyAssets() []*entity.Asset {
	var readyAssets []*entity.Asset
	for _, asset := range b.assets {
		if asset.IsReady() {
			readyAssets = append(readyAssets, asset)
		}
	}
	return readyAssets
}

func (b *Bucket) GetAssetsWithVideos() []*entity.Asset {
	var assetsWithVideos []*entity.Asset
	for _, asset := range b.assets {
		if asset.HasVideo() {
			assetsWithVideos = append(assetsWithVideos, asset)
		}
	}
	return assetsWithVideos
}

func (b *Bucket) GetAssetsWithImages() []*entity.Asset {
	var assetsWithImages []*entity.Asset
	for _, asset := range b.assets {
		if asset.HasImage() {
			assetsWithImages = append(assetsWithImages, asset)
		}
	}
	return assetsWithImages
}
