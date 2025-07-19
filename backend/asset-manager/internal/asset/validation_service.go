package asset

import (
	"regexp"

	apperrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (s *ValidationService) ValidateAsset(a *Asset) error {
	log := logger.Get().WithService("validation-service")

	if a.Slug != "" {
		if !s.isValidSlug(a.Slug) {
			log.Error("Invalid slug format", "slug", a.Slug)
			return apperrors.NewValidationError("invalid slug format", nil)
		}
	}

	if a.Type != nil {
		validTypes := []string{
			AssetTypeMovie, AssetTypeSeries, AssetTypeSeason, AssetTypeEpisode,
			AssetTypeDocumentary, AssetTypeMusic, AssetTypePodcast,
			AssetTypeTrailer, AssetTypeBehindTheScenes, AssetTypeInterview,
		}
		if !s.contains(validTypes, *a.Type) {
			log.Error("Invalid asset type", "type", *a.Type, "valid_types", validTypes)
			return apperrors.NewValidationError("invalid type value", nil)
		}
	}

	if a.Genre != nil {
		validGenres := []string{
			AssetGenreAction, AssetGenreDrama, AssetGenreComedy, AssetGenreHorror,
			AssetGenreSciFi, AssetGenreRomance, AssetGenreThriller, AssetGenreFantasy,
			AssetGenreDocumentary, AssetGenreMusic, AssetGenreNews,
			AssetGenreSports, AssetGenreKids, AssetGenreEducational,
		}
		if !s.contains(validGenres, *a.Genre) {
			log.Error("Invalid primary genre", "genre", *a.Genre, "valid_genres", validGenres)
			return apperrors.NewValidationError("invalid genre value", nil)
		}
	}

	if len(a.Genres) > 0 {
		validGenres := []string{
			AssetGenreAction, AssetGenreDrama, AssetGenreComedy, AssetGenreHorror,
			AssetGenreSciFi, AssetGenreRomance, AssetGenreThriller, AssetGenreFantasy,
			AssetGenreDocumentary, AssetGenreMusic, AssetGenreNews,
			AssetGenreSports, AssetGenreKids, AssetGenreEducational,
		}
		for _, genre := range a.Genres {
			if !s.contains(validGenres, genre) {
				log.Error("Invalid additional genre", "genre", genre, "valid_genres", validGenres)
				return apperrors.NewValidationError("invalid genre value in genres array", nil)
			}
		}
	}

	if a.PublishRule != nil {
		if !a.PublishRule.PublishAt.IsZero() && !a.PublishRule.UnpublishAt.IsZero() {
			if a.PublishRule.PublishAt.After(a.PublishRule.UnpublishAt) {
				log.Error("Publish date cannot be after unpublish date")
				return apperrors.NewValidationError("publish date cannot be after unpublish date", nil)
			}
		}
	}

	return nil
}

func (s *ValidationService) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (s *ValidationService) isValidSlug(slug string) bool {
	if slug == "" {
		return false
	}

	matched, _ := regexp.MatchString(`^[a-z0-9-_]+$`, slug)
	return matched
}
