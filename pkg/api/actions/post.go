package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// POST Action Helpers
// ─────────────────────────────────────────────────────────────────────────────

// PostAndExpectOK sends a POST request with a JSON body pointer, expects 200 OK.
func PostAndExpectOK(tc *tests.TestContext, c *client.Client, path string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("POST", path, nil, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("POST %s failed: %v", path, err)
		tc.Fatalf("POST %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for POST %s", path)
}

// PostAndExpectCreated sends a POST request with a JSON body pointer, accepts 200 or 201 as success.
func PostAndExpectCreated(tc *tests.TestContext, c *client.Client, path string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("POST", path, nil, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("POST %s failed: %v", path, err)
		tc.Fatalf("POST %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 2xx Created for POST %s", path)
}

// PostAndExpectStatus sends a POST request and expects a specific HTTP status code.
func PostAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, body interface{}, wantStatus int) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("POST", path, nil, body, &dummy, nil)
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

// PostWithHeadersAndExpectOK sends a POST request with custom headers and body pointer, expects 200 OK.
func PostWithHeadersAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("POST", path, headers, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("POST %s failed: %v", path, err)
		tc.Fatalf("POST %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for POST %s", path)
}

// PostWithHeadersAndExpectStatus sends a POST request with headers and expects a specific status.
func PostWithHeadersAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, body interface{}, wantStatus int) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("POST", path, headers, body, &dummy, nil)
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

// AuthenticatedPost sends a POST request with auth credentials, expects 200/201.
func AuthenticatedPost(tc *tests.TestContext, c *client.Client, path string, auth client.Authenticator, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("POST", path, nil, body, resp, auth)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("Authenticated POST %s failed: %v", path, err)
		tc.Errorf("Authenticated POST %s failed: %v", path, err)
	} else {
		tc.Actual = fmt.Sprintf("HTTP 200 OK (authenticated) for POST %s", path)
	}
}
