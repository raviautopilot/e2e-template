package actions

import (
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Delete executes a DELETE request and expects HTTP 2xx success.
func Delete(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "DELETE", path, headers, reqBody, respBody, auth)
}

// DeleteAndExpectOK is an alias for Delete — executes a DELETE request and expects HTTP 2xx success.
func DeleteAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "DELETE", path, headers, reqBody, respBody, auth)
}

// DeleteAndExpectStatus executes a DELETE request and asserts the status code matches wantStatus.
func DeleteAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	sendHttpRequestAndExpectStatus(tc, c, "DELETE", path, headers, reqBody, respBody, auth, wantStatus)
}
