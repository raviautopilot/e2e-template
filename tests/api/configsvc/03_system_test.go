package configsvc_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_System_03_HealthCheck verifies the service health endpoint.
func TestAPI_System_03_HealthCheck(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "GET /health returns healthy status",
		"Get the health status of the microservice",
		"Health Check — HTTP 200 OK",
		func(tc *tests.TestContext) {
			var resp HealthResponse
			actions.GetAndExpectOK(tc, c, "/health", &resp)
			actions.AssertNotEmpty(tc, "status", resp.Status)
		},
	)

}
