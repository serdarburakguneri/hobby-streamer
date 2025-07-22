package bucket

import (
	"errors"
	"time"
)

type Bucket struct {
	id          string
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

func NewBucket(name string, key string, description *string, ownerID *string, status *string) (*Bucket, error) {
	if name == "" {
		return nil, errors.New("bucket name cannot be empty")
	}

	if !isValidKey(key) {
		return nil, errors.New("invalid bucket key format")
	}

	now := time.Now().UTC()
	return &Bucket{
		id:          generateID(),
		key:         key,
		name:        name,
		description: description,
		bucketType:  nil,
		status:      status,
		ownerID:     ownerID,
		metadata:    make(map[string]interface{}),
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstructBucket(
	id string,
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

func (b *Bucket) ID() string {
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
		return errors.New("bucket name cannot be empty")
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
