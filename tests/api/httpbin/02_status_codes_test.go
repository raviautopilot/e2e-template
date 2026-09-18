package httpbin_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

// TestAPI_HttpBin_02_StatusCodes verifies standard HTTP status responses (200, 404, 500).
func TestAPI_HttpBin_02_StatusCodes(t *testing.T) {
	tests.RunAPITestWithClients(t, "GET /status/{code} returns expected HTTP status codes",
		"Verifies that httpbin accurately returns requested status codes: 200 OK, 404 Not Found, 500 Server Error.",
		"HTTP 200, 404, and 500 status codes returned respectively",
		apiClient, client2,
		func(tc *tests.TestContext) {
			actions.GetAndExpectStatus(tc, apiClient, "/status/200", 200)
			actions.GetAndExpectStatus(tc, apiClient, "/status/404", 404)
			actions.GetAndExpectStatus(tc, apiClient, "/status/500", 500)
		},
	)
}
