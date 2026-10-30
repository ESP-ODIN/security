package implementation

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"security/internal/dto"
	"security/internal/model"
	"security/internal/repository"
	"security/internal/service"
)

type pipelineService struct {
	repo repository.PipelineRepository
}

// NewPipelineService returns a PipelineService.
func NewPipelineService(repo repository.PipelineRepository) service.PipelineService {
	return &pipelineService{repo: repo}
}

func (s *pipelineService) List(ctx context.Context, status string) (dto.PipelineRunsResponse, error) {
	if status == "" {
		status = "manual_review"
	}
	runs, err := s.repo.ListByStatus(ctx, status)
	if err != nil {
		return dto.PipelineRunsResponse{}, translate(err)
	}
	if runs == nil {
		runs = []model.PipelineRun{}
	}
	return dto.PipelineRunsResponse{Status: status, Runs: runs}, nil
}

func (s *pipelineService) Decide(ctx context.Context, runID, adminID uuid.UUID, req dto.DecisionRequest) (dto.DecisionResponse, error) {
	req.Decision = strings.TrimSpace(req.Decision)
	req.Reason = strings.TrimSpace(req.Reason)
	if req.Decision == "" {
		return dto.DecisionResponse{}, service.ErrDecisionRequired
	}
	if req.Reason == "" {
		return dto.DecisionResponse{}, service.ErrReasonRequired
	}
	if req.Decision != service.DecisionApproved && req.Decision != service.DecisionRejected {
		return dto.DecisionResponse{}, service.ErrInvalidDecision
	}
	run, err := s.repo.SetDecision(ctx, runID, req.Decision, req.Reason, adminID)
	if err != nil {
		return dto.DecisionResponse{}, translate(err)
	}
	decidedBy := adminID
	if run.DecidedBy != nil {
		decidedBy = *run.DecidedBy
	}
	return dto.DecisionResponse{
		RunID:     run.ID,
		Decision:  run.Decision,
		Reason:    run.Reason,
		DecidedBy: decidedBy,
	}, nil
}
