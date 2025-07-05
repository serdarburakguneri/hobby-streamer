package job

import (
	"context"
	"encoding/json"
)

type JobRunner interface {
	Run(ctx context.Context, payload json.RawMessage) error
}

type Registry struct {
	runners map[string]JobRunner
}

func NewRegistry() *Registry {
	return &Registry{runners: make(map[string]JobRunner)}
}

func (r *Registry) Register(t string, runner JobRunner) {
	r.runners[t] = runner
}

func (r *Registry) Get(t string) (JobRunner, bool) {
	runner, ok := r.runners[t]
	return runner, ok
}