package repository

import (
	"context"

	"security/model"
)

// VersionRepository persists agent versions.
type VersionRepository interface {
	Create(ctx context.Context, version model.Version, idempotencyKey string) (model.Version, error)
	GetByID(ctx context.Context, agentID, versionID string) (model.Version, error)
	GetPermissions(ctx context.Context, agentID, versionID string) ([]model.Permission, error)
}
