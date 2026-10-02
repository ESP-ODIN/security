package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"security/dto"
	"security/service"
)

// PipelineHandler handles admin pipeline-run endpoints.
type PipelineHandler struct {
	svc service.PipelineService
}

// NewPipelineHandler returns a PipelineHandler.
func NewPipelineHandler(svc service.PipelineService) *PipelineHandler {
	return &PipelineHandler{svc: svc}
}

// List handles GET /api/v1/admin/security/pipeline-runs.
func (h *PipelineHandler) List(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.List(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// Decide handles PUT /api/v1/admin/security/pipeline-runs/{runId}/decision.
func (h *PipelineHandler) Decide(w http.ResponseWriter, r *http.Request) {
	var body dto.DecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.svc.Decide(r.Context(), r.PathValue("runId"), body)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "required") || strings.Contains(msg, "must be") {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeError(w, http.StatusInternalServerError, msg)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
