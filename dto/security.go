package dto

// SecurityStatusResponse is the pipeline status for a version.
type SecurityStatusResponse struct {
	AgentID   string  `json:"agentId"`
	VersionID string  `json:"versionId"`
	Step      *string `json:"step"`
	Decision  *string `json:"decision"`
}

// SecurityReportResponse is the detailed analysis report.
type SecurityReportResponse struct {
	AgentID         string `json:"agentId"`
	VersionID       string `json:"versionId"`
	Secrets         []any  `json:"secrets"`
	Vulnerabilities []any  `json:"vulnerabilities"`
	SandboxAccess   []any  `json:"sandboxAccess"`
}

// PermissionsResponse lists permissions declared in a version manifest.
type PermissionsResponse struct {
	AgentID     string `json:"agentId"`
	VersionID   string `json:"versionId"`
	Permissions []any  `json:"permissions"`
}
