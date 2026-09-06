package example_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

// TestAPI_Example_01_HealthCheck tests the health/ping endpoint of your custom backend.
func TestAPI_Example_01_HealthCheck(t *testing.T) {
	if !isServerRunning() {
		t.Skipf("Skipping: no server running at %s (set baseUrl in config.json)", tests.GlobalConfig.BaseURL)
	}

	tests.RunAPITestWithDetails(
		t,
		"Health Check Endpoint Returns 200 OK",
		"Verifies that the target application's health endpoint responds with 200 OK.",
		"HTTP 200 OK with healthy status",
		func(tc *tests.TestContext) {
			var resp HealthResponse
			actions.GetAndExpectOK(tc, tc.Client, "/health", &resp)
			actions.AssertNotEmpty(tc, "status", resp.Status)
		},
	)
}
