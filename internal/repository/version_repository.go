package repository

import (
	"context"

	"github.com/google/uuid"

	"security/internal/model"
)

// VersionRepository persists agent versions.
type VersionRepository interface {
	Create(ctx context.Context, version model.Version, idempotencyKey string) (model.Version, error)
	// GetByID returns ErrNotFound when the version does not exist.
	GetByID(ctx context.Context, agentID, versionID uuid.UUID) (model.Version, error)
	// GetPermissions returns ErrNotFound when the version does not exist.
	GetPermissions(ctx context.Context, agentID, versionID uuid.UUID) ([]model.Permission, error)
}
