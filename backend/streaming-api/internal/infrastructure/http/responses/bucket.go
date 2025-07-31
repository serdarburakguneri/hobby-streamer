package responses

import (
	"time"

	assetentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/entity"
	bucketentity "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/entity"
)

type BucketsResponse struct {
	Buckets []BucketResponse `json:"buckets"`
	Count   int              `json:"count"`
	Limit   int              `json:"limit"`
	NextKey *string          `json:"nextKey,omitempty"`
}

type BucketResponse struct {
	ID          string          `json:"id"`
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Type        string          `json:"type"`
	Status      *string         `json:"status,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Assets      []AssetResponse `json:"assets,omitempty"`
}

func NewBucketResponse(b *bucketentity.Bucket) BucketResponse {
	var description *string
	if b.Description() != nil {
		desc := b.Description().Value()
		description = &desc
	}

	var status *string
	if b.Status() != nil {
		statusVal := b.Status().Value()
		status = &statusVal
	}

	return BucketResponse{
		ID:          b.ID().Value(),
		Key:         b.Key().Value(),
		Name:        b.Name().Value(),
		Description: description,
		Type:        b.Type().Value(),
		Status:      status,
		CreatedAt:   b.CreatedAt().Value(),
		UpdatedAt:   b.UpdatedAt().Value(),
		Assets:      convertAssetsToResponse(b.Assets()),
	}
}

func convertAssetsToResponse(assets []*assetentity.Asset) []AssetResponse {
	responses := make([]AssetResponse, len(assets))
	for i, a := range assets {
		responses[i] = NewAssetResponse(a)
	}
	return responses
}
