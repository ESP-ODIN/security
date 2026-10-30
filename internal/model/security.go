package model

import "github.com/google/uuid"

// SecurityReport holds analysis results for a version.
type SecurityReport struct {
	AgentID         uuid.UUID `json:"agent_id" format:"uuid" swaggertype:"string"`
	VersionID       uuid.UUID `json:"version_id" format:"uuid" swaggertype:"string"`
	Secrets         []Finding `json:"secrets"`
	Vulnerabilities []Finding `json:"vulnerabilities"`
	SandboxAccess   []Finding `json:"sandbox_access"`
}

// Finding is a single analysis result item.
type Finding struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Permission is a declared permission from a version manifest.
type Permission struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
}
