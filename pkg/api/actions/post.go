package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Post executes a POST request without asserting the response status.
// It is up to the caller to validate the returned error or response.
func Post(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) error {
	return sendHttpRequest(c, "POST", path, headers, reqBody, respBody, auth)
}

// PostAndExpectOK executes a POST request and asserts HTTP 2xx success.
func PostAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	if err := Post(tc, c, path, headers, reqBody, respBody, auth); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Fatalf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP 2xx OK for POST %s", path)
	}
}

// PostAndExpectStatus executes a POST request and asserts the status code matches wantStatus.
func PostAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	if err := sendHttpRequestAndExpectStatus(c, "POST", path, headers, reqBody, respBody, auth, wantStatus); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Errorf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP %d as expected for POST %s", wantStatus, path)
	}
}

// PostAndExpectStatusCode is an alias for PostAndExpectStatus.
func PostAndExpectStatusCode(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	PostAndExpectStatus(tc, c, path, headers, reqBody, respBody, auth, wantStatus)
}
