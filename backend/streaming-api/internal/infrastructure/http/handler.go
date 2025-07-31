package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/security"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/application"
	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
	bucketvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/bucket/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/infrastructure/http/responses"
)

type Handler struct {
	assetService  application.AssetServiceInterface
	bucketService application.BucketServiceInterface
	logger        *logger.Logger
	config        *config.BaseConfig
}

func NewHandler(assetService application.AssetServiceInterface, bucketService application.BucketServiceInterface, config *config.BaseConfig) *Handler {
	return &Handler{
		assetService:  assetService,
		bucketService: bucketService,
		logger:        logger.Get().WithService("streaming-handler"),
		config:        config,
	}
}

func (h *Handler) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	corsCfg := h.config.Security.CORS
	allowedOrigins := corsCfg.AllowedOrigins
	allowedMethods := corsCfg.AllowedMethods
	allowedHeaders := corsCfg.AllowedHeaders

	router.Use(security.CORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders))
	router.Use(logger.CompressionMiddleware)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/buckets", h.GetBuckets).Methods("GET")
	api.HandleFunc("/buckets/{key}", h.GetBucket).Methods("GET")
	api.HandleFunc("/buckets/{key}/assets", h.GetAssetsInBucket).Methods("GET")

	api.HandleFunc("/assets", h.GetAssets).Methods("GET")
	api.HandleFunc("/assets/{slug}", h.GetAsset).Methods("GET")

	router.HandleFunc("/health", h.HealthCheck).Methods("GET")

	return router
}

func (h *Handler) GetBuckets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	var nextKey *string
	if nk := r.URL.Query().Get("nextKey"); nk != "" {
		nextKey = &nk
	}

	buckets, err := h.bucketService.GetBuckets(ctx, limit, nextKey)
	if err != nil {
		h.handleError(w, err, "Failed to get buckets")
		return
	}

	var bucketResponses []responses.BucketResponse
	for _, bucket := range buckets {
		bucketResponses = append(bucketResponses, responses.NewBucketResponse(bucket))
	}

	response := responses.BucketsResponse{
		Buckets: bucketResponses,
		Count:   len(bucketResponses),
		Limit:   limit,
		NextKey: nextKey,
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetBucket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	key := vars["key"]

	bucketKeyVO, err := bucketvalueobjects.NewBucketKey(key)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid bucket key")
		return
	}

	bucket, err := h.bucketService.GetBucket(ctx, *bucketKeyVO)
	if err != nil {
		h.handleError(w, err, "Failed to get bucket")
		return
	}

	if bucket == nil {
		h.writeError(w, http.StatusNotFound, "Bucket not found")
		return
	}

	bucketResponse := responses.NewBucketResponse(bucket)
	w.Header().Set("Cache-Control", "public, max-age=1800")
	h.writeJSON(w, http.StatusOK, bucketResponse)
}

func (h *Handler) GetAssetsInBucket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	key := vars["key"]

	bucketKeyVO, err := bucketvalueobjects.NewBucketKey(key)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid bucket key")
		return
	}

	assets, err := h.assetService.GetAssetsInBucket(ctx, *bucketKeyVO)
	if err != nil {
		h.handleError(w, err, "Failed to get assets in bucket")
		return
	}

	var assetResponses []responses.AssetResponse
	for _, asset := range assets {
		assetResponses = append(assetResponses, responses.NewAssetResponse(asset))
	}

	response := responses.AssetsResponse{
		Assets: assetResponses,
		Count:  len(assetResponses),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetAssets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	assets, err := h.assetService.GetAssets(ctx)
	if err != nil {
		h.handleError(w, err, "Failed to get assets")
		return
	}

	var assetResponses []responses.AssetResponse
	for _, asset := range assets {
		assetResponses = append(assetResponses, responses.NewAssetResponse(asset))
	}

	response := responses.AssetsResponse{
		Assets: assetResponses,
		Count:  len(assetResponses),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	slug := vars["slug"]

	slugVO, err := assetvalueobjects.NewSlug(slug)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid asset slug")
		return
	}

	asset, err := h.assetService.GetAsset(ctx, *slugVO)
	if err != nil {
		h.handleError(w, err, "Failed to get asset")
		return
	}

	if asset == nil {
		h.writeError(w, http.StatusNotFound, "Asset not found")
		return
	}

	assetResponse := responses.NewAssetResponse(asset)
	h.writeJSON(w, http.StatusOK, assetResponse)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := responses.HealthResponse{
		Status:  "healthy",
		Service: "streaming-api",
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) handleError(w http.ResponseWriter, err error, defaultMessage string) {
	if errors.IsAppError(err) {
		appErr := err.(*errors.AppError)
		h.logger.WithError(err).Error("Application error", "error_type", appErr.Type, "context", appErr.Context)

		status := appErr.HTTPStatus()
		message := appErr.Message
		if appErr.Type == errors.ErrorTypeCircuitBreaker {
			message = "Service temporarily unavailable"
		}

		h.writeError(w, status, message)
		return
	}

	h.logger.WithError(err).Error("Unexpected error")
	h.writeError(w, http.StatusInternalServerError, defaultMessage)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.WithError(err).Error("Failed to encode JSON response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	response := responses.ErrorResponse{
		Error: message,
	}

	h.writeJSON(w, status, response)
}
