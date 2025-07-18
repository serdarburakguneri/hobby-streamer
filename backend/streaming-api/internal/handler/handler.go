package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/model"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/service"
)

type Handler struct {
	service service.ServiceInterface
	logger  *logger.Logger
}

func NewHandler(service service.ServiceInterface) *Handler {
	return &Handler{
		service: service,
		logger:  logger.Get().WithService("streaming-handler"),
	}
}

func (h *Handler) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

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

	buckets, err := h.service.GetBuckets(ctx)
	if err != nil {
		h.handleError(w, err, "Failed to get buckets")
		return
	}

	response := model.BucketsResponse{
		Buckets: buckets,
		Count:   len(buckets),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetBucket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	key := vars["key"]

	bucket, err := h.service.GetBucket(ctx, key)
	if err != nil {
		h.handleError(w, err, "Failed to get bucket")
		return
	}

	if bucket == nil {
		h.writeError(w, http.StatusNotFound, "Bucket not found")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=1800")
	h.writeJSON(w, http.StatusOK, bucket)
}

func (h *Handler) GetAssetsInBucket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	key := vars["key"]

	assets, err := h.service.GetAssetsInBucket(ctx, key)
	if err != nil {
		h.handleError(w, err, "Failed to get assets in bucket")
		return
	}

	response := model.AssetsResponse{
		Assets: assets,
		Count:  len(assets),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetAssets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	assets, err := h.service.GetAssets(ctx)
	if err != nil {
		h.handleError(w, err, "Failed to get assets")
		return
	}

	response := model.AssetsResponse{
		Assets: assets,
		Count:  len(assets),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	slug := vars["slug"]

	asset, err := h.service.GetAsset(ctx, slug)
	if err != nil {
		h.handleError(w, err, "Failed to get asset")
		return
	}

	if asset == nil {
		h.writeError(w, http.StatusNotFound, "Asset not found")
		return
	}

	h.writeJSON(w, http.StatusOK, asset)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := model.HealthResponse{
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
	response := model.ErrorResponse{
		Error: message,
	}

	h.writeJSON(w, status, response)
}
