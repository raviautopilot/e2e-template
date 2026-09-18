package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Delete executes a DELETE request without asserting the response status.
// It is up to the caller to validate the returned error or response.
func Delete(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) error {
	return sendHttpRequest(c, "DELETE", path, headers, reqBody, respBody, auth)
}

// DeleteAndExpectOK executes a DELETE request and asserts HTTP 2xx success.
func DeleteAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	if err := Delete(tc, c, path, headers, reqBody, respBody, auth); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Fatalf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP 2xx OK for DELETE %s", path)
	}
}

// DeleteAndExpectStatus executes a DELETE request and asserts the status code matches wantStatus.
func DeleteAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	if err := sendHttpRequestAndExpectStatus(c, "DELETE", path, headers, reqBody, respBody, auth, wantStatus); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Errorf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP %d as expected for DELETE %s", wantStatus, path)
	}
}

// DeleteAndExpectStatusCode is an alias for DeleteAndExpectStatus.
func DeleteAndExpectStatusCode(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	DeleteAndExpectStatus(tc, c, path, headers, reqBody, respBody, auth, wantStatus)
}
