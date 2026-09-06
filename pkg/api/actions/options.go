package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// ─────────────────────────────────────────────────────────────────────────────
// OPTIONS Action Helpers
// ─────────────────────────────────────────────────────────────────────────────

// OptionsAndExpectOK sends an OPTIONS request and expects HTTP 200 OK.
func OptionsAndExpectOK(tc *tests.TestContext, c *client.Client, path string, resp interface{}) {
	err := c.SendHttpRequest("OPTIONS", path, nil, nil, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("OPTIONS %s failed: %v", path, err)
		tc.Fatalf("OPTIONS %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for OPTIONS %s", path)
}

// OptionsAndExpectStatus sends an OPTIONS request and expects a specific status code.
func OptionsAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, wantStatus int) {
	var dummy map[string]interface{}
	err := c.SendHttpRequest("OPTIONS", path, nil, nil, &dummy, nil)
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

// OptionsWithHeadersAndExpectOK sends an OPTIONS request with custom headers, expects 200 OK.
func OptionsWithHeadersAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, resp interface{}) {
	err := c.SendHttpRequest("OPTIONS", path, headers, nil, resp, nil)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("OPTIONS %s failed: %v", path, err)
		tc.Fatalf("OPTIONS %s failed: %v", path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 200 OK for OPTIONS %s", path)
}

// AuthenticatedOptions sends an OPTIONS request with auth credentials, expects 200 OK.
func AuthenticatedOptions(tc *tests.TestContext, c *client.Client, path string, auth client.Authenticator, resp interface{}) {
	err := c.SendHttpRequest("OPTIONS", path, nil, nil, resp, auth)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("Authenticated OPTIONS %s failed: %v", path, err)
		tc.Errorf("Authenticated OPTIONS %s failed: %v", path, err)
	} else {
		tc.Actual = fmt.Sprintf("HTTP 200 OK (authenticated) for OPTIONS %s", path)
	}
}
