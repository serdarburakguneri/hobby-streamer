package asset

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/config"
)

type mockRepository struct {
	assets map[string]*Asset
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		assets: make(map[string]*Asset),
	}
}

func (m *mockRepository) GetAssetByID(ctx context.Context, id string) (*Asset, error) {
	if asset, exists := m.assets[id]; exists {
		return asset, nil
	}
	return nil, errors.New("asset not found")
}

func (m *mockRepository) GetAssetBySlug(ctx context.Context, slug string) (*Asset, error) {
	for _, asset := range m.assets {
		if asset.Slug == slug {
			return asset, nil
		}
	}
	return nil, errors.New("asset not found")
}

func (m *mockRepository) ListAssets(ctx context.Context, limit int, lastKey map[string]interface{}) (*AssetPage, error) {
	var assets []Asset
	count := 0
	for _, asset := range m.assets {
		if count >= limit {
			break
		}
		assets = append(assets, *asset)
		count++
	}
	return &AssetPage{
		Items: assets,
	}, nil
}

func (m *mockRepository) SaveAsset(ctx context.Context, asset *Asset) error {
	// Simulate the timestamp setting behavior from the real repository
	now := time.Now().UTC()
	if asset.CreatedAt.IsZero() {
		asset.CreatedAt = now
	}
	asset.UpdatedAt = now

	m.assets[asset.ID] = asset
	return nil
}

func (m *mockRepository) PatchAsset(ctx context.Context, id string, patch map[string]interface{}) error {
	if _, exists := m.assets[id]; !exists {
		return errors.New("asset not found")
	}
	return nil
}

func (m *mockRepository) DeleteAsset(ctx context.Context, id string) error {
	if _, exists := m.assets[id]; !exists {
		return errors.New("asset not found")
	}
	delete(m.assets, id)
	return nil
}

func (m *mockRepository) GetParent(ctx context.Context, childID string) (*Asset, error) {
	return nil, nil
}

func (m *mockRepository) SearchAssets(ctx context.Context, query string, limit int, lastKey map[string]interface{}) (*AssetPage, error) {
	var assets []Asset
	count := 0
	for _, asset := range m.assets {
		if count >= limit {
			break
		}
		if asset.Title != nil && *asset.Title == query {
			assets = append(assets, *asset)
			count++
		}
	}
	return &AssetPage{
		Items: assets,
	}, nil
}

func (m *mockRepository) GetChildren(ctx context.Context, parentID string) ([]Asset, error) {
	return []Asset{}, nil
}

func (m *mockRepository) GetAssetsByTypeAndGenre(ctx context.Context, assetType, genre string) ([]Asset, error) {
	return []Asset{}, nil
}

func (m *mockRepository) GetAssetsByIDs(ctx context.Context, ids []string) ([]Asset, error) {
	var assets []Asset
	for _, id := range ids {
		if asset, exists := m.assets[id]; exists {
			assets = append(assets, *asset)
		}
	}
	return assets, nil
}

func dummyDynamicConfig() *config.DynamicConfig {
	base := &config.BaseConfig{
		Environment: config.Test,
		Service:     "asset-manager",
		Components:  map[string]interface{}{},
	}
	return config.NewDynamicConfig(base)
}

func TestService_CreateAsset(t *testing.T) {
	tests := []struct {
		name    string
		asset   *Asset
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid asset creation",
			asset: &Asset{
				Title:       stringPtr("Test Movie"),
				Description: stringPtr("A test movie"),
				Type:        stringPtr(AssetTypeMovie),
				Genre:       stringPtr("action"),
			},
			wantErr: false,
		},
		{
			name: "asset with ID should fail",
			asset: &Asset{
				ID:    "123",
				Title: stringPtr("Test Movie"),
				Type:  stringPtr(AssetTypeMovie),
			},
			wantErr: true,
			errMsg:  "validation: asset ID should not be set during creation",
		},
		{
			name: "invalid type should fail",
			asset: &Asset{
				Title: stringPtr("Test Movie"),
				Type:  stringPtr("invalid"),
			},
			wantErr: true,
			errMsg:  "validation: invalid type value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			service := NewService(repo, dummyDynamicConfig())
			ctx := context.Background()

			result, err := service.CreateAsset(ctx, tt.asset)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateAsset() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("CreateAsset() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateAsset() unexpected error = %v", err)
				return
			}

			if result.ID == "" {
				t.Errorf("CreateAsset() expected ID to be set")
			}

			// Status is computed based on PublishRule, so we don't check it here

			if result.CreatedAt.IsZero() {
				t.Errorf("CreateAsset() expected CreatedAt to be set")
			}

			if result.UpdatedAt.IsZero() {
				t.Errorf("CreateAsset() expected UpdatedAt to be set")
			}
		})
	}
}

func TestService_UpdateAsset(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		asset   *Asset
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid asset update",
			id:   "123",
			asset: &Asset{
				ID:          "123",
				Title:       stringPtr("Updated Movie"),
				Description: stringPtr("An updated movie"),
				Type:        stringPtr(AssetTypeMovie),
			},
			wantErr: false,
		},
		{
			name: "ID mismatch should fail",
			id:   "123",
			asset: &Asset{
				ID:    "456",
				Title: stringPtr("Updated Movie"),
				Type:  stringPtr(AssetTypeMovie),
			},
			wantErr: true,
			errMsg:  "validation: asset ID mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepository()
			service := NewService(repo, dummyDynamicConfig())
			ctx := context.Background()

			err := service.UpdateAsset(ctx, tt.id, tt.asset)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateAsset() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("UpdateAsset() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateAsset() unexpected error = %v", err)
			}
		})
	}
}

func TestService_GetAssetByID(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo, dummyDynamicConfig())
	ctx := context.Background()

	// Create a test asset
	testAsset := &Asset{
		ID:          "123",
		Title:       stringPtr("Test Movie"),
		Description: stringPtr("A test movie"),
		Type:        stringPtr(AssetTypeMovie),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.assets["123"] = testAsset

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "existing asset",
			id:      "123",
			wantErr: false,
		},
		{
			name:    "non-existing asset",
			id:      "456",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetAssetByID(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAssetByID() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetAssetByID() unexpected error = %v", err)
				return
			}

			if result.ID != tt.id {
				t.Errorf("GetAssetByID() got ID = %v, want %v", result.ID, tt.id)
			}
		})
	}
}

func TestService_ListAssets(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo, dummyDynamicConfig())
	ctx := context.Background()

	// Create test assets
	assets := []*Asset{
		{
			ID:    "1",
			Title: stringPtr("Movie 1"),
			Type:  stringPtr(AssetTypeMovie),
		},
		{
			ID:    "2",
			Title: stringPtr("Movie 2"),
			Type:  stringPtr(AssetTypeMovie),
		},
		{
			ID:    "3",
			Title: stringPtr("Movie 3"),
			Type:  stringPtr(AssetTypeMovie),
		},
	}

	for _, asset := range assets {
		repo.assets[asset.ID] = asset
	}

	tests := []struct {
		name     string
		limit    int
		expected int
	}{
		{
			name:     "limit 2",
			limit:    2,
			expected: 2,
		},
		{
			name:     "limit 5",
			limit:    5,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ListAssets(ctx, tt.limit, nil)

			if err != nil {
				t.Errorf("ListAssets() unexpected error = %v", err)
				return
			}

			if len(result.Items) != tt.expected {
				t.Errorf("ListAssets() got %v items, want %v", len(result.Items), tt.expected)
			}
		})
	}
}

func TestService_DeleteAsset(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo, dummyDynamicConfig())
	ctx := context.Background()

	// Create a test asset
	testAsset := &Asset{
		ID:    "123",
		Title: stringPtr("Test Movie"),
		Type:  stringPtr(AssetTypeMovie),
	}
	repo.assets["123"] = testAsset

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "existing asset",
			id:      "123",
			wantErr: false,
		},
		{
			name:    "non-existing asset",
			id:      "456",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteAsset(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteAsset() expected error but got none")
				}
				return
			}

			// Verify asset was deleted
			if _, exists := repo.assets[tt.id]; exists {
				t.Errorf("DeleteAsset() asset still exists after deletion")
			}
		})
	}
}

func TestService_AddImage(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo, dummyDynamicConfig())
	ctx := context.Background()

	// Create a test asset
	testAsset := &Asset{
		ID:     "123",
		Title:  stringPtr("Test Movie"),
		Type:   stringPtr(AssetTypeMovie),
		Images: []Image{},
	}
	repo.assets["123"] = testAsset

	tests := []struct {
		name    string
		id      string
		image   *Image
		wantErr bool
	}{
		{
			name: "add new image",
			id:   "123",
			image: &Image{
				FileName: "poster.jpg",
				URL:      "https://example.com/poster.jpg",
				Width:    1920,
				Height:   1080,
			},
			wantErr: false,
		},
		{
			name: "add duplicate image should fail",
			id:   "123",
			image: &Image{
				FileName: "poster.jpg",
				URL:      "https://example.com/poster2.jpg",
				Width:    1920,
				Height:   1080,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddImage(ctx, tt.id, tt.image)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddImage() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("AddImage() unexpected error = %v", err)
				return
			}

			// Verify image was added
			asset, _ := repo.GetAssetByID(ctx, tt.id)
			found := false
			for _, img := range asset.Images {
				if img.FileName == tt.image.FileName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("AddImage() image was not added")
			}
		})
	}
}

func TestService_AddVideo(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo, dummyDynamicConfig())
	ctx := context.Background()

	// Create a test asset
	testAsset := &Asset{
		ID:     "123",
		Title:  stringPtr("Test Movie"),
		Type:   stringPtr(AssetTypeMovie),
		Videos: []Video{},
	}
	repo.assets["123"] = testAsset

	tests := []struct {
		name    string
		id      string
		video   *Video
		wantErr bool
	}{
		{
			name: "add new video",
			id:   "123",
			video: &Video{
				Type:   VideoTypeMain,
				Format: "mp4",
				StorageLocation: S3Object{
					Bucket: "test-bucket",
					Key:    "video.mp4",
					URL:    "https://example.com/video.mp4",
				},
				Width:       1920,
				Height:      1080,
				Duration:    120.5,
				ContentType: "video/mp4",
			},
			wantErr: false,
		},
		{
			name: "update existing video",
			id:   "123",
			video: &Video{
				Type:   VideoTypeMain,
				Format: "mp4",
				StorageLocation: S3Object{
					Bucket: "test-bucket",
					Key:    "video2.mp4",
					URL:    "https://example.com/video2.mp4",
				},
				Width:       1920,
				Height:      1080,
				Duration:    125.0,
				ContentType: "video/mp4",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddVideo(ctx, tt.id, tt.video)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddVideo() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("AddVideo() unexpected error = %v", err)
				return
			}

			// Verify video was added/updated
			asset, _ := repo.GetAssetByID(ctx, tt.id)
			found := false
			for _, video := range asset.Videos {
				if video.Type == tt.video.Type && video.Format == tt.video.Format && video.StorageLocation.URL == tt.video.StorageLocation.URL {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("AddVideo() video was not added")
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func (m *mockRepository) Target() url.URL {
	return url.URL{}
}
