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

// Submit handles POST /api/v1/agents/{id}/versions.
func (h *VersionHandler) Submit(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.Submit(r.Context(), r.PathValue("id"), r.Header.Get("Idempotency-Key"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, resp)
}

// GetPermissions handles GET /api/v1/agents/{id}/versions/{versionId}/permissions.
func (h *VersionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.GetPermissions(r.Context(), r.PathValue("id"), r.PathValue("versionId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
