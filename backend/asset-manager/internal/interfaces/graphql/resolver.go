package graphql

import (
	appasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset"
	appbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket"
)

type Resolver struct {
	assetAppService  *appasset.ApplicationService
	bucketAppService *appbucket.ApplicationService
}

func NewResolver(assetAppService *appasset.ApplicationService, bucketAppService *appbucket.ApplicationService) *Resolver {
	return &Resolver{
		assetAppService:  assetAppService,
		bucketAppService: bucketAppService,
	}
}
