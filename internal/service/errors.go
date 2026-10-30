package service

import "errors"

const (
	DecisionApproved = "approved"
	DecisionRejected = "rejected"
)

var (
	// ErrNotFound is returned when a version or a pipeline run does not exist.
	ErrNotFound = errors.New("not found")
	// ErrDecisionRequired is returned when a manual decision omits decision.
	ErrDecisionRequired = errors.New("decision is required")
	// ErrReasonRequired is returned when a manual decision omits reason.
	ErrReasonRequired = errors.New("reason is required")
	// ErrInvalidDecision is returned when decision is not approved or rejected.
	ErrInvalidDecision = errors.New("decision must be approved or rejected")
)
