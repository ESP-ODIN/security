package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"security/internal/dto"
)

// Pinger is implemented by the database pool.
type Pinger interface {
	Ping(context.Context) error
}

// Health godoc
//
//	@Summary	Liveness check
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	dto.HealthResponse
//	@Router		/health [get]
func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, dto.HealthResponse{Status: "ok"})
}

// Ready godoc
//
//	@Summary		Readiness check
//	@Description	Pings the database on every request.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	dto.HealthResponse
//	@Failure		503	{object}	dto.HealthResponse
//	@Router			/ready [get]
func Ready(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		w.Header().Set("Cache-Control", "no-store")
		if err := db.Ping(ctx); err != nil {
			slog.Warn("readiness check failed", "error", err)
			writeJSON(w, http.StatusServiceUnavailable, dto.HealthResponse{Status: "unavailable"})
			return
		}
		writeJSON(w, http.StatusOK, dto.HealthResponse{Status: "ok"})
	}
}
