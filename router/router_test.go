package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"security/handler"
	"security/repository"
	"security/router"
	"security/service"
)

func testRouter() http.Handler {
	return router.New(router.Deps{
		Version:  handler.NewVersionHandler(service.NewVersionService(repository.NewVersionRepository())),
		Security: handler.NewSecurityHandler(service.NewSecurityService(repository.NewSecurityRepository())),
		Pipeline: handler.NewPipelineHandler(service.NewPipelineService(repository.NewPipelineRepository())),
	})
}

func TestHealth(t *testing.T) {
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{http.MethodGet, "/health", http.StatusOK},
		{http.MethodPost, "/health", http.StatusMethodNotAllowed},
		{http.MethodGet, "/unknown", http.StatusNotFound},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			testRouter().ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d", rec.Code, tc.status)
			}
			if tc.status == http.StatusOK {
				if got := rec.Header().Get("Content-Type"); got != "application/json" {
					t.Fatalf("Content-Type = %q", got)
				}
				var body struct {
					Status string `json:"status"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Status != "ok" {
					t.Fatalf("invalid health response: %s (error: %v)", rec.Body, err)
				}
			}
		})
	}
}

func TestSecurityRoutesRequireAuth(t *testing.T) {
	paths := []struct {
		method, path string
	}{
		{http.MethodPost, "/api/v1/agents/a1/versions"},
		{http.MethodGet, "/api/v1/agents/a1/versions/v1/security"},
		{http.MethodGet, "/api/v1/agents/a1/versions/v1/security/report"},
		{http.MethodGet, "/api/v1/agents/a1/versions/v1/permissions"},
	}
	for _, tc := range paths {
		t.Run("unauth "+tc.method+" "+tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			testRouter().ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestSecurityRoutesAuthenticated(t *testing.T) {
	cases := []struct {
		method, path string
		status       int
		headers      map[string]string
	}{
		{http.MethodPost, "/api/v1/agents/a1/versions", http.StatusAccepted, map[string]string{
			"Idempotency-Key": "key-1",
		}},
		{http.MethodGet, "/api/v1/agents/a1/versions/v1/security", http.StatusOK, nil},
		{http.MethodGet, "/api/v1/agents/a1/versions/v1/security/report", http.StatusOK, nil},
		{http.MethodGet, "/api/v1/agents/a1/versions/v1/permissions", http.StatusOK, nil},
	}
	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Authorization", "Bearer test")
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()
			testRouter().ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.status, rec.Body)
			}
		})
	}
}

func TestAdminRoutes(t *testing.T) {
	t.Run("list without auth", func(t *testing.T) {
		rec := httptest.NewRecorder()
		testRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/admin/security/pipeline-runs?status=manual_review", nil))
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("list without admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/security/pipeline-runs?status=manual_review", nil)
		req.Header.Set("Authorization", "Bearer test")
		rec := httptest.NewRecorder()
		testRouter().ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("list as admin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/security/pipeline-runs?status=manual_review", nil)
		req.Header.Set("Authorization", "Bearer test")
		req.Header.Set("X-Admin", "true")
		rec := httptest.NewRecorder()
		testRouter().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body)
		}
	})

	t.Run("decision requires reason", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/security/pipeline-runs/run-1/decision",
			strings.NewReader(`{"decision":"validé"}`))
		req.Header.Set("Authorization", "Bearer test")
		req.Header.Set("X-Admin", "true")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		testRouter().ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body)
		}
	})

	t.Run("decision ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/security/pipeline-runs/run-1/decision",
			strings.NewReader(`{"decision":"validé","reason":"revue manuelle OK"}`))
		req.Header.Set("Authorization", "Bearer test")
		req.Header.Set("X-Admin", "true")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		testRouter().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body)
		}
	})
}
