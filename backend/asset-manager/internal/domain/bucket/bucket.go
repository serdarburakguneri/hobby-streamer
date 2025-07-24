package bucket

import (
	"regexp"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type Bucket struct {
	id          BucketID
	key         string
	name        string
	description *string
	bucketType  *string
	status      *string
	ownerID     *string
	metadata    map[string]interface{}
	createdAt   time.Time
	updatedAt   time.Time
}

func NewBucket(name, key string) (*Bucket, error) {
	if name == "" {
		return nil, pkgerrors.NewValidationError("bucket name cannot be empty", nil)
	}

	if !isValidKey(key) {
		return nil, pkgerrors.NewValidationError("invalid bucket key format", nil)
	}

	now := time.Now().UTC()
	id, _ := NewBucketID(generateID())
	return &Bucket{
		id:          *id,
		key:         key,
		name:        name,
		description: nil,
		bucketType:  nil,
		status:      nil,
		ownerID:     nil,
		metadata:    make(map[string]interface{}),
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstructBucket(
	id BucketID,
	name string,
	description *string,
	key string,
	ownerID *string,
	status *string,
	metadata map[string]interface{},
	createdAt time.Time,
	updatedAt time.Time,
) *Bucket {
	return &Bucket{
		id:          id,
		key:         key,
		name:        name,
		description: description,
		bucketType:  nil,
		status:      status,
		ownerID:     ownerID,
		metadata:    metadata,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (b *Bucket) ID() BucketID {
	return b.id
}

func (b *Bucket) Name() string {
	return b.name
}

func (b *Bucket) Description() *string {
	return b.description
}

func (b *Bucket) Key() string {
	return b.key
}

func (b *Bucket) Type() *string {
	return b.bucketType
}

func (b *Bucket) Status() *string {
	return b.status
}

func (b *Bucket) OwnerID() *string {
	return b.ownerID
}

func (b *Bucket) Metadata() map[string]interface{} {
	return b.metadata
}

func (b *Bucket) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Bucket) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *Bucket) UpdateName(name string) error {
	if name == "" {
		return pkgerrors.NewValidationError("bucket name cannot be empty", nil)
	}
	b.name = name
	b.updatedAt = time.Now().UTC()
	return nil
}

func (b *Bucket) UpdateDescription(description *string) {
	b.description = description
	b.updatedAt = time.Now().UTC()
}

func (b *Bucket) SetOwnerID(ownerID *string) {
	b.ownerID = ownerID
	b.updatedAt = time.Now().UTC()
}

func (b *Bucket) SetMetadata(metadata map[string]interface{}) {
	b.metadata = metadata
	b.updatedAt = time.Now().UTC()
}

func (b *Bucket) SetStatus(status *string) {
	b.status = status
	b.updatedAt = time.Now().UTC()
}

func (b *Bucket) ValidateName() error {
	if b.name == "" {
		return pkgerrors.NewValidationError("bucket name cannot be empty", nil)
	}
	if len(b.name) > 100 {
		return pkgerrors.NewValidationError("bucket name too long", nil)
	}
	return nil
}

func (b *Bucket) ValidateKey() error {
	if !isValidKey(b.key) {
		return pkgerrors.NewValidationError("invalid bucket key format", nil)
	}
	return nil
}

func (b *Bucket) CanAddAsset(assetID string, hasAssetFunc func(bucketID, assetID string) (bool, error)) error {
	exists, err := hasAssetFunc(b.id.Value(), assetID)
	if err != nil {
		return pkgerrors.WithContext(err, map[string]interface{}{"operation": "CanAddAsset", "bucketID": b.id.Value(), "assetID": assetID})
	}
	if exists {
		return pkgerrors.NewValidationError("asset already exists in bucket", nil)
	}
	return nil
}

func (b *Bucket) ValidateNotEmpty(assetCountFunc func(bucketID string) (int, error)) error {
	count, err := assetCountFunc(b.id.Value())
	if err != nil {
		return pkgerrors.WithContext(err, map[string]interface{}{"operation": "ValidateNotEmpty", "bucketID": b.id.Value()})
	}
	if count == 0 {
		return pkgerrors.NewValidationError("cannot perform operation on empty bucket", nil)
	}
	return nil
}

func (b *Bucket) ValidateOwnership(ownerID string) error {
	// if b.ownerID == nil || *b.ownerID != ownerID {
	//     return pkgerrors.NewForbiddenError("unauthorized access to bucket", nil)
	// }
	return nil
}

func isValidKey(key string) bool {
	if len(key) < 3 || len(key) > 50 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, key)
	return matched
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
