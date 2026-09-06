package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// API Test Action Helpers
//
// These functions extract the repetitive HTTP call + assertion boilerplate
// into callable actions, so test bodies become 2-3 lines instead of 10-15.
//
// Usage in tests:
//   var resp MyType
//   actions.GetAndExpectOK(tc, c, "/endpoint", &resp)
//   actions.AssertNotEmpty(tc, "fieldName", resp.Field)
// ─────────────────────────────────────────────────────────────────────────────

// GetAndExpectOK sends a GET request and expects HTTP 200 OK.
func GetAndExpectOK(tc *tests.TestContext, c *client.Client, path string, resp interface{}) {
	err := c.SendHttpRequest("GET", path, nil, nil, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("GET %s failed: %v", path, err)
		tc.Fatalf("GET %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for GET %s", path)
}

// GetAndExpectStatus sends a GET request and expects a specific HTTP status code.
func GetAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, wantStatus int) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("GET", path, nil, nil, &dummy, nil)
	if wantStatus >= 200 && wantStatus < 300 {
		// Expecting success
		if err != nil {
			tc.FailureReason = fmt.Sprintf("Expected %d, got error: %v", wantStatus, err)
			tc.Errorf("Expected %d, got error: %v", wantStatus, err)
		} else {
			tc.Actual = fmt.Sprintf("HTTP %d OK as expected", wantStatus)
		}
	} else {
		// Expecting error status
		if err == nil {
			tc.FailureReason = fmt.Sprintf("Expected %d, got 200 OK", wantStatus)
			tc.Errorf("Expected %d, got 200 OK", wantStatus)
		} else if err.StatusCode() != wantStatus {
			tc.FailureReason = fmt.Sprintf("Expected %d, got %d", wantStatus, err.StatusCode())
			tc.Errorf("Expected %d, got %d", wantStatus, err.StatusCode())
		} else {
			tc.Actual = fmt.Sprintf("HTTP %d as expected", wantStatus)
		}
	}
}

// GetAndExpectNotCrash sends a GET request and fails only if it returns HTTP 500.
func GetAndExpectNotCrash(tc *tests.TestContext, c *client.Client, path string) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("GET", path, nil, nil, &dummy, nil)
	if err != nil && err.StatusCode() == 500 {
		tc.FailureReason = fmt.Sprintf("Server crashed with 500 for GET %s", path)
		tc.Errorf("Server crashed with 500 for GET %s: %v", path, err)
	} else {
		tc.Actual = "Handled gracefully without a 500 crash"
	}
}

// GetWithHeadersAndExpectOK sends a GET request with custom headers, expects 200 OK.
func GetWithHeadersAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, resp interface{}) {
	err := c.SendHttpRequest("GET", path, headers, nil, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("GET %s failed: %v", path, err)
		tc.Fatalf("GET %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for GET %s", path)
}

// GetWithHeadersAndExpectStatus sends a GET request with headers and expects a specific status.
func GetWithHeadersAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, wantStatus int) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("GET", path, headers, nil, &dummy, nil)
	if err == nil {
		if wantStatus >= 300 {
			tc.FailureReason = fmt.Sprintf("Expected %d, got 200 OK", wantStatus)
			tc.Errorf("Expected %d, got 200 OK", wantStatus)
		} else {
			tc.Actual = "HTTP 200 OK as expected"
		}
	} else if err.StatusCode() != wantStatus {
		tc.FailureReason = fmt.Sprintf("Expected %d, got %d", wantStatus, err.StatusCode())
		tc.Errorf("Expected %d, got %d", wantStatus, err.StatusCode())
	} else {
		tc.Actual = fmt.Sprintf("HTTP %d as expected", wantStatus)
	}
}

// PostAndExpectOK sends a POST request with a JSON body, expects 200 OK.
func PostAndExpectOK(tc *tests.TestContext, c *client.Client, path string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("POST", path, nil, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("POST %s failed: %v", path, err)
		tc.Fatalf("POST %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for POST %s", path)
}

// PostAndExpectCreated sends a POST request, accepts either 200 or 201 as success.
func PostAndExpectCreated(tc *tests.TestContext, c *client.Client, path string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("POST", path, nil, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("POST %s failed: %v", path, err)
		tc.Fatalf("POST %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 2xx Created for POST %s", path)
}

// AuthenticatedGet sends a GET request with auth credentials, expects 200 OK.
func AuthenticatedGet(tc *tests.TestContext, c *client.Client, path string, auth client.Authenticator, resp interface{}) {
	err := c.SendHttpRequest("GET", path, nil, nil, resp, auth)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("Authenticated GET %s failed: %v", path, err)
		tc.Errorf("Authenticated GET %s failed: %v", path, err)
	} else {
		tc.Actual = fmt.Sprintf("HTTP 200 OK (authenticated) for GET %s", path)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Assertion Helpers
// ─────────────────────────────────────────────────────────────────────────────

// AssertNotEmpty fails the test if the value is empty.
func AssertNotEmpty(tc *tests.TestContext, fieldName string, value string) {
	if value == "" {
		tc.FailureReason = fmt.Sprintf("Expected non-empty %s", fieldName)
		tc.Errorf("Expected non-empty %s", fieldName)
	}
}

// AssertEquals fails the test if got != want.
func AssertEquals(tc *tests.TestContext, fieldName string, got, want string) {
	if got != want {
		tc.FailureReason = fmt.Sprintf("Expected %s=%q, got %q", fieldName, want, got)
		tc.Errorf("Expected %s=%q, got %q", fieldName, want, got)
	}
}

// AssertIntEquals fails the test if got != want.
func AssertIntEquals(tc *tests.TestContext, fieldName string, got, want int) {
	if got != want {
		tc.FailureReason = fmt.Sprintf("Expected %s=%d, got %d", fieldName, want, got)
		tc.Errorf("Expected %s=%d, got %d", fieldName, want, got)
	}
}

// AssertNotZero fails the test if the value is zero.
func AssertNotZero(tc *tests.TestContext, fieldName string, value int) {
	if value == 0 {
		tc.FailureReason = fmt.Sprintf("Expected non-zero %s", fieldName)
		tc.Errorf("Expected non-zero %s", fieldName)
	}
}

// AssertNonEmptyList fails the test if the list is empty.
func AssertNonEmptyList(tc *tests.TestContext, fieldName string, length int) {
	if length == 0 {
		tc.FailureReason = fmt.Sprintf("Expected non-empty %s list", fieldName)
		tc.Errorf("Expected non-empty %s list", fieldName)
	}
}

// AssertListLength fails the test if the list length doesn't match.
func AssertListLength(tc *tests.TestContext, fieldName string, got, want int) {
	if got != want {
		tc.FailureReason = fmt.Sprintf("Expected %s count=%d, got %d", fieldName, want, got)
		tc.Errorf("Expected %s count=%d, got %d", fieldName, want, got)
	}
}
