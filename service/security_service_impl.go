package service

import (
	"context"

	"security/dto"
	"security/repository"
)

type securityService struct {
	repo repository.SecurityRepository
}

// NewSecurityService returns a SecurityService.
func NewSecurityService(repo repository.SecurityRepository) SecurityService {
	return &securityService{repo: repo}
}

func (s *securityService) GetStatus(ctx context.Context, agentID, versionID string) (dto.SecurityStatusResponse, error) {
	run, err := s.repo.GetStatus(ctx, agentID, versionID)
	if err != nil {
		return dto.SecurityStatusResponse{}, err
	}
	resp := dto.SecurityStatusResponse{
		AgentID:   run.AgentID,
		VersionID: run.VersionID,
	}
	if run.Step != "" {
		resp.Step = &run.Step
	}
	if run.Decision != "" {
		resp.Decision = &run.Decision
	}
	return resp, nil
}

func (s *securityService) GetReport(ctx context.Context, agentID, versionID string) (dto.SecurityReportResponse, error) {
	report, err := s.repo.GetReport(ctx, agentID, versionID)
	if err != nil {
		return dto.SecurityReportResponse{}, err
	}
	return dto.SecurityReportResponse{
		AgentID:         report.AgentID,
		VersionID:       report.VersionID,
		Secrets:         toAnySlice(report.Secrets),
		Vulnerabilities: toAnySlice(report.Vulnerabilities),
		SandboxAccess:   toAnySlice(report.SandboxAccess),
	}, nil
}

func toAnySlice[T any](items []T) []any {
	out := make([]any, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	return out
}
