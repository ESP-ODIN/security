package repository

import (
	"context"

	"github.com/google/uuid"

	"security/internal/model"
)

// SecurityRepository loads security analysis data.
type SecurityRepository interface {
	// GetStatus returns the pipeline run for a version.
	// It returns ErrNotFound when the version or its run does not exist.
	GetStatus(ctx context.Context, agentID, versionID uuid.UUID) (model.PipelineRun, error)
	// GetReport returns the analysis report for a version.
	// It returns ErrNotFound when the version does not exist.
	GetReport(ctx context.Context, agentID, versionID uuid.UUID) (model.SecurityReport, error)
}
