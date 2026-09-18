package actions

import (
	"fmt"

	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Options executes an OPTIONS request without asserting the response status.
// It is up to the caller to validate the returned error or response.
func Options(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) error {
	return sendHttpRequest(c, "OPTIONS", path, headers, reqBody, respBody, auth)
}

// OptionsAndExpectOK executes an OPTIONS request and asserts HTTP 2xx success.
func OptionsAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	if err := Options(tc, c, path, headers, reqBody, respBody, auth); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Fatalf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP 2xx OK for OPTIONS %s", path)
	}
}

// OptionsAndExpectStatus executes an OPTIONS request and asserts the status code matches wantStatus.
func OptionsAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	if err := sendHttpRequestAndExpectStatus(c, "OPTIONS", path, headers, reqBody, respBody, auth, wantStatus); err != nil {
		if tc != nil {
			tc.FailureReason = err.Error()
			tc.Errorf("%v", err)
		}
		return
	}
	if tc != nil {
		tc.Actual = fmt.Sprintf("HTTP %d as expected for OPTIONS %s", wantStatus, path)
	}
}

// OptionsAndExpectStatusCode is an alias for OptionsAndExpectStatus.
func OptionsAndExpectStatusCode(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	OptionsAndExpectStatus(tc, c, path, headers, reqBody, respBody, auth, wantStatus)
}
