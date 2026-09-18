package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Head executes a HEAD request without asserting the response status.
// It is up to the caller to validate the returned error or response.
func Head(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) error {
	return sendHttpRequest(c, "HEAD", path, headers, reqBody, respBody, auth)
}

// HeadAndExpectOK executes a HEAD request and asserts HTTP 2xx success.
func HeadAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	if err := Head(tc, c, path, headers, reqBody, respBody, auth); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Fatalf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP 2xx OK for HEAD %s", path)
	}
}

// HeadAndExpectStatus executes a HEAD request and asserts the status code matches wantStatus.
func HeadAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	if err := sendHttpRequestAndExpectStatus(c, "HEAD", path, headers, reqBody, respBody, auth, wantStatus); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Errorf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP %d as expected for HEAD %s", wantStatus, path)
	}
}

// HeadAndExpectStatusCode is an alias for HeadAndExpectStatus.
func HeadAndExpectStatusCode(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	HeadAndExpectStatus(tc, c, path, headers, reqBody, respBody, auth, wantStatus)
}
