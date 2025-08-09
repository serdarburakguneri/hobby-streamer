package outbox

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Record struct {
	ID        string
	Topic     string
	Payload   []byte
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Store interface {
	Enqueue(ctx context.Context, topic string, payload []byte, headers map[string]string) (string, error)
	DequeueBatch(ctx context.Context, limit int) ([]Record, error)
	MarkDispatched(ctx context.Context, id string) error
}

type Neo4jStore struct {
	driver neo4j.Driver
}

func NewNeo4jStore(driver neo4j.Driver) *Neo4jStore {
	return &Neo4jStore{driver: driver}
}

func (s *Neo4jStore) Enqueue(ctx context.Context, topic string, payload []byte, headers map[string]string) (string, error) {
	session := s.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	id := time.Now().UTC().Format("20060102150405.000000000")
	_, err := session.Run(`
        MERGE (o:Outbox {id: $id})
        SET o.topic = $topic, o.payload = $payload, o.status = 'pending', o.createdAt = timestamp(), o.updatedAt = timestamp()
    `, map[string]interface{}{
		"id":      id,
		"topic":   topic,
		"payload": string(payload),
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Neo4jStore) DequeueBatch(ctx context.Context, limit int) ([]Record, error) {
	session := s.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	result, err := session.Run(`
        MATCH (o:Outbox {status: 'pending'})
        WITH o ORDER BY o.createdAt ASC LIMIT $limit
        SET o.status = 'processing', o.updatedAt = timestamp()
        RETURN o.id as id, o.topic as topic, o.payload as payload
    `, map[string]interface{}{"limit": limit})
	if err != nil {
		return nil, err
	}
	records := make([]Record, 0)
	for result.Next() {
		rec := Record{ID: result.Record().Values[0].(string), Topic: result.Record().Values[1].(string)}
		switch v := result.Record().Values[2].(type) {
		case []byte:
			rec.Payload = v
		case string:
			rec.Payload = []byte(v)
		}
		records = append(records, rec)
	}
	return records, nil
}

func (s *Neo4jStore) MarkDispatched(ctx context.Context, id string) error {
	session := s.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	_, err := session.Run(`
        MATCH (o:Outbox {id: $id})
        SET o.status = 'dispatched', o.updatedAt = timestamp()
    `, map[string]interface{}{"id": id})
	return err
}
