package actions

import (
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// Patch executes a PATCH request and expects HTTP 2xx success.
func Patch(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "PATCH", path, headers, reqBody, respBody, auth)
}

// PatchAndExpectOK is an alias for Patch — executes a PATCH request and expects HTTP 2xx success.
func PatchAndExpectOK(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) {
	sendHttpRequest(tc, c, "PATCH", path, headers, reqBody, respBody, auth)
}

// PatchAndExpectStatus executes a PATCH request and asserts the status code matches wantStatus.
func PatchAndExpectStatus(tc *tests.TestContext, c *client.Client, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) {
	sendHttpRequestAndExpectStatus(tc, c, "PATCH", path, headers, reqBody, respBody, auth, wantStatus)
}
