package handler

import (
	"net/http"

	"security/service"
)

// SecurityHandler handles security status and report endpoints.
type SecurityHandler struct {
	svc service.SecurityService
}

// NewSecurityHandler returns a SecurityHandler.
func NewSecurityHandler(svc service.SecurityService) *SecurityHandler {
	return &SecurityHandler{svc: svc}
}

// GetStatus godoc
//
//	@Summary		Get security pipeline status
//	@Description	Returns the pipeline status for this version (current step, decision)
//	@Tags			security
//	@Produce		json
//	@Param			id			path		string	true	"Agent ID"
//	@Param			versionId	path		string	true	"Version ID"
//	@Success		200			{object}	dto.SecurityStatusResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/agents/{id}/versions/{versionId}/security [get]
func (h *SecurityHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.GetStatus(r.Context(), r.PathValue("id"), r.PathValue("versionId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetReport godoc
//
//	@Summary		Get security analysis report
//	@Description	Returns detailed analysis results (secrets, vulnerabilities, observed sandbox access)
//	@Tags			security
//	@Produce		json
//	@Param			id			path		string	true	"Agent ID"
//	@Param			versionId	path		string	true	"Version ID"
//	@Success		200			{object}	dto.SecurityReportResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/agents/{id}/versions/{versionId}/security/report [get]
func (h *SecurityHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.GetReport(r.Context(), r.PathValue("id"), r.PathValue("versionId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
