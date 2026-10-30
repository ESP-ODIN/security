package repository

import (
	"context"

	"github.com/google/uuid"

	"security/internal/model"
)

// PipelineRepository manages security pipeline runs.
type PipelineRepository interface {
	ListByStatus(ctx context.Context, status string) ([]model.PipelineRun, error)
	// SetDecision stores a manual decision, its reason, and the admin who made it.
	// It returns ErrNotFound when the run does not exist.
	SetDecision(ctx context.Context, runID uuid.UUID, decision, reason string, decidedBy uuid.UUID) (model.PipelineRun, error)
}
