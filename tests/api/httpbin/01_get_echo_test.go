package httpbin_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

type httpbinGetResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Args    map[string]string `json:"args"`
}

// TestAPI_HttpBin_01_GetEcho verifies GET /get echoes back request metadata and parameters.
func TestAPI_HttpBin_01_GetEcho(t *testing.T) {
	tests.RunAPITestWithClients(t, "GET /get returns 200 OK with request echo",
		"httpbin.org/get echoes the request back. Verifies 200 OK and URL field is non-empty.",
		"HTTP 200 OK with non-empty url field in response",
		apiClient, client2,
		func(tc *tests.TestContext) {
			var resp httpbinGetResponse
			actions.Get(tc, apiClient, "/get", nil, nil, &resp, nil)
			actions.AssertNotEmpty(tc, "url", resp.URL)

			var paramResp httpbinGetResponse
			actions.Get(tc, apiClient, "/get?search=golang&limit=10", nil, nil, &paramResp, nil)
			actions.AssertEquals(tc, "search", paramResp.Args["search"], "golang")
			actions.AssertEquals(tc, "limit", paramResp.Args["limit"], "10")
		},
	)
}
