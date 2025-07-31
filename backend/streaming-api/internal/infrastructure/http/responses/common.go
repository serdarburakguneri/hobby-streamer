package responses

import (
	"time"

	assetvalueobjects "github.com/serdarburakguneri/hobby-streamer/backend/streaming-api/internal/domain/asset/valueobjects"
)

type S3ObjectResponse struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	URL    string `json:"url"`
}

type StreamInfoResponse struct {
	DownloadURL *string `json:"downloadUrl,omitempty"`
	CDNPrefix   *string `json:"cdnPrefix,omitempty"`
	URL         *string `json:"url,omitempty"`
}

type PublishRuleResponse struct {
	PublishAt   *time.Time `json:"publishAt,omitempty"`
	UnpublishAt *time.Time `json:"unpublishAt,omitempty"`
	Regions     []string   `json:"regions,omitempty"`
	AgeRating   *string    `json:"ageRating,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func convertS3ObjectToResponse(s3Obj assetvalueobjects.S3ObjectValue) S3ObjectResponse {
	return S3ObjectResponse{
		Bucket: s3Obj.Bucket(),
		Key:    s3Obj.Key(),
		URL:    s3Obj.URL(),
	}
}

func convertS3ObjectToResponsePtr(s3Obj *assetvalueobjects.S3ObjectValue) *S3ObjectResponse {
	if s3Obj == nil {
		return nil
	}
	return &S3ObjectResponse{
		Bucket: s3Obj.Bucket(),
		Key:    s3Obj.Key(),
		URL:    s3Obj.URL(),
	}
}

func convertStreamInfoToResponse(streamInfo *assetvalueobjects.StreamInfoValue) *StreamInfoResponse {
	if streamInfo == nil {
		return nil
	}
	return &StreamInfoResponse{
		DownloadURL: streamInfo.DownloadURL(),
		CDNPrefix:   streamInfo.CDNPrefix(),
		URL:         streamInfo.URL(),
	}
}

func convertPublishRuleToResponse(rule *assetvalueobjects.PublishRuleValue) *PublishRuleResponse {
	if rule == nil {
		return nil
	}

	return &PublishRuleResponse{
		PublishAt:   rule.PublishAt(),
		UnpublishAt: rule.UnpublishAt(),
		Regions:     rule.Regions(),
		AgeRating:   rule.AgeRating(),
	}
}
