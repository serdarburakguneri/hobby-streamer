package pipeline

import (
	"context"
	"encoding/json"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	domain "github.com/serdarburakguneri/hobby-streamer/backend/asset-manager/internal/domain/pipeline/entity"
)

type Repository struct {
	driver neo4j.Driver
}

func NewRepository(driver neo4j.Driver) *Repository { return &Repository{driver: driver} }

const upsertQuery = `
MERGE (p:Pipeline {assetId: $assetId, videoId: $videoId})
ON CREATE SET p.createdAt = datetime($now), p.updatedAt = datetime($now), p.steps = $steps
ON MATCH SET p.updatedAt = datetime($now), p.steps = $steps
RETURN p
`

func (r *Repository) Upsert(ctx context.Context, p *domain.Pipeline) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	stepsJSON, _ := json.Marshal(p.Steps)
	_, err := session.Run(upsertQuery, map[string]interface{}{
		"assetId": p.AssetID,
		"videoId": p.VideoID,
		"now":     time.Now().UTC().Format(time.RFC3339),
		"steps":   string(stepsJSON),
	})
	return err
}

const getQuery = `
MATCH (p:Pipeline {assetId: $assetId, videoId: $videoId}) RETURN p
`

func (r *Repository) Get(ctx context.Context, assetID, videoID string) (*domain.Pipeline, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	res, err := session.Run(getQuery, map[string]interface{}{"assetId": assetID, "videoId": videoID})
	if err != nil {
		return nil, err
	}
	if !res.Next() {
		return nil, nil
	}
	rec := res.Record()
	node := rec.Values[0].(neo4j.Node)
	stepsStr, _ := node.Props["steps"].(string)
	steps := map[string]domain.StepState{}
	if stepsStr != "" {
		_ = json.Unmarshal([]byte(stepsStr), &steps)
	}
	p := &domain.Pipeline{AssetID: node.Props["assetId"].(string), VideoID: node.Props["videoId"].(string), Steps: steps}
	return p, nil
}
