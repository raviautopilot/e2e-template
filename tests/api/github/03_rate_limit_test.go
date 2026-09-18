package github_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

// TestAPI_GitHub_03_RateLimit verifies GitHub API rate_limit endpoint responses.
func TestAPI_GitHub_03_RateLimit(t *testing.T) {
	tests.RunAPITestWithClients(t, "GET /rate_limit shows current rate limit status",
		"Verifies the rate_limit endpoint returns valid info without authentication.",
		"HTTP 200 OK with rate limit information",
		apiClient, client2,
		func(tc *tests.TestContext) {
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, apiClient, "/rate_limit", githubHeaders, nil, &resp, nil)
			if resp["rate"] == nil {
				tc.Errorf("Expected non-nil rate limit response")
			}
		},
	)
}
