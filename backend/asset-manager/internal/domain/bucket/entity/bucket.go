package entity

import (
	"errors"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/bucket/valueobjects"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/operations"
)

type Bucket struct {
	id          valueobjects.BucketID
	version     int
	key         valueobjects.BucketKey
	name        valueobjects.BucketName
	description *valueobjects.BucketDescription
	bucketType  *valueobjects.BucketType
	status      *valueobjects.BucketStatus
	ownerID     *valueobjects.OwnerID
	metadata    map[string]interface{}
	createdAt   valueobjects.CreatedAt
	updatedAt   valueobjects.UpdatedAt
	assetIDs    []string
}

func NewBucket(name, key string) (*Bucket, error) {
	bucketName, err := valueobjects.NewBucketName(name)
	if err != nil {
		return nil, err
	}

	bucketKey, err := valueobjects.NewBucketKey(key)
	if err != nil {
		return nil, err
	}

	bucketID, err := valueobjects.NewBucketID(operations.GenerateID())
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Bucket{
		id:          *bucketID,
		version:     0,
		key:         *bucketKey,
		name:        *bucketName,
		description: nil,
		bucketType:  nil,
		status:      nil,
		ownerID:     nil,
		metadata:    make(map[string]interface{}),
		createdAt:   valueobjects.NewCreatedAt(now),
		updatedAt:   valueobjects.NewUpdatedAt(now),
		assetIDs:    make([]string, 0),
	}, nil
}

func NewBucketWithProperties(name, key string, description *string, bucketType *string, status *string, ownerID *string, metadata map[string]interface{}) (*Bucket, error) {
	bucketName, err := valueobjects.NewBucketName(name)
	if err != nil {
		return nil, err
	}

	bucketKey, err := valueobjects.NewBucketKey(key)
	if err != nil {
		return nil, err
	}

	bucketID, err := valueobjects.NewBucketID(operations.GenerateID())
	if err != nil {
		return nil, err
	}

	var bucketDescription *valueobjects.BucketDescription
	if description != nil {
		bucketDescription, err = valueobjects.NewBucketDescription(*description)
		if err != nil {
			return nil, err
		}
	}

	var bucketTypeVO *valueobjects.BucketType
	if bucketType != nil {
		bucketTypeVO, err = valueobjects.NewBucketType(*bucketType)
		if err != nil {
			return nil, err
		}
	}

	var bucketStatus *valueobjects.BucketStatus
	if status != nil {
		bucketStatus, err = valueobjects.NewBucketStatus(*status)
		if err != nil {
			return nil, err
		}
	}

	var ownerIDVO *valueobjects.OwnerID
	if ownerID != nil {
		ownerIDVO, err = valueobjects.NewOwnerID(*ownerID)
		if err != nil {
			return nil, err
		}
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	now := time.Now().UTC()
	return &Bucket{
		id:          *bucketID,
		version:     0,
		key:         *bucketKey,
		name:        *bucketName,
		description: bucketDescription,
		bucketType:  bucketTypeVO,
		status:      bucketStatus,
		ownerID:     ownerIDVO,
		metadata:    metadata,
		createdAt:   valueobjects.NewCreatedAt(now),
		updatedAt:   valueobjects.NewUpdatedAt(now),
		assetIDs:    make([]string, 0),
	}, nil
}

func ReconstructBucket(
	id valueobjects.BucketID,
	name valueobjects.BucketName,
	description *valueobjects.BucketDescription,
	key valueobjects.BucketKey,
	bucketType *valueobjects.BucketType,
	ownerID *valueobjects.OwnerID,
	status *valueobjects.BucketStatus,
	metadata map[string]interface{},
	createdAt valueobjects.CreatedAt,
	updatedAt valueobjects.UpdatedAt,
	assetIDs []string,
) *Bucket {
	return &Bucket{
		id:          id,
		version:     0,
		key:         key,
		name:        name,
		description: description,
		bucketType:  bucketType,
		status:      status,
		ownerID:     ownerID,
		metadata:    metadata,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		assetIDs:    assetIDs,
	}
}

func (b *Bucket) ID() valueobjects.BucketID {
	return b.id
}

func (b *Bucket) Version() int     { return b.version }
func (b *Bucket) SetVersion(v int) { b.version = v }

func (b *Bucket) Name() valueobjects.BucketName {
	return b.name
}

func (b *Bucket) Description() *valueobjects.BucketDescription {
	return b.description
}

func (b *Bucket) Key() valueobjects.BucketKey {
	return b.key
}

func (b *Bucket) Type() *valueobjects.BucketType {
	return b.bucketType
}

func (b *Bucket) Status() *valueobjects.BucketStatus {
	return b.status
}

func (b *Bucket) OwnerID() *valueobjects.OwnerID {
	return b.ownerID
}

func (b *Bucket) Metadata() map[string]interface{} {
	return b.metadata
}

func (b *Bucket) CreatedAt() valueobjects.CreatedAt {
	return b.createdAt
}

func (b *Bucket) UpdatedAt() valueobjects.UpdatedAt {
	return b.updatedAt
}

func (b *Bucket) AssetIDs() []string {
	return b.assetIDs
}

func (b *Bucket) UpdateName(name valueobjects.BucketName) {
	b.name = name
	b.touch()
}

func (b *Bucket) UpdateDescription(description *valueobjects.BucketDescription) {
	b.description = description
	b.touch()
}

func (b *Bucket) UpdateStatus(status *valueobjects.BucketStatus) {
	b.status = status
	b.touch()
}

func (b *Bucket) UpdateType(bucketType *valueobjects.BucketType) {
	b.bucketType = bucketType
	b.touch()
}

func (b *Bucket) UpdateOwnerID(ownerID *valueobjects.OwnerID) {
	b.ownerID = ownerID
	b.touch()
}

func (b *Bucket) UpdateMetadata(metadata map[string]interface{}) {
	b.metadata = metadata
	b.touch()
}

func (b *Bucket) Activate() error {
	activeStatus, err := valueobjects.NewBucketStatus(constants.StatusActive)
	if err != nil {
		return err
	}
	b.status = activeStatus
	b.touch()
	return nil
}

func (b *Bucket) Deactivate() error {
	inactiveStatus, err := valueobjects.NewBucketStatus(constants.StatusPending)
	if err != nil {
		return err
	}
	b.status = inactiveStatus
	b.touch()
	return nil
}

func (b *Bucket) AddAsset(assetID string) error {
	if b.HasAsset(assetID) {
		return errors.New("asset already exists in bucket")
	}
	b.assetIDs = append(b.assetIDs, assetID)
	b.touch()
	return nil
}

func (b *Bucket) RemoveAsset(assetID string) error {
	for i, id := range b.assetIDs {
		if id == assetID {
			b.assetIDs = append(b.assetIDs[:i], b.assetIDs[i+1:]...)
			b.touch()
			return nil
		}
	}
	return errors.New("asset not found in bucket")
}

func (b *Bucket) HasAsset(assetID string) bool {
	for _, id := range b.assetIDs {
		if id == assetID {
			return true
		}
	}
	return false
}

func (b *Bucket) IsActive() bool {
	return b.status != nil && b.status.Value() == constants.StatusActive
}

func (b *Bucket) IsOwnedBy(userID string) bool {
	return b.ownerID != nil && b.ownerID.Value() == userID
}

func (b *Bucket) touch() {
	b.updatedAt = valueobjects.NewUpdatedAt(time.Now().UTC())
}
