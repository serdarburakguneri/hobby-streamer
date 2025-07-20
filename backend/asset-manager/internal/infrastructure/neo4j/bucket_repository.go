package neo4j

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	domainbucket "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket"
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

	query := `
		CREATE (b:Bucket {
			id: $id,
			name: $name,
			description: $description,
			key: $key,
			ownerID: $ownerID,
			metadata: $metadata,
			createdAt: $createdAt,
			updatedAt: $updatedAt
		})
		RETURN b
	`

	params := map[string]interface{}{
		"id":          bucket.ID(),
		"name":        bucket.Name(),
		"description": bucket.Description(),
		"key":         bucket.Key(),
		"ownerID":     bucket.OwnerID(),
		"metadata":    bucket.Metadata(),
		"createdAt":   bucket.CreatedAt(),
		"updatedAt":   bucket.UpdatedAt(),
	}

	_, err := session.Run(query, params)
	return err
}

func (r *BucketRepository) GetByID(ctx context.Context, id string) (*domainbucket.Bucket, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $id})
		OPTIONAL MATCH (b)-[:CONTAINS]->(a:Asset)
		RETURN b, collect(a) as assets
	`

	result, err := session.Run(query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, fmt.Errorf("bucket not found")
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
		return nil, fmt.Errorf("failed to get bucket by slug: %w", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, fmt.Errorf("bucket not found")
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
			b.metadata = $metadata,
			b.updatedAt = $updatedAt
		RETURN b
	`

	params := map[string]interface{}{
		"id":          bucket.ID(),
		"name":        bucket.Name(),
		"description": bucket.Description(),
		"ownerID":     bucket.OwnerID(),
		"metadata":    bucket.Metadata(),
		"updatedAt":   bucket.UpdatedAt(),
	}

	_, err := session.Run(query, params)
	return err
}

func (r *BucketRepository) Delete(ctx context.Context, id string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $id})
		OPTIONAL MATCH (b)-[r:CONTAINS]->(a:Asset)
		DELETE r, b
	`

	_, err := session.Run(query, map[string]interface{}{"id": id})
	return err
}

func (r *BucketRepository) List(ctx context.Context, limit *int, lastKey map[string]interface{}) (*domainbucket.BucketPage, error) {
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
		return nil, fmt.Errorf("failed to list buckets: %w", err)
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
		return nil, fmt.Errorf("failed to search buckets: %w", err)
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
		return nil, fmt.Errorf("failed to get buckets by owner: %w", err)
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

func (r *BucketRepository) AddAsset(ctx context.Context, bucketID string, assetID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $bucketID})
		MATCH (a:Asset {id: $assetID})
		MERGE (b)-[:CONTAINS]->(a)
		RETURN b, a
	`

	params := map[string]interface{}{
		"bucketID": bucketID,
		"assetID":  assetID,
	}

	_, err := session.Run(query, params)
	return err
}

func (r *BucketRepository) RemoveAsset(ctx context.Context, bucketID string, assetID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $bucketID})-[r:CONTAINS]->(a:Asset {id: $assetID})
		DELETE r
	`

	params := map[string]interface{}{
		"bucketID": bucketID,
		"assetID":  assetID,
	}

	_, err := session.Run(query, params)
	return err
}

func (r *BucketRepository) GetAssetIDs(ctx context.Context, bucketID string, limit *int, lastKey map[string]interface{}) ([]string, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (b:Bucket {id: $bucketID})-[:CONTAINS]->(a:Asset)
		RETURN a.id
		ORDER BY a.createdAt DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"bucketID": bucketID,
		"limit":    limit,
	}

	result, err := session.Run(query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket asset IDs: %w", err)
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
		return nil, fmt.Errorf("failed to get bucket by key: %w", err)
	}

	record, err := result.Single()
	if err != nil {
		return nil, fmt.Errorf("bucket not found")
	}

	return r.recordToBucket(record)
}

func (r *BucketRepository) recordToBucket(record *neo4j.Record) (*domainbucket.Bucket, error) {
	bucketNode, ok := record.Get("b")
	if !ok {
		return nil, fmt.Errorf("bucket node not found in record")
	}

	bucketProps := bucketNode.(neo4j.Node).Props

	var assetIDs []string
	if assetsList, ok := record.Get("assets"); ok {
		if assetsList != nil {
			assetsSlice := assetsList.([]interface{})
			for _, assetNode := range assetsSlice {
				if assetNode != nil {
					if props := assetNode.(neo4j.Node).Props; props != nil {
						if id, ok := props["id"].(string); ok {
							assetIDs = append(assetIDs, id)
						}
					}
				}
			}
		}
	}

	var descriptionPtr, ownerIDPtr *string
	description, _ := bucketProps["description"].(string)
	ownerID, _ := bucketProps["ownerID"].(string)

	if description != "" {
		descriptionPtr = &description
	}
	if ownerID != "" {
		ownerIDPtr = &ownerID
	}

	var metadata map[string]interface{}
	if metadataInterface, exists := bucketProps["metadata"]; exists {
		if metadataMap, ok := metadataInterface.(map[string]interface{}); ok {
			metadata = metadataMap
		}
	}

	return domainbucket.ReconstructBucket(
		bucketProps["id"].(string),
		bucketProps["name"].(string),
		descriptionPtr,
		bucketProps["key"].(string),
		ownerIDPtr,
		assetIDs,
		metadata,
		bucketProps["createdAt"].(time.Time),
		bucketProps["updatedAt"].(time.Time),
	), nil
}
