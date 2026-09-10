package jsonplaceholder_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

type jsonPlaceholderPost struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

// TestAPI_JSONPlaceholder_01_Posts verifies GET /posts list and single post retrieval.
func TestAPI_JSONPlaceholder_01_Posts(t *testing.T) {
	c := client.NewClient("https://jsonplaceholder.typicode.com", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "GET /posts and GET /posts/1 validation",
		"Verifies fetching list of 100 posts, single post fields, and 404 for missing post.",
		"HTTP 200 OK with 100 items, valid single post, and HTTP 404 for id 9999",
		func(tc *tests.TestContext) {
			var posts []jsonPlaceholderPost
			actions.GetAndExpectOK(tc, c, "/posts", &posts)
			actions.AssertListLength(tc, "posts", len(posts), 100)

			var post jsonPlaceholderPost
			actions.GetAndExpectOK(tc, c, "/posts/1", &post)
			actions.AssertIntEquals(tc, "id", post.ID, 1)
			actions.AssertNotEmpty(tc, "title", post.Title)

			actions.GetAndExpectStatus(tc, c, "/posts/9999", 404)
		},
	)
}
