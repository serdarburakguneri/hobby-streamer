package http

import (
	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/bucket"
)

func NewRouter(assetService *asset.Service, bucketService *bucket.Service) *mux.Router {
	r := mux.NewRouter()

	assetHandler := &asset.AssetHandler{Service: assetService}
	bucketHandler := &bucket.BucketHandler{Service: bucketService}

	// Core asset routes
	r.HandleFunc("/assets", assetHandler.ListAssets).Methods("GET")
	r.HandleFunc("/assets/{id}", assetHandler.GetAsset).Methods("GET")
	r.HandleFunc("/assets", assetHandler.CreateAsset).Methods("POST")
	r.HandleFunc("/assets/{id}", assetHandler.PatchAsset).Methods("PATCH")

	// Publish rule routes
	r.HandleFunc("/assets/{id}/publishRule", assetHandler.GetPublishRule).Methods("GET")
	r.HandleFunc("/assets/{id}/publishRule", assetHandler.PatchPublishRule).Methods("PATCH")

	// Video routes
	r.HandleFunc("/assets/{id}/videos/{label}", assetHandler.SetVideoVariant).Methods("POST")
	r.HandleFunc("/assets/{id}/videos/{label}", assetHandler.DeleteVideoVariant).Methods("DELETE")

	// Image routes
	r.HandleFunc("/assets/{id}/images", assetHandler.AddImage).Methods("POST")
	r.HandleFunc("/assets/{id}/images/{filename}", assetHandler.DeleteImage).Methods("DELETE")

	// Bucket routes
	r.HandleFunc("/buckets", bucketHandler.ListBuckets).Methods("GET")
	r.HandleFunc("/buckets/{id}", bucketHandler.GetBucket).Methods("GET")
	r.HandleFunc("/buckets", bucketHandler.CreateBucket).Methods("POST")
	r.HandleFunc("/buckets/{id}", bucketHandler.PatchBucket).Methods("PATCH")
	r.HandleFunc("/buckets/{id}/assets", bucketHandler.AddAssetToBucket).Methods("POST")
	r.HandleFunc("/buckets/{id}/assets/{assetId}", bucketHandler.RemoveAssetFromBucket).Methods("DELETE")

	return r
}
