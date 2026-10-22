package implementation

import (
	"context"

	"github.com/google/uuid"

	"security/internal/model"
	"security/internal/repository"
)

// versionRepository is a stand-in until the database layer is ready.
// Reads return ErrNotFound because nothing is persisted yet.
type versionRepository struct{}

// NewVersionRepository returns a VersionRepository.
func NewVersionRepository() repository.VersionRepository {
	return &versionRepository{}
}

func (r *versionRepository) Create(_ context.Context, version model.Version, _ string) (model.Version, error) {
	if version.ID == uuid.Nil {
		version.ID = uuid.New()
	}
	return version, nil
}

func (r *versionRepository) GetByID(context.Context, uuid.UUID, uuid.UUID) (model.Version, error) {
	return model.Version{}, repository.ErrNotFound
}

func (r *versionRepository) GetPermissions(context.Context, uuid.UUID, uuid.UUID) ([]model.Permission, error) {
	return nil, repository.ErrNotFound
}
