package service

import (
	"context"

	"security/dto"
)

// PipelineService handles admin pipeline-run operations.
type PipelineService interface {
	List(ctx context.Context, status string) (dto.PipelineRunsResponse, error)
	Decide(ctx context.Context, runID string, req dto.DecisionRequest) (dto.DecisionResponse, error)
}
