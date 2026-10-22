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

// List godoc
//
//	@Summary		List pipeline runs
//	@Description	Lists security pipeline runs, typically those awaiting manual review
//	@Tags			admin
//	@Produce		json
//	@Param			status	query		string	false	"Filter by status"	default(manual_review)
//	@Success		200		{object}	dto.PipelineRunsResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Security		AdminAuth
//	@Router			/api/v1/admin/security/pipeline-runs [get]
func (h *PipelineHandler) List(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.List(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// Decide godoc
//
//	@Summary		Decide on a pipeline run
//	@Description	Records a manual decision for complementary analysis. Body must include decision (validé/rejeté) and reason (free text, audited).
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			runId	path		string				true	"Pipeline run ID"
//	@Param			body	body		dto.DecisionRequest	true	"Manual decision"
//	@Success		200		{object}	dto.DecisionResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Security		AdminAuth
//	@Router			/api/v1/admin/security/pipeline-runs/{runId}/decision [put]
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
