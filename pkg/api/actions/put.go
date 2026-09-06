package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// PUT Action Helpers
// ─────────────────────────────────────────────────────────────────────────────

// PutAndExpectOK sends a PUT request with a JSON body pointer, expects 200 OK.
func PutAndExpectOK(tc *tests.TestContext, c *client.Client, path string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("PUT", path, nil, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("PUT %s failed: %v", path, err)
		tc.Fatalf("PUT %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for PUT %s", path)
}

// PutAndExpectStatus sends a PUT request and expects a specific HTTP status code.
func PutAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, body interface{}, wantStatus int) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("PUT", path, nil, body, &dummy, nil)
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

// PutWithHeadersAndExpectOK sends a PUT request with custom headers, expects 200 OK.
func PutWithHeadersAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("PUT", path, headers, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("PUT %s failed: %v", path, err)
		tc.Fatalf("PUT %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for PUT %s", path)
}

// AuthenticatedPut sends a PUT request with auth credentials, expects 200 OK.
func AuthenticatedPut(tc *tests.TestContext, c *client.Client, path string, auth client.Authenticator, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("PUT", path, nil, body, resp, auth)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("Authenticated PUT %s failed: %v", path, err)
		tc.Errorf("Authenticated PUT %s failed: %v", path, err)
	} else {
		tc.Actual = fmt.Sprintf("HTTP 200 OK (authenticated) for PUT %s", path)
	}
}
