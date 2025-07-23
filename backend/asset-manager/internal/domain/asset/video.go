package asset

import (
	"fmt"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type VideoType string
type VideoFormat string
type VideoStatus string
type VideoQuality string

type Video struct {
	id                 string
	label              string
	videoType          VideoType
	format             VideoFormat
	storageLocation    S3Object
	width              int
	height             int
	duration           float64
	bitrate            int
	codec              string
	size               int64
	contentType        string
	streamInfo         *StreamInfo
	metadata           map[string]string
	status             VideoStatus
	thumbnail          *Image
	transcodingInfo    TranscodingInfo
	createdAt          time.Time
	updatedAt          time.Time
	segmentCount       int
	videoCodec         string
	audioCodec         string
	avgSegmentDuration float64
	segments           []string
	frameRate          string
	audioChannels      int
	audioSampleRate    int
}

func NewVideoFormat(value string) (*VideoFormat, error) {
	switch VideoFormat(value) {
	case VideoFormat(constants.VideoStreamingFormatRaw), VideoFormat(constants.VideoStreamingFormatHLS), VideoFormat(constants.VideoStreamingFormatDASH):
		format := VideoFormat(value)
		return &format, nil
	default:
		return nil, fmt.Errorf("invalid video format: %s", value)
	}
}

func NewVideo(label string, format *VideoFormat, storageLocation S3Object, segmentCount int, videoCodec, audioCodec string, avgSegmentDuration float64, segments []string) *Video {
	now := time.Now().UTC()
	video := &Video{
		id:                 generateID(),
		label:              label,
		videoType:          VideoType(constants.VideoTypeMain),
		storageLocation:    storageLocation,
		status:             VideoStatus(constants.VideoStatusPending),
		metadata:           make(map[string]string),
		createdAt:          now,
		updatedAt:          now,
		segmentCount:       segmentCount,
		videoCodec:         videoCodec,
		audioCodec:         audioCodec,
		avgSegmentDuration: avgSegmentDuration,
		segments:           segments,
	}

	if format != nil {
		video.format = *format
	}

	return video
}

func ReconstructVideo(
	id string,
	label string,
	videoType VideoType,
	format VideoFormat,
	storageLocation S3Object,
	width int,
	height int,
	duration float64,
	bitrate int,
	codec string,
	size int64,
	contentType string,
	status VideoStatus,
	createdAt time.Time,
	updatedAt time.Time,
	segmentCount int,
	videoCodec string,
	audioCodec string,
	avgSegmentDuration float64,
	segments []string,
	frameRate string,
	audioChannels int,
	audioSampleRate int,
) *Video {
	return &Video{
		id:                 id,
		label:              label,
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
		status:             status,
		metadata:           make(map[string]string),
		createdAt:          createdAt,
		updatedAt:          updatedAt,
		segmentCount:       segmentCount,
		videoCodec:         videoCodec,
		audioCodec:         audioCodec,
		avgSegmentDuration: avgSegmentDuration,
		segments:           segments,
		frameRate:          frameRate,
		audioChannels:      audioChannels,
		audioSampleRate:    audioSampleRate,
	}
}

func (v *Video) ID() string {
	return v.id
}

func (v *Video) Label() string {
	return v.label
}

func (v *Video) Type() VideoType {
	return v.videoType
}

func (v *Video) Format() VideoFormat {
	return v.format
}

func (v *Video) StorageLocation() S3Object {
	return v.storageLocation
}

func (v *Video) Width() int {
	return v.width
}

func (v *Video) Height() int {
	return v.height
}

func (v *Video) Duration() float64 {
	return v.duration
}

func (v *Video) Bitrate() int {
	return v.bitrate
}

func (v *Video) Codec() string {
	return v.codec
}

func (v *Video) Size() int64 {
	return v.size
}

func (v *Video) ContentType() string {
	return v.contentType
}

func (v *Video) StreamInfo() *StreamInfo {
	return v.streamInfo
}

func (v *Video) Metadata() map[string]string {
	return v.metadata
}

func (v *Video) Status() VideoStatus {
	return v.status
}

func (v *Video) Thumbnail() *Image {
	return v.thumbnail
}

func (v *Video) TranscodingInfo() TranscodingInfo {
	return v.transcodingInfo
}

func (v *Video) CreatedAt() time.Time {
	return v.createdAt
}

func (v *Video) UpdatedAt() time.Time {
	return v.updatedAt
}

func (v *Video) UpdateStatus(status VideoStatus) {
	v.status = status
	v.updatedAt = time.Now().UTC()
}

func (v *Video) UpdateTranscodingInfo(info TranscodingInfo) {
	v.transcodingInfo = info
	v.updatedAt = time.Now().UTC()
}

func (v *Video) UpdateDimensions(width, height int) {
	v.width = width
	v.height = height
	v.updatedAt = time.Now().UTC()
}

func (v *Video) UpdateDuration(duration float64) {
	v.duration = duration
	v.updatedAt = time.Now().UTC()
}

func (v *Video) UpdateBitrate(bitrate int) {
	v.bitrate = bitrate
	v.updatedAt = time.Now().UTC()
}

func (v *Video) UpdateCodec(codec string) {
	v.codec = codec
	v.updatedAt = time.Now().UTC()
}

func (v *Video) UpdateSize(size int64) {
	v.size = size
	v.updatedAt = time.Now().UTC()
}

func (v *Video) UpdateContentType(contentType string) {
	v.contentType = contentType
	v.updatedAt = time.Now().UTC()
}

func (v *Video) SetStreamInfo(streamInfo *StreamInfo) {
	v.streamInfo = streamInfo
	v.updatedAt = time.Now().UTC()
}

func (v *Video) SetMetadata(metadata map[string]string) {
	v.metadata = metadata
	v.updatedAt = time.Now().UTC()
}

func (v *Video) SetThumbnail(thumbnail *Image) {
	v.thumbnail = thumbnail
	v.updatedAt = time.Now().UTC()
}

func (v *Video) SetType(videoType VideoType) {
	v.videoType = videoType
	v.updatedAt = time.Now().UTC()
}

func (v *Video) SetStorageLocation(s3obj S3Object) {
	v.storageLocation = s3obj
	v.updatedAt = time.Now().UTC()
}

func (v *Video) Quality() VideoQuality {
	if v.width >= 3840 {
		return VideoQuality(constants.VideoQuality4K)
	} else if v.width >= 1920 {
		return VideoQuality(constants.VideoQualityFHD)
	} else if v.width >= 1280 {
		return VideoQuality(constants.VideoQualityHD)
	}
	return VideoQuality(constants.VideoQualitySD)
}

func (v *Video) IsReady() bool {
	return v.status == constants.VideoStatusReady
}

func (v *Video) IsProcessing() bool {
	return v.status == constants.VideoStatusAnalyzing || v.status == constants.VideoStatusTranscoding
}

func (v *Video) IsFailed() bool {
	return v.status == constants.VideoStatusFailed
}

func (v *Video) SegmentCount() int {
	return v.segmentCount
}
func (v *Video) SetSegmentCount(count int) {
	v.segmentCount = count
	v.updatedAt = time.Now().UTC()
}
func (v *Video) VideoCodec() string {
	return v.videoCodec
}
func (v *Video) SetVideoCodec(codec string) {
	v.videoCodec = codec
	v.updatedAt = time.Now().UTC()
}
func (v *Video) AudioCodec() string {
	return v.audioCodec
}
func (v *Video) SetAudioCodec(codec string) {
	v.audioCodec = codec
	v.updatedAt = time.Now().UTC()
}
func (v *Video) AvgSegmentDuration() float64 {
	return v.avgSegmentDuration
}
func (v *Video) SetAvgSegmentDuration(dur float64) {
	v.avgSegmentDuration = dur
	v.updatedAt = time.Now().UTC()
}
func (v *Video) Segments() []string {
	return v.segments
}
func (v *Video) SetSegments(segs []string) {
	v.segments = segs
	v.updatedAt = time.Now().UTC()
}

func (v *Video) FrameRate() string {
	return v.frameRate
}

func (v *Video) SetFrameRate(frameRate string) {
	v.frameRate = frameRate
	v.updatedAt = time.Now().UTC()
}

func (v *Video) AudioChannels() int {
	return v.audioChannels
}

func (v *Video) SetAudioChannels(channels int) {
	v.audioChannels = channels
	v.updatedAt = time.Now().UTC()
}

func (v *Video) AudioSampleRate() int {
	return v.audioSampleRate
}

func (v *Video) SetAudioSampleRate(sampleRate int) {
	v.audioSampleRate = sampleRate
	v.updatedAt = time.Now().UTC()
}
