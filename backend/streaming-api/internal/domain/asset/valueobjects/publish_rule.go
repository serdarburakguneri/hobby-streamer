package valueobjects

import (
	"regexp"
	"strings"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type PublishRuleValue struct {
	publishAt   *time.Time
	unpublishAt *time.Time
	regions     []string
	ageRating   *string
}

func NewPublishRuleValue(publishAt, unpublishAt *time.Time, regions []string, ageRating *string) (*PublishRuleValue, error) {
	if publishAt != nil && unpublishAt != nil {
		if publishAt.After(*unpublishAt) {
			return nil, ErrInvalidPublishDates
		}
	}

	if len(regions) > 50 {
		return nil, ErrTooManyRegions
	}

	validatedRegions := make([]string, 0, len(regions))
	for _, region := range regions {
		if len(region) > 10 {
			return nil, ErrInvalidRegion
		}

		regionRegex := regexp.MustCompile(`^[A-Z]{2,3}$`)
		if !regionRegex.MatchString(region) {
			return nil, ErrInvalidRegion
		}

		validatedRegions = append(validatedRegions, strings.ToUpper(region))
	}

	if ageRating != nil {
		validRatings := map[string]bool{
			constants.AgeRatingG: true, constants.AgeRatingPG: true, constants.AgeRatingPG13: true, constants.AgeRatingR: true, constants.AgeRatingNC17: true,
			constants.AgeRatingTVY: true, constants.AgeRatingTVY7: true, constants.AgeRatingTVG: true, constants.AgeRatingTVPG: true, constants.AgeRatingTV14: true, constants.AgeRatingTVMA: true,
		}

		if !validRatings[*ageRating] {
			return nil, ErrInvalidAgeRating
		}
	}

	return &PublishRuleValue{
		publishAt:   publishAt,
		unpublishAt: unpublishAt,
		regions:     validatedRegions,
		ageRating:   ageRating,
	}, nil
}

func (p PublishRuleValue) PublishAt() *time.Time {
	return p.publishAt
}

func (p PublishRuleValue) UnpublishAt() *time.Time {
	return p.unpublishAt
}

func (p PublishRuleValue) Regions() []string {
	return p.regions
}

func (p PublishRuleValue) AgeRating() *string {
	return p.ageRating
}

func (p PublishRuleValue) Equals(other PublishRuleValue) bool {
	return p.publishAt == other.publishAt &&
		p.unpublishAt == other.unpublishAt &&
		p.ageRating == other.ageRating
}

var (
	ErrInvalidPublishDates = pkgerrors.NewValidationError("invalid publish dates", nil)
	ErrTooManyRegions      = pkgerrors.NewValidationError("too many regions", nil)
	ErrInvalidRegion       = pkgerrors.NewValidationError("invalid region", nil)
	ErrInvalidAgeRating    = pkgerrors.NewValidationError("invalid age rating", nil)
)
