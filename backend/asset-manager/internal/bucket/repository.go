package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type BucketRepository interface {
	GetBucketByID(ctx context.Context, id string) (*Bucket, error)
	ListBuckets(ctx context.Context, limit int) (*BucketPage, error)
	CreateBucket(ctx context.Context, bucket *Bucket) error
	UpdateBucket(ctx context.Context, bucket *Bucket) error
	PatchBucket(ctx context.Context, id string, patch map[string]interface{}) error
	DeleteBucket(ctx context.Context, id string) error
	GetBucketsByType(ctx context.Context, bucketType string) ([]Bucket, error)
	GetBucketsByAsset(ctx context.Context, assetID string) ([]Bucket, error)
	GetAssetsInBucket(ctx context.Context, bucketID string) ([]string, error)
	GetBucketByKey(ctx context.Context, key string) (*Bucket, error)
}

type Repository struct {
	driver neo4j.Driver
	logger *logger.Logger
}

func NewRepository(driver neo4j.Driver) *Repository {
	return &Repository{
		driver: driver,
		logger: logger.WithService("bucket-neo4j-repository"),
	}
}

func (r *Repository) CreateBucket(ctx context.Context, b *Bucket) error {
	log := r.logger.WithContext(ctx)

	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		CREATE (b:Bucket {
			id: $id,
			key: $key,
			name: $name,
			description: $description,
			type: $type,
			status: $status,
			assetIds: $assetIds,
			createdAt: $createdAt,
			updatedAt: $updatedAt
		})
		RETURN b
	`

	params := map[string]interface{}{
		"id":          b.ID,
		"key":         b.Key,
		"name":        b.Name,
		"description": b.Description,
		"type":        b.Type,
		"status":      b.Status,
		"assetIds":    b.AssetIDs,
		"createdAt":   b.CreatedAt,
		"updatedAt":   b.UpdatedAt,
	}

	_, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to create bucket", "bucket_id", b.ID)
		return apperrors.NewInternalError("failed to create bucket", err)
	}

	log.Debug("Bucket created successfully", "bucket_id", b.ID, "name", b.Name)
	return nil
}

func (r *Repository) UpdateBucket(ctx context.Context, b *Bucket) error {
	log := r.logger.WithContext(ctx)

	b.UpdatedAt = time.Now().UTC()

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $id})
		SET b.name = $name,
			b.description = $description,
			b.type = $type,
			b.status = $status,
			b.assetIds = $assetIds,
			b.updatedAt = $updatedAt
		RETURN b
	`

	params := map[string]interface{}{
		"id":          b.ID,
		"name":        b.Name,
		"description": b.Description,
		"type":        b.Type,
		"status":      b.Status,
		"assetIds":    b.AssetIDs,
		"updatedAt":   b.UpdatedAt,
	}

	_, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to update bucket", "bucket_id", b.ID)
		return apperrors.NewInternalError("failed to update bucket", err)
	}

	log.Debug("Bucket updated successfully", "bucket_id", b.ID, "name", b.Name)
	return nil
}

func (r *Repository) GetBucketByID(ctx context.Context, id string) (*Bucket, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $id})
		RETURN b
	`

	params := map[string]interface{}{
		"id": id,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to get bucket", "bucket_id", id)
		return nil, apperrors.NewInternalError("failed to get bucket", err)
	}

	record, err := result.Single()
	if err != nil {
		log.Debug("Bucket not found", "bucket_id", id)
		return nil, apperrors.NewNotFoundError("bucket not found", err)
	}

	bucket, err := r.recordToBucket(record)
	if err != nil {
		log.WithError(err).Error("Failed to convert Neo4j record to bucket", "bucket_id", id)
		return nil, apperrors.NewInternalError("convert record to bucket failed", err)
	}

	log.Debug("Bucket retrieved successfully", "bucket_id", id, "name", bucket.Name)
	return bucket, nil
}

func (r *Repository) ListBuckets(ctx context.Context, limit int) (*BucketPage, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket)
		RETURN b
		ORDER BY b.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"limit": limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to list buckets")
		return nil, apperrors.NewInternalError("failed to list buckets", err)
	}

	var buckets []Bucket
	for result.Next() {
		bucket, err := r.recordToBucket(result.Record())
		if err != nil {
			log.WithError(err).Warn("Failed to convert record to bucket, skipping")
			continue
		}
		buckets = append(buckets, *bucket)
	}

	if err := result.Err(); err != nil {
		log.WithError(err).Error("Error iterating through Neo4j results")
		return nil, apperrors.NewInternalError("iterate results failed", err)
	}

	log.Debug("Buckets listed successfully from Neo4j", "count", len(buckets))
	return &BucketPage{
		Items: buckets,
	}, nil
}

func (r *Repository) PatchBucket(ctx context.Context, id string, patch map[string]interface{}) error {
	log := r.logger.WithContext(ctx)

	if len(patch) == 0 {
		return nil
	}

	if _, hasID := patch["id"]; hasID {
		return apperrors.NewValidationError("cannot modify id field", nil)
	}
	if _, hasKey := patch["key"]; hasKey {
		return apperrors.NewValidationError("cannot modify key field", nil)
	}

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	setClause := "SET b.updatedAt = $updatedAt"
	params := map[string]interface{}{
		"id":        id,
		"updatedAt": time.Now().UTC(),
	}

	for key, value := range patch {
		paramKey := fmt.Sprintf("param_%s", key)
		setClause += fmt.Sprintf(", b.%s = $%s", key, paramKey)
		params[paramKey] = value
	}

	query := fmt.Sprintf(`
		MATCH (b:Bucket {id: $id})
		%s
		RETURN b
	`, setClause)

	_, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to patch bucket", "bucket_id", id)
		return apperrors.NewInternalError("failed to patch bucket", err)
	}

	log.Debug("Bucket patched successfully", "bucket_id", id)
	return nil
}

func (r *Repository) DeleteBucket(ctx context.Context, id string) error {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $id})
		DETACH DELETE b
	`

	result, err := session.Run(query, map[string]interface{}{"id": id})
	if err != nil {
		log.WithError(err).Error("Failed to delete bucket from Neo4j", "bucket_id", id)
		return apperrors.NewInternalError("delete bucket failed", err)
	}

	summary, err := result.Consume()
	if err != nil {
		log.WithError(err).Error("Failed to consume delete result", "bucket_id", id)
		return apperrors.NewInternalError("consume delete result failed", err)
	}

	if summary.Counters().NodesDeleted() == 0 {
		log.Debug("Bucket not found for deletion", "bucket_id", id)
		return apperrors.NewNotFoundError("bucket not found", nil)
	}

	log.Debug("Bucket deleted successfully from Neo4j", "bucket_id", id)
	return nil
}

func (r *Repository) recordToBucket(record *neo4j.Record) (*Bucket, error) {
	node, ok := record.Get("b")
	if !ok {
		return nil, apperrors.NewInternalError("no 'b' field in record", nil)
	}

	neo4jNode, ok := node.(neo4j.Node)
	if !ok {
		return nil, apperrors.NewInternalError("field 'b' is not a node", nil)
	}

	props := neo4jNode.Props

	getString := func(key string) string {
		if val, exists := props[key]; exists && val != nil {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	getTime := func(key string) time.Time {
		if val, exists := props[key]; exists && val != nil {
			if t, ok := val.(time.Time); ok {
				return t
			}
		}
		return time.Time{}
	}

	bucket := &Bucket{
		ID:          getString("id"),
		Key:         getString("key"),
		Name:        getString("name"),
		Description: getString("description"),
		Type:        getString("type"),
		Status:      getString("status"),
		CreatedAt:   getTime("createdAt"),
		UpdatedAt:   getTime("updatedAt"),
	}

	if assetIDs, ok := props["assetIds"].([]interface{}); ok {
		for _, assetID := range assetIDs {
			if assetIDStr, ok := assetID.(string); ok {
				bucket.AssetIDs = append(bucket.AssetIDs, assetIDStr)
			}
		}
	}

	return bucket, nil
}

func (r *Repository) GetBucketsByType(ctx context.Context, bucketType string) ([]Bucket, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {type: $type})
		RETURN b
		ORDER BY b.createdAt DESC
	`

	result, err := session.Run(query, map[string]interface{}{"type": bucketType})
	if err != nil {
		log.WithError(err).Error("Failed to get buckets by type from Neo4j", "type", bucketType)
		return nil, apperrors.NewInternalError("get buckets by type failed", err)
	}

	var buckets []Bucket
	for result.Next() {
		bucket, err := r.recordToBucket(result.Record())
		if err != nil {
			log.WithError(err).Warn("Failed to convert record to bucket, skipping")
			continue
		}
		buckets = append(buckets, *bucket)
	}

	if err := result.Err(); err != nil {
		log.WithError(err).Error("Error iterating through Neo4j bucket type results")
		return nil, apperrors.NewInternalError("iterate bucket type results failed", err)
	}

	log.Debug("Buckets by type retrieved successfully from Neo4j", "type", bucketType, "count", len(buckets))
	return buckets, nil
}

func (r *Repository) GetBucketsByAsset(ctx context.Context, assetID string) ([]Bucket, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket)
		WHERE $assetID IN b.assetIds
		RETURN b
		ORDER BY b.createdAt DESC
	`

	result, err := session.Run(query, map[string]interface{}{"assetID": assetID})
	if err != nil {
		log.WithError(err).Error("Failed to get buckets by asset from Neo4j", "asset_id", assetID)
		return nil, apperrors.NewInternalError("get buckets by asset failed", err)
	}

	var buckets []Bucket
	for result.Next() {
		bucket, err := r.recordToBucket(result.Record())
		if err != nil {
			log.WithError(err).Warn("Failed to convert record to bucket, skipping")
			continue
		}
		buckets = append(buckets, *bucket)
	}

	if err := result.Err(); err != nil {
		log.WithError(err).Error("Error iterating through Neo4j bucket asset results")
		return nil, apperrors.NewInternalError("iterate bucket asset results failed", err)
	}

	log.Debug("Buckets by asset retrieved successfully from Neo4j", "asset_id", assetID, "count", len(buckets))
	return buckets, nil
}

func (r *Repository) GetAssetsInBucket(ctx context.Context, bucketID string) ([]string, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $bucketID})
		RETURN b.assetIds
	`

	neo4jResult, err := session.Run(query, map[string]interface{}{"bucketID": bucketID})
	if err != nil {
		log.WithError(err).Error("Failed to get assets in bucket from Neo4j", "bucket_id", bucketID)
		return nil, apperrors.NewInternalError("get assets in bucket failed", err)
	}

	record, err := neo4jResult.Single()
	if err != nil {
		log.Debug("Bucket not found in Neo4j", "bucket_id", bucketID)
		return nil, apperrors.NewNotFoundError("bucket not found", err)
	}

	assetIDsInterface, ok := record.Get("b.assetIds")
	if !ok {
		log.Debug("No asset IDs found in bucket", "bucket_id", bucketID)
		return []string{}, nil
	}

	assetIDs, ok := assetIDsInterface.([]interface{})
	if !ok {
		log.Debug("Asset IDs is not a slice", "bucket_id", bucketID)
		return []string{}, nil
	}

	var assetIDList []string
	for _, assetID := range assetIDs {
		if assetIDStr, ok := assetID.(string); ok {
			assetIDList = append(assetIDList, assetIDStr)
		}
	}

	log.Debug("Assets in bucket retrieved successfully from Neo4j", "bucket_id", bucketID, "count", len(assetIDList))
	return assetIDList, nil
}

func (r *Repository) GetBucketByKey(ctx context.Context, key string) (*Bucket, error) {
	log := r.logger.WithContext(ctx)

	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {key: $key})
		RETURN b
	`

	params := map[string]interface{}{
		"key": key,
	}

	result, err := session.Run(query, params)
	if err != nil {
		log.WithError(err).Error("Failed to get bucket by key from Neo4j", "key", key)
		return nil, apperrors.NewInternalError("get bucket by key failed", err)
	}

	record, err := result.Single()
	if err != nil {
		log.Debug("Bucket not found by key", "key", key)
		return nil, apperrors.NewNotFoundError("bucket not found by key", err)
	}

	bucket, err := r.recordToBucket(record)
	if err != nil {
		log.WithError(err).Error("Failed to convert Neo4j record to bucket by key", "key", key)
		return nil, apperrors.NewInternalError("convert record to bucket by key failed", err)
	}

	log.Debug("Bucket retrieved by key successfully", "key", key, "bucket_id", bucket.ID)
	return bucket, nil
}
