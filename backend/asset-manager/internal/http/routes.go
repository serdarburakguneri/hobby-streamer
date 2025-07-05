package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/bucket"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

func NewRouter(assetService *asset.Service, bucketService *bucket.Service, authMiddleware *auth.AuthMiddleware) *mux.Router {
	r := mux.NewRouter()

	assetHandler := &asset.AssetHandler{Service: assetService}
	bucketHandler := &bucket.BucketHandler{Service: bucketService}

	// Health check endpoint (no auth required)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "asset-manager"})
	}).Methods("GET")

	// Core asset routes
	r.HandleFunc("/assets", authMiddleware.RequireAuth(
		authMiddleware.RequireAnyRole([]string{constants.RoleAdmin, constants.RoleUser})(assetHandler.ListAssets),
	)).Methods("GET")

	r.HandleFunc("/assets/{id}", authMiddleware.RequireAuth(
		authMiddleware.RequireAnyRole([]string{constants.RoleAdmin, constants.RoleUser})(assetHandler.GetAsset),
	)).Methods("GET")

	r.HandleFunc("/assets", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(assetHandler.CreateAsset),
	)).Methods("POST")

	r.HandleFunc("/assets/{id}", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(assetHandler.PatchAsset),
	)).Methods("PATCH")

	// Publish rule routes
	r.HandleFunc("/assets/{id}/publishRule", authMiddleware.RequireAuth(
		authMiddleware.RequireAnyRole([]string{constants.RoleAdmin, constants.RoleUser})(assetHandler.GetPublishRule),
	)).Methods("GET")

	r.HandleFunc("/assets/{id}/publishRule", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(assetHandler.PatchPublishRule),
	)).Methods("PATCH")

	// Video routes
	r.HandleFunc("/assets/{id}/videos/{label}", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(assetHandler.SetVideoVariant),
	)).Methods("POST")

	r.HandleFunc("/assets/{id}/videos/{label}", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(assetHandler.DeleteVideoVariant),
	)).Methods("DELETE")

	// Image routes
	r.HandleFunc("/assets/{id}/images", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(assetHandler.AddImage),
	)).Methods("POST")

	r.HandleFunc("/assets/{id}/images/{filename}", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(assetHandler.DeleteImage),
	)).Methods("DELETE")

	// Bucket routes
	r.HandleFunc("/buckets", authMiddleware.RequireAuth(
		authMiddleware.RequireAnyRole([]string{constants.RoleAdmin, constants.RoleUser})(bucketHandler.ListBuckets),
	)).Methods("GET")

	r.HandleFunc("/buckets/{id}", authMiddleware.RequireAuth(
		authMiddleware.RequireAnyRole([]string{constants.RoleAdmin, constants.RoleUser})(bucketHandler.GetBucket),
	)).Methods("GET")

	r.HandleFunc("/buckets", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(bucketHandler.CreateBucket),
	)).Methods("POST")

	r.HandleFunc("/buckets/{id}", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(bucketHandler.PatchBucket),
	)).Methods("PATCH")

	r.HandleFunc("/buckets/{id}/assets", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(bucketHandler.AddAssetToBucket),
	)).Methods("POST")

	r.HandleFunc("/buckets/{id}/assets/{assetId}", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(constants.RoleAdmin)(bucketHandler.RemoveAssetFromBucket),
	)).Methods("DELETE")

	return r
}
