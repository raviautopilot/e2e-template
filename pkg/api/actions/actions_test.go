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

		case "/echo-put":
			var req echoRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			resp := echoResponse{Greeting: "updated " + req.Name}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(resp)

		case "/echo-patch":
			var req echoRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			resp := echoResponse{Greeting: "patched " + req.Name}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(resp)

		case "/status/200":
			w.WriteHeader(http.StatusOK)

		case "/status/204":
			w.WriteHeader(http.StatusNoContent)

		case "/status/404":
			w.WriteHeader(http.StatusNotFound)

		case "/head-ok":
			w.WriteHeader(http.StatusOK)

		case "/options-ok":
			w.Header().Set("Allow", "GET, POST, OPTIONS")
			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
}

func TestGet(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Get with headers and auth", func(t *testing.T) {
		auth := &client.BearerTokenAuth{Token: "test-token-123"}
		headers := map[string]string{"X-Custom": "world"}
		var resp echoResponse

		Get(tc, c, "/ok", headers, nil, &resp, auth)

		if resp.Greeting != "hello world" {
			t.Errorf("expected 'hello world', got %q", resp.Greeting)
		}
		if resp.AuthUser != "Bearer test-token-123" {
			t.Errorf("expected 'Bearer test-token-123', got %q", resp.AuthUser)
		}
	})

	t.Run("GetAndExpectOK", func(t *testing.T) {
		var resp echoResponse
		GetAndExpectOK(tc, c, "/ok", nil, nil, &resp, nil)
		if resp.Greeting != "hello " {
			t.Errorf("expected 'hello ', got %q", resp.Greeting)
		}
	})

	t.Run("GetAndExpectStatus 404", func(t *testing.T) {
		GetAndExpectStatus(tc, c, "/status/404", nil, nil, nil, nil, 404)
	})

	t.Run("GetAndExpectStatus 200", func(t *testing.T) {
		GetAndExpectStatus(tc, c, "/status/200", nil, nil, nil, nil, 200)
	})
}

func TestPost(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Post with body and auth", func(t *testing.T) {
		auth := &client.BasicAuth{Username: "admin", Password: "pwd"}
		req := echoRequest{Name: "tester"}
		var resp echoResponse

		Post(tc, c, "/echo-post", nil, &req, &resp, auth)

		if resp.Greeting != "created tester" {
			t.Errorf("expected 'created tester', got %q", resp.Greeting)
		}
		if resp.AuthUser == "" {
			t.Errorf("expected Basic auth header to be present")
		}
	})

	t.Run("PostAndExpectOK", func(t *testing.T) {
		req := echoRequest{Name: "tester"}
		var resp echoResponse
		PostAndExpectOK(tc, c, "/echo-post", nil, &req, &resp, nil)
		if resp.Greeting != "created tester" {
			t.Errorf("expected 'created tester', got %q", resp.Greeting)
		}
	})
}

func TestPut(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Put with body", func(t *testing.T) {
		req := echoRequest{Name: "resource"}
		var resp echoResponse
		Put(tc, c, "/echo-put", nil, &req, &resp, nil)
		if resp.Greeting != "updated resource" {
			t.Errorf("expected 'updated resource', got %q", resp.Greeting)
		}
	})
}

func TestPatch(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Patch with body", func(t *testing.T) {
		req := echoRequest{Name: "field"}
		var resp echoResponse
		Patch(tc, c, "/echo-patch", nil, &req, &resp, nil)
		if resp.Greeting != "patched field" {
			t.Errorf("expected 'patched field', got %q", resp.Greeting)
		}
	})
}

func TestDelete(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Delete 204 No Content", func(t *testing.T) {
		Delete(tc, c, "/status/204", nil, nil, nil, nil)
	})

	t.Run("DeleteAndExpectStatus 404", func(t *testing.T) {
		DeleteAndExpectStatus(tc, c, "/status/404", nil, nil, nil, nil, 404)
	})
}

func TestHead(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Head OK", func(t *testing.T) {
		Head(tc, c, "/head-ok", nil, nil, nil, nil)
	})

	t.Run("HeadAndExpectStatus 404", func(t *testing.T) {
		HeadAndExpectStatus(tc, c, "/status/404", nil, nil, nil, nil, 404)
	})
}

func TestOptions(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Options OK", func(t *testing.T) {
		Options(tc, c, "/options-ok", nil, nil, nil, nil)
	})
}
