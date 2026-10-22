package dto

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"

	"security/internal/model"
)

func TestJSONUsesSnakeCase(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	step := "analyze"
	body, err := json.Marshal(SecurityStatusResponse{AgentID: id, VersionID: id, Step: &step})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "agentId") || !strings.Contains(string(body), `"agent_id"`) {
		t.Fatalf("status = %s", body)
	}

	admin := id
	runBody, err := json.Marshal(model.PipelineRun{ID: id, AgentID: id, DecidedBy: &admin, Reason: "ok"})
	if err != nil {
		t.Fatal(err)
	}
	raw := string(runBody)
	for _, key := range []string{`"id"`, `"agent_id"`, `"decided_by"`, `"reason"`} {
		if !strings.Contains(raw, key) {
			t.Fatalf("run JSON missing %s: %s", key, raw)
		}
	}
	if strings.Contains(raw, `"ID"`) || strings.Contains(raw, `"AgentID"`) {
		t.Fatalf("run JSON uses Go names: %s", raw)
	}

	report, err := json.Marshal(SecurityReportResponse{
		AgentID: id,
		Secrets: []model.Finding{{Type: "secret", Message: "key"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(report), `"sandbox_access"`) || !strings.Contains(string(report), `"type"`) {
		t.Fatalf("report = %s", report)
	}
}
