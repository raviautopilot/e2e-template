package actions

import (
	"fmt"

	"e2e-template/pkg/client"
)

// sendHttpRequest executes an HTTP request and expects HTTP 2xx success.
// It operates independently of TestContext and returns an error if the request fails or returns a non-2xx status code.
func sendHttpRequest(c *client.Client, method string, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator) error {
	err := c.SendHttpRequest(method, path, headers, reqBody, respBody, auth)
	if err != nil {
		return fmt.Errorf("%s %s failed: %w", method, path, err)
	}
	return nil
}

// sendHttpRequestAndExpectStatus executes an HTTP request and asserts that the returned status code matches wantStatus.
// It operates independently of TestContext and returns an error if the status code does not match wantStatus.
func sendHttpRequestAndExpectStatus(c *client.Client, method string, path string, headers map[string]string, reqBody interface{}, respBody interface{}, auth client.Authenticator, wantStatus int) error {
	err := c.SendHttpRequest(method, path, headers, reqBody, respBody, auth)
	if wantStatus >= 200 && wantStatus < 300 {
		if err != nil {
			return fmt.Errorf("expected %d for %s %s, got error: %w", wantStatus, method, path, err)
		}
		return nil
	}
	if err == nil {
		return fmt.Errorf("expected %d for %s %s, got 2xx success", wantStatus, method, path)
	}
	if err.StatusCode() != wantStatus {
		return fmt.Errorf("expected %d for %s %s, got %d (body: %s)", wantStatus, method, path, err.StatusCode(), err.ResponseBody())
	}
	return nil
}
