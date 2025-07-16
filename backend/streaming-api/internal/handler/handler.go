package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

	router.Use(h.corsMiddleware)
	router.Use(h.loggingMiddleware)

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
		h.logger.WithError(err).Error("Failed to get buckets")
		h.writeError(w, http.StatusInternalServerError, "Failed to get buckets")
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
		h.logger.WithError(err).Error("Failed to get bucket", "key", key)
		h.writeError(w, http.StatusInternalServerError, "Failed to get bucket")
		return
	}

	if bucket == nil {
		h.writeError(w, http.StatusNotFound, "Bucket not found")
		return
	}

	h.writeJSON(w, http.StatusOK, bucket)
}

func (h *Handler) GetAssetsInBucket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	key := vars["key"]

	assets, err := h.service.GetAssetsInBucket(ctx, key)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get assets in bucket", "key", key)
		h.writeError(w, http.StatusInternalServerError, "Failed to get assets in bucket")
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
		h.logger.WithError(err).Error("Failed to get assets")
		h.writeError(w, http.StatusInternalServerError, "Failed to get assets")
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
		h.logger.WithError(err).Error("Failed to get asset", "slug", slug)
		h.writeError(w, http.StatusInternalServerError, "Failed to get asset")
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

func (h *Handler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(w, r)
	})
}
