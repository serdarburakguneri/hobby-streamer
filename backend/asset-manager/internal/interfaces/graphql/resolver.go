package graphql

import (
	appasset "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/asset"
	appbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/bucket"
	cdn "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/application/cdn"
)

type Resolver struct {
	assetCommandService  *appasset.CommandService
	assetQueryService    *appasset.QueryService
	bucketCommandService *appbucket.CommandService
	bucketQueryService   *appbucket.QueryService
	cdnService           cdn.Service
}

func NewResolver(
	assetCommandService *appasset.CommandService,
	assetQueryService *appasset.QueryService,
	bucketCommandService *appbucket.CommandService,
	bucketQueryService *appbucket.QueryService,
	cdnService cdn.Service,
) *Resolver {
	return &Resolver{
		assetCommandService:  assetCommandService,
		assetQueryService:    assetQueryService,
		bucketCommandService: bucketCommandService,
		bucketQueryService:   bucketQueryService,
		cdnService:           cdnService,
	}
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }
func (r *Resolver) Bucket() BucketResolver     { return &bucketResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type bucketResolver struct{ *Resolver }
