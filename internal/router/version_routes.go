package router

import "net/http"

func registerVersionRoutes(mux *http.ServeMux, deps Deps) {
	auth := deps.Authenticator
	mux.Handle("POST /api/v1/agents/{id}/versions", auth.RequireAuth(http.HandlerFunc(deps.Version.Submit)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/security", auth.RequireAuth(http.HandlerFunc(deps.Security.GetStatus)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/security/report", auth.RequireAuth(http.HandlerFunc(deps.Security.GetReport)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/permissions", auth.RequireAuth(http.HandlerFunc(deps.Version.GetPermissions)))
}
