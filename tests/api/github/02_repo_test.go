package github_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

type githubRepo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

// TestAPI_GitHub_02_Repo verifies GitHub public repository lookup and 404 behavior.
func TestAPI_GitHub_02_Repo(t *testing.T) {
	tests.RunAPITestWithClients(t, "GET /repos/octocat/Hello-World returns public repo",
		"Verifies GitHub API returns public repository details, and 404 for non-existent repo.",
		"HTTP 200 OK with repository full_name, and 404 for missing repo",
		apiClient, client2,
		func(tc *tests.TestContext) {
			var repo githubRepo
			actions.GetWithHeadersAndExpectOK(tc, apiClient, "/repos/octocat/Hello-World", githubHeaders, &repo)
			actions.AssertEquals(tc, "full_name", repo.FullName, "octocat/Hello-World")

			actions.GetWithHeadersAndExpectStatus(tc, apiClient, "/repos/nonexistent-org-xyz/nonexistent-repo-abc", githubHeaders, 404)
		},
	)
}
