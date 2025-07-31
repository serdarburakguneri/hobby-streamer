package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	CloudEventsVersion = "1.0"
	ContentTypeJSON    = "application/json"
)

type Event struct {
	SpecVersion     string                 `json:"specversion"`
	ID              string                 `json:"id"`
	Source          string                 `json:"source"`
	Type            string                 `json:"type"`
	DataContentType string                 `json:"datacontenttype,omitempty"`
	Time            time.Time              `json:"time"`
	Data            interface{}            `json:"data,omitempty"`
	CorrelationID   string                 `json:"correlationid,omitempty"`
	CausationID     string                 `json:"causationid,omitempty"`
	Extensions      map[string]interface{} `json:"-"`
}

func NewEvent(eventType string, data interface{}) *Event {
	return &Event{
		SpecVersion:     CloudEventsVersion,
		ID:              uuid.New().String(),
		Type:            eventType,
		DataContentType: ContentTypeJSON,
		Time:            time.Now().UTC(),
		Data:            data,
		Extensions:      make(map[string]interface{}),
	}
}

func (e *Event) SetSource(source string) *Event {
	e.Source = source
	return e
}

func (e *Event) SetCorrelationID(correlationID string) *Event {
	e.CorrelationID = correlationID
	return e
}

func (e *Event) SetCausationID(causationID string) *Event {
	e.CausationID = causationID
	return e
}

func (e *Event) AddExtension(key string, value interface{}) *Event {
	if e.Extensions == nil {
		e.Extensions = make(map[string]interface{})
	}
	e.Extensions[key] = value
	return e
}

func (e *Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	base := (*Alias)(e)

	eventMap := make(map[string]interface{})
	baseBytes, err := json.Marshal(base)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(baseBytes, &eventMap); err != nil {
		return nil, err
	}

	for key, value := range e.Extensions {
		eventMap[key] = value
	}

	return json.Marshal(eventMap)
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.Extensions = make(map[string]interface{})
	baseFields := map[string]bool{
		"specversion":     true,
		"id":              true,
		"source":          true,
		"type":            true,
		"datacontenttype": true,
		"time":            true,
		"data":            true,
		"correlationid":   true,
		"causationid":     true,
	}

	for key, value := range raw {
		if !baseFields[key] {
			e.Extensions[key] = value
		}
	}

	return nil
}

func (e *Event) Validate() error {
	if e.SpecVersion == "" {
		return fmt.Errorf("specversion is required")
	}
	if e.ID == "" {
		return fmt.Errorf("id is required")
	}
	if e.Source == "" {
		return fmt.Errorf("source is required")
	}
	if e.Type == "" {
		return fmt.Errorf("type is required")
	}
	if e.Time.IsZero() {
		return fmt.Errorf("time is required")
	}
	return nil
}

func (e *Event) GetDataAs(target interface{}) error {
	if e.Data == nil {
		return fmt.Errorf("event has no data")
	}

	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	return json.Unmarshal(dataBytes, target)
}

func (e *Event) String() string {
	bytes, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("Event{ID: %s, Type: %s, Source: %s}", e.ID, e.Type, e.Source)
	}
	return string(bytes)
}
