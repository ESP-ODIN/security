package dto

import (
	"github.com/google/uuid"

	"security/internal/model"
)

// PipelineRunsResponse lists pipeline runs filtered by status.
type PipelineRunsResponse struct {
	Status string              `json:"status"`
	Runs   []model.PipelineRun `json:"runs"`
}

// DecisionRequest is the body for a manual pipeline decision.
type DecisionRequest struct {
	Decision string `json:"decision"` // approved | rejected
	Reason   string `json:"reason"`
}

// DecisionResponse confirms a manual decision was recorded.
type DecisionResponse struct {
	RunID     uuid.UUID `json:"run_id" format:"uuid" swaggertype:"string"`
	Decision  string    `json:"decision"`
	Reason    string    `json:"reason"`
	DecidedBy uuid.UUID `json:"decided_by" format:"uuid" swaggertype:"string"`
}

// ErrorResponse is a generic API error payload.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid agent id"`
}
