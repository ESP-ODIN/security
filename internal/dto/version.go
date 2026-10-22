package dto

// SubmitVersionResponse is returned after submitting a new version.
type SubmitVersionResponse struct {
	AgentID        string `json:"agentId"`
	IdempotencyKey string `json:"idempotencyKey"`
	Status         string `json:"status"`
	Message        string `json:"message"`
}
