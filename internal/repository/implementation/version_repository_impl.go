package repository

import (
	"context"

	"security/model"
)

// versionRepository is a stub implementation until the database layer is ready.
type versionRepository struct{}

// NewVersionRepository returns a VersionRepository.
func NewVersionRepository() VersionRepository {
	return &versionRepository{}
}

func (r *versionRepository) Create(_ context.Context, version model.Version, _ string) (model.Version, error) {
	return version, nil
}

func (r *versionRepository) GetByID(_ context.Context, agentID, versionID string) (model.Version, error) {
	return model.Version{ID: versionID, AgentID: agentID}, nil
}

func (r *versionRepository) GetPermissions(_ context.Context, _, _ string) ([]model.Permission, error) {
	return []model.Permission{}, nil
}
