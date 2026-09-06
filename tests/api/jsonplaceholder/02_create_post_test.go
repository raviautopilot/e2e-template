package jsonplaceholder_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_JSONPlaceholder_02_CreatePost verifies POST /posts creates a new resource (201 Created).
func TestAPI_JSONPlaceholder_02_CreatePost(t *testing.T) {
	c := client.NewClient("https://jsonplaceholder.typicode.com", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "POST /posts creates a new post (201 Created)",
		"Sends a JSON body to POST /posts. Verifies HTTP 201 Created and generated ID.",
		"HTTP 201 Created with generated post ID",
		func(tc *tests.TestContext) {
			newPost := map[string]interface{}{
				"title":  "E2E Test Post",
				"body":   "This post was created during automated E2E testing.",
				"userId": 1,
			}
			var created jsonPlaceholderPost
			actions.PostAndExpectCreated(tc, c, "/posts", &newPost, &created)
			actions.AssertNotZero(tc, "id", created.ID)
			actions.AssertEquals(tc, "title", created.Title, "E2E Test Post")
		},
	)
}
