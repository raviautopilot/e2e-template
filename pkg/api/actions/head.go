package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// HEAD Action Helpers
// ─────────────────────────────────────────────────────────────────────────────

// HeadAndExpectOK sends a HEAD request and expects HTTP 200 OK.
func HeadAndExpectOK(tc *tests.TestContext, c *client.Client, path string) {
	err := c.SendHttpRequest("HEAD", path, nil, nil, nil, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("HEAD %s failed: %v", path, err)
		tc.Fatalf("HEAD %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for HEAD %s", path)
}

// HeadAndExpectStatus sends a HEAD request and expects a specific status code.
func HeadAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, wantStatus int) {
	err := c.SendHttpRequest("HEAD", path, nil, nil, nil, nil)
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

// HeadWithHeadersAndExpectOK sends a HEAD request with custom headers, expects 200 OK.
func HeadWithHeadersAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string) {
	err := c.SendHttpRequest("HEAD", path, headers, nil, nil, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("HEAD %s failed: %v", path, err)
		tc.Fatalf("HEAD %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for HEAD %s", path)
}

// AuthenticatedHead sends a HEAD request with auth credentials, expects 200 OK.
func AuthenticatedHead(tc *tests.TestContext, c *client.Client, path string, auth client.Authenticator) {
	err := c.SendHttpRequest("HEAD", path, nil, nil, nil, auth)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("Authenticated HEAD %s failed: %v", path, err)
		tc.Errorf("Authenticated HEAD %s failed: %v", path, err)
	} else {
		tc.Actual = fmt.Sprintf("HTTP 200 OK (authenticated) for HEAD %s", path)
	}
}
