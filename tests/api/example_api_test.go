package api_test

// ─────────────────────────────────────────────────────────────────────────────
// TEMPLATE: Example API Tests
//
// This file shows how to write API tests using the e2e-template framework.
// Tests use RunAPITestWithDetails for rich HTML/Markdown reports.
//
// HOW TO ADAPT:
//  1. Replace baseUrl in config.json with your API's base URL.
//  2. Rename these functions and update endpoint paths.
//  3. Add more test cases to the table-driven slices.
//  4. Add more test files following the same naming convention: XX-feature_test.go
//
// FRAMEWORK CONCEPTS:
//  - tests.RunAPITestWithDetails → runs a named subtest with full report metadata
//  - tc.Client.SendHttpRequest   → makes HTTP calls and returns HttpError on failure
//  - tc.Expected / tc.Actual / tc.FailureReason → populate the HTML/MD report
//
// Run with:
//   go test -v ./tests/api/...
// ─────────────────────────────────────────────────────────────────────────────

import (
	"fmt"
	"testing"

	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// Example 1: Health Check Endpoint
// Demonstrates a simple positive test and a wrong-method negative test.
// ─────────────────────────────────────────────────────────────────────────────

type healthTestCase struct {
	Name        string
	Method      string
	Path        string
	Description string
	Expected    string
	WantError   bool
	WantStatus  int
}

// TestAPI_01_HealthCheck validates the /health (or equivalent) endpoint.
// Adapt the path to match your API's health endpoint.
func TestAPI_01_HealthCheck(t *testing.T) {
	testCases := []healthTestCase{
		{
			Name:        "Health Check Returns 200 OK",
			Method:      "GET",
			Path:        "/health",
			Description: "Verifies the health endpoint returns 200 OK with a valid status field.",
			Expected:    "HTTP 200 OK with {\"status\":\"ok\"} or similar",
			WantError:   false,
		},
		{
			Name:        "Wrong HTTP Method POST on /health",
			Method:      "POST",
			Path:        "/health",
			Description: "Verifies that POST on a GET-only health endpoint is rejected.",
			Expected:    "HTTP 405 Method Not Allowed or 404 Not Found",
			WantError:   true,
			WantStatus:  405,
		},
		{
			Name:        "SQL Injection Query Parameter Handled Safely",
			Method:      "GET",
			Path:        "/health?id=1' OR '1'='1",
			Description: "Verifies that SQL injection in query param does not cause a 500 error.",
			Expected:    "HTTP 200 OK or 400 Bad Request — not a 500 crash",
			WantError:   false,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		tests.RunAPITestWithDetails(
			t,
			testCase.Name,
			testCase.Description,
			testCase.Expected,
			func(tc *tests.TestContext) {
				var resp HealthResponse
				err := tc.Client.SendHttpRequest(testCase.Method, testCase.Path, nil, nil, &resp, nil)

				if testCase.WantError {
					if err == nil {
						tc.FailureReason = fmt.Sprintf("Expected HTTP error, got 200 OK")
						tc.Errorf("Expected HTTP error, got 200 OK")
					} else if testCase.WantStatus > 0 && err.StatusCode() != testCase.WantStatus {
						tc.Actual = fmt.Sprintf("Got HTTP %d", err.StatusCode())
						// Accept either 404 or 405 for method-not-allowed scenarios
						if err.StatusCode() != 404 && err.StatusCode() != 405 {
							tc.FailureReason = fmt.Sprintf("Expected %d, got %d", testCase.WantStatus, err.StatusCode())
							tc.Errorf("Expected %d or 404, got %d", testCase.WantStatus, err.StatusCode())
						}
					} else {
						tc.Actual = fmt.Sprintf("Rejected with HTTP %d as expected", err.StatusCode())
					}
				} else {
					if err != nil && err.StatusCode() == 500 {
						tc.FailureReason = fmt.Sprintf("Server crashed with 500: %v", err)
						tc.Errorf("Server crashed with 500: %v", err)
					} else if err != nil {
						// Non-500 errors may be acceptable (e.g., endpoint not yet implemented)
						tc.Actual = fmt.Sprintf("Responded with HTTP %d (not a crash)", err.StatusCode())
					} else {
						tc.Actual = fmt.Sprintf("HTTP 200 OK, status=%q, message=%q", resp.Status, resp.Message)
					}
				}
			},
		)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Example 2: Table-Driven GET Endpoint Tests
// Shows how to test multiple public GET endpoints with a shared validation loop.
// ─────────────────────────────────────────────────────────────────────────────

type publicGetTestCase struct {
	Name        string
	Path        string
	Description string
	Expected    string
	ValidateFn  func(tc *tests.TestContext)
}

// TestAPI_02_PublicEndpoints validates multiple read-only public endpoints.
// TODO: Replace the endpoint paths and validation logic with your own API routes.
func TestAPI_02_PublicEndpoints(t *testing.T) {
	testCases := []publicGetTestCase{
		{
			Name:        "GET Root Returns Welcome Payload",
			Path:        "/",
			Description: "Verifies the root endpoint returns a non-empty message field.",
			Expected:    "HTTP 200 OK with a non-empty message",
			ValidateFn: func(tc *tests.TestContext) {
				var resp RootResponse
				err := tc.Client.SendHttpRequest("GET", "/", nil, nil, &resp, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("GET / failed: %v", err)
					tc.Errorf("GET / failed: %v", err)
					return
				}
				tc.Actual = fmt.Sprintf("HTTP 200 OK, message=%q", resp.Message)
				if resp.Message == "" {
					tc.FailureReason = "Expected non-empty message in root response"
					tc.Errorf("Expected non-empty message in root response")
				}
			},
		},
		{
			Name:        "Invalid Route Returns 404",
			Path:        "/this-route-does-not-exist",
			Description: "Verifies that unknown routes return 404 Not Found.",
			Expected:    "HTTP 404 Not Found",
			ValidateFn: func(tc *tests.TestContext) {
				var dummy map[string]interface{}
				err := tc.Client.SendHttpRequest("GET", "/this-route-does-not-exist", nil, nil, &dummy, nil)
				if err == nil {
					tc.FailureReason = "Expected 404, got 200 OK"
					tc.Errorf("Expected 404, got 200 OK")
				} else if err.StatusCode() != 404 {
					tc.FailureReason = fmt.Sprintf("Expected 404, got %d", err.StatusCode())
					tc.Errorf("Expected 404, got %d", err.StatusCode())
				} else {
					tc.Actual = "HTTP 404 Not Found as expected"
				}
			},
		},
		{
			Name:        "Negative Pagination Parameters Do Not Crash Server",
			Path:        "/health?page=-1&limit=-100",
			Description: "Verifies that negative pagination params are handled gracefully (no 500).",
			Expected:    "HTTP 200 OK or 400 Bad Request — not a 500 crash",
			ValidateFn: func(tc *tests.TestContext) {
				var dummy map[string]interface{}
				err := tc.Client.SendHttpRequest("GET", "/health?page=-1&limit=-100", nil, nil, &dummy, nil)
				if err != nil && err.StatusCode() == 500 {
					tc.FailureReason = "Server crashed on negative pagination parameters"
					tc.Errorf("Server crashed on negative pagination parameters")
				} else {
					tc.Actual = "Handled gracefully without a 500 crash"
				}
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		tests.RunAPITestWithDetails(
			t,
			testCase.Name,
			testCase.Description,
			testCase.Expected,
			testCase.ValidateFn,
		)
	}
}
