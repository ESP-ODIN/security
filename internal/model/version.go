package model

import "github.com/google/uuid"

// Version is an agent version submitted for the security pipeline.
type Version struct {
	ID      uuid.UUID `json:"id" format:"uuid" swaggertype:"string"`
	AgentID uuid.UUID `json:"agent_id" format:"uuid" swaggertype:"string"`
	Status  string    `json:"status"`
}
