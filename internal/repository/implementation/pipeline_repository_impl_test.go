package implementation

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"security/internal/model"
	"security/internal/repository"
)

func TestSetDecisionStoresAudit(t *testing.T) {
	repo := &pipelineRepository{runs: map[uuid.UUID]model.PipelineRun{}}
	runID := uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	adminID := uuid.MustParse("7e733b0f-136e-4cd9-923d-19c733cd466a")
	repo.runs[runID] = model.PipelineRun{ID: runID, Status: "manual_review"}

	got, err := repo.SetDecision(context.Background(), runID, "approved", "revue manuelle OK", adminID)
	if err != nil {
		t.Fatalf("SetDecision: %v", err)
	}
	if got.Decision != "approved" || got.Reason != "revue manuelle OK" || got.DecidedBy == nil || *got.DecidedBy != adminID || got.Status != "approved" {
		t.Fatalf("stored run = %+v", got)
	}

	listed, err := repo.ListByStatus(context.Background(), "approved")
	if err != nil || len(listed) != 1 || listed[0].Reason != "revue manuelle OK" || listed[0].DecidedBy == nil || *listed[0].DecidedBy != adminID {
		t.Fatalf("listed = %+v, err = %v", listed, err)
	}
	stillReview, err := repo.ListByStatus(context.Background(), "manual_review")
	if err != nil || len(stillReview) != 0 {
		t.Fatalf("manual_review = %+v, err = %v", stillReview, err)
	}
}

func TestSetDecisionUnknownRun(t *testing.T) {
	repo := NewPipelineRepository()
	_, err := repo.SetDecision(context.Background(), uuid.New(), "rejected", "not safe", uuid.New())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("error = %v, want ErrNotFound", err)
	}
}
