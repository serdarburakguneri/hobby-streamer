package valueobjects

import (
	"errors"
	"time"
)

type PublishRule struct {
	publishAt   *time.Time
	unpublishAt *time.Time
	regions     []string
	ageRating   *string
}

func NewPublishRule(publishAt, unpublishAt *time.Time, regions []string, ageRating *string) (*PublishRule, error) {
	if publishAt != nil && unpublishAt != nil {
		if publishAt.After(*unpublishAt) {
			return nil, errors.New("publish date cannot be after unpublish date")
		}
	}

	if len(regions) > 50 {
		return nil, errors.New("too many regions")
	}

	for _, region := range regions {
		if len(region) > 10 {
			return nil, errors.New("region code too long")
		}
	}

	if ageRating != nil && len(*ageRating) > 10 {
		return nil, errors.New("age rating too long")
	}

	return &PublishRule{
		publishAt:   publishAt,
		unpublishAt: unpublishAt,
		regions:     regions,
		ageRating:   ageRating,
	}, nil
}

func (pr PublishRule) PublishAt() *time.Time {
	return pr.publishAt
}

func (pr PublishRule) UnpublishAt() *time.Time {
	return pr.unpublishAt
}

func (pr PublishRule) Regions() []string {
	return pr.regions
}

func (pr PublishRule) AgeRating() *string {
	return pr.ageRating
}

func (pr PublishRule) IsPublished() bool {
	if pr.publishAt == nil {
		return false
	}
	now := time.Now().UTC()
	if now.Before(*pr.publishAt) {
		return false
	}
	if pr.unpublishAt != nil && now.After(*pr.unpublishAt) {
		return false
	}
	return true
}
