package main

import (
	"encoding/json"
	"net/http"
)

// POST /api/v1/agents/{id}/versions
// Soumission d'une nouvelle version — déclenche le pipeline de sécurité.
// Idempotent via header Idempotency-Key.
func submitVersion(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("id")
	idempotencyKey := r.Header.Get("Idempotency-Key")

	writeJSON(w, http.StatusAccepted, map[string]any{
		"agentId":         agentID,
		"idempotencyKey":  idempotencyKey,
		"status":          "accepted",
		"message":         "version submitted; security pipeline will start",
	})
}

// GET /api/v1/agents/{id}/versions/{versionId}/security
// Statut du pipeline pour cette version (étape en cours, décision).
func getSecurityStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"agentId":   r.PathValue("id"),
		"versionId": r.PathValue("versionId"),
		"step":      nil,
		"decision":  nil,
	})
}

// GET /api/v1/agents/{id}/versions/{versionId}/security/report
// Détail des résultats (secrets, vulnérabilités, accès sandbox observés).
func getSecurityReport(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"agentId":         r.PathValue("id"),
		"versionId":       r.PathValue("versionId"),
		"secrets":         []any{},
		"vulnerabilities": []any{},
		"sandboxAccess":   []any{},
	})
}

// GET /api/v1/agents/{id}/versions/{versionId}/permissions
// Permissions déclarées dans le manifeste de cette version.
func getPermissions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"agentId":     r.PathValue("id"),
		"versionId":   r.PathValue("versionId"),
		"permissions": []any{},
	})
}

// GET /api/v1/admin/security/pipeline-runs?status=manual_review
// Liste des runs en attente de revue manuelle.
func listPipelineRuns(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "manual_review"
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status": status,
		"runs":   []any{},
	})
}

type decisionRequest struct {
	Decision string `json:"decision"` // validé | rejeté
	Reason   string `json:"reason"`
}

// PUT /api/v1/admin/security/pipeline-runs/{runId}/decision
// Décision manuelle (analyse complémentaire). Body: decision + reason.
func decidePipelineRun(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("runId")

	var body decisionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.Decision == "" {
		writeError(w, http.StatusBadRequest, "decision is required")
		return
	}
	if body.Reason == "" {
		writeError(w, http.StatusBadRequest, "reason is required")
		return
	}
	if body.Decision != "validé" && body.Decision != "rejeté" {
		writeError(w, http.StatusBadRequest, "decision must be validé or rejeté")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"runId":    runID,
		"decision": body.Decision,
		"reason":   body.Reason,
	})
}
