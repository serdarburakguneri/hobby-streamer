package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
)

func (c *TranscoderEventConsumer) unmarshalEventData(event *events.Event, target interface{}) error {
	var dataBytes []byte
	switch v := event.Data.(type) {
	case []byte:
		dataBytes = v
	case map[string]interface{}:
		var err error
		dataBytes, err = json.Marshal(v)
		if err != nil {
			c.logger.WithError(err).Error("Failed to marshal event data")
			return err
		}
	default:
		c.logger.Error("Unexpected event data type", "type", fmt.Sprintf("%T", event.Data))
		return fmt.Errorf("unexpected event data type: %T", event.Data)
	}
	return json.Unmarshal(dataBytes, target)
}
