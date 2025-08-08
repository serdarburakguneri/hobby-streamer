package entity

import (
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/asset/valueobjects"
)

type Video struct {
	id                 valueobjects.ID
	label              valueobjects.ValidatedString
	videoType          valueobjects.VideoType
	format             valueobjects.VideoFormat
	storageLocation    valueobjects.S3Object
	width              int
	height             int
	duration           float64
	bitrate            int
	codec              string
	size               int64
	contentType        valueobjects.ContentType
	status             valueobjects.VideoStatus
	timestamps         *valueobjects.Timestamps
	segmentCount       int
	videoCodec         string
	audioCodec         string
	avgSegmentDuration float64
	segments           []string
	frameRate          string
	audioChannels      int
	audioSampleRate    int
	streamInfo         *valueobjects.StreamInfo
}

func NewVideo(
	label string,
	format *valueobjects.VideoFormat,
	storageLocation valueobjects.S3Object,
	width, height int,
	duration float64,
	bitrate int,
	codec string,
	size int64,
	contentType string,
	videoCodec, audioCodec string,
	frameRate string,
	audioChannels, audioSampleRate int,
	streamInfo *valueobjects.StreamInfo,
) (*Video, error) {
	videoLabel, err := valueobjects.NewValidatedString(label, 100, "video label")
	if err != nil {
		return nil, err
	}

	contentTypeVO, err := valueobjects.NewContentType(contentType)
	if err != nil {
		return nil, err
	}

	videoIDPtr, err := valueobjects.GenerateVideoID()
	if err != nil {
		return nil, err
	}
	videoID := *videoIDPtr
	timestamps := valueobjects.NewTimestamps()

	return &Video{
		id:              videoID,
		label:           *videoLabel,
		videoType:       valueobjects.VideoTypeMain,
		format:          *format,
		storageLocation: storageLocation,
		width:           width,
		height:          height,
		duration:        duration,
		bitrate:         bitrate,
		codec:           codec,
		size:            size,
		contentType:     *contentTypeVO,
		status:          valueobjects.VideoStatusReady,
		timestamps:      timestamps,
		segments:        make([]string, 0),
		videoCodec:      videoCodec,
		audioCodec:      audioCodec,
		frameRate:       frameRate,
		audioChannels:   audioChannels,
		audioSampleRate: audioSampleRate,
		streamInfo:      streamInfo,
	}, nil
}

func ReconstructVideo(
	id valueobjects.ID,
	label valueobjects.ValidatedString,
	videoType valueobjects.VideoType,
	format valueobjects.VideoFormat,
	storageLocation valueobjects.S3Object,
	width, height int,
	duration float64,
	bitrate int,
	codec string,
	size int64,
	contentType valueobjects.ContentType,
	status valueobjects.VideoStatus,
	timestamps *valueobjects.Timestamps,
	segmentCount int,
	videoCodec, audioCodec string,
	avgSegmentDuration float64,
	segments []string,
	frameRate string,
	audioChannels, audioSampleRate int,
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
		timestamps:         timestamps,
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

func (v *Video) ID() valueobjects.ID                    { return v.id }
func (v *Video) Label() valueobjects.ValidatedString    { return v.label }
func (v *Video) Type() valueobjects.VideoType           { return v.videoType }
func (v *Video) Format() valueobjects.VideoFormat       { return v.format }
func (v *Video) StorageLocation() valueobjects.S3Object { return v.storageLocation }
func (v *Video) Width() int                             { return v.width }
func (v *Video) Height() int                            { return v.height }
func (v *Video) Duration() float64                      { return v.duration }
func (v *Video) Bitrate() int                           { return v.bitrate }
func (v *Video) Codec() string                          { return v.codec }
func (v *Video) Size() int64                            { return v.size }
func (v *Video) ContentType() valueobjects.ContentType  { return v.contentType }
func (v *Video) Status() valueobjects.VideoStatus       { return v.status }
func (v *Video) Timestamps() *valueobjects.Timestamps   { return v.timestamps }
func (v *Video) SegmentCount() int                      { return v.segmentCount }
func (v *Video) VideoCodec() string                     { return v.videoCodec }
func (v *Video) AudioCodec() string                     { return v.audioCodec }
func (v *Video) AvgSegmentDuration() float64            { return v.avgSegmentDuration }
func (v *Video) Segments() []string                     { return v.segments }
func (v *Video) FrameRate() string                      { return v.frameRate }
func (v *Video) AudioChannels() int                     { return v.audioChannels }
func (v *Video) AudioSampleRate() int                   { return v.audioSampleRate }
func (v *Video) StreamInfo() *valueobjects.StreamInfo   { return v.streamInfo }
func (v *Video) CreatedAt() time.Time                   { return v.timestamps.CreatedAt() }
func (v *Video) UpdatedAt() time.Time                   { return v.timestamps.UpdatedAt() }

func (v *Video) SetStreamInfo(streamInfo *valueobjects.StreamInfo) {
	v.streamInfo = streamInfo
	v.timestamps.Update()
}

func (v *Video) UpdateStatus(status valueobjects.VideoStatus) {
	v.status = status
	v.timestamps.Update()
}

func (v *Video) UpdateSize(size int64) {
	v.size = size
	v.timestamps.Update()
}

func (v *Video) UpdateMediaInfo(info valueobjects.TranscodingInfo) {
	v.width = info.Width()
	v.height = info.Height()
	v.duration = info.Duration()
	v.bitrate = info.Bitrate()
	v.codec = info.Codec()
	v.size = info.Size()
	v.contentType = info.ContentType()
	v.videoCodec = info.VideoCodec()
	v.audioCodec = info.AudioCodec()
	v.frameRate = info.FrameRate()
	v.audioChannels = info.AudioChannels()
	v.audioSampleRate = info.AudioSampleRate()
	v.timestamps.Update()
}

func (v *Video) IsReady() bool      { return v.status.IsReady() }
func (v *Video) IsFailed() bool     { return v.status.IsFailed() }
func (v *Video) IsProcessing() bool { return v.status.IsProcessing() }

func (v *Video) UpdateStorageLocation(location valueobjects.S3Object) {
	v.storageLocation = location
	v.timestamps.Update()
}

// removed UpdateTechnicalDetails in favor of UpdateMediaInfo
