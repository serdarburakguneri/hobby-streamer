package config

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs"
)

type ServicesConfig struct {
	AssetService  *asset.Service
	BucketService *bucket.Service
}

func NewServicesConfig(driver neo4j.Driver, sqsProducer *sqs.Producer, dynamicCfg *config.DynamicConfig) *ServicesConfig {
	assetRepo := asset.NewRepository(driver)
	bucketRepo := bucket.NewRepository(driver)

	assetService := asset.NewServiceWithSQS(assetRepo, sqsProducer, dynamicCfg)
	bucketService := bucket.NewService(bucketRepo)

	return &ServicesConfig{
		AssetService:  assetService,
		BucketService: bucketService,
	}
}
