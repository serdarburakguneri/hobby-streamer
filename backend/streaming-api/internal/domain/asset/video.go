package asset

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type Video struct {
	id                 VideoID
	videoType          *VideoType
	format             *VideoFormat
	storageLocation    S3ObjectValue
	width              *int
	height             *int
	duration           *float64
	bitrate            *int
	codec              *string
	size               *int
	contentType        *string
	streamInfo         *StreamInfoValue
	metadata           *string
	status             *string
	thumbnail          *Image
	createdAt          time.Time
	updatedAt          time.Time
	quality            *VideoQuality
	isReady            bool
	isProcessing       bool
	isFailed           bool
	segmentCount       *int
	videoCodec         *string
	audioCodec         *string
	avgSegmentDuration *float64
	segments           []string
	frameRate          *string
	audioChannels      *int
	audioSampleRate    *int
	transcodingInfo    *TranscodingInfo
}

func NewVideo(
	id VideoID,
	videoType *VideoType,
	format *VideoFormat,
	storageLocation S3ObjectValue,
	width *int,
	height *int,
	duration *float64,
	bitrate *int,
	codec *string,
	size *int,
	contentType *string,
	streamInfo *StreamInfoValue,
	metadata *string,
	status *string,
	thumbnail *Image,
	createdAt time.Time,
	updatedAt time.Time,
	quality *VideoQuality,
	isReady bool,
	isProcessing bool,
	isFailed bool,
	segmentCount *int,
	videoCodec *string,
	audioCodec *string,
	avgSegmentDuration *float64,
	segments []string,
	frameRate *string,
	audioChannels *int,
	audioSampleRate *int,
	transcodingInfo *TranscodingInfo,
) *Video {
	return &Video{
		id:                 id,
		videoType:          videoType,
		format:             format,
		storageLocation:    storageLocation,
		width:              width,
		height:             height,
		duration:           duration,
		bitrate:            bitrate,
		codec:              codec,
		size:               size,
		contentType:        contentType,
		streamInfo:         streamInfo,
		metadata:           metadata,
		status:             status,
		thumbnail:          thumbnail,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
		quality:            quality,
		isReady:            isReady,
		isProcessing:       isProcessing,
		isFailed:           isFailed,
		segmentCount:       segmentCount,
		videoCodec:         videoCodec,
		audioCodec:         audioCodec,
		avgSegmentDuration: avgSegmentDuration,
		segments:           segments,
		frameRate:          frameRate,
		audioChannels:      audioChannels,
		audioSampleRate:    audioSampleRate,
		transcodingInfo:    transcodingInfo,
	}
}

func (v *Video) ID() VideoID {
	return v.id
}

func (v *Video) Type() *VideoType {
	return v.videoType
}

func (v *Video) Format() *VideoFormat {
	return v.format
}

func (v *Video) StorageLocation() S3ObjectValue {
	return v.storageLocation
}

func (v *Video) Width() *int {
	return v.width
}

func (v *Video) Height() *int {
	return v.height
}

func (v *Video) Duration() *float64 {
	return v.duration
}

func (v *Video) Bitrate() *int {
	return v.bitrate
}

func (v *Video) Codec() *string {
	return v.codec
}

func (v *Video) Size() *int {
	return v.size
}

func (v *Video) ContentType() *string {
	return v.contentType
}

func (v *Video) StreamInfo() *StreamInfoValue {
	return v.streamInfo
}

func (v *Video) Metadata() *string {
	return v.metadata
}

func (v *Video) Status() *string {
	return v.status
}

func (v *Video) Thumbnail() *Image {
	return v.thumbnail
}

func (v *Video) CreatedAt() time.Time {
	return v.createdAt
}

func (v *Video) UpdatedAt() time.Time {
	return v.updatedAt
}

func (v *Video) IsReady() bool {
	return v.status != nil && *v.status == constants.VideoStatusReady
}

func (v *Video) IsMain() bool {
	return v.videoType != nil && v.videoType.Value() == constants.VideoTypeMain
}

func (v *Video) IsTrailer() bool {
	return v.videoType != nil && v.videoType.Value() == "trailer"
}

func (v *Video) HasStreamInfo() bool {
	return v.streamInfo != nil
}

func (v *Video) GetAspectRatio() *float64 {
	if v.width == nil || v.height == nil || *v.width == 0 || *v.height == 0 {
		return nil
	}

	ratio := float64(*v.width) / float64(*v.height)
	return &ratio
}

func (v *Video) GetDurationInMinutes() *float64 {
	if v.duration == nil {
		return nil
	}

	minutes := *v.duration / 60.0
	return &minutes
}

func (v *Video) IsHD() bool {
	if v.width == nil || v.height == nil {
		return false
	}
	return *v.width >= 1280 && *v.height >= 720
}

func (v *Video) Is4K() bool {
	if v.width == nil || v.height == nil {
		return false
	}
	return *v.width >= 3840 && *v.height >= 2160
}

func (v *Video) Quality() *VideoQuality            { return v.quality }
func (v *Video) IsReadyFlag() bool                 { return v.isReady }
func (v *Video) IsProcessing() bool                { return v.isProcessing }
func (v *Video) IsFailed() bool                    { return v.isFailed }
func (v *Video) SegmentCount() *int                { return v.segmentCount }
func (v *Video) VideoCodec() *string               { return v.videoCodec }
func (v *Video) AudioCodec() *string               { return v.audioCodec }
func (v *Video) AvgSegmentDuration() *float64      { return v.avgSegmentDuration }
func (v *Video) Segments() []string                { return v.segments }
func (v *Video) FrameRate() *string                { return v.frameRate }
func (v *Video) AudioChannels() *int               { return v.audioChannels }
func (v *Video) AudioSampleRate() *int             { return v.audioSampleRate }
func (v *Video) TranscodingInfo() *TranscodingInfo { return v.transcodingInfo }
