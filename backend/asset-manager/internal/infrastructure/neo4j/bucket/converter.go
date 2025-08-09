package bucket

import (
	"encoding/json"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

func RecordToBucket(record *neo4j.Record) (*entity.Bucket, error) {
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
	bucketID, err := valueobjects.NewBucketID(id)
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

	descriptionPtr, ownerIDPtr, statusPtr, typePtr := (*string)(nil), (*string)(nil), (*string)(nil), (*string)(nil)
	description, _ := bucketProps["description"].(string)
	ownerID, _ := bucketProps["ownerID"].(string)
	status, _ := bucketProps["status"].(string)
	bucketType, _ := bucketProps["type"].(string)
	if description != "" {
		descriptionPtr = &description
	}
	if ownerID != "" {
		ownerIDPtr = &ownerID
	}
	if status != "" {
		statusPtr = &status
	}
	if bucketType != "" {
		typePtr = &bucketType
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

	bucketName, err := valueobjects.NewBucketName(name)
	if err != nil {
		return nil, err
	}

	bucketKey, err := valueobjects.NewBucketKey(key)
	if err != nil {
		return nil, err
	}

	var bucketDescription *valueobjects.BucketDescription
	if descriptionPtr != nil {
		desc, err := valueobjects.NewBucketDescription(*descriptionPtr)
		if err != nil {
			return nil, err
		}
		bucketDescription = desc
	}

	var bucketOwnerID *valueobjects.OwnerID
	if ownerIDPtr != nil {
		ownerID, err := valueobjects.NewOwnerID(*ownerIDPtr)
		if err != nil {
			return nil, err
		}
		bucketOwnerID = ownerID
	}

	var bucketStatus *valueobjects.BucketStatus
	if statusPtr != nil {
		status, err := valueobjects.NewBucketStatus(*statusPtr)
		if err != nil {
			return nil, err
		}
		bucketStatus = status
	}

	var bucketTypeVO *valueobjects.BucketType
	if typePtr != nil {
		bucketTypeVO, err = valueobjects.NewBucketType(*typePtr)
		if err != nil {
			return nil, err
		}
	}

	createdAtVO := valueobjects.NewCreatedAt(createdAt)
	updatedAtVO := valueobjects.NewUpdatedAt(updatedAt)

	bucket := entity.ReconstructBucket(
		*bucketID,
		*bucketName,
		bucketDescription,
		*bucketKey,
		bucketTypeVO,
		bucketOwnerID,
		bucketStatus,
		metadata,
		createdAtVO,
		updatedAtVO,
		nil,
	)

	if v, ok := bucketProps["version"].(int64); ok {
		bucket.SetVersion(int(v))
	}

	return bucket, nil
}
