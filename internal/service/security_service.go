package service

import (
	"context"

	"github.com/google/uuid"

	"security/internal/dto"
)

// SecurityService exposes pipeline status and analysis reports.
type SecurityService interface {
	GetStatus(ctx context.Context, agentID, versionID uuid.UUID) (dto.SecurityStatusResponse, error)
	GetReport(ctx context.Context, agentID, versionID uuid.UUID) (dto.SecurityReportResponse, error)
}
