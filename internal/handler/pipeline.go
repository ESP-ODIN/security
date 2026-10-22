package handler

import (
	"net/http"

	"github.com/google/uuid"

	"security/internal/dto"
	"security/internal/service"
)

// PipelineHandler handles admin pipeline-run endpoints.
type PipelineHandler struct {
	svc    service.PipelineService
	admins map[uuid.UUID]struct{}
}

// NewPipelineHandler returns a PipelineHandler.
// adminUserIDs lists Auth user ids allowed to call admin routes.
func NewPipelineHandler(svc service.PipelineService, adminUserIDs []uuid.UUID) *PipelineHandler {
	admins := make(map[uuid.UUID]struct{}, len(adminUserIDs))
	for _, id := range adminUserIDs {
		admins[id] = struct{}{}
	}
	return &PipelineHandler{svc: svc, admins: admins}
}

// List godoc
//
//	@Summary		List pipeline runs
//	@Description	Lists security pipeline runs, typically those awaiting manual review. The caller must be listed in ADMIN_USER_IDS.
//	@Tags			admin
//	@Produce		json
//	@Param			status	query		string	false	"Filter by status"	default(manual_review)
//	@Success		200		{object}	dto.PipelineRunsResponse
//	@Failure		401		{string}	string	"unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/admin/security/pipeline-runs [get]
func (h *PipelineHandler) List(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.adminID(w, r); !ok {
		return
	}
	resp, err := h.svc.List(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// Decide godoc
//
//	@Summary		Decide on a pipeline run
//	@Description	Records a manual decision for complementary analysis. Body must include decision (approved or rejected) and reason. The reason and the admin id are stored for audit.
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			runId	path		string				true	"Pipeline run ID"	format(uuid)
//	@Param			body	body		dto.DecisionRequest	true	"Manual decision"
//	@Success		200		{object}	dto.DecisionResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{string}	string	"unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		413		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/admin/security/pipeline-runs/{runId}/decision [put]
func (h *PipelineHandler) Decide(w http.ResponseWriter, r *http.Request) {
	adminID, ok := h.adminID(w, r)
	if !ok {
		return
	}
	runID, ok := parsePathID(w, r.PathValue("runId"), "run id")
	if !ok {
		return
	}
	var body dto.DecisionRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	resp, err := h.svc.Decide(r.Context(), runID, adminID, body)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *PipelineHandler) adminID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	identity, ok := authenticatedIdentity(w, r)
	if !ok {
		return uuid.Nil, false
	}
	if _, allowed := h.admins[identity.ID]; !allowed {
		writeError(w, http.StatusForbidden, "admin access required")
		return uuid.Nil, false
	}
	return identity.ID, true
}
