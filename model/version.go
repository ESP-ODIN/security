package model

// Version is an agent version submitted for the security pipeline.
type Version struct {
	ID      string
	AgentID string
	Status  string
}
