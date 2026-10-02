package dto

// PipelineRunsResponse lists pipeline runs filtered by status.
type PipelineRunsResponse struct {
	Status string `json:"status"`
	Runs   []any  `json:"runs"`
}

// DecisionRequest is the body for a manual pipeline decision.
type DecisionRequest struct {
	Decision string `json:"decision"` // validé | rejeté
	Reason   string `json:"reason"`
}

// DecisionResponse confirms a manual decision was recorded.
type DecisionResponse struct {
	RunID    string `json:"runId"`
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// ErrorResponse is a generic API error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}
