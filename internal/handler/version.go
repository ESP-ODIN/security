package handler

import (
	"net/http"

	"security/service"
)

// VersionHandler handles version-related endpoints.
type VersionHandler struct {
	svc service.VersionService
}

// NewVersionHandler returns a VersionHandler.
func NewVersionHandler(svc service.VersionService) *VersionHandler {
	return &VersionHandler{svc: svc}
}

// Submit godoc
//
//	@Summary		Submit a new agent version
//	@Description	Submits a new version and starts the security pipeline (normalize, quarantine, analyze, smoke test, sandbox). Idempotent via Idempotency-Key: the same key does not start a second run.
//	@Tags			versions
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string	true	"Agent ID"
//	@Param			Idempotency-Key	header		string	false	"Idempotency key to avoid duplicate pipeline runs"
//	@Success		202				{object}	dto.SubmitVersionResponse
//	@Failure		401				{object}	dto.ErrorResponse
//	@Failure		500				{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/agents/{id}/versions [post]
func (h *VersionHandler) Submit(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.Submit(r.Context(), r.PathValue("id"), r.Header.Get("Idempotency-Key"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, resp)
}

// GetPermissions godoc
//
//	@Summary		List version permissions
//	@Description	Returns permissions declared in this version's manifest (useful to compare with the previous version)
//	@Tags			versions
//	@Produce		json
//	@Param			id			path		string	true	"Agent ID"
//	@Param			versionId	path		string	true	"Version ID"
//	@Success		200			{object}	dto.PermissionsResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/agents/{id}/versions/{versionId}/permissions [get]
func (h *VersionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.GetPermissions(r.Context(), r.PathValue("id"), r.PathValue("versionId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
