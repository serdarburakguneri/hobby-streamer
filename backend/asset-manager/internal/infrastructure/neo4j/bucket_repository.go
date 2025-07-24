package neo4j

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type BucketRepository struct {
	driver neo4j.Driver
}

func NewBucketRepository(driver neo4j.Driver) *BucketRepository {
	return &BucketRepository{
		driver: driver,
	}
}

func (r *BucketRepository) Create(ctx context.Context, bucket *domainbucket.Bucket) error {
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

	query := `
		CREATE (b:Bucket {
			id: $id,
			name: $name,
			description: $description,
			key: $key,
			ownerID: $ownerID,
			status: $status,
			type: $type,
			metadata: $metadata,
			createdAt: $createdAt,
			updatedAt: $updatedAt
		})
		RETURN b
	`

	params := map[string]interface{}{
		"id":          bucket.ID().Value(),
		"name":        bucket.Name(),
		"description": bucket.Description(),
		"key":         bucket.Key(),
		"ownerID":     bucket.OwnerID(),
		"status":      bucket.Status(),
		"type":        bucket.Type(),
		"metadata":    metadataJSON,
		"createdAt":   bucket.CreatedAt().Format(time.RFC3339),
		"updatedAt":   bucket.UpdatedAt().Format(time.RFC3339),
	}

	log.Info(fmt.Sprintf("Creating bucket with params: %+v", params))
	result, err := session.Run(query, params)
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

func (r *BucketRepository) GetByID(ctx context.Context, id domainbucket.BucketID) (*domainbucket.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $id})
		OPTIONAL MATCH (b)-[:CONTAINS]->(a:Asset)
		RETURN b, collect(a) as assets
	`

	result, err := session.Run(query, map[string]interface{}{"id": id.Value()})
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to get bucket", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, pkgerrors.NewNotFoundError("bucket not found", nil)
	}

	return r.recordToBucket(record)
}

func (r *BucketRepository) GetBySlug(ctx context.Context, slug string) (*domainbucket.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {slug: $slug})
		OPTIONAL MATCH (b)-[:CONTAINS]->(a:Asset)
		RETURN b, collect(a) as assets
	`

	result, err := session.Run(query, map[string]interface{}{"slug": slug})
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to get bucket by slug", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, pkgerrors.NewNotFoundError("bucket not found", nil)
	}

	return r.recordToBucket(record)
}

func (r *BucketRepository) Update(ctx context.Context, bucket *domainbucket.Bucket) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $id})
		SET b.name = $name,
			b.description = $description,
			b.ownerID = $ownerID,
			b.status = $status,
			b.type = $type,
			b.metadata = $metadata,
			b.updatedAt = $updatedAt
		RETURN b
	`

	params := map[string]interface{}{
		"id":          bucket.ID().Value(),
		"name":        bucket.Name(),
		"description": bucket.Description(),
		"ownerID":     bucket.OwnerID(),
		"status":      bucket.Status(),
		"type":        bucket.Type(),
		"metadata":    bucket.Metadata(),
		"updatedAt":   bucket.UpdatedAt().Format(time.RFC3339),
	}

	_, err := session.Run(query, params)
	return err
}

func (r *BucketRepository) Delete(ctx context.Context, id domainbucket.BucketID) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $id})
		OPTIONAL MATCH (b)-[r:CONTAINS]->(a:Asset)
		DELETE r, b
	`

	_, err := session.Run(query, map[string]interface{}{"id": id.Value()})
	return err
}

const defaultLimit = 10

func (r *BucketRepository) List(ctx context.Context, limit *int, lastKey map[string]interface{}) (*domainbucket.BucketPage, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	if limit == nil {
		limit = new(int)
		*limit = defaultLimit
	}

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
		return nil, pkgerrors.NewInternalError("failed to list buckets", err)
	}

	var buckets []*domainbucket.Bucket
	for result.Next() {
		record := result.Record()
		bucket, err := r.recordToBucket(record)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return &domainbucket.BucketPage{
		Items:   buckets,
		HasMore: len(buckets) == *limit,
		Total:   len(buckets),
	}, nil
}

func (r *BucketRepository) Search(ctx context.Context, query string, limit *int, lastKey map[string]interface{}) (*domainbucket.BucketPage, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	searchQuery := fmt.Sprintf("(?i).*%s.*", query)
	cypherQuery := `
		MATCH (b:Bucket)
		WHERE b.name =~ $searchQuery OR b.description =~ $searchQuery
		RETURN b
		ORDER BY b.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"searchQuery": searchQuery,
		"limit":       limit,
	}

	result, err := session.Run(cypherQuery, params)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to search buckets", err)
	}

	var buckets []*domainbucket.Bucket
	for result.Next() {
		record := result.Record()
		bucket, err := r.recordToBucket(record)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return &domainbucket.BucketPage{
		Items:   buckets,
		HasMore: len(buckets) == *limit,
		Total:   len(buckets),
	}, nil
}

func (r *BucketRepository) GetByOwnerID(ctx context.Context, ownerID string, limit *int, lastKey map[string]interface{}) (*domainbucket.BucketPage, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {ownerID: $ownerID})
		RETURN b
		ORDER BY b.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"ownerID": ownerID,
		"limit":   limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to get buckets by owner", err)
	}

	var buckets []*domainbucket.Bucket
	for result.Next() {
		record := result.Record()
		bucket, err := r.recordToBucket(record)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}

	return &domainbucket.BucketPage{
		Items:   buckets,
		HasMore: len(buckets) == *limit,
		Total:   len(buckets),
	}, nil
}

func (r *BucketRepository) AddAsset(ctx context.Context, bucketID domainbucket.BucketID, assetID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	log := logger.WithService("neo4j-bucket-repository").WithContext(ctx)

	query := `
		MATCH (b:Bucket {id: $bucketID})
		MATCH (a:Asset {id: $assetID})
		MERGE (b)-[:CONTAINS]->(a)
		RETURN b, a
	`

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
		"assetID":  assetID,
	}

	log.Info(fmt.Sprintf("AddAsset: bucketID=%s assetID=%s query=%s", bucketID.Value(), assetID, query))
	result, err := session.Run(query, params)
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

func (r *BucketRepository) RemoveAsset(ctx context.Context, bucketID domainbucket.BucketID, assetID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $bucketID})-[r:CONTAINS]->(a:Asset {id: $assetID})
		DELETE r
	`

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
		"assetID":  assetID,
	}

	_, err := session.Run(query, params)
	return err
}

func (r *BucketRepository) GetAssetIDs(ctx context.Context, bucketID domainbucket.BucketID, limit *int, lastKey map[string]interface{}) ([]string, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	log := logger.WithService("neo4j-bucket-repository").WithContext(ctx)

	if limit == nil {
		limit = new(int)
		*limit = defaultLimit
	}

	query := `
		MATCH (b:Bucket {id: $bucketID})-[:CONTAINS]->(a:Asset)
		RETURN a.id
		ORDER BY a.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
		"limit":    *limit,
	}

	log.Info(fmt.Sprintf("GetAssetIDs: bucketID=%s query=%s params=%+v", bucketID.Value(), query, params))
	result, err := session.Run(query, params)
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

func (r *BucketRepository) GetByKey(ctx context.Context, key string) (*domainbucket.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {key: $key})
		OPTIONAL MATCH (b)-[:CONTAINS]->(a:Asset)
		RETURN b, collect(a) as assets
	`

	result, err := session.Run(query, map[string]interface{}{"key": key})
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to get bucket by key", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, domainbucket.ErrBucketNotFound
	}

	return r.recordToBucket(record)
}

func (r *BucketRepository) recordToBucket(record *neo4j.Record) (*domainbucket.Bucket, error) {
	bucketNode, ok := record.Get("b")
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket node not found in record", nil)
	}

	bucketProps := bucketNode.(neo4j.Node).Props

	idRaw, ok := bucketProps["id"]
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket id missing", nil)
	}
	id, ok := idRaw.(string)
	if !ok || id == "" {
		return nil, pkgerrors.NewInternalError("bucket id is not a string or is empty", nil)
	}
	bucketID, err := domainbucket.NewBucketID(id)
	if err != nil {
		return nil, err
	}

	nameRaw, ok := bucketProps["name"]
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket name missing", nil)
	}
	name, ok := nameRaw.(string)
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket name is not a string", nil)
	}

	keyRaw, ok := bucketProps["key"]
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket key missing", nil)
	}
	key, ok := keyRaw.(string)
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket key is not a string", nil)
	}

	descriptionPtr, ownerIDPtr, statusPtr := (*string)(nil), (*string)(nil), (*string)(nil)
	description, _ := bucketProps["description"].(string)
	ownerID, _ := bucketProps["ownerID"].(string)
	status, _ := bucketProps["status"].(string)
	if description != "" {
		descriptionPtr = &description
	}
	if ownerID != "" {
		ownerIDPtr = &ownerID
	}
	if status != "" {
		statusPtr = &status
	}

	var metadata map[string]interface{}
	if metadataInterface, exists := bucketProps["metadata"]; exists {
		if metadataMap, ok := metadataInterface.(string); ok {
			if err := json.Unmarshal([]byte(metadataMap), &metadata); err != nil {
				return nil, pkgerrors.NewInternalError("failed to unmarshal metadata", err)
			}
		}
	}

	createdAtRaw, ok := bucketProps["createdAt"]
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket createdAt missing", nil)
	}
	createdAtStr, ok := createdAtRaw.(string)
	if !ok || createdAtStr == "" {
		return nil, pkgerrors.NewInternalError("bucket createdAt is not a string or is empty", nil)
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to parse createdAt: "+createdAtStr, err)
	}

	updatedAtRaw, ok := bucketProps["updatedAt"]
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket updatedAt missing", nil)
	}
	updatedAtStr, ok := updatedAtRaw.(string)
	if !ok || updatedAtStr == "" {
		return nil, pkgerrors.NewInternalError("bucket updatedAt is not a string or is empty", nil)
	}
	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return nil, pkgerrors.NewInternalError("failed to parse updatedAt: "+updatedAtStr, err)
	}

	bucket := domainbucket.ReconstructBucket(
		*bucketID,
		name,
		descriptionPtr,
		key,
		ownerIDPtr,
		statusPtr,
		metadata,
		createdAt,
		updatedAt,
	)

	typeRaw, ok := bucketProps["type"]
	if !ok {
		return nil, pkgerrors.NewInternalError("bucket type missing", nil)
	}
	typeStr, ok := typeRaw.(string)
	if !ok || typeStr == "" {
		return nil, pkgerrors.NewInternalError("bucket type is not a string or is empty", nil)
	}
	if err := bucket.SetType(typeStr); err != nil {
		return nil, err
	}

	return bucket, nil
}

func (r *BucketRepository) HasAsset(ctx context.Context, bucketID domainbucket.BucketID, assetID string) (bool, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $bucketID})-[:CONTAINS]->(a:Asset {id: $assetID})
		RETURN count(a) as count
	`

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
		"assetID":  assetID,
	}

	result, err := session.Run(query, params)
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

func (r *BucketRepository) AssetCount(ctx context.Context, bucketID domainbucket.BucketID) (int, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $bucketID})-[:CONTAINS]->(a:Asset)
		RETURN count(a) as count
	`

	params := map[string]interface{}{
		"bucketID": bucketID.Value(),
	}

	result, err := session.Run(query, params)
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
