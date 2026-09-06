package github_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_GitHub_03_RateLimit verifies GitHub API rate_limit endpoint responses.
func TestAPI_GitHub_03_RateLimit(t *testing.T) {
	c := client.NewClient("https://api.github.com", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "GET /rate_limit shows current rate limit status",
		"Verifies the rate_limit endpoint returns valid info without authentication.",
		"HTTP 200 OK with rate limit information",
		func(tc *tests.TestContext) {
			var resp map[string]interface{}
			actions.GetWithHeadersAndExpectOK(tc, c, "/rate_limit", githubHeaders, &resp)
			if resp["rate"] == nil {
				tc.Errorf("Expected non-nil rate limit response")
			}
		},
	)
}
