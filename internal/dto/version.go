package dto

import "github.com/google/uuid"

// SubmitVersionResponse is returned after submitting a new version.
type SubmitVersionResponse struct {
	AgentID        uuid.UUID `json:"agent_id" format:"uuid" swaggertype:"string"`
	IdempotencyKey string    `json:"idempotency_key"`
	Status         string    `json:"status"`
	Message        string    `json:"message"`
}
