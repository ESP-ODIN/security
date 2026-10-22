package repository

import (
	"context"

	"security/model"
)

// PipelineRepository manages security pipeline runs.
type PipelineRepository interface {
	ListByStatus(ctx context.Context, status string) ([]model.PipelineRun, error)
	SetDecision(ctx context.Context, runID, decision, reason string) (model.PipelineRun, error)
}
