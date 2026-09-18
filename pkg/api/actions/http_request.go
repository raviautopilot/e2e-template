package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// sendHttpRequest executes an HTTP request and expects HTTP 2xx success.
// This is the internal engine — callers should use the verb-specific helpers (Get, Post, Put, etc.).
func sendHttpRequest(tc *tests.TestContext, c *client.Client, method string, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	err := c.SendHttpRequest(method, path, headers, reqBody, respBody, auth)
	if err != nil {
		tc.FailureReason = fmt.Sprintf("%s %s failed: %v", method, path, err)
		tc.Fatalf("%s %s failed: %v", method, path, err)
	}
	tc.Actual = fmt.Sprintf("HTTP 2xx OK for %s %s", method, path)
}

// sendHttpRequestAndExpectStatus executes an HTTP request and asserts the returned status code matches wantStatus.
// This is the internal engine — callers should use the verb-specific helpers (GetAndExpectStatus, PostAndExpectStatus, etc.).
func sendHttpRequestAndExpectStatus(tc *tests.TestContext, c *client.Client, method string, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	err := c.SendHttpRequest(method, path, headers, reqBody, respBody, auth)
	if wantStatus >= 200 && wantStatus < 300 {
		if err != nil {
			tc.FailureReason = fmt.Sprintf("Expected %d for %s %s, got error: %v", wantStatus, method, path, err)
			tc.Errorf("Expected %d for %s %s, got error: %v", wantStatus, method, path, err)
		} else {
			tc.Actual = fmt.Sprintf("HTTP %d OK as expected for %s %s", wantStatus, method, path)
		}
	} else {
		if err == nil {
			tc.FailureReason = fmt.Sprintf("Expected %d for %s %s, got 2xx success", wantStatus, method, path)
			tc.Errorf("Expected %d for %s %s, got 2xx success", wantStatus, method, path)
		} else if err.StatusCode() != wantStatus {
			tc.FailureReason = fmt.Sprintf("Expected %d for %s %s, got %d (body: %s)", wantStatus, method, path, err.StatusCode(), err.ResponseBody())
			tc.Errorf("Expected %d for %s %s, got %d (body: %s)", wantStatus, method, path, err.StatusCode(), err.ResponseBody())
		} else {
			tc.Actual = fmt.Sprintf("HTTP %d as expected for %s %s", wantStatus, method, path)
		}
	}
}
