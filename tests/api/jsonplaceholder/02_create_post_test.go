package jsonplaceholder_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

// TestAPI_JSONPlaceholder_02_CreatePost verifies POST /posts creates a new resource (201 Created).
func TestAPI_JSONPlaceholder_02_CreatePost(t *testing.T) {
	tests.RunAPITestWithClients(t, "POST /posts creates a new post (201 Created)",
		"Sends a JSON body to POST /posts. Verifies HTTP 201 Created and generated ID.",
		"HTTP 201 Created with generated post ID",
		apiClient, client2,
		func(tc *tests.TestContext) {
			newPost := map[string]interface{}{
				"title":  "E2E Test Post",
				"body":   "This post was created during automated E2E testing.",
				"userId": 1,
			}
			var created jsonPlaceholderPost
			actions.Post(tc, apiClient, "/posts", nil, &newPost, &created, nil)
			actions.AssertNotZero(tc, "id", created.ID)
			actions.AssertEquals(tc, "title", created.Title, "E2E Test Post")
		},
	)
}
