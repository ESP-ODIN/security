package model

// PipelineRun tracks a security pipeline execution for a version.
type PipelineRun struct {
	ID        string
	AgentID   string
	VersionID string
	Status    string // e.g. running, manual_review, validated, rejected
	Step      string
	Decision  string
}
