package job

import (
	"context"

	"github.com/serdarburakguneri/hobby-streamer/backend/transcoder/internal/domain/job/entity"
)

type DomainService interface {
	ProcessJob(ctx context.Context, job *entity.Job) (interface{}, error)
}
