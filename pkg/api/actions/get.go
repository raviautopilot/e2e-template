package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// GET Action Helpers
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
		if err != nil {
			tc.FailureReason = fmt.Sprintf("Expected %d, got error: %v", wantStatus, err)
			tc.Errorf("Expected %d, got error: %v", wantStatus, err)
		} else {
			tc.Actual = fmt.Sprintf("HTTP %d OK as expected", wantStatus)
		}
	} else {
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
