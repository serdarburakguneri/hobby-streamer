package bucket

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset"
)

type Bucket struct {
	id          BucketID
	key         BucketKey
	name        BucketName
	description *BucketDescription
	bucketType  BucketType
	status      *BucketStatus
	assetIDs    *AssetIDs
	createdAt   CreatedAt
	updatedAt   UpdatedAt
	assets      []*asset.Asset
}

func NewBucket(
	id BucketID,
	key BucketKey,
	name BucketName,
	description *BucketDescription,
	bucketType BucketType,
	status *BucketStatus,
	assetIDs *AssetIDs,
	createdAt CreatedAt,
	updatedAt UpdatedAt,
	assets []*asset.Asset,
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

func (b *Bucket) ID() BucketID {
	return b.id
}

func (b *Bucket) Key() BucketKey {
	return b.key
}

func (b *Bucket) Name() BucketName {
	return b.name
}

func (b *Bucket) Description() *BucketDescription {
	return b.description
}

func (b *Bucket) Type() BucketType {
	return b.bucketType
}

func (b *Bucket) Status() *BucketStatus {
	return b.status
}

func (b *Bucket) AssetIDs() *AssetIDs {
	return b.assetIDs
}

func (b *Bucket) CreatedAt() CreatedAt {
	return b.createdAt
}

func (b *Bucket) UpdatedAt() UpdatedAt {
	return b.updatedAt
}

func (b *Bucket) Assets() []*asset.Asset {
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

func (b *Bucket) GetPublicAssets() []*asset.Asset {
	var publicAssets []*asset.Asset
	for _, asset := range b.assets {
		if asset.IsPublished() {
			publicAssets = append(publicAssets, asset)
		}
	}
	return publicAssets
}

func (b *Bucket) GetAssetsByType(assetType asset.AssetType) []*asset.Asset {
	var filteredAssets []*asset.Asset
	for _, asset := range b.assets {
		if asset.Type().Equals(assetType) {
			filteredAssets = append(filteredAssets, asset)
		}
	}
	return filteredAssets
}

func (b *Bucket) GetAssetsByGenre(genre asset.Genre) []*asset.Asset {
	var filteredAssets []*asset.Asset
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

func (b *Bucket) GetReadyAssets() []*asset.Asset {
	var readyAssets []*asset.Asset
	for _, asset := range b.assets {
		if asset.IsReady() {
			readyAssets = append(readyAssets, asset)
		}
	}
	return readyAssets
}

func (b *Bucket) GetAssetsWithVideos() []*asset.Asset {
	var assetsWithVideos []*asset.Asset
	for _, asset := range b.assets {
		if asset.HasVideo() {
			assetsWithVideos = append(assetsWithVideos, asset)
		}
	}
	return assetsWithVideos
}

func (b *Bucket) GetAssetsWithImages() []*asset.Asset {
	var assetsWithImages []*asset.Asset
	for _, asset := range b.assets {
		if asset.HasImage() {
			assetsWithImages = append(assetsWithImages, asset)
		}
	}
	return assetsWithImages
}
