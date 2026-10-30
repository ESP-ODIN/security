package dto

import (
	"github.com/google/uuid"

	"security/internal/model"
)

// SecurityStatusResponse is the pipeline status for a version.
type SecurityStatusResponse struct {
	AgentID   uuid.UUID `json:"agent_id" format:"uuid" swaggertype:"string"`
	VersionID uuid.UUID `json:"version_id" format:"uuid" swaggertype:"string"`
	Step      *string   `json:"step"`
	Decision  *string   `json:"decision"`
}

// SecurityReportResponse is the detailed analysis report.
type SecurityReportResponse struct {
	AgentID         uuid.UUID       `json:"agent_id" format:"uuid" swaggertype:"string"`
	VersionID       uuid.UUID       `json:"version_id" format:"uuid" swaggertype:"string"`
	Secrets         []model.Finding `json:"secrets"`
	Vulnerabilities []model.Finding `json:"vulnerabilities"`
	SandboxAccess   []model.Finding `json:"sandbox_access"`
}

// PermissionsResponse lists permissions declared in a version manifest.
type PermissionsResponse struct {
	AgentID     uuid.UUID          `json:"agent_id" format:"uuid" swaggertype:"string"`
	VersionID   uuid.UUID          `json:"version_id" format:"uuid" swaggertype:"string"`
	Permissions []model.Permission `json:"permissions"`
}
