package job

import (
	"context"
	"testing"
)

func TestJobDomainService_ProcessJob_Errors(t *testing.T) {
	ds := NewJobDomainService()
	ctx := context.Background()

	tests := []struct {
		name    string
		job     *Job
		wantErr bool
	}{
		{
			name:    "analyze job missing assetID",
			job:     NewAnalyzeJob(AssetID{""}, VideoID{"vid"}, "input.mp4"),
			wantErr: true,
		},
		{
			name:    "analyze job missing videoID",
			job:     NewAnalyzeJob(AssetID{"aid"}, VideoID{""}, "input.mp4"),
			wantErr: true,
		},
		{
			name:    "analyze job missing input",
			job:     NewAnalyzeJob(AssetID{"aid"}, VideoID{"vid"}, ""),
			wantErr: true,
		},
		{
			name:    "transcode job missing output",
			job:     NewTranscodeJob(AssetID{"aid"}, VideoID{"vid"}, "input.mp4", "", JobFormatHLS),
			wantErr: true,
		},
		{
			name:    "transcode job missing format",
			job:     NewTranscodeJob(AssetID{"aid"}, VideoID{"vid"}, "input.mp4", "s3://bucket/key", ""),
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
