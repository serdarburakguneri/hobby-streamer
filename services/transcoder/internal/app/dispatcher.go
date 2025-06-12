package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/services/transcoder/internal/job"
	"github.com/serdarburakguneri/hobby-streamer/services/transcoder/internal/queue"
)

type Dispatcher struct {
	r *job.Registry
}

func NewDispatcher(r *job.Registry) *Dispatcher {
	return &Dispatcher{r: r}
}

func (d *Dispatcher) HandleMessage(msg queue.QueueMessage) error {
	runner, ok := d.r.Get(msg.Type)
	if !ok {
		return fmt.Errorf("no runner registered for type: %s", msg.Type)
	}
	return runner.Run(context.Background(), msg.Payload)
}