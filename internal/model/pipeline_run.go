package model

import "github.com/google/uuid"

// PipelineRun tracks a security pipeline execution for a version.
type PipelineRun struct {
	ID        uuid.UUID  `json:"id" format:"uuid" swaggertype:"string"`
	AgentID   uuid.UUID  `json:"agent_id" format:"uuid" swaggertype:"string"`
	VersionID uuid.UUID  `json:"version_id" format:"uuid" swaggertype:"string"`
	Status    string     `json:"status"`
	Step      string     `json:"step"`
	Decision  string     `json:"decision"`
	Reason    string     `json:"reason"`
	DecidedBy *uuid.UUID `json:"decided_by,omitempty" format:"uuid" swaggertype:"string"`
}
