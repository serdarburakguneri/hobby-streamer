package graph

import (
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bucket"
)

type Resolver struct {
	AssetService  asset.AssetService
	BucketService bucket.BucketService
}

func NewResolver(assetService asset.AssetService, bucketService bucket.BucketService) *Resolver {
	return &Resolver{
		AssetService:  assetService,
		BucketService: bucketService,
	}
}
