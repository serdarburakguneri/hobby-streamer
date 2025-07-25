package sqs

import (
	"context"
	"encoding/json"
	"time"

	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type BucketEventPublisher struct {
	producer *sqs.Producer
}

func NewBucketEventPublisher(producer *sqs.Producer) *BucketEventPublisher {
	return &BucketEventPublisher{
		producer: producer,
	}
}

func (p *BucketEventPublisher) PublishBucketCreated(ctx context.Context, bucket *domainbucket.Bucket) error {
	event := BucketCreatedEvent{
		EventType: "bucket.created",
		BucketID:  bucket.ID().Value(),
		Data: BucketEventData{
			ID:          bucket.ID().Value(),
			Name:        bucket.Name(),
			Description: stringPtrToString(bucket.Description()),
			Key:         bucket.Key(),
			OwnerID:     stringPtrToString(bucket.OwnerID()),
			Metadata:    bucket.Metadata(),
			CreatedAt:   bucket.CreatedAt().Format(time.RFC3339),
			UpdatedAt:   bucket.UpdatedAt().Format(time.RFC3339),
		},
	}

	return p.publishEvent(ctx, event)
}

func (p *BucketEventPublisher) PublishBucketUpdated(ctx context.Context, bucket *domainbucket.Bucket) error {
	event := BucketUpdatedEvent{
		EventType: "bucket.updated",
		BucketID:  bucket.ID().Value(),
		Data: BucketEventData{
			ID:          bucket.ID().Value(),
			Name:        bucket.Name(),
			Description: stringPtrToString(bucket.Description()),
			Key:         bucket.Key(),
			OwnerID:     stringPtrToString(bucket.OwnerID()),
			Metadata:    bucket.Metadata(),
			CreatedAt:   bucket.CreatedAt().Format(time.RFC3339),
			UpdatedAt:   bucket.UpdatedAt().Format(time.RFC3339),
		},
	}

	return p.publishEvent(ctx, event)
}

func (p *BucketEventPublisher) PublishBucketDeleted(ctx context.Context, bucketID string) error {
	event := BucketDeletedEvent{
		EventType: "bucket.deleted",
		BucketID:  bucketID,
		Data: BucketDeletedEventData{
			BucketID: bucketID,
		},
	}

	return p.publishEvent(ctx, event)
}

func (p *BucketEventPublisher) PublishAssetAddedToBucket(ctx context.Context, bucketID string, assetID string) error {
	event := AssetAddedToBucketEvent{
		EventType: "bucket.asset.added",
		BucketID:  bucketID,
		Data: AssetAddedToBucketEventData{
			BucketID: bucketID,
			AssetID:  assetID,
		},
	}

	return p.publishEvent(ctx, event)
}

func (p *BucketEventPublisher) PublishAssetRemovedFromBucket(ctx context.Context, bucketID string, assetID string) error {
	event := AssetRemovedFromBucketEvent{
		EventType: "bucket.asset.removed",
		BucketID:  bucketID,
		Data: AssetRemovedFromBucketEventData{
			BucketID: bucketID,
			AssetID:  assetID,
		},
	}

	return p.publishEvent(ctx, event)
}

func (p *BucketEventPublisher) publishEvent(ctx context.Context, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return pkgerrors.NewInternalError("failed to marshal event", err)
	}

	err = p.producer.SendMessage(ctx, string(payload), "bucket-event")
	if err != nil {
		return pkgerrors.NewInternalError("failed to send bucket event", err)
	}

	return nil
}

func stringPtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

type BucketCreatedEvent struct {
	EventType string          `json:"eventType"`
	BucketID  string          `json:"bucketID"`
	Data      BucketEventData `json:"data"`
}

type BucketUpdatedEvent struct {
	EventType string          `json:"eventType"`
	BucketID  string          `json:"bucketID"`
	Data      BucketEventData `json:"data"`
}

type BucketDeletedEvent struct {
	EventType string                 `json:"eventType"`
	BucketID  string                 `json:"bucketID"`
	Data      BucketDeletedEventData `json:"data"`
}

type AssetAddedToBucketEvent struct {
	EventType string                      `json:"eventType"`
	BucketID  string                      `json:"bucketID"`
	Data      AssetAddedToBucketEventData `json:"data"`
}

type AssetRemovedFromBucketEvent struct {
	EventType string                          `json:"eventType"`
	BucketID  string                          `json:"bucketID"`
	Data      AssetRemovedFromBucketEventData `json:"data"`
}

type BucketEventData struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Key         string                 `json:"key"`
	OwnerID     string                 `json:"ownerID"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
}

type BucketDeletedEventData struct {
	BucketID string `json:"bucketID"`
}

type AssetAddedToBucketEventData struct {
	BucketID string `json:"bucketID"`
	AssetID  string `json:"assetID"`
}

type AssetRemovedFromBucketEventData struct {
	BucketID string `json:"bucketID"`
	AssetID  string `json:"assetID"`
}
