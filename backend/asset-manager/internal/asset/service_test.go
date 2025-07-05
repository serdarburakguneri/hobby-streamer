package asset

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type MockRepository struct {
	SaveAssetFunc    func(ctx context.Context, a *Asset) error
	GetAssetByIDFunc func(ctx context.Context, id int) (*Asset, error)
}

func (m *MockRepository) SaveAsset(ctx context.Context, a *Asset) error {
	if m.SaveAssetFunc != nil {
		return m.SaveAssetFunc(ctx, a)
	}
	return nil
}
func (m *MockRepository) GetAssetByID(ctx context.Context, id int) (*Asset, error) {
	if m.GetAssetByIDFunc != nil {
		return m.GetAssetByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) ListAssets(ctx context.Context, limit int, lastKey map[string]types.AttributeValue) (*AssetPage, error) {
	return nil, nil
}

func (m *MockRepository) PatchAsset(ctx context.Context, id int, patch map[string]interface{}) error {
	return nil
}

func TestCreateAsset_Success(t *testing.T) {
	mockRepo := &MockRepository{
		SaveAssetFunc: func(ctx context.Context, a *Asset) error {
			return nil
		},
	}
	svc := NewService(mockRepo)

	asset := &Asset{}
	got, err := svc.CreateAsset(context.Background(), asset)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID == 0 {
		t.Errorf("expected asset ID to be set, got 0")
	}
}

func TestCreateAsset_WithIDSet(t *testing.T) {
	mockRepo := &MockRepository{}
	svc := NewService(mockRepo)

	asset := &Asset{ID: 123}
	_, err := svc.CreateAsset(context.Background(), asset)
	if err == nil {
		t.Fatalf("expected error when ID is set, got nil")
	}
}

func TestAddImage_Duplicate(t *testing.T) {
	mockRepo := &MockRepository{
		GetAssetByIDFunc: func(ctx context.Context, id int) (*Asset, error) {
			return &Asset{
				Images: []Image{{FileName: "foo.jpg"}},
			}, nil
		},
	}
	svc := NewService(mockRepo)
	img := &Image{FileName: "foo.jpg"}
	err := svc.AddImage(context.Background(), 1, img)
	if err == nil || err.Error() != "image with same filename already exists" {
		t.Errorf("expected duplicate error, got %v", err)
	}
}
