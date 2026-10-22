package router

import "net/http"

func registerAdminRoutes(mux *http.ServeMux, deps Deps) {
	auth := deps.Authenticator
	mux.Handle("GET /api/v1/admin/security/pipeline-runs", auth.RequireAuth(http.HandlerFunc(deps.Pipeline.List)))
	mux.Handle("PUT /api/v1/admin/security/pipeline-runs/{runId}/decision", auth.RequireAuth(http.HandlerFunc(deps.Pipeline.Decide)))
}
