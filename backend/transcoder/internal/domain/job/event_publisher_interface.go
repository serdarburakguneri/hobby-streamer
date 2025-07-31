package job

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/events"
)

type EventPublisher interface {
	PublishJobCompleted(ctx context.Context, event *events.JobCompletedEvent) error
}
