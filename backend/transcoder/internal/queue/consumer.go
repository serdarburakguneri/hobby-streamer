package queue

import "context"

type QueueConsumer interface {
	Start(ctx context.Context, handle func(QueueMessage) error)
}
