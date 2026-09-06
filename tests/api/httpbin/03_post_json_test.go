package httpbin_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

type httpbinPostResponse struct {
	URL  string            `json:"url"`
	JSON map[string]string `json:"json"`
	Data string            `json:"data"`
}

// TestAPI_HttpBin_03_PostJSON verifies POST /post with JSON body echoes payload back.
func TestAPI_HttpBin_03_PostJSON(t *testing.T) {
	c := client.NewClient("https://httpbin.org", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "POST /post with JSON body is echoed back",
		"POST to httpbin.org/post with a JSON payload. Verifies 200 OK and payload echo.",
		"HTTP 200 OK with payload matched in response",
		func(tc *tests.TestContext) {
			payload := map[string]string{
				"framework": "go-e2e",
				"purpose":   "testing",
			}
			var resp httpbinPostResponse
			actions.PostAndExpectOK(tc, c, "/post", &payload, &resp)
			actions.AssertEquals(tc, "framework", resp.JSON["framework"], "go-e2e")
			actions.AssertEquals(tc, "purpose", resp.JSON["purpose"], "testing")
		},
	)
}
