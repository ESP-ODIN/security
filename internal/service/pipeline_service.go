package service

import (
	"context"

	"github.com/google/uuid"

	"security/internal/dto"
)

// PipelineService handles admin pipeline-run operations.
type PipelineService interface {
	List(ctx context.Context, status string) (dto.PipelineRunsResponse, error)
	Decide(ctx context.Context, runID, adminID uuid.UUID, req dto.DecisionRequest) (dto.DecisionResponse, error)
}
