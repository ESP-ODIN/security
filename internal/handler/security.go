package handler

import (
	"net/http"

	"security/internal/dto"
	"security/internal/service"
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
//	@Param			id			path		string	true	"Agent ID"		format(uuid)
//	@Param			versionId	path		string	true	"Version ID"	format(uuid)
//	@Success		200			{object}	dto.SecurityStatusResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{string}	string	"unauthorized"
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/agents/{id}/versions/{versionId}/security [get]
func (h *SecurityHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if _, ok := authenticatedIdentity(w, r); !ok {
		return
	}
	agentID, ok := parsePathID(w, r.PathValue("id"), "agent id")
	if !ok {
		return
	}
	versionID, ok := parsePathID(w, r.PathValue("versionId"), "version id")
	if !ok {
		return
	}
	resp, err := h.svc.GetStatus(r.Context(), agentID, versionID)
	if err != nil {
		writeServiceError(w, err)
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
//	@Param			id			path		string	true	"Agent ID"		format(uuid)
//	@Param			versionId	path		string	true	"Version ID"	format(uuid)
//	@Success		200			{object}	dto.SecurityReportResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{string}	string	"unauthorized"
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/agents/{id}/versions/{versionId}/security/report [get]
func (h *SecurityHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	if _, ok := authenticatedIdentity(w, r); !ok {
		return
	}
	agentID, ok := parsePathID(w, r.PathValue("id"), "agent id")
	if !ok {
		return
	}
	versionID, ok := parsePathID(w, r.PathValue("versionId"), "version id")
	if !ok {
		return
	}
	var resp dto.SecurityReportResponse
	var err error
	resp, err = h.svc.GetReport(r.Context(), agentID, versionID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
