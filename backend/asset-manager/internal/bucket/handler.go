package bucket

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/shared"
)

type BucketHandler struct {
	Service BucketService
}

func (h *BucketHandler) ListBuckets(w http.ResponseWriter, r *http.Request) {
	limit := 10
	if val := r.URL.Query().Get("limit"); val != "" {
		if l, err := strconv.Atoi(val); err == nil && l > 0 {
			limit = l
		}
	}

	var scanKey map[string]types.AttributeValue
	if token := r.URL.Query().Get("nextKey"); token != "" {
		decoded, err := shared.DecodeLastEvaluatedKey(token)
		if err != nil {
			shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
			return
		}
		scanKey, err = shared.ToDynamoKey(decoded)
		if err != nil {
			shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
			return
		}
	}

	page, err := h.Service.ListBuckets(r.Context(), limit, scanKey)
	if err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, "Could not list buckets")
		return
	}

	shared.JSON(w, http.StatusOK, BuildPaginatedResponse(page), "")
}

func BuildPaginatedResponse(page *BucketPage) map[string]interface{} {
	return map[string]interface{}{
		"items":   page.Items,
		"nextKey": shared.EncodeLastEvaluatedKey(page.LastEvaluatedKey),
	}
}

func (h *BucketHandler) GetBucket(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
		return
	}
	b, err := h.Service.GetBucketByID(r.Context(), id)
	if err != nil {
		shared.JSON(w, http.StatusNotFound, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusOK, b, "")
}

func (h *BucketHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	var b Bucket
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
		return
	}
	created, err := h.Service.CreateBucket(r.Context(), &b)
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusCreated, created, "")
}

func (h *BucketHandler) PatchBucket(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
		return
	}
	var patch map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
		return
	}
	if err := h.Service.PatchBucket(r.Context(), id, patch); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}

func (h *BucketHandler) AddAssetToBucket(w http.ResponseWriter, r *http.Request) {
	bucketID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
		return
	}
	var payload struct {
		AssetID int `json:"assetId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.AssetID == 0 {
		shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
		return
	}
	if err := h.Service.AddAssetToBucket(r.Context(), bucketID, payload.AssetID); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}

func (h *BucketHandler) RemoveAssetFromBucket(w http.ResponseWriter, r *http.Request) {
	bucketID, err := strconv.Atoi(mux.Vars(r)["id"])
	assetID, err2 := strconv.Atoi(mux.Vars(r)["assetId"])
	if err != nil || err2 != nil {
		shared.JSON(w, http.StatusBadRequest, nil, constants.ErrInvalidRequestBody)
		return
	}
	if err := h.Service.RemoveAssetFromBucket(r.Context(), bucketID, assetID); err != nil {
		shared.JSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	shared.JSON(w, http.StatusNoContent, nil, "")
}
