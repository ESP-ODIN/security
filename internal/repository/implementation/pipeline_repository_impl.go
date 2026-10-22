package repository

import (
	"context"

	"security/model"
)

type pipelineRepository struct{}

// NewPipelineRepository returns a PipelineRepository.
func NewPipelineRepository() PipelineRepository {
	return &pipelineRepository{}
}

func (r *pipelineRepository) ListByStatus(_ context.Context, _ string) ([]model.PipelineRun, error) {
	return []model.PipelineRun{}, nil
}

func (r *pipelineRepository) SetDecision(_ context.Context, runID, decision, _ string) (model.PipelineRun, error) {
	return model.PipelineRun{ID: runID, Decision: decision, Status: decision}, nil
}
