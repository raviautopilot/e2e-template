package actions

import (
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Head executes a HEAD request and expects HTTP 2xx success.
func Head(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "HEAD", path, headers, reqBody, respBody, auth)
}

// HeadAndExpectOK is an alias for Head — executes a HEAD request and expects HTTP 2xx success.
func HeadAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "HEAD", path, headers, reqBody, respBody, auth)
}

// HeadAndExpectStatus executes a HEAD request and asserts the status code matches wantStatus.
func HeadAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	sendHttpRequestAndExpectStatus(tc, c, "HEAD", path, headers, reqBody, respBody, auth, wantStatus)
}
