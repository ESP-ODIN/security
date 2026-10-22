package service

import (
	"context"

	"security/dto"
)

// SecurityService exposes pipeline status and analysis reports.
type SecurityService interface {
	GetStatus(ctx context.Context, agentID, versionID string) (dto.SecurityStatusResponse, error)
	GetReport(ctx context.Context, agentID, versionID string) (dto.SecurityReportResponse, error)
}
