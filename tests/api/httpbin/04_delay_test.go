package httpbin_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

// TestAPI_HttpBin_04_Delay verifies the client gracefully waits for delayed responses.
func TestAPI_HttpBin_04_Delay(t *testing.T) {
	tests.RunAPITestWithClients(t, "GET /delay/1 responds within timeout",
		"Verifies the client waits for delayed responses without premature timeouts.",
		"HTTP 200 OK received within client timeout",
		apiClient, client2,
		func(tc *tests.TestContext) {
			actions.GetAndExpectOK(tc, apiClient, "/delay/1", nil)
		},
	)
}
