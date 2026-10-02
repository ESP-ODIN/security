package handler

import (
	"net/http"

	"security/dto"
)

// Health godoc
//
//	@Summary		Health check
//	@Description	Returns API health status
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	dto.HealthResponse
//	@Router			/health [get]
func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, dto.HealthResponse{Status: "ok"})
}
