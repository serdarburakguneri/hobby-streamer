package bucket

import (
	"context"

	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
)

type EventPublisher interface {
	PublishBucketCreated(ctx context.Context, bucket *domainbucket.Bucket) error
	PublishBucketUpdated(ctx context.Context, bucket *domainbucket.Bucket) error
	PublishBucketDeleted(ctx context.Context, bucketID string) error
	PublishAssetAddedToBucket(ctx context.Context, bucketID string, assetID string) error
	PublishAssetRemovedFromBucket(ctx context.Context, bucketID string, assetID string) error
}
