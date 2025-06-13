package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/asset"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/pkg/response"
)

type Handler struct {
	Repo *asset.Repository
}

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, nil, "Invalid ID")
		return
	}

	a, err := h.Repo.GetAssetByID(r.Context(), id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, nil, "Asset not found")
		return
	}
	response.JSON(w, http.StatusOK, a, "")
}

func (h *Handler) SaveAsset(w http.ResponseWriter, r *http.Request) {
	var a asset.Asset
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		response.JSON(w, http.StatusBadRequest, nil, "Invalid body")
		return
	}

	if a.Id == 0 {
		response.JSON(w, http.StatusBadRequest, nil, "Missing ID")
		return
	}

	if err := h.Repo.SaveAsset(r.Context(), &a); err != nil {
		response.JSON(w, http.StatusInternalServerError, nil, "Save failed")
		return
	}

	response.JSON(w, http.StatusOK, a, "")
}

func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if val := r.URL.Query().Get("limit"); val != "" {
		if l, err := strconv.Atoi(val); err == nil && l > 0 {
			limit = l
		}
	}

	var scanKey map[string]map[string]string
	if token := r.URL.Query().Get("nextKey"); token != "" {
		decoded, err := asset.DecodeLastEvaluatedKey(token)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil, "Invalid nextKey")
			return
		}
		scanKey, err = asset.ToDynamoKey(decoded)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil, "Invalid scan key")
			return
		}
	}

	page, err := h.Repo.ListAssets(r.Context(), limit, scanKey)
	if err != nil {
		log.Printf("ListAssets failed: %v", err)
		response.JSON(w, http.StatusInternalServerError, nil, "Could not list assets")
		return
	}

	response.JSON(w, http.StatusOK, asset.BuildPaginatedResponse(page), "")
}