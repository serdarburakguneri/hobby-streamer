package job

import (
	"context"
	"testing"

	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/valueobjects"
)

type TestEventPublisher struct{}

func (m *TestEventPublisher) PublishJobCompleted(ctx context.Context, event events.CompletedEvent) error {
	return nil
}

// noop storage implementation for testing
type nopStorage struct{}

func (nopStorage) Download(ctx context.Context, input string) (string, error) { return "", nil }
func (nopStorage) CreateDir(path string) error                                { return nil }
func (nopStorage) Remove(path string) error                                   { return nil }
func (nopStorage) RemoveAll(path string) error                                { return nil }
func (nopStorage) Upload(ctx context.Context, localDir, s3Path string) error  { return nil }

// noop registry and strategy for testing
type nopRegistry struct{}

func (nopRegistry) Get(format string) TranscodeStrategy { return &nopStrategy{} }

type nopStrategy struct{}

func (n *nopStrategy) ValidateInput(ctx context.Context, job *entity.Job) error { return nil }
func (n *nopStrategy) Transcode(ctx context.Context, job *entity.Job, localPath, outputDir string) (string, error) {
	return "", nil
}
func (n *nopStrategy) ExtractMetadata(ctx context.Context, filePath string, job *entity.Job) (*valueobjects.TranscodeMetadata, error) {
	return &valueobjects.TranscodeMetadata{}, nil
}
func (n *nopStrategy) ValidateOutput(job *entity.Job) error { return nil }

func TestJobDomainService_ProcessJob_Errors(t *testing.T) {
	testPublisher := &TestEventPublisher{}
	ds := NewDomainService(nopStorage{}, nopRegistry{}, testPublisher)
	ctx := context.Background()

	assetID, _ := valueobjects.NewAssetID("aid")
	videoID, _ := valueobjects.NewVideoID("vid")

	tests := []struct {
		name    string
		job     *entity.Job
		wantErr bool
	}{
		{
			name:    "analyze job missing input",
			job:     entity.NewAnalyzeJob(*assetID, *videoID, ""),
			wantErr: true,
		},
		{
			name:    "transcode job missing output",
			job:     entity.NewTranscodeJob(*assetID, *videoID, "input.mp4", "", "1080p", valueobjects.JobFormatHLS),
			wantErr: true,
		},
		{
			name:    "transcode job missing format",
			job:     entity.NewTranscodeJob(*assetID, *videoID, "input.mp4", "s3://bucket/key", "1080p", valueobjects.JobFormat("")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ds.ProcessJob(ctx, tt.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
