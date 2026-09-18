package actions

import (
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Get executes a GET request and expects HTTP 2xx success.
func Get(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "GET", path, headers, reqBody, respBody, auth)
}

// GetAndExpectOK is an alias for Get — executes a GET request and expects HTTP 2xx success.
func GetAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "GET", path, headers, reqBody, respBody, auth)
}

// GetAndExpectStatus executes a GET request and asserts the status code matches wantStatus.
func GetAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	sendHttpRequestAndExpectStatus(tc, c, "GET", path, headers, reqBody, respBody, auth, wantStatus)
}
