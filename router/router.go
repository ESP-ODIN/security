package router

import (
	"net/http"

	"security/handler"
	"security/middleware"
)

// Deps groups HTTP handlers wired into the router.
type Deps struct {
	Version  *handler.VersionHandler
	Security *handler.SecurityHandler
	Pipeline *handler.PipelineHandler
}

// New builds the HTTP mux with all API routes.
func New(deps Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)

	// Authentifié
	mux.Handle("POST /api/v1/agents/{id}/versions", middleware.RequireAuth(http.HandlerFunc(deps.Version.Submit)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/security", middleware.RequireAuth(http.HandlerFunc(deps.Security.GetStatus)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/security/report", middleware.RequireAuth(http.HandlerFunc(deps.Security.GetReport)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/permissions", middleware.RequireAuth(http.HandlerFunc(deps.Version.GetPermissions)))

	// Admin
	mux.Handle("GET /api/v1/admin/security/pipeline-runs", middleware.RequireAdmin(http.HandlerFunc(deps.Pipeline.List)))
	mux.Handle("PUT /api/v1/admin/security/pipeline-runs/{runId}/decision", middleware.RequireAdmin(http.HandlerFunc(deps.Pipeline.Decide)))

	return mux
}
