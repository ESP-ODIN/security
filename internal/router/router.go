package router

import (
	"net/http"

	"github.com/ESP-ODIN/authkit-go"

	"security/internal/handler"
)

// Deps groups HTTP handlers wired into the router.
type Deps struct {
	Version       *handler.VersionHandler
	Security      *handler.SecurityHandler
	Pipeline      *handler.PipelineHandler
	Authenticator *authkit.Authenticator
}

// New builds the HTTP mux with all API routes.
func New(db handler.Pinger, deps Deps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /ready", handler.Ready(db))
	registerVersionRoutes(mux, deps)
	registerAdminRoutes(mux, deps)
	registerSwaggerRoutes(mux)
	return mux
}
