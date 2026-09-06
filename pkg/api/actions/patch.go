package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// PATCH Action Helpers
// ─────────────────────────────────────────────────────────────────────────────

// PatchAndExpectOK sends a PATCH request with a JSON body pointer, expects 200 OK.
func PatchAndExpectOK(tc *tests.TestContext, c *client.Client, path string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("PATCH", path, nil, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("PATCH %s failed: %v", path, err)
		tc.Fatalf("PATCH %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for PATCH %s", path)
}

// PatchAndExpectStatus sends a PATCH request and expects a specific HTTP status code.
func PatchAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, body interface{}, wantStatus int) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("PATCH", path, nil, body, &dummy, nil)
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

// PatchWithHeadersAndExpectOK sends a PATCH request with custom headers, expects 200 OK.
func PatchWithHeadersAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("PATCH", path, headers, body, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("PATCH %s failed: %v", path, err)
		tc.Fatalf("PATCH %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for PATCH %s", path)
}

// AuthenticatedPatch sends a PATCH request with auth credentials, expects 200 OK.
func AuthenticatedPatch(tc *tests.TestContext, c *client.Client, path string, auth client.Authenticator, body interface{}, resp interface{}) {
	err := c.SendHttpRequest("PATCH", path, nil, body, resp, auth)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("Authenticated PATCH %s failed: %v", path, err)
		tc.Errorf("Authenticated PATCH %s failed: %v", path, err)
	} else {
		tc.Actual = fmt.Sprintf("HTTP 200 OK (authenticated) for PATCH %s", path)
	}
}
