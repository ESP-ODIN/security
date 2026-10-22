package implementation

import (
	"context"

	"github.com/google/uuid"

	"security/internal/dto"
	"security/internal/model"
	"security/internal/repository"
	"security/internal/service"
)

type securityService struct {
	repo repository.SecurityRepository
}

// NewSecurityService returns a SecurityService.
func NewSecurityService(repo repository.SecurityRepository) service.SecurityService {
	return &securityService{repo: repo}
}

func (s *securityService) GetStatus(ctx context.Context, agentID, versionID uuid.UUID) (dto.SecurityStatusResponse, error) {
	run, err := s.repo.GetStatus(ctx, agentID, versionID)
	if err != nil {
		return dto.SecurityStatusResponse{}, translate(err)
	}
	resp := dto.SecurityStatusResponse{AgentID: agentID, VersionID: versionID}
	if run.Step != "" {
		resp.Step = &run.Step
	}
	if run.Decision != "" {
		resp.Decision = &run.Decision
	}
	return resp, nil
}

func (s *securityService) GetReport(ctx context.Context, agentID, versionID uuid.UUID) (dto.SecurityReportResponse, error) {
	report, err := s.repo.GetReport(ctx, agentID, versionID)
	if err != nil {
		return dto.SecurityReportResponse{}, translate(err)
	}
	return dto.SecurityReportResponse{
		AgentID:         agentID,
		VersionID:       versionID,
		Secrets:         findingsOrEmpty(report.Secrets),
		Vulnerabilities: findingsOrEmpty(report.Vulnerabilities),
		SandboxAccess:   findingsOrEmpty(report.SandboxAccess),
	}, nil
}

func findingsOrEmpty(items []model.Finding) []model.Finding {
	if items == nil {
		return []model.Finding{}
	}
	return items
}
