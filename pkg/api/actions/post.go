package actions

import (
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Post executes a POST request and expects HTTP 2xx success.
func Post(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "POST", path, headers, reqBody, respBody, auth)
}

// PostAndExpectOK is an alias for Post — executes a POST request and expects HTTP 2xx success.
func PostAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "POST", path, headers, reqBody, respBody, auth)
}

// PostAndExpectStatus executes a POST request and asserts the status code matches wantStatus.
func PostAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	sendHttpRequestAndExpectStatus(tc, c, "POST", path, headers, reqBody, respBody, auth, wantStatus)
}
