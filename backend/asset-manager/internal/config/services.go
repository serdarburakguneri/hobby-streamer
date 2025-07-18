package config

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
)

type ServicesConfig struct {
	AssetService  *asset.Service
	BucketService *bucket.Service
}

func NewServicesConfig(driver neo4j.Driver, dynamicCfg *config.DynamicConfig) *ServicesConfig {
	assetRepo := asset.NewRepository(driver)
	bucketRepo := bucket.NewRepository(driver)

	assetService := asset.NewService(assetRepo, dynamicCfg)
	bucketService := bucket.NewService(bucketRepo)

	return &ServicesConfig{
		AssetService:  assetService,
		BucketService: bucketService,
	}
}
