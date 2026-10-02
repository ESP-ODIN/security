package repository

import (
	"context"

	"security/model"
)

// SecurityRepository loads security analysis data.
type SecurityRepository interface {
	GetStatus(ctx context.Context, agentID, versionID string) (model.PipelineRun, error)
	GetReport(ctx context.Context, agentID, versionID string) (model.SecurityReport, error)
}
