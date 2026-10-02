package model

// SecurityReport holds analysis results for a version.
type SecurityReport struct {
	AgentID         string
	VersionID       string
	Secrets         []Finding
	Vulnerabilities []Finding
	SandboxAccess   []Finding
}

// Finding is a single analysis result item.
type Finding struct {
	Type    string
	Message string
}

// Permission is a declared permission from a version manifest.
type Permission struct {
	Name  string
	Scope string
}
