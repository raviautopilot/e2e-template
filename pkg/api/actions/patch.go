package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Patch executes a PATCH request without asserting the response status.
// It is up to the caller to validate the returned error or response.
func Patch(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) error {
	return sendHttpRequest(c, "PATCH", path, headers, reqBody, respBody, auth)
}

// PatchAndExpectOK executes a PATCH request and asserts HTTP 2xx success.
func PatchAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	if err := Patch(tc, c, path, headers, reqBody, respBody, auth); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Fatalf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP 2xx OK for PATCH %s", path)
	}
}

// PatchAndExpectStatus executes a PATCH request and asserts the status code matches wantStatus.
func PatchAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	if err := sendHttpRequestAndExpectStatus(c, "PATCH", path, headers, reqBody, respBody, auth, wantStatus); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Errorf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP %d as expected for PATCH %s", wantStatus, path)
	}
}

// PatchAndExpectStatusCode is an alias for PatchAndExpectStatus.
func PatchAndExpectStatusCode(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	PatchAndExpectStatus(tc, c, path, headers, reqBody, respBody, auth, wantStatus)
}
