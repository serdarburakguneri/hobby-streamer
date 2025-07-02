package asset

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
)

type MockService struct {
	ListAssetsFunc   func(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*AssetPage, error)
	GetAssetByIDFunc func(ctx context.Context, id int) (*Asset, error)
	CreateAssetFunc  func(ctx context.Context, a *Asset) (*Asset, error)
	PatchAssetFunc   func(ctx context.Context, id int, patch map[string]interface{}) error
	AddImageFunc     func(ctx context.Context, id int, img *Image) error
	DeleteImageFunc  func(ctx context.Context, id int, filename string) error
	AddVideoFunc     func(ctx context.Context, id int, label string, video *Video) error
	DeleteVideoFunc  func(ctx context.Context, id int, label string) error
}

func (m *MockService) ListAssets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*AssetPage, error) {
	if m.ListAssetsFunc != nil {
		return m.ListAssetsFunc(ctx, limit, lastKey)
	}
	return &AssetPage{}, nil
}
func (m *MockService) GetAssetByID(ctx context.Context, id int) (*Asset, error) {
	if m.GetAssetByIDFunc != nil {
		return m.GetAssetByIDFunc(ctx, id)
	}
	return &Asset{ID: id}, nil
}
func (m *MockService) CreateAsset(ctx context.Context, a *Asset) (*Asset, error) {
	if m.CreateAssetFunc != nil {
		return m.CreateAssetFunc(ctx, a)
	}
	a.ID = 1
	return a, nil
}
func (m *MockService) PatchAsset(ctx context.Context, id int, patch map[string]interface{}) error {
	if m.PatchAssetFunc != nil {
		return m.PatchAssetFunc(ctx, id, patch)
	}
	return nil
}
func (m *MockService) AddImage(ctx context.Context, id int, img *Image) error {
	if m.AddImageFunc != nil {
		return m.AddImageFunc(ctx, id, img)
	}
	return nil
}
func (m *MockService) DeleteImage(ctx context.Context, id int, filename string) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, id, filename)
	}
	return nil
}
func (m *MockService) AddVideo(ctx context.Context, id int, label string, video *Video) error {
	if m.AddVideoFunc != nil {
		return m.AddVideoFunc(ctx, id, label, video)
	}
	return nil
}
func (m *MockService) DeleteVideo(ctx context.Context, id int, label string) error {
	if m.DeleteVideoFunc != nil {
		return m.DeleteVideoFunc(ctx, id, label)
	}
	return nil
}

func TestListAssetsHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{}}
	req := httptest.NewRequest(http.MethodGet, "/assets", nil)
	w := httptest.NewRecorder()
	h.ListAssets(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestGetAssetHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{
		GetAssetByIDFunc: func(ctx context.Context, id int) (*Asset, error) {
			return &Asset{ID: id, Title: ptr("Test Asset")}, nil
		},
	}}
	req := httptest.NewRequest(http.MethodGet, "/assets/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	h.GetAsset(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestCreateAssetHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{}}
	payload := AssetCreateDTO{Title: ptr("Test Asset"), Type: ptr("video")}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/assets", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateAsset(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Result().StatusCode)
	}
}

func TestPatchAssetHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{
		PatchAssetFunc: func(ctx context.Context, id int, patch map[string]interface{}) error {
			return nil
		},
	}}
	patch := map[string]interface{}{"title": "new title"}
	body, _ := json.Marshal(patch)
	req := httptest.NewRequest(http.MethodPatch, "/assets/1", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	h.PatchAsset(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}

func TestGetPublishRuleHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{
		GetAssetByIDFunc: func(ctx context.Context, id int) (*Asset, error) {
			return &Asset{ID: id, PublishRule: &PublishRule{}}, nil
		},
	}}
	req := httptest.NewRequest(http.MethodGet, "/assets/1/publishRule", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	h.GetPublishRule(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestPatchPublishRuleHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{
		PatchAssetFunc: func(ctx context.Context, id int, patch map[string]interface{}) error {
			return nil
		},
	}}
	patch := map[string]interface{}{"publishRule": map[string]interface{}{"regionLock": []string{"US"}}}
	body, _ := json.Marshal(patch)
	req := httptest.NewRequest(http.MethodPatch, "/assets/1/publishRule", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	h.PatchPublishRule(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}

func TestSetVideoVariantHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{
		AddVideoFunc: func(ctx context.Context, id int, label string, video *Video) error {
			return nil
		},
	}}
	video := Video{FileName: "test.mp4"}
	body, _ := json.Marshal(video)
	req := httptest.NewRequest(http.MethodPost, "/assets/1/videos/hls", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "label": "hls"})
	w := httptest.NewRecorder()
	h.SetVideoVariant(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}

func TestDeleteVideoVariantHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{
		DeleteVideoFunc: func(ctx context.Context, id int, label string) error {
			return nil
		},
	}}
	req := httptest.NewRequest(http.MethodDelete, "/assets/1/videos/hls", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1", "label": "hls"})
	w := httptest.NewRecorder()
	h.DeleteVideoVariant(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}

func TestAddImageHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{
		AddImageFunc: func(ctx context.Context, id int, img *Image) error {
			return nil
		},
	}}
	img := Image{FileName: "foo.jpg"}
	body, _ := json.Marshal(img)
	req := httptest.NewRequest(http.MethodPost, "/assets/1/images", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	h.AddImage(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}

func TestDeleteImageHandler(t *testing.T) {
	h := &AssetHandler{Service: &MockService{
		DeleteImageFunc: func(ctx context.Context, id int, filename string) error {
			return nil
		},
	}}
	req := httptest.NewRequest(http.MethodDelete, "/assets/1/images/foo.jpg", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1", "filename": "foo.jpg"})
	w := httptest.NewRecorder()
	h.DeleteImage(w, req)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Result().StatusCode)
	}
}

// Helper to get pointer to string
func ptr(s string) *string { return &s }
