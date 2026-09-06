package api_test

// ─────────────────────────────────────────────────────────────────────────────
// TEMPLATE: Example API Tests (YOUR APPLICATION)
//
// This file shows how to write tests against YOUR local/staging API server.
// Tests are skipped automatically when no server is running at baseUrl.
//
// See public_api_test.go for fully working examples against real public APIs.
//
// HOW TO ADAPT:
//  1. Replace baseUrl in config.json with your API's base URL.
//  2. Rename these functions and update endpoint paths.
//  3. Remove the isServerRunning() skip guard once your server is configured.
//  4. Add more test files: tests/api/XX-feature_test.go
//
// Run with:
//   go test -v ./tests/api/... -run TestAPI_01
// ─────────────────────────────────────────────────────────────────────────────

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

// isServerRunning does a quick connectivity check against the configured baseUrl.
func isServerRunning() bool {
	if tests.GlobalConfig == nil || tests.GlobalConfig.BaseURL == "" {
		return false
	}
	httpClient := &http.Client{Timeout: 2 * time.Second}
	resp, err := httpClient.Get(tests.GlobalConfig.BaseURL)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

// ─────────────────────────────────────────────────────────────────────────────
// Example 1: Health Check Endpoint
// ─────────────────────────────────────────────────────────────────────────────

func TestAPI_01_HealthCheck(t *testing.T) {
	if !isServerRunning() {
		t.Skipf("Skipping: no server running at %s (set baseUrl in config.json)", tests.GlobalConfig.BaseURL)
	}

	type testCase struct {
		Name        string
		Description string
		Expected    string
		RunFn       func(tc *tests.TestContext)
	}

	testCases := []testCase{
		{
			Name:        "Health Check Returns 200 OK",
			Description: "Verifies the health endpoint returns 200 OK.",
			Expected:    "HTTP 200 OK with status field",
			RunFn: func(tc *tests.TestContext) {
				var resp HealthResponse
				actions.GetAndExpectOK(tc, tc.Client, "/health", &resp)
				tc.Actual = fmt.Sprintf("status=%q, message=%q", resp.Status, resp.Message)
			},
		},
		{
			Name:        "Wrong HTTP Method POST on /health",
			Description: "Verifies that POST on a GET-only endpoint is rejected.",
			Expected:    "HTTP 405 Method Not Allowed or 404 Not Found",
			RunFn: func(tc *tests.TestContext) {
				var resp HealthResponse
				err := tc.Client.SendHttpRequest("POST", "/health", nil, nil, &resp, nil)
				if err == nil {
					tc.FailureReason = "Expected HTTP error, got 200 OK"
					tc.Errorf("Expected HTTP error, got 200 OK")
				} else if err.StatusCode() != 404 && err.StatusCode() != 405 {
					tc.FailureReason = fmt.Sprintf("Expected 404 or 405, got %d", err.StatusCode())
					tc.Errorf("Expected 404 or 405, got %d", err.StatusCode())
				} else {
					tc.Actual = fmt.Sprintf("Rejected with HTTP %d as expected", err.StatusCode())
				}
			},
		},
		{
			Name:        "SQL Injection Query Parameter Handled Safely",
			Description: "Verifies SQL injection in query param does not cause a 500 error.",
			Expected:    "HTTP 200 OK or 400 Bad Request — not a 500 crash",
			RunFn: func(tc *tests.TestContext) {
				actions.GetAndExpectNotCrash(tc, tc.Client, "/health?id=1' OR '1'='1")
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		tests.RunAPITestWithDetails(t, testCase.Name, testCase.Description, testCase.Expected, testCase.RunFn)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Example 2: Table-Driven GET Endpoint Tests
// ─────────────────────────────────────────────────────────────────────────────

func TestAPI_02_PublicEndpoints(t *testing.T) {
	if !isServerRunning() {
		t.Skipf("Skipping: no server running at %s (set baseUrl in config.json)", tests.GlobalConfig.BaseURL)
	}

	type testCase struct {
		Name        string
		Description string
		Expected    string
		RunFn       func(tc *tests.TestContext)
	}

	testCases := []testCase{
		{
			Name:        "GET Root Returns Welcome Payload",
			Description: "Verifies the root endpoint returns a non-empty message.",
			Expected:    "HTTP 200 OK with a non-empty message",
			RunFn: func(tc *tests.TestContext) {
				var resp RootResponse
				actions.GetAndExpectOK(tc, tc.Client, "/", &resp)
				actions.AssertNotEmpty(tc, "message", resp.Message)
			},
		},
		{
			Name:        "Invalid Route Returns 404",
			Description: "Verifies that unknown routes return 404 Not Found.",
			Expected:    "HTTP 404 Not Found",
			RunFn: func(tc *tests.TestContext) {
				actions.GetAndExpectStatus(tc, tc.Client, "/this-route-does-not-exist", 404)
			},
		},
		{
			Name:        "Negative Pagination Parameters Do Not Crash Server",
			Description: "Verifies that negative pagination params are handled gracefully.",
			Expected:    "HTTP 200 OK or 400 Bad Request — not a 500 crash",
			RunFn: func(tc *tests.TestContext) {
				actions.GetAndExpectNotCrash(tc, tc.Client, "/health?page=-1&limit=-100")
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		tests.RunAPITestWithDetails(t, testCase.Name, testCase.Description, testCase.Expected, testCase.RunFn)
	}
}
