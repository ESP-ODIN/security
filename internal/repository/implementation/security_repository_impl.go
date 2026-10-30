package implementation

import (
	"context"

	"github.com/google/uuid"

	"security/internal/model"
	"security/internal/repository"
)

// securityRepository is a stand-in until the database layer is ready.
// Lookups return ErrNotFound because nothing is persisted yet.
type securityRepository struct{}

// NewSecurityRepository returns a SecurityRepository.
func NewSecurityRepository() repository.SecurityRepository {
	return &securityRepository{}
}

func (r *securityRepository) GetStatus(context.Context, uuid.UUID, uuid.UUID) (model.PipelineRun, error) {
	return model.PipelineRun{}, repository.ErrNotFound
}

func (r *securityRepository) GetReport(context.Context, uuid.UUID, uuid.UUID) (model.SecurityReport, error) {
	return model.SecurityReport{}, repository.ErrNotFound
}
