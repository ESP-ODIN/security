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

// GetStatus handles GET /api/v1/agents/{id}/versions/{versionId}/security.
func (h *SecurityHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.GetStatus(r.Context(), r.PathValue("id"), r.PathValue("versionId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetReport handles GET /api/v1/agents/{id}/versions/{versionId}/security/report.
func (h *SecurityHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.GetReport(r.Context(), r.PathValue("id"), r.PathValue("versionId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
