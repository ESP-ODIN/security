package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "security/docs"
)

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.Handle("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
}
