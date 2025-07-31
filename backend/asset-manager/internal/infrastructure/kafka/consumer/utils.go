package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/events"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger"
)

func unmarshalEvent[T any](logger *logger.Logger, event *events.Event) (*T, error) {
	var result T
	if err := unmarshalEventData(logger, event, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func unmarshalEventData(logger *logger.Logger, event *events.Event, target interface{}) error {
	var dataBytes []byte
	switch v := event.Data.(type) {
	case []byte:
		dataBytes = v
	case map[string]interface{}:
		var err error
		dataBytes, err = json.Marshal(v)
		if err != nil {
			logger.WithError(err).Error("Failed to marshal event data")
			return err
		}
	default:
		logger.Error("Unexpected event data type", "type", fmt.Sprintf("%T", event.Data))
		return fmt.Errorf("unexpected event data type: %T", event.Data)
	}

	return json.Unmarshal(dataBytes, target)
}
