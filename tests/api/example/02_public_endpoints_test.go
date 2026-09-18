package example_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

// TestAPI_Example_02_PublicEndpoints tests public/unauthenticated endpoints of your custom backend.
func TestAPI_Example_02_PublicEndpoints(t *testing.T) {
	if !isServerRunning() {
		t.Skipf("Skipping: no server running at %s (set baseUrl in config.json)", tests.GlobalConfig.BaseURL)
	}

	tests.RunAPITestWithClients(
		t,
		"Public Root Endpoint Returns 200 OK",
		"Verifies that the target application root responds with 200 OK.",
		"HTTP 200 OK with welcome message",
		apiClient,
		client2,
		func(tc *tests.TestContext) {
			var resp RootResponse
			actions.SendHttpRequest(tc, tc.Client, "GET", "/", nil, nil, &resp, nil)
			actions.AssertNotEmpty(tc, "message", resp.Message)
		},
	)
}
