package main

import "net/http"

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Authentifié
	mux.Handle("POST /api/v1/agents/{id}/versions", requireAuth(http.HandlerFunc(submitVersion)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/security", requireAuth(http.HandlerFunc(getSecurityStatus)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/security/report", requireAuth(http.HandlerFunc(getSecurityReport)))
	mux.Handle("GET /api/v1/agents/{id}/versions/{versionId}/permissions", requireAuth(http.HandlerFunc(getPermissions)))

	// Admin
	mux.Handle("GET /api/v1/admin/security/pipeline-runs", requireAdmin(http.HandlerFunc(listPipelineRuns)))
	mux.Handle("PUT /api/v1/admin/security/pipeline-runs/{runId}/decision", requireAdmin(http.HandlerFunc(decidePipelineRun)))

	return mux
}
