package events

const (
	EventNamespace = "com.hobbystreamer"

	AssetCreatedEventType   = EventNamespace + ".asset.created"
	AssetUpdatedEventType   = EventNamespace + ".asset.updated"
	AssetDeletedEventType   = EventNamespace + ".asset.deleted"
	AssetPublishedEventType = EventNamespace + ".asset.published"

	VideoAddedEventType         = EventNamespace + ".video.added"
	VideoRemovedEventType       = EventNamespace + ".video.removed"
	VideoStatusUpdatedEventType = EventNamespace + ".video.status.updated"

	BucketCreatedEventType      = EventNamespace + ".bucket.created"
	BucketUpdatedEventType      = EventNamespace + ".bucket.updated"
	BucketDeletedEventType      = EventNamespace + ".bucket.deleted"
	BucketAssetAddedEventType   = EventNamespace + ".bucket.asset.added"
	BucketAssetRemovedEventType = EventNamespace + ".bucket.asset.removed"

	JobAnalyzeRequestedEventType   = EventNamespace + ".job.analyze.requested"
	JobTranscodeRequestedEventType = EventNamespace + ".job.transcode.requested"
	JobAnalyzeCompletedEventType   = EventNamespace + ".job.analyze.completed"
	JobTranscodeCompletedEventType = EventNamespace + ".job.transcode.completed"

	ContentAnalysisRequestedEventType = EventNamespace + ".content.analysis.requested"
	ContentAnalysisCompletedEventType = EventNamespace + ".content.analysis.completed"
	ContentAnalysisFailedEventType    = EventNamespace + ".content.analysis.failed"
)

const (
	AssetEventsTopic  = "asset-events"
	BucketEventsTopic = "bucket-events"

	ContentAnalysisTopic          = "content-analysis"
	ContentAnalysisRequestedTopic = "content.analysis.requested"
	ContentAnalysisCompletedTopic = "content.analysis.completed"
	ContentAnalysisFailedTopic    = "content.analysis.failed"

	RawVideoUploadedTopic = "raw-video-uploaded"

	AnalyzeJobRequestedTopic = "analyze.job.requested"
	HLSJobRequestedTopic     = "hls.job.requested"
	DASHJobRequestedTopic    = "dash.job.requested"

	AnalyzeJobCompletedTopic = "analyze.job.completed"
	HLSJobCompletedTopic     = "hls.job.completed"
	DASHJobCompletedTopic    = "dash.job.completed"
)

const (
	AssetManagerGroupID    = "asset-manager-group"
	TranscoderGroupID      = "transcoder-group"
	ContentAnalyzerGroupID = "content-analyzer-group"
)

func NewAssetCreatedEvent(assetID, slug, title, assetType string) *Event {
	return NewEvent(AssetCreatedEventType, map[string]interface{}{
		"assetId": assetID,
		"slug":    slug,
		"title":   title,
		"type":    assetType,
	})
}

func NewAssetUpdatedEvent(assetID, slug, title, assetType string) *Event {
	return NewEvent(AssetUpdatedEventType, map[string]interface{}{
		"assetId": assetID,
		"slug":    slug,
		"title":   title,
		"type":    assetType,
	})
}

func NewAssetDeletedEvent(assetID, slug string) *Event {
	return NewEvent(AssetDeletedEventType, map[string]interface{}{
		"assetId": assetID,
		"slug":    slug,
	})
}

func NewAssetPublishedEvent(assetID, slug string) *Event {
	return NewEvent(AssetPublishedEventType, map[string]interface{}{
		"assetId": assetID,
		"slug":    slug,
	})
}

func NewVideoAddedEvent(assetID, videoID, label, format string) *Event {
	return NewEvent(VideoAddedEventType, map[string]interface{}{
		"assetId": assetID,
		"videoId": videoID,
		"label":   label,
		"format":  format,
	})
}

func NewVideoRemovedEvent(assetID, videoID string) *Event {
	return NewEvent(VideoRemovedEventType, map[string]interface{}{
		"assetId": assetID,
		"videoId": videoID,
	})
}

func NewVideoStatusUpdatedEvent(assetID, videoID, status string) *Event {
	return NewEvent(VideoStatusUpdatedEventType, map[string]interface{}{
		"assetId": assetID,
		"videoId": videoID,
		"status":  status,
	})
}

func NewBucketCreatedEvent(bucketID, name, key string) *Event {
	return NewEvent(BucketCreatedEventType, map[string]interface{}{
		"bucketId": bucketID,
		"name":     name,
		"key":      key,
	})
}

func NewBucketUpdatedEvent(bucketID, name, key string) *Event {
	return NewEvent(BucketUpdatedEventType, map[string]interface{}{
		"bucketId": bucketID,
		"name":     name,
		"key":      key,
	})
}

func NewBucketDeletedEvent(bucketID string) *Event {
	return NewEvent(BucketDeletedEventType, map[string]interface{}{
		"bucketId": bucketID,
	})
}

func NewBucketAssetAddedEvent(bucketID, assetID string) *Event {
	return NewEvent(BucketAssetAddedEventType, map[string]interface{}{
		"bucketId": bucketID,
		"assetId":  assetID,
	})
}

func NewBucketAssetRemovedEvent(bucketID, assetID string) *Event {
	return NewEvent(BucketAssetRemovedEventType, map[string]interface{}{
		"bucketId": bucketID,
		"assetId":  assetID,
	})
}

func NewJobAnalyzeRequestedEvent(assetID, videoID, input string) *Event {
	return NewEvent(JobAnalyzeRequestedEventType, map[string]interface{}{
		"assetId": assetID,
		"videoId": videoID,
		"input":   input,
		"jobType": "analyze",
	})
}

func NewJobTranscodeRequestedEvent(assetID, videoID, input, format, outputBucket, outputKey string) *Event {
	return NewEvent(JobTranscodeRequestedEventType, map[string]interface{}{
		"assetId":      assetID,
		"videoId":      videoID,
		"input":        input,
		"format":       format,
		"outputBucket": outputBucket,
		"outputKey":    outputKey,
		"jobType":      "transcode",
	})
}

func NewJobAnalyzeCompletedEvent(assetID, videoID string, success bool, metadata map[string]interface{}, errorMsg string) *Event {
	data := map[string]interface{}{
		"assetId": assetID,
		"videoId": videoID,
		"success": success,
		"jobType": "analyze",
	}

	if success && metadata != nil {
		for key, value := range metadata {
			data[key] = value
		}
	} else if !success {
		data["error"] = errorMsg
	}

	return NewEvent(JobAnalyzeCompletedEventType, data)
}

func NewJobTranscodeCompletedEvent(assetID, videoID, format string, success bool, metadata map[string]interface{}, errorMsg string) *Event {
	data := map[string]interface{}{
		"assetId": assetID,
		"videoId": videoID,
		"format":  format,
		"success": success,
		"jobType": "transcode",
	}

	if success && metadata != nil {
		for key, value := range metadata {
			data[key] = value
		}
	} else if !success {
		data["error"] = errorMsg
	}

	return NewEvent(JobTranscodeCompletedEventType, data)
}

func NewContentAnalysisRequestedEvent(assetID, videoID string) *Event {
	return NewEvent(ContentAnalysisRequestedEventType, map[string]interface{}{
		"assetId": assetID,
		"videoId": videoID,
	})
}

func NewContentAnalysisCompletedEvent(assetID, videoID string, analysis map[string]interface{}) *Event {
	return NewEvent(ContentAnalysisCompletedEventType, map[string]interface{}{
		"assetId":  assetID,
		"videoId":  videoID,
		"analysis": analysis,
	})
}

func NewContentAnalysisFailedEvent(assetID, videoID, errorMsg string) *Event {
	return NewEvent(ContentAnalysisFailedEventType, map[string]interface{}{
		"assetId": assetID,
		"videoId": videoID,
		"error":   errorMsg,
	})
}
