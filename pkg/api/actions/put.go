package actions

import (
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Put executes a PUT request and expects HTTP 2xx success.
func Put(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "PUT", path, headers, reqBody, respBody, auth)
}

// PutAndExpectOK is an alias for Put — executes a PUT request and expects HTTP 2xx success.
func PutAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "PUT", path, headers, reqBody, respBody, auth)
}

// PutAndExpectStatus executes a PUT request and asserts the status code matches wantStatus.
func PutAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	sendHttpRequestAndExpectStatus(tc, c, "PUT", path, headers, reqBody, respBody, auth, wantStatus)
}
