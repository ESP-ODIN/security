package main

import "net/http"

// requireAuth ensures the request is authenticated.
// Placeholder: real JWT / session checks will replace this.
func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requireAdmin ensures the request has admin privileges.
// Placeholder: real role checks will replace this.
func requireAdmin(next http.Handler) http.Handler {
	return requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Admin") != "true" {
			writeError(w, http.StatusForbidden, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	}))
}
