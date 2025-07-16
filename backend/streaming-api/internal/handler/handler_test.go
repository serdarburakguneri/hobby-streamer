package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
	"github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/model"
)

type mockService struct {
	buckets      []model.Bucket
	assets       []model.Asset
	bucket       *model.Bucket
	asset        *model.Asset
	assetsInBuck []model.Asset
	err          error
}

func (m *mockService) GetBuckets(_ context.Context) ([]model.Bucket, error) { return m.buckets, m.err }
func (m *mockService) GetBucket(_ context.Context, key string) (*model.Bucket, error) {
	return m.bucket, m.err
}
func (m *mockService) GetAssets(_ context.Context) ([]model.Asset, error) { return m.assets, m.err }
func (m *mockService) GetAsset(_ context.Context, slug string) (*model.Asset, error) {
	return m.asset, m.err
}
func (m *mockService) GetAssetsInBucket(_ context.Context, key string) ([]model.Asset, error) {
	return m.assetsInBuck, m.err
}

func newTestHandler(s *mockService) *Handler {
	h := &Handler{
		service: s,
		logger:  logger.Get().WithService("test-logger"),
	}
	return h
}

func TestGetBuckets(t *testing.T) {
	svc := &mockService{buckets: []model.Bucket{{ID: "1", Key: "k", Name: "n"}}}
	h := newTestHandler(svc)
	r := h.SetupRoutes()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/buckets", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if int(resp["count"].(float64)) != 1 {
		t.Errorf("expected count 1, got %v", resp["count"])
	}
}

func TestGetBucket_NotFound(t *testing.T) {
	svc := &mockService{bucket: nil}
	h := newTestHandler(svc)
	r := h.SetupRoutes()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/buckets/doesnotexist", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetAssets(t *testing.T) {
	svc := &mockService{assets: []model.Asset{{ID: "1", Slug: "s"}}}
	h := newTestHandler(svc)
	r := h.SetupRoutes()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/assets", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if int(resp["count"].(float64)) != 1 {
		t.Errorf("expected count 1, got %v", resp["count"])
	}
}

func TestGetAsset_NotFound(t *testing.T) {
	svc := &mockService{asset: nil}
	h := newTestHandler(svc)
	r := h.SetupRoutes()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/assets/doesnotexist", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetAssetsInBucket(t *testing.T) {
	svc := &mockService{assetsInBuck: []model.Asset{{ID: "1", Slug: "s"}}}
	h := newTestHandler(svc)
	r := h.SetupRoutes()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/buckets/bucket1/assets", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if int(resp["count"].(float64)) != 1 {
		t.Errorf("expected count 1, got %v", resp["count"])
	}
}

func TestHealthCheck(t *testing.T) {
	svc := &mockService{}
	h := newTestHandler(svc)
	r := h.SetupRoutes()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["status"] != "healthy" {
		t.Errorf("expected healthy, got %v", resp["status"])
	}
}

func TestGetBuckets_Error(t *testing.T) {
	svc := &mockService{err: errors.New("fail")}
	h := newTestHandler(svc)
	r := h.SetupRoutes()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/buckets", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
