package graphql

import (
	appasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset"
	appbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket"
)

type Resolver struct {
	assetCommandService  *appasset.CommandService
	assetQueryService    *appasset.QueryService
	bucketCommandService *appbucket.CommandService
	bucketQueryService   *appbucket.QueryService
}

func NewResolver(
	assetCommandService *appasset.CommandService,
	assetQueryService *appasset.QueryService,
	bucketCommandService *appbucket.CommandService,
	bucketQueryService *appbucket.QueryService,
) *Resolver {
	return &Resolver{
		assetCommandService:  assetCommandService,
		assetQueryService:    assetQueryService,
		bucketCommandService: bucketCommandService,
		bucketQueryService:   bucketQueryService,
	}
}
