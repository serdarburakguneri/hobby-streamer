package asset

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/shared"
)

type AssetService interface {
	ListAssets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*AssetPage, error)
	GetAssetByID(ctx context.Context, id int) (*Asset, error)
	CreateAsset(ctx context.Context, a *Asset) (*Asset, error)
	PatchAsset(ctx context.Context, id int, patch map[string]interface{}) error
	AddImage(ctx context.Context, id int, img *Image) error
	DeleteImage(ctx context.Context, id int, filename string) error
	AddVideo(ctx context.Context, id int, label string, video *Video) error
	DeleteVideo(ctx context.Context, id int, label string) error
}

type AssetHandler struct {
	Service AssetService
}

func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if val := r.URL.Query().Get("limit"); val != "" {
		if l, err := strconv.Atoi(val); err == nil && l > 0 {
			limit = l
		}
	}

	var flatKey map[string]types.AttributeValue
	if token := r.URL.Query().Get("nextKey"); token != "" {
		decoded, err := shared.DecodeLastEvaluatedKey(token)
		if err != nil {
			shared.JSON(w, http.StatusBadRequest, nil, "Invalid nextKey")
			return
		}
		flatKey, err = shared.ToDynamoKey(decoded)
		if err != nil {
			shared.JSON(w, http.StatusBadRequest, nil, "Invalid scan key")
			return
		}
	} else {
		flatKey = make(map[string]types.AttributeValue)
	}

	page, err := h.Service.ListAssets(r.Context(), limit, flatKey)
	if err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, "Could not list assets")
		return
	}

	shared.JSON(w, http.StatusOK, BuildPaginatedResponse(page), "")
}

func (h *AssetHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid ID")
		return
	}
	asset, err := h.Service.GetAssetByID(r.Context(), id)
	if err != nil {
		shared.JSON(w, http.StatusNotFound, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusOK, asset, "")
}

func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var dto AssetCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid body")
		return
	}

	a := &Asset{
		Title:       dto.Title,
		Description: dto.Description,
		Type:        dto.Type,
		Category:    dto.Category,
		Genres:      dto.Genres,
		Tags:        dto.Tags,
		Credits:     dto.Credits,
		PublishRule: dto.PublishRule,
		Attributes:  dto.Attributes,
		Status:      AssetStatusPending,
	}

	created, err := h.Service.CreateAsset(r.Context(), a)
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusCreated, created, "")
}

func (h *AssetHandler) PatchAsset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid ID")
		return
	}
	var patch map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid patch")
		return
	}
	if err := h.Service.PatchAsset(r.Context(), id, patch); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}

func (h *AssetHandler) GetPublishRule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid ID")
		return
	}
	asset, err := h.Service.GetAssetByID(r.Context(), id)
	if err != nil {
		shared.JSON(w, http.StatusNotFound, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusOK, asset.PublishRule, "")
}

func (h *AssetHandler) PatchPublishRule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid ID")
		return
	}
	var patch map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid patch")
		return
	}
	if err := h.Service.PatchAsset(r.Context(), id, map[string]interface{}{"publishRule": patch}); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}

func (h *AssetHandler) SetVideoVariant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	label := mux.Vars(r)["label"]
	if err != nil || label == "" {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid path params")
		return
	}
	var v Video
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid video body")
		return
	}
	if err := h.Service.AddVideo(r.Context(), id, label, &v); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}

func (h *AssetHandler) DeleteVideoVariant(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	label := mux.Vars(r)["label"]
	if err != nil || label == "" {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid path params")
		return
	}
	if err := h.Service.DeleteVideo(r.Context(), id, label); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}

func (h *AssetHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid ID")
		return
	}
	var image Image
	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid image")
		return
	}
	if err := h.Service.AddImage(r.Context(), id, &image); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}

func (h *AssetHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	filename := mux.Vars(r)["filename"]
	if err != nil || filename == "" {
		shared.JSON(w, http.StatusBadRequest, nil, "Invalid path params")
		return
	}
	if err := h.Service.DeleteImage(r.Context(), id, filename); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}
