package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Get executes a GET request without asserting the response status.
// It is up to the caller to validate the returned error or response.
func Get(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) error {
	return sendHttpRequest(c, "GET", path, headers, reqBody, respBody, auth)
}

// GetAndExpectOK executes a GET request and asserts HTTP 2xx success.
func GetAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	if err := Get(tc, c, path, headers, reqBody, respBody, auth); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Fatalf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP 2xx OK for GET %s", path)
	}
}

// GetAndExpectStatus executes a GET request and asserts the status code matches wantStatus.
func GetAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	if err := sendHttpRequestAndExpectStatus(c, "GET", path, headers, reqBody, respBody, auth, wantStatus); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Errorf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP %d as expected for GET %s", wantStatus, path)
	}
}

// GetAndExpectStatusCode is an alias for GetAndExpectStatus.
func GetAndExpectStatusCode(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	GetAndExpectStatus(tc, c, path, headers, reqBody, respBody, auth, wantStatus)
}
