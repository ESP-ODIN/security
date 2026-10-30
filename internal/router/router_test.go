package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"security/internal/handler"
	repositoryimpl "security/internal/repository/implementation"
	serviceimpl "security/internal/service/implementation"
	"security/internal/testutil"
)

var (
	agentID   = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	versionID = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	runID     = uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	adminID   = uuid.MustParse("7e733b0f-136e-4cd9-923d-19c733cd466a")
	userID    = uuid.MustParse("11111111-1111-1111-1111-111111111111")
)

type fakePinger struct{ err error }

func (p fakePinger) Ping(context.Context) error { return p.err }

func testRouter(t *testing.T, admins []uuid.UUID, pingErr error) (*testutil.Auth, http.Handler) {
	t.Helper()
	authn := testutil.NewAuth(t)
	h := New(fakePinger{err: pingErr}, Deps{
		Version:       handler.NewVersionHandler(serviceimpl.NewVersionService(repositoryimpl.NewVersionRepository())),
		Security:      handler.NewSecurityHandler(serviceimpl.NewSecurityService(repositoryimpl.NewSecurityRepository())),
		Pipeline:      handler.NewPipelineHandler(serviceimpl.NewPipelineService(repositoryimpl.NewPipelineRepository()), admins),
		Authenticator: authn.Authenticator,
	})
	return authn, h
}

func authed(t *testing.T, authn *testutil.Auth, id uuid.UUID, method, path, body string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+authn.Token(t, authn.Claims(id)))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid error body: %s", rec.Body)
	}
	return body.Error
}

func TestHealthAndReady(t *testing.T) {
	for _, tc := range []struct {
		name, method, path string
		pingErr            error
		status             int
		bodyStatus         string
	}{
		{"health", http.MethodGet, "/health", nil, http.StatusOK, "ok"},
		{"health method", http.MethodPost, "/health", nil, http.StatusMethodNotAllowed, ""},
		{"ready", http.MethodGet, "/ready", nil, http.StatusOK, "ok"},
		{"database unavailable", http.MethodGet, "/ready", errors.New("database unavailable"), http.StatusServiceUnavailable, "unavailable"},
		{"ready method", http.MethodPost, "/ready", nil, http.StatusMethodNotAllowed, ""},
		{"unknown", http.MethodGet, "/unknown", nil, http.StatusNotFound, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, h := testRouter(t, nil, tc.pingErr)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d", rec.Code, tc.status)
			}
			if tc.bodyStatus == "" {
				return
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q", got)
			}
			if tc.path == "/ready" {
				if got := rec.Header().Get("Cache-Control"); got != "no-store" {
					t.Fatalf("Cache-Control = %q", got)
				}
			}
			var body struct {
				Status string `json:"status"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Status != tc.bodyStatus {
				t.Fatalf("invalid response: %s (error: %v)", rec.Body, err)
			}
		})
	}
}

func TestSecurityRoutesRequireAuth(t *testing.T) {
	_, h := testRouter(t, nil, nil)
	paths := []struct {
		method, path string
	}{
		{http.MethodPost, "/api/v1/agents/" + agentID.String() + "/versions"},
		{http.MethodGet, "/api/v1/agents/" + agentID.String() + "/versions/" + versionID.String() + "/security"},
		{http.MethodGet, "/api/v1/agents/" + agentID.String() + "/versions/" + versionID.String() + "/security/report"},
		{http.MethodGet, "/api/v1/agents/" + agentID.String() + "/versions/" + versionID.String() + "/permissions"},
		{http.MethodGet, "/api/v1/admin/security/pipeline-runs"},
	}
	for _, tc := range paths {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Authorization", "Bearer not-a-jwt")
			req.Header.Set("X-Admin", "true")
			h.ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestInvalidIDs(t *testing.T) {
	authn, h := testRouter(t, []uuid.UUID{adminID}, nil)
	for _, tc := range []struct {
		method, path, want string
	}{
		{http.MethodPost, "/api/v1/agents/not-a-uuid/versions", "invalid agent id"},
		{http.MethodGet, "/api/v1/agents/" + agentID.String() + "/versions/v1/security", "invalid version id"},
		{http.MethodPut, "/api/v1/admin/security/pipeline-runs/run-1/decision", "invalid run id"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			req := authed(t, authn, adminID, tc.method, tc.path, "")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest || decodeError(t, rec) != tc.want {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestAuthenticatedRoutes(t *testing.T) {
	authn, h := testRouter(t, []uuid.UUID{adminID}, nil)

	t.Run("submit", func(t *testing.T) {
		req := authed(t, authn, userID, http.MethodPost, "/api/v1/agents/"+agentID.String()+"/versions", "")
		req.Header.Set("Idempotency-Key", "key-1")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusAccepted {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
		}
		var body struct {
			AgentID        string `json:"agent_id"`
			IdempotencyKey string `json:"idempotency_key"`
			Status         string `json:"status"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.AgentID != agentID.String() || body.IdempotencyKey != "key-1" || body.Status != "accepted" {
			t.Fatalf("body = %+v", body)
		}
		if strings.Contains(rec.Body.String(), "agentId") {
			t.Fatalf("response uses camelCase: %s", rec.Body)
		}
	})

	for _, path := range []string{
		"/api/v1/agents/" + agentID.String() + "/versions/" + versionID.String() + "/security",
		"/api/v1/agents/" + agentID.String() + "/versions/" + versionID.String() + "/security/report",
		"/api/v1/agents/" + agentID.String() + "/versions/" + versionID.String() + "/permissions",
	} {
		t.Run("missing "+path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, authed(t, authn, userID, http.MethodGet, path, ""))
			if rec.Code != http.StatusNotFound || decodeError(t, rec) != "not found" {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestAdminRoutes(t *testing.T) {
	authn, h := testRouter(t, []uuid.UUID{adminID}, nil)

	t.Run("list without auth", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/admin/security/pipeline-runs?status=manual_review", nil))
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("x-admin header is ignored", func(t *testing.T) {
		req := authed(t, authn, userID, http.MethodGet, "/api/v1/admin/security/pipeline-runs", "")
		req.Header.Set("X-Admin", "true")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden || decodeError(t, rec) != "admin access required" {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("list as configured admin", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, authed(t, authn, adminID, http.MethodGet, "/api/v1/admin/security/pipeline-runs?status=manual_review", ""))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
		}
		var body struct {
			Status string            `json:"status"`
			Runs   []json.RawMessage `json:"runs"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Status != "manual_review" || body.Runs == nil {
			t.Fatalf("body = %s, err = %v", rec.Body, err)
		}
	})

	t.Run("decision requires reason", func(t *testing.T) {
		req := authed(t, authn, adminID, http.MethodPut, "/api/v1/admin/security/pipeline-runs/"+runID.String()+"/decision", `{"decision":"approved"}`)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest || decodeError(t, rec) != "reason is required" {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("legacy decision is rejected", func(t *testing.T) {
		req := authed(t, authn, adminID, http.MethodPut, "/api/v1/admin/security/pipeline-runs/"+runID.String()+"/decision", `{"decision":"validé","reason":"ok"}`)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest || decodeError(t, rec) != "decision must be approved or rejected" {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		req := authed(t, authn, adminID, http.MethodPut, "/api/v1/admin/security/pipeline-runs/"+runID.String()+"/decision", `{"decision":"approved","reason":"ok","extra":true}`)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("body too large", func(t *testing.T) {
		payload := `{"decision":"approved","reason":"` + strings.Repeat("a", 1<<20) + `"}`
		req := authed(t, authn, adminID, http.MethodPut, "/api/v1/admin/security/pipeline-runs/"+runID.String()+"/decision", payload)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusRequestEntityTooLarge || decodeError(t, rec) != "request body too large" {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("missing run", func(t *testing.T) {
		req := authed(t, authn, adminID, http.MethodPut, "/api/v1/admin/security/pipeline-runs/"+runID.String()+"/decision", `{"decision":"approved","reason":"revue OK"}`)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound || decodeError(t, rec) != "not found" {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
		}
	})
}

func TestSwagger(t *testing.T) {
	_, h := testRouter(t, nil, nil)
	for _, path := range []string{"/swagger/index.html", "/swagger/doc.json"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200", path, rec.Code)
		}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil))
	var spec struct {
		Paths map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatalf("invalid spec: %v", err)
	}
	for _, path := range []string{"/health", "/ready", "/api/v1/agents/{id}/versions", "/api/v1/admin/security/pipeline-runs"} {
		if _, ok := spec.Paths[path]; !ok {
			t.Fatalf("spec is missing %s; regenerate it with `go tool swag init -g cmd/api/main.go -o docs --parseInternal`", path)
		}
	}
	if strings.Contains(rec.Body.String(), "X-Admin") || strings.Contains(rec.Body.String(), "agentId") {
		t.Fatalf("spec still documents X-Admin or camelCase fields")
	}
}
