package service

import (
	"context"
	"fmt"

	"security/dto"
	"security/repository"
)

type pipelineService struct {
	repo repository.PipelineRepository
}

// NewPipelineService returns a PipelineService.
func NewPipelineService(repo repository.PipelineRepository) PipelineService {
	return &pipelineService{repo: repo}
}

func (s *pipelineService) List(ctx context.Context, status string) (dto.PipelineRunsResponse, error) {
	if status == "" {
		status = "manual_review"
	}
	runs, err := s.repo.ListByStatus(ctx, status)
	if err != nil {
		return dto.PipelineRunsResponse{}, err
	}
	items := make([]any, 0, len(runs))
	for _, run := range runs {
		items = append(items, run)
	}
	return dto.PipelineRunsResponse{Status: status, Runs: items}, nil
}

func (s *pipelineService) Decide(ctx context.Context, runID string, req dto.DecisionRequest) (dto.DecisionResponse, error) {
	if req.Decision == "" {
		return dto.DecisionResponse{}, fmt.Errorf("decision is required")
	}
	if req.Reason == "" {
		return dto.DecisionResponse{}, fmt.Errorf("reason is required")
	}
	if req.Decision != "validé" && req.Decision != "rejeté" {
		return dto.DecisionResponse{}, fmt.Errorf("decision must be validé or rejeté")
	}
	if _, err := s.repo.SetDecision(ctx, runID, req.Decision, req.Reason); err != nil {
		return dto.DecisionResponse{}, err
	}
	return dto.DecisionResponse{
		RunID:    runID,
		Decision: req.Decision,
		Reason:   req.Reason,
	}, nil
}
