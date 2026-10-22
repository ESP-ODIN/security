package repository

import (
	"context"

	"security/model"
)

type securityRepository struct{}

// NewSecurityRepository returns a SecurityRepository.
func NewSecurityRepository() SecurityRepository {
	return &securityRepository{}
}

func (r *securityRepository) GetStatus(_ context.Context, agentID, versionID string) (model.PipelineRun, error) {
	return model.PipelineRun{AgentID: agentID, VersionID: versionID}, nil
}

func (r *securityRepository) GetReport(_ context.Context, agentID, versionID string) (model.SecurityReport, error) {
	return model.SecurityReport{
		AgentID:         agentID,
		VersionID:       versionID,
		Secrets:         []model.Finding{},
		Vulnerabilities: []model.Finding{},
		SandboxAccess:   []model.Finding{},
	}, nil
}
