package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

type Publisher struct {
	store    Store
	fallback *events.Producer
	logger   *logger.Logger
}

func NewPublisher(store Store, fallback *events.Producer) *Publisher {
	return &Publisher{store: store, fallback: fallback, logger: logger.WithService("outbox-publisher")}
}

func (p *Publisher) Publish(ctx context.Context, topic string, ev *events.Event) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if p.store == nil {
		return p.fallback.SendEvent(ctx, topic, ev)
	}
	_, err = p.store.Enqueue(ctx, topic, b, nil)
	return err
}

type Dispatcher struct {
	store  Store
	prod   *events.Producer
	logger *logger.Logger
	quitCh chan struct{}
	closed bool
}

func NewDispatcher(store Store, prod *events.Producer) *Dispatcher {
	return &Dispatcher{store: store, prod: prod, logger: logger.WithService("outbox-dispatcher"), quitCh: make(chan struct{}, 1)}
}

func (d *Dispatcher) Start(ctx context.Context) {
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-d.quitCh:
				return
			case <-t.C:
				recs, err := d.store.DequeueBatch(ctx, 50)
				if err != nil {
					d.logger.WithError(err).Error("outbox dequeue failed")
					continue
				}
				for _, r := range recs {
					var ev events.Event
					if err := json.Unmarshal(r.Payload, &ev); err != nil {
						d.logger.WithError(err).Error("outbox payload unmarshal failed")
						continue
					}
					if err := d.prod.SendEvent(ctx, r.Topic, &ev); err != nil {
						d.logger.WithError(err).Error("outbox publish failed", "topic", r.Topic)
						continue
					}
					_ = d.store.MarkDispatched(ctx, r.ID)
				}
			}
		}
	}()
}

func (d *Dispatcher) Stop() {
	if !d.closed {
		d.closed = true
		d.quitCh <- struct{}{}
	}
}
