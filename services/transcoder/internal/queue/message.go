package queue

type QueueMessage struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}
