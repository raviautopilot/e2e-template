package httpbin_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_HttpBin_04_Delay verifies the client gracefully waits for delayed responses.
func TestAPI_HttpBin_04_Delay(t *testing.T) {
	c := client.NewClient("https://httpbin.org", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "GET /delay/1 responds within timeout",
		"Verifies the client waits for delayed responses without premature timeouts.",
		"HTTP 200 OK received within client timeout",
		func(tc *tests.TestContext) {
			actions.GetAndExpectOK(tc, c, "/delay/1", nil)
		},
	)
}
