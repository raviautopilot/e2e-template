package httpbin_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_HttpBin_02_StatusCodes verifies standard HTTP status responses (200, 404, 500).
func TestAPI_HttpBin_02_StatusCodes(t *testing.T) {
	c := client.NewClient("https://httpbin.org", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "GET /status/{code} returns expected HTTP status codes",
		"Verifies that httpbin accurately returns requested status codes: 200 OK, 404 Not Found, 500 Server Error.",
		"HTTP 200, 404, and 500 status codes returned respectively",
		func(tc *tests.TestContext) {
			actions.GetAndExpectStatus(tc, c, "/status/200", 200)
			actions.GetAndExpectStatus(tc, c, "/status/404", 404)
			actions.GetAndExpectStatus(tc, c, "/status/500", 500)
		},
	)
}
