package implementation

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"security/internal/dto"
	"security/internal/model"
	"security/internal/repository"
	"security/internal/service"
)

type fakePipelineRepo struct {
	missing bool
	called  bool
	admin   uuid.UUID
	reason  string
}

func (f *fakePipelineRepo) ListByStatus(context.Context, string) ([]model.PipelineRun, error) {
	return []model.PipelineRun{}, nil
}

func (f *fakePipelineRepo) SetDecision(_ context.Context, runID uuid.UUID, decision, reason string, decidedBy uuid.UUID) (model.PipelineRun, error) {
	f.called = true
	f.admin = decidedBy
	f.reason = reason
	if f.missing {
		return model.PipelineRun{}, repository.ErrNotFound
	}
	admin := decidedBy
	return model.PipelineRun{
		ID:        runID,
		Decision:  decision,
		Reason:    reason,
		Status:    decision,
		DecidedBy: &admin,
	}, nil
}

func TestDecide(t *testing.T) {
	runID := uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	adminID := uuid.MustParse("7e733b0f-136e-4cd9-923d-19c733cd466a")

	t.Run("requires decision and reason", func(t *testing.T) {
		repo := &fakePipelineRepo{}
		svc := NewPipelineService(repo)
		_, err := svc.Decide(context.Background(), runID, adminID, dto.DecisionRequest{Reason: "ok"})
		if !errors.Is(err, service.ErrDecisionRequired) || repo.called {
			t.Fatalf("error = %v, called = %v", err, repo.called)
		}
		_, err = svc.Decide(context.Background(), runID, adminID, dto.DecisionRequest{Decision: "approved", Reason: "  "})
		if !errors.Is(err, service.ErrReasonRequired) || repo.called {
			t.Fatalf("error = %v, called = %v", err, repo.called)
		}
	})

	t.Run("rejects legacy decision values", func(t *testing.T) {
		repo := &fakePipelineRepo{}
		svc := NewPipelineService(repo)
		_, err := svc.Decide(context.Background(), runID, adminID, dto.DecisionRequest{Decision: "validé", Reason: "ok"})
		if !errors.Is(err, service.ErrInvalidDecision) || repo.called {
			t.Fatalf("error = %v, called = %v", err, repo.called)
		}
	})

	t.Run("missing run", func(t *testing.T) {
		repo := &fakePipelineRepo{missing: true}
		svc := NewPipelineService(repo)
		_, err := svc.Decide(context.Background(), runID, adminID, dto.DecisionRequest{Decision: "rejected", Reason: "secrets"})
		if !errors.Is(err, service.ErrNotFound) || !repo.called || repo.admin != adminID || repo.reason != "secrets" {
			t.Fatalf("error = %v, called = %v, admin = %s", err, repo.called, repo.admin)
		}
	})

	t.Run("stores reason and admin", func(t *testing.T) {
		repo := &fakePipelineRepo{}
		svc := NewPipelineService(repo)
		resp, err := svc.Decide(context.Background(), runID, adminID, dto.DecisionRequest{Decision: " approved ", Reason: " revue OK "})
		if err != nil {
			t.Fatal(err)
		}
		if resp.RunID != runID || resp.Decision != "approved" || resp.Reason != "revue OK" || resp.DecidedBy != adminID {
			t.Fatalf("response = %+v", resp)
		}
		if repo.admin != adminID || repo.reason != "revue OK" {
			t.Fatalf("stored admin = %s, reason = %q", repo.admin, repo.reason)
		}
	})
}
