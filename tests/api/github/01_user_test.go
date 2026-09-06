package github_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

type githubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Type  string `json:"type"`
}

var githubHeaders = map[string]string{
	"Accept":               "application/vnd.github.v3+json",
	"X-GitHub-Api-Version": "2022-11-28",
}

// TestAPI_GitHub_01_User verifies GitHub user profile lookup and 404 behavior.
func TestAPI_GitHub_01_User(t *testing.T) {
	c := client.NewClient("https://api.github.com", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "GET /users/octocat returns valid GitHub user",
		"Verifies GitHub public API returns octocat's profile, and 404 for non-existent users.",
		"HTTP 200 OK for existing user, HTTP 404 for missing user",
		func(tc *tests.TestContext) {
			var user githubUser
			actions.GetWithHeadersAndExpectOK(tc, c, "/users/octocat", githubHeaders, &user)
			actions.AssertEquals(tc, "login", user.Login, "octocat")
			actions.AssertNotZero(tc, "id", user.ID)

			actions.GetWithHeadersAndExpectStatus(tc, c, "/users/nonexistent-user-xyz-abc-123", githubHeaders, 404)
		},
	)
}
