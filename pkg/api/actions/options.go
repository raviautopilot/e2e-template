package actions

import (
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Options executes an OPTIONS request and expects HTTP 2xx success.
func Options(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "OPTIONS", path, headers, reqBody, respBody, auth)
}

// OptionsAndExpectOK is an alias for Options — executes an OPTIONS request and expects HTTP 2xx success.
func OptionsAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "OPTIONS", path, headers, reqBody, respBody, auth)
}

// OptionsAndExpectStatus executes an OPTIONS request and asserts the status code matches wantStatus.
func OptionsAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	sendHttpRequestAndExpectStatus(tc, c, "OPTIONS", path, headers, reqBody, respBody, auth, wantStatus)
}
