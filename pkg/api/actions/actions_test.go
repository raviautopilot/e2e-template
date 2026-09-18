package actions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

type echoRequest struct {
	Name string `json:"name"`
}

type echoResponse struct {
	Greeting string `json:"greeting"`
	AuthUser string `json:"auth_user,omitempty"`
}

func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			resp := echoResponse{Greeting: "hello " + r.Header.Get("X-Custom")}
			if auth := r.Header.Get("Authorization"); auth != "" {
				resp.AuthUser = auth
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(resp)

		case "/echo-post":
			var req echoRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			resp := echoResponse{Greeting: "created " + req.Name}
			if auth := r.Header.Get("Authorization"); auth != "" {
				resp.AuthUser = auth
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(resp)

		case "/status/204":
			w.WriteHeader(http.StatusNoContent)

		case "/status/404":
			w.WriteHeader(http.StatusNotFound)

		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
}

func TestSendHttpRequest(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("GET with headers and auth", func(t *testing.T) {
		auth := &client.BearerTokenAuth{Token: "test-token-123"}
		headers := map[string]string{"X-Custom": "world"}
		var resp echoResponse

		SendHttpRequest(tc, c, "GET", "/ok", headers, nil, &resp, auth)

		if resp.Greeting != "hello world" {
			t.Errorf("expected 'hello world', got %q", resp.Greeting)
		}
		if resp.AuthUser != "Bearer test-token-123" {
			t.Errorf("expected 'Bearer test-token-123', got %q", resp.AuthUser)
		}
	})

	t.Run("POST with body and auth", func(t *testing.T) {
		auth := &client.BasicAuth{Username: "admin", Password: "pwd"}
		req := echoRequest{Name: "tester"}
		var resp echoResponse

		SendHttpRequest(tc, c, "POST", "/echo-post", nil, &req, &resp, auth)

		if resp.Greeting != "created tester" {
			t.Errorf("expected 'created tester', got %q", resp.Greeting)
		}
		if resp.AuthUser == "" {
			t.Errorf("expected Basic auth header to be present")
		}
	})

	t.Run("DELETE with 204 No Content", func(t *testing.T) {
		SendHttpRequest(tc, c, "DELETE", "/status/204", nil, nil, nil, nil)
	})
}

func TestSendHttpRequestAndExpectStatus(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Expect 404", func(t *testing.T) {
		SendHttpRequestAndExpectStatus(tc, c, "GET", "/status/404", nil, nil, nil, nil, 404)
	})

	t.Run("Expect 200", func(t *testing.T) {
		SendHttpRequestAndExpectStatus(tc, c, "GET", "/ok", nil, nil, nil, nil, 200)
	})
}
