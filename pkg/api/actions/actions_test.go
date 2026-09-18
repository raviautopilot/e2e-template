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

func TestSendHttpRequestInternal(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())

	t.Run("sendHttpRequest success without tc", func(t *testing.T) {
		var resp echoResponse
		err := sendHttpRequest(c, "GET", "/ok", nil, nil, &resp, nil)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("sendHttpRequest error without tc", func(t *testing.T) {
		err := sendHttpRequest(c, "GET", "/status/404", nil, nil, nil, nil)
		if err == nil {
			t.Fatalf("expected error for 404, got nil")
		}
	})

	t.Run("sendHttpRequestAndExpectStatus match without tc", func(t *testing.T) {
		err := sendHttpRequestAndExpectStatus(c, "GET", "/status/404", nil, nil, nil, nil, 404)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("sendHttpRequestAndExpectStatus mismatch without tc", func(t *testing.T) {
		err := sendHttpRequestAndExpectStatus(c, "GET", "/ok", nil, nil, nil, nil, 404)
		if err == nil {
			t.Fatalf("expected error for status mismatch, got nil")
		}
	})
}

func TestGet(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Get returns error without failing tc", func(t *testing.T) {
		// Get on 404 returns an error directly to the caller, does not fail tc
		err := Get(tc, c, "/status/404", nil, nil, nil, nil)
		if err == nil {
			t.Errorf("expected error from Get on 404, got nil")
		}
		if tc.FailureReason != "" {
			t.Errorf("expected Get to not set tc.FailureReason, got %q", tc.FailureReason)
		}
	})

	t.Run("Get with headers and auth returns nil without validation", func(t *testing.T) {
		auth := &client.BearerTokenAuth{Token: "test-token-123"}
		headers := map[string]string{"X-Custom": "world"}
		var resp echoResponse

		err := Get(tc, c, "/ok", headers, nil, &resp, auth)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Greeting != "hello world" {
			t.Errorf("expected 'hello world', got %q", resp.Greeting)
		}
		if resp.AuthUser != "Bearer test-token-123" {
			t.Errorf("expected 'Bearer test-token-123', got %q", resp.AuthUser)
		}
	})

	t.Run("GetAndExpectOK validates and sets tc.Actual", func(t *testing.T) {
		var resp echoResponse
		GetAndExpectOK(tc, c, "/ok", nil, nil, &resp, nil)
		if resp.Greeting != "hello " {
			t.Errorf("expected 'hello ', got %q", resp.Greeting)
		}
		if tc.Actual == "" {
			t.Errorf("expected tc.Actual to be set by GetAndExpectOK")
		}
	})

	t.Run("GetAndExpectStatus 404", func(t *testing.T) {
		GetAndExpectStatus(tc, c, "/status/404", nil, nil, nil, nil, 404)
		if tc.Actual != "HTTP 404 as expected for GET /status/404" {
			t.Errorf("unexpected tc.Actual: %q", tc.Actual)
		}
	})

	t.Run("GetAndExpectStatusCode 200", func(t *testing.T) {
		GetAndExpectStatusCode(tc, c, "/status/200", nil, nil, nil, nil, 200)
		if tc.Actual != "HTTP 200 as expected for GET /status/200" {
			t.Errorf("unexpected tc.Actual: %q", tc.Actual)
		}
	})
}

func TestPost(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Post returns error without failing tc", func(t *testing.T) {
		auth := &client.BasicAuth{Username: "admin", Password: "pwd"}
		req := echoRequest{Name: "tester"}
		var resp echoResponse

		err := Post(tc, c, "/echo-post", nil, &req, &resp, auth)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Greeting != "created tester" {
			t.Errorf("expected 'created tester', got %q", resp.Greeting)
		}
		if resp.AuthUser == "" {
			t.Errorf("expected Basic auth header to be present")
		}
	})

	t.Run("Post on 404 returns error without setting tc.FailureReason", func(t *testing.T) {
		testTc := &tests.TestContext{T: t, Client: c}
		err := Post(testTc, c, "/status/404", nil, nil, nil, nil)
		if err == nil {
			t.Errorf("expected error on 404, got nil")
		}
		if testTc.FailureReason != "" {
			t.Errorf("expected Post to not set FailureReason, got %q", testTc.FailureReason)
		}
	})

	t.Run("PostAndExpectOK validates and sets tc.Actual", func(t *testing.T) {
		req := echoRequest{Name: "tester"}
		var resp echoResponse
		PostAndExpectOK(tc, c, "/echo-post", nil, &req, &resp, nil)
		if resp.Greeting != "created tester" {
			t.Errorf("expected 'created tester', got %q", resp.Greeting)
		}
		if tc.Actual == "" {
			t.Errorf("expected tc.Actual to be set by PostAndExpectOK")
		}
	})

	t.Run("PostAndExpectStatusCode", func(t *testing.T) {
		PostAndExpectStatusCode(tc, c, "/status/404", nil, nil, nil, nil, 404)
	})
}

func TestPut(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Put with body returns error without failing tc", func(t *testing.T) {
		req := echoRequest{Name: "resource"}
		var resp echoResponse
		err := Put(tc, c, "/echo-put", nil, &req, &resp, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Greeting != "updated resource" {
			t.Errorf("expected 'updated resource', got %q", resp.Greeting)
		}
	})

	t.Run("PutAndExpectOK validates", func(t *testing.T) {
		req := echoRequest{Name: "resource"}
		var resp echoResponse
		PutAndExpectOK(tc, c, "/echo-put", nil, &req, &resp, nil)
		if tc.Actual == "" {
			t.Errorf("expected tc.Actual to be set by PutAndExpectOK")
		}
	})
}

func TestPatch(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Patch with body returns error without failing tc", func(t *testing.T) {
		req := echoRequest{Name: "field"}
		var resp echoResponse
		err := Patch(tc, c, "/echo-patch", nil, &req, &resp, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Greeting != "patched field" {
			t.Errorf("expected 'patched field', got %q", resp.Greeting)
		}
	})

	t.Run("PatchAndExpectOK validates", func(t *testing.T) {
		req := echoRequest{Name: "field"}
		var resp echoResponse
		PatchAndExpectOK(tc, c, "/echo-patch", nil, &req, &resp, nil)
		if tc.Actual == "" {
			t.Errorf("expected tc.Actual to be set by PatchAndExpectOK")
		}
	})
}

func TestDelete(t *testing.T) {
	srv := setupMockServer()
	defer srv.Close()

	c := client.NewClient(srv.URL, 5*time.Second, t.TempDir())
	tc := &tests.TestContext{T: t, Client: c}

	t.Run("Delete 204 No Content returns error without failing tc", func(t *testing.T) {
		err := Delete(tc, c, "/status/204", nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("DeleteAndExpectOK validates", func(t *testing.T) {
		DeleteAndExpectOK(tc, c, "/status/204", nil, nil, nil, nil)
		if tc.Actual == "" {
			t.Errorf("expected tc.Actual to be set by DeleteAndExpectOK")
		}
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

	t.Run("Head OK returns error without failing tc", func(t *testing.T) {
		err := Head(tc, c, "/head-ok", nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("HeadAndExpectOK validates", func(t *testing.T) {
		HeadAndExpectOK(tc, c, "/head-ok", nil, nil, nil, nil)
		if tc.Actual == "" {
			t.Errorf("expected tc.Actual to be set by HeadAndExpectOK")
		}
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

	t.Run("Options OK returns error without failing tc", func(t *testing.T) {
		err := Options(tc, c, "/options-ok", nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("OptionsAndExpectOK validates", func(t *testing.T) {
		OptionsAndExpectOK(tc, c, "/options-ok", nil, nil, nil, nil)
		if tc.Actual == "" {
			t.Errorf("expected tc.Actual to be set by OptionsAndExpectOK")
		}
	})
}
