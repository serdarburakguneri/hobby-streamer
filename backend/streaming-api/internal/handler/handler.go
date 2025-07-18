package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
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

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"buckets": buckets,
		"count":   len(buckets),
	})
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

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"assets": assets,
		"count":  len(assets),
	})
}

func (h *Handler) GetAssets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	assets, err := h.service.GetAssets(ctx)
	if err != nil {
		h.handleError(w, err, "Failed to get assets")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"assets": assets,
		"count":  len(assets),
	})
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
	h.writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "streaming-api",
	})
}

func (h *Handler) handleError(w http.ResponseWriter, err error, defaultMessage string) {
	if errors.IsAppError(err) {
		appErr := err.(*errors.AppError)
		h.logger.WithError(err).Error("Application error", "error_type", appErr.Type, "context", appErr.Context)

		switch appErr.Type {
		case errors.ErrorTypeNotFound:
			h.writeError(w, http.StatusNotFound, appErr.Message)
			return
		case errors.ErrorTypeValidation:
			h.writeError(w, http.StatusBadRequest, appErr.Message)
			return
		case errors.ErrorTypeUnauthorized:
			h.writeError(w, http.StatusUnauthorized, appErr.Message)
			return
		case errors.ErrorTypeForbidden:
			h.writeError(w, http.StatusForbidden, appErr.Message)
			return
		case errors.ErrorTypeConflict:
			h.writeError(w, http.StatusConflict, appErr.Message)
			return
		case errors.ErrorTypeTransient:
			h.writeError(w, http.StatusServiceUnavailable, appErr.Message)
			return
		case errors.ErrorTypeTimeout:
			h.writeError(w, http.StatusGatewayTimeout, appErr.Message)
			return
		case errors.ErrorTypeCircuitBreaker:
			h.writeError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		case errors.ErrorTypeExternal:
			h.writeError(w, http.StatusBadGateway, appErr.Message)
			return
		default:
			h.writeError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
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
	h.writeJSON(w, status, map[string]string{
		"error": message,
	})
}
