package implementation

import (
	"context"
	"sort"
	"sync"

	"github.com/google/uuid"

	"security/internal/model"
	"security/internal/repository"
)

// pipelineRepository is an in-memory stand-in until PostgreSQL is wired.
// It starts empty, so a decision on an unknown run returns ErrNotFound.
type pipelineRepository struct {
	mu   sync.Mutex
	runs map[uuid.UUID]model.PipelineRun
}

// NewPipelineRepository returns a PipelineRepository.
func NewPipelineRepository() repository.PipelineRepository {
	return &pipelineRepository{runs: map[uuid.UUID]model.PipelineRun{}}
}

func (r *pipelineRepository) ListByStatus(_ context.Context, status string) ([]model.PipelineRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]model.PipelineRun, 0)
	for _, run := range r.runs {
		if run.Status == status {
			out = append(out, run)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID.String() < out[j].ID.String() })
	return out, nil
}

func (r *pipelineRepository) SetDecision(_ context.Context, runID uuid.UUID, decision, reason string, decidedBy uuid.UUID) (model.PipelineRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	run, ok := r.runs[runID]
	if !ok {
		return model.PipelineRun{}, repository.ErrNotFound
	}
	admin := decidedBy
	run.Decision = decision
	run.Reason = reason
	run.Status = decision
	run.DecidedBy = &admin
	r.runs[runID] = run
	return run, nil
}
