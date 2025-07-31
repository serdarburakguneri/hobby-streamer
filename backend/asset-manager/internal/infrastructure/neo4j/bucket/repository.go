package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Repository struct {
	driver neo4j.Driver
}

func NewRepository(driver neo4j.Driver) *Repository {
	return &Repository{
		driver: driver,
	}
}

func (r *Repository) Create(ctx context.Context, bucket *entity.Bucket) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	log := logger.WithService("neo4j-bucket-repository").WithContext(ctx)

	metadataJSON := ""
	if bucket.Metadata() != nil {
		b, err := json.Marshal(bucket.Metadata())
		if err != nil {
			log.Error(fmt.Sprintf("Failed to marshal metadata: %v", err))
			return pkgerrors.NewInternalError("failed to marshal metadata", err)
		}
		metadataJSON = string(b)
	}

	description := ""
	if bucket.Description() != nil {
		description = bucket.Description().Value()
	}

	ownerID := ""
	if bucket.OwnerID() != nil {
		ownerID = bucket.OwnerID().Value()
	}

	status := ""
	if bucket.Status() != nil {
		status = bucket.Status().Value()
	}

	bucketType := ""
	if bucket.Type() != nil {
		bucketType = bucket.Type().Value()
	}

	params := map[string]interface{}{
		"id":          bucket.ID().Value(),
		"name":        bucket.Name().Value(),
		"description": description,
		"key":         bucket.Key().Value(),
		"ownerID":     ownerID,
		"status":      status,
		"type":        bucketType,
		"metadata":    metadataJSON,
		"createdAt":   bucket.CreatedAt().Value().Format(time.RFC3339),
		"updatedAt":   bucket.UpdatedAt().Value().Format(time.RFC3339),
	}

	log.Info(fmt.Sprintf("Creating bucket with params: %+v", params))
	result, err := session.Run(createQuery, params)
	if err != nil {
		log.Error(fmt.Sprintf("Create error: %v", err))
		return pkgerrors.NewInternalError("neo4j create error", err)
	}
	if !result.Next() {
		if result.Err() != nil {
			log.Error(fmt.Sprintf("Create result error: %v", result.Err()))
			return pkgerrors.NewInternalError("neo4j create result error", result.Err())
		}
		log.Warn("Failed to create bucket: no result returned")
		return pkgerrors.NewInternalError("failed to create bucket: no result returned", nil)
	}
	log.Info(fmt.Sprintf("Bucket created successfully: %v", params["id"]))
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id valueobjects.BucketID) (*entity.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(getByIDQuery, map[string]interface{}{"id": id.Value()})
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to get bucket", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, pkgerrors.NewNotFoundError("bucket not found", nil)
	}

	return RecordToBucket(record)
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*entity.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(getBySlugQuery, map[string]interface{}{"slug": slug})
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to get bucket by slug", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, pkgerrors.NewNotFoundError("bucket not found", nil)
	}

	return RecordToBucket(record)
}

func (r *Repository) Update(ctx context.Context, bucket *entity.Bucket) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	log := logger.WithService("neo4j-bucket-repository").WithContext(ctx)

	metadataJSON := ""
	if bucket.Metadata() != nil {
		b, err := json.Marshal(bucket.Metadata())
		if err != nil {
			log.Error(fmt.Sprintf("Failed to marshal metadata: %v", err))
			return pkgerrors.NewInternalError("failed to marshal metadata", err)
		}
		metadataJSON = string(b)
	}

	description := ""
	if bucket.Description() != nil {
		description = bucket.Description().Value()
	}

	ownerID := ""
	if bucket.OwnerID() != nil {
		ownerID = bucket.OwnerID().Value()
	}

	status := ""
	if bucket.Status() != nil {
		status = bucket.Status().Value()
	}

	bucketType := ""
	if bucket.Type() != nil {
		bucketType = bucket.Type().Value()
	}

	params := map[string]interface{}{
		"id":          bucket.ID().Value(),
		"name":        bucket.Name().Value(),
		"description": description,
		"ownerID":     ownerID,
		"status":      status,
		"type":        bucketType,
		"metadata":    metadataJSON,
		"updatedAt":   bucket.UpdatedAt().Value().Format(time.RFC3339),
	}

	log.Info(fmt.Sprintf("Updating bucket with params: %+v", params))
	result, err := session.Run(updateQuery, params)
	if err != nil {
		log.Error(fmt.Sprintf("Update error: %v", err))
		return pkgerrors.NewInternalError("neo4j update error", err)
	}
	if !result.Next() {
		if result.Err() != nil {
			log.Error(fmt.Sprintf("Update result error: %v", result.Err()))
			return pkgerrors.NewInternalError("neo4j update result error", result.Err())
		}
		log.Warn("Failed to update bucket: no result returned")
		return pkgerrors.NewInternalError("failed to update bucket: no result returned", nil)
	}
	log.Info(fmt.Sprintf("Bucket updated successfully: %v", params["id"]))
	return nil
}

func (r *Repository) Delete(ctx context.Context, id valueobjects.BucketID) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.Run(deleteQuery, map[string]interface{}{"id": id.Value()})
	return err
}

func (r *Repository) List(ctx context.Context, limit *int, offset *int) ([]*entity.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	limitVal := 10
	if limit != nil {
		limitVal = *limit
	}

	params := map[string]interface{}{
		"limit": limitVal,
	}

	result, err := session.Run(listQuery, params)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to list buckets", err)
	}

	var buckets []*entity.Bucket
	for result.Next() {
		record := result.Record()
		bucket, err := RecordToBucket(record)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func (r *Repository) Search(ctx context.Context, query string, limit *int, offset *int) ([]*entity.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	limitVal := 10
	if limit != nil {
		limitVal = *limit
	}

	params := map[string]interface{}{
		"query": query,
		"limit": limitVal,
	}

	result, err := session.Run(searchQuery, params)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to search buckets", err)
	}

	var buckets []*entity.Bucket
	for result.Next() {
		record := result.Record()
		bucket, err := RecordToBucket(record)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func (r *Repository) GetByOwnerID(ctx context.Context, ownerID string, limit *int, lastKey map[string]interface{}) (*entity.BucketPage, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	params := map[string]interface{}{
		"ownerID": ownerID,
		"limit":   limit,
	}

	result, err := session.Run(getByOwnerIDQuery, params)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to get buckets by owner", err)
	}

	var buckets []*entity.Bucket
	for result.Next() {
		record := result.Record()
		bucket, err := RecordToBucket(record)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return &entity.BucketPage{
		Items:   buckets,
		HasMore: len(buckets) == *limit,
	}, nil
}

func (r *Repository) AddAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	log := logger.WithService("neo4j-bucket-repository").WithContext(ctx)

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
		"assetID":  assetID,
	}

	log.Info(fmt.Sprintf("AddAsset: bucketID=%s assetID=%s query=%s", bucketID.Value(), assetID, addAssetQuery))
	result, err := session.Run(addAssetQuery, params)
	if err != nil {
		log.Error(fmt.Sprintf("AddAsset error: %v", err))
		return pkgerrors.NewInternalError("add asset error", err)
	}
	if !result.Next() {
		if result.Err() != nil {
			log.Error(fmt.Sprintf("AddAsset result error: %v", result.Err()))
			return pkgerrors.NewInternalError("add asset result error", result.Err())
		}
		log.Warn("AddAsset: no result returned")
		return pkgerrors.NewInternalError("add asset: no result returned", nil)
	}
	log.Info("AddAsset: relationship created successfully")
	return nil
}

func (r *Repository) RemoveAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
		"assetID":  assetID,
	}

	_, err := session.Run(removeAssetQuery, params)
	return err
}

func (r *Repository) GetAssetIDs(ctx context.Context, bucketID valueobjects.BucketID, limit *int, lastKey map[string]interface{}) ([]string, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	log := logger.WithService("neo4j-bucket-repository").WithContext(ctx)

	if limit == nil {
		limit = new(int)
		*limit = 10
	}

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
		"limit":    *limit,
	}

	log.Info(fmt.Sprintf("GetAssetIDs: bucketID=%s query=%s params=%+v", bucketID.Value(), getAssetIDsQuery, params))
	result, err := session.Run(getAssetIDsQuery, params)
	if err != nil {
		log.Error(fmt.Sprintf("GetAssetIDs error: %v", err))
		return nil, pkgerrors.NewInternalError("failed to get bucket asset IDs", err)
	}

	var assetIDs []string
	for result.Next() {
		record := result.Record()
		if assetID, ok := record.Get("a.id"); ok {
			if id, ok := assetID.(string); ok {
				assetIDs = append(assetIDs, id)
			}
		}
	}

	log.Info(fmt.Sprintf("GetAssetIDs: found %d asset IDs", len(assetIDs)))
	return assetIDs, nil
}

func (r *Repository) GetByKey(ctx context.Context, key string) (*entity.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(getByKeyQuery, map[string]interface{}{"key": key})
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to get bucket by key", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, pkgerrors.NewNotFoundError("bucket not found", nil)
	}

	return RecordToBucket(record)
}

func (r *Repository) HasAsset(ctx context.Context, bucketID valueobjects.BucketID, assetID string) (bool, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
		"assetID":  assetID,
	}

	result, err := session.Run(hasAssetQuery, params)
	if err != nil {
		return false, pkgerrors.NewInternalError("failed to check if asset is in bucket", err)
	}

	if result.Next() {
		count, _ := result.Record().Get("count")
		if cnt, ok := count.(int64); ok && cnt > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (r *Repository) AssetCount(ctx context.Context, bucketID valueobjects.BucketID) (int, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
	}

	result, err := session.Run(assetCountQuery, params)
	if err != nil {
		return 0, pkgerrors.NewInternalError("failed to get asset count for bucket", err)
	}

	if result.Next() {
		count, _ := result.Record().Get("count")
		if cnt, ok := count.(int64); ok {
			return int(cnt), nil
		}
	}
	return 0, nil
}

func (r *Repository) FindByType(ctx context.Context, bucketType valueobjects.BucketType, limit *int, offset *int) ([]*entity.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	limitVal := 10
	if limit != nil {
		limitVal = *limit
	}

	params := map[string]interface{}{
		"type":  bucketType.Value(),
		"limit": limitVal,
	}

	result, err := session.Run(findByTypeQuery, params)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to find buckets by type", err)
	}

	var buckets []*entity.Bucket
	for result.Next() {
		record := result.Record()
		bucket, err := RecordToBucket(record)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func (r *Repository) FindByStatus(ctx context.Context, status valueobjects.BucketStatus, limit *int, offset *int) ([]*entity.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	limitVal := 10
	if limit != nil {
		limitVal = *limit
	}

	params := map[string]interface{}{
		"status": status.Value(),
		"limit":  limitVal,
	}

	result, err := session.Run(findByStatusQuery, params)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to find buckets by status", err)
	}

	var buckets []*entity.Bucket
	for result.Next() {
		record := result.Record()
		bucket, err := RecordToBucket(record)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(countQuery, nil)
	if err != nil {
		return 0, pkgerrors.NewInternalError("failed to count buckets", err)
	}

	if result.Next() {
		count, _ := result.Record().Get("count")
		if cnt, ok := count.(int64); ok {
			return cnt, nil
		}
	}
	return 0, nil
}

func (r *Repository) CountByOwnerID(ctx context.Context, ownerID valueobjects.OwnerID) (int64, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	params := map[string]interface{}{
		"ownerID": ownerID.Value(),
	}

	result, err := session.Run(countByOwnerIDQuery, params)
	if err != nil {
		return 0, pkgerrors.NewInternalError("failed to count buckets by owner", err)
	}

	if result.Next() {
		count, _ := result.Record().Get("count")
		if cnt, ok := count.(int64); ok {
			return cnt, nil
		}
	}
	return 0, nil
}

func (r *Repository) CountByType(ctx context.Context, bucketType valueobjects.BucketType) (int64, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	params := map[string]interface{}{
		"type": bucketType.Value(),
	}

	result, err := session.Run(countByTypeQuery, params)
	if err != nil {
		return 0, pkgerrors.NewInternalError("failed to count buckets by type", err)
	}

	if result.Next() {
		count, _ := result.Record().Get("count")
		if cnt, ok := count.(int64); ok {
			return cnt, nil
		}
	}
	return 0, nil
}

func (r *Repository) Exists(ctx context.Context, id valueobjects.BucketID) (bool, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	params := map[string]interface{}{
		"id": id.Value(),
	}

	result, err := session.Run(existsQuery, params)
	if err != nil {
		return false, pkgerrors.NewInternalError("failed to check if bucket exists", err)
	}

	if result.Next() {
		count, _ := result.Record().Get("count")
		if cnt, ok := count.(int64); ok && cnt > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (r *Repository) ExistsByKey(ctx context.Context, key valueobjects.BucketKey) (bool, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	params := map[string]interface{}{
		"key": key.Value(),
	}

	result, err := session.Run(existsByKeyQuery, params)
	if err != nil {
		return false, pkgerrors.NewInternalError("failed to check if bucket exists by key", err)
	}

	if result.Next() {
		count, _ := result.Record().Get("count")
		if cnt, ok := count.(int64); ok && cnt > 0 {
			return true, nil
		}
	}
	return false, nil
}
