package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// DELETE Action Helpers
// ─────────────────────────────────────────────────────────────────────────────

// DeleteAndExpectOK sends a DELETE request and expects HTTP 200 OK.
func DeleteAndExpectOK(tc *tests.TestContext, c *client.Client, path string, resp interface{}) {
	err := c.SendHttpRequest("DELETE", path, nil, nil, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("DELETE %s failed: %v", path, err)
		tc.Fatalf("DELETE %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for DELETE %s", path)
}

// DeleteAndExpectNoContent sends a DELETE request and expects HTTP 204 No Content (or 200).
func DeleteAndExpectNoContent(tc *tests.TestContext, c *client.Client, path string) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("DELETE", path, nil, nil, &dummy, nil)
	if err != nil && err.StatusCode() != 204 && err.StatusCode() != 200 {
		tc.FailureReason = fmt.Sprintf("DELETE %s expected 204/200, got: %v", path, err)
		tc.Fatalf("DELETE %s expected 204/200, got: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 204/200 for DELETE %s", path)
}

// DeleteAndExpectStatus sends a DELETE request and expects a specific HTTP status code.
func DeleteAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, wantStatus int) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("DELETE", path, nil, nil, &dummy, nil)
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

// DeleteWithHeadersAndExpectOK sends a DELETE request with custom headers, expects 200 OK.
func DeleteWithHeadersAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, resp interface{}) {
	err := c.SendHttpRequest("DELETE", path, headers, nil, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("DELETE %s failed: %v", path, err)
		tc.Fatalf("DELETE %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for DELETE %s", path)
}

// AuthenticatedDelete sends a DELETE request with auth credentials, expects 200/204.
func AuthenticatedDelete(tc *tests.TestContext, c *client.Client, path string, auth client.Authenticator, resp interface{}) {
	err := c.SendHttpRequest("DELETE", path, nil, nil, resp, auth)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("Authenticated DELETE %s failed: %v", path, err)
		tc.Errorf("Authenticated DELETE %s failed: %v", path, err)
	} else {
		tc.Actual = fmt.Sprintf("HTTP 200/204 OK (authenticated) for DELETE %s", path)
	}
}
