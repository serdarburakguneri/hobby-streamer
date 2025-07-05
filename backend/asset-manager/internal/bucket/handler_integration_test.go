package bucket_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
	"github.com/serdarburakguneri/hobby-streamer/services/asset-manager/internal/bucket"
)

type MockService struct {
	ListBucketsFunc           func(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*bucket.BucketPage, error)
	GetBucketByIDFunc         func(ctx context.Context, id int) (*bucket.Bucket, error)
	CreateBucketFunc          func(ctx context.Context, b *bucket.Bucket) (*bucket.Bucket, error)
	PatchBucketFunc           func(ctx context.Context, id int, patch map[string]interface{}) error
	AddAssetToBucketFunc      func(ctx context.Context, id int, assetID int) error
	RemoveAssetFromBucketFunc func(ctx context.Context, id int, assetID int) error
	UpdateBucketFunc          func(ctx context.Context, id int, b *bucket.Bucket) error
}

func (m *MockService) ListBuckets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*bucket.BucketPage, error) {
	if m.ListBucketsFunc != nil {
		return m.ListBucketsFunc(ctx, limit, lastKey)
	}
	return &bucket.BucketPage{}, nil
}
func (m *MockService) GetBucketByID(ctx context.Context, id int) (*bucket.Bucket, error) {
	if m.GetBucketByIDFunc != nil {
		return m.GetBucketByIDFunc(ctx, id)
	}
	return &bucket.Bucket{ID: strconv.Itoa(id)}, nil
}
func (m *MockService) CreateBucket(ctx context.Context, b *bucket.Bucket) (*bucket.Bucket, error) {
	if m.CreateBucketFunc != nil {
		return m.CreateBucketFunc(ctx, b)
	}
	b.ID = "1"
	return b, nil
}
func (m *MockService) PatchBucket(ctx context.Context, id int, patch map[string]interface{}) error {
	if m.PatchBucketFunc != nil {
		return m.PatchBucketFunc(ctx, id, patch)
	}
	return nil
}
func (m *MockService) AddAssetToBucket(ctx context.Context, id int, assetID int) error {
	if m.AddAssetToBucketFunc != nil {
		return m.AddAssetToBucketFunc(ctx, id, assetID)
	}
	return nil
}
func (m *MockService) RemoveAssetFromBucket(ctx context.Context, id int, assetID int) error {
	if m.RemoveAssetFromBucketFunc != nil {
		return m.RemoveAssetFromBucketFunc(ctx, id, assetID)
	}
	return nil
}
func (m *MockService) UpdateBucket(ctx context.Context, id int, b *bucket.Bucket) error {
	if m.UpdateBucketFunc != nil {
		return m.UpdateBucketFunc(ctx, id, b)
	}
	return nil
}

func TestListBucketsHandler(t *testing.T) {
	h := &bucket.BucketHandler{Service: &MockService{}}
	req := httptest.NewRequest(http.MethodGet, "/buckets", nil)
	w := httptest.NewRecorder()
	h.ListBuckets(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestGetBucketHandler(t *testing.T) {
	h := &bucket.BucketHandler{Service: &MockService{
		GetBucketByIDFunc: func(ctx context.Context, id int) (*bucket.Bucket, error) {
			return &bucket.Bucket{ID: strconv.Itoa(id), Name: "Test Bucket"}, nil
		},
	}}
	req := httptest.NewRequest(http.MethodGet, "/buckets/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	h.GetBucket(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestCreateBucketHandler(t *testing.T) {
	h := &bucket.BucketHandler{Service: &MockService{}}
	payload := bucket.Bucket{Name: "Test Bucket"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/buckets", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateBucket(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Result().StatusCode)
	}
}

func TestPatchBucketHandler(t *testing.T) {
	h := &bucket.BucketHandler{Service: &MockService{
		PatchBucketFunc: func(ctx context.Context, id int, patch map[string]interface{}) error {
			return nil
		},
	}}
	patch := map[string]interface{}{"name": "new name"}
	body, _ := json.Marshal(patch)
	req := httptest.NewRequest(http.MethodPatch, "/buckets/1", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	h.PatchBucket(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}

func TestAddAssetToBucketHandler(t *testing.T) {
	h := &bucket.BucketHandler{Service: &MockService{
		AddAssetToBucketFunc: func(ctx context.Context, id int, assetID int) error {
			return nil
		},
	}}
	payload := struct {
		AssetID int `json:"assetId"`
	}{AssetID: 42}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/buckets/1/assets", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	h.AddAssetToBucket(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}

func TestRemoveAssetFromBucketHandler(t *testing.T) {
	h := &bucket.BucketHandler{Service: &MockService{
		RemoveAssetFromBucketFunc: func(ctx context.Context, id int, assetID int) error {
			return nil
		},
	}}
	req := httptest.NewRequest(http.MethodDelete, "/buckets/1/assets/42", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1", "assetId": "42"})
	w := httptest.NewRecorder()
	h.RemoveAssetFromBucket(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}
