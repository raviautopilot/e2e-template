package api_test

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC SITE API TESTS
//
// These tests target real, publicly-available HTTP APIs so the template works
// out-of-the-box without any local server running.
//
// APIs used:
//   - https://httpbin.org  — echo/inspect service (HTTP mechanics testing)
//   - https://api.github.com — GitHub public REST API (no auth needed)
//   - https://jsonplaceholder.typicode.com — Fake REST API (CRUD patterns)
//
// Run with:
//   go test -v ./tests/api/... -run TestAPI_Public
// ─────────────────────────────────────────────────────────────────────────────

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/pkg/logger"
	"e2e-template/tests"
)

// newPublicClient creates a Client pointed at an external base URL.
func newPublicClient(baseURL string, logDir string) *client.Client {
	return client.NewClient(baseURL, 15*time.Second, logDir)
}

// ─────────────────────────────────────────────────────────────────────────────
// 01 · httpbin.org — HTTP mechanics testing
// ─────────────────────────────────────────────────────────────────────────────

type httpbinGetResponse struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Args    map[string]string `json:"args"`
}

type httpbinPostResponse struct {
	URL  string            `json:"url"`
	JSON map[string]string `json:"json"`
	Data string            `json:"data"`
}

func TestAPI_Public_01_HttpBin(t *testing.T) {
	const baseURL = "https://httpbin.org"
	logDir := tests.ExecutionLogDir

	type testCase struct {
		Name        string
		Description string
		Expected    string
		RunFn       func(tc *tests.TestContext, c *client.Client)
	}

	testCases := []testCase{
		{
			Name:        "GET /get returns 200 OK with request echo",
			Description: "httpbin.org/get echoes the request back. Verifies 200 OK and URL field is non-empty.",
			Expected:    "HTTP 200 OK with non-empty url field in response",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var resp httpbinGetResponse
				actions.GetAndExpectOK(tc, c, "/get", &resp)
				actions.AssertNotEmpty(tc, "url", resp.URL)
			},
		},
		{
			Name:        "GET /get with custom query params are echoed back",
			Description: "Verifies that query parameters are correctly echoed back by httpbin.",
			Expected:    "HTTP 200 OK with args.foo=bar and args.baz=42",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var resp httpbinGetResponse
				actions.GetAndExpectOK(tc, c, "/get?foo=bar&baz=42", &resp)
				tc.Actual = fmt.Sprintf("args=%v", resp.Args)
				actions.AssertEquals(tc, "args.foo", resp.Args["foo"], "bar")
				actions.AssertEquals(tc, "args.baz", resp.Args["baz"], "42")
			},
		},
		{
			Name:        "GET /status/200 returns 200 OK",
			Description: "httpbin.org/status/:code returns the requested HTTP status code.",
			Expected:    "HTTP 200 OK",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				actions.GetAndExpectStatus(tc, c, "/status/200", 200)
			},
		},
		{
			Name:        "GET /status/404 returns 404 Not Found",
			Description: "Verifies that requesting a 404 status code returns 404 Not Found.",
			Expected:    "HTTP 404 Not Found",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				actions.GetAndExpectStatus(tc, c, "/status/404", 404)
			},
		},
		{
			Name:        "GET /status/500 returns 500 Internal Server Error",
			Description: "Verifies that a 500 status is correctly returned and handled (not swallowed).",
			Expected:    "HTTP 500 Internal Server Error returned",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				actions.GetAndExpectStatus(tc, c, "/status/500", 500)
			},
		},
		{
			Name:        "POST /post with JSON body is echoed back",
			Description: "Verifies POST with a JSON payload is accepted and body is echoed.",
			Expected:    "HTTP 200 OK with json.key=value echoed in response",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				reqBody := map[string]string{"key": "value", "framework": "e2e-template"}
				var resp httpbinPostResponse
				actions.PostAndExpectOK(tc, c, "/post", &reqBody, &resp)
				actions.AssertEquals(tc, "json.key", resp.JSON["key"], "value")
			},
		},
		{
			Name:        "GET /delay/1 responds within timeout",
			Description: "Verifies a 1-second delayed response completes within the client timeout.",
			Expected:    "HTTP 200 OK within 15-second client timeout",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var resp httpbinGetResponse
				actions.GetAndExpectOK(tc, c, "/delay/1", &resp)
			},
		},
		{
			Name:        "Bearer Token Auth via /bearer endpoint",
			Description: "Verifies that a valid Bearer token is accepted by httpbin's /bearer endpoint.",
			Expected:    "HTTP 200 OK when Authorization: Bearer <token> header is sent",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				auth := &client.BearerTokenAuth{Token: "test-framework-token"}
				var resp map[string]interface{}
				actions.AuthenticatedGet(tc, c, "/bearer", auth, &resp)
			},
		},
		{
			Name:        "Basic Auth via /basic-auth endpoint",
			Description: "Verifies that correct Basic Auth credentials are accepted by httpbin.",
			Expected:    "HTTP 200 OK with correct Basic Auth user=user, password=pass",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				auth := &client.BasicAuth{Username: "user", Password: "pass"}
				var resp map[string]interface{}
				actions.AuthenticatedGet(tc, c, "/basic-auth/user/pass", auth, &resp)
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		c := newPublicClient(baseURL, logDir)
		tests.RunAPITestWithDetails(t, testCase.Name, testCase.Description, testCase.Expected,
			func(tc *tests.TestContext) { testCase.RunFn(tc, c) },
		)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 02 · GitHub Public REST API
// ─────────────────────────────────────────────────────────────────────────────

type githubUser struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
}

type githubRepo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

func TestAPI_Public_02_GitHub(t *testing.T) {
	const baseURL = "https://api.github.com"
	logDir := tests.ExecutionLogDir
	githubHeaders := map[string]string{
		"Accept":               "application/vnd.github.v3+json",
		"X-GitHub-Api-Version": "2022-11-28",
	}

	type testCase struct {
		Name        string
		Description string
		Expected    string
		RunFn       func(tc *tests.TestContext, c *client.Client)
	}

	testCases := []testCase{
		{
			Name:        "GET /users/octocat returns valid GitHub user",
			Description: "Verifies the GitHub API returns octocat's public profile.",
			Expected:    "HTTP 200 OK with login=octocat, id>0",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var user githubUser
				actions.GetWithHeadersAndExpectOK(tc, c, "/users/octocat", githubHeaders, &user)
				actions.AssertEquals(tc, "login", user.Login, "octocat")
				actions.AssertNotZero(tc, "id", user.ID)
			},
		},
		{
			Name:        "GET /users/nonexistent-user-xyz-abc-123 returns 404",
			Description: "Verifies that fetching a non-existent GitHub user returns 404.",
			Expected:    "HTTP 404 Not Found",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				actions.GetWithHeadersAndExpectStatus(tc, c, "/users/nonexistent-user-xyz-abc-123", githubHeaders, 404)
			},
		},
		{
			Name:        "GET /repos/octocat/Hello-World returns valid public repo",
			Description: "Verifies the GitHub API returns the Hello-World public repo.",
			Expected:    "HTTP 200 OK with full_name=octocat/Hello-World",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var repo githubRepo
				actions.GetWithHeadersAndExpectOK(tc, c, "/repos/octocat/Hello-World", githubHeaders, &repo)
				actions.AssertEquals(tc, "full_name", repo.FullName, "octocat/Hello-World")
			},
		},
		{
			Name:        "GET /repos/nonexistent-org-xyz/nonexistent-repo-abc returns 404",
			Description: "Verifies that a non-existent repository returns 404.",
			Expected:    "HTTP 404 Not Found",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				actions.GetWithHeadersAndExpectStatus(tc, c, "/repos/nonexistent-org-xyz/nonexistent-repo-abc", githubHeaders, 404)
			},
		},
		{
			Name:        "GET /rate_limit shows current rate limit status",
			Description: "Verifies the rate_limit endpoint returns valid info (no auth needed).",
			Expected:    "HTTP 200 OK with rate limit information",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var resp map[string]interface{}
				actions.GetWithHeadersAndExpectOK(tc, c, "/rate_limit", githubHeaders, &resp)
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		c := newPublicClient(baseURL, logDir)
		tests.RunAPITestWithDetails(t, testCase.Name, testCase.Description, testCase.Expected,
			func(tc *tests.TestContext) { testCase.RunFn(tc, c) },
		)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 03 · JSONPlaceholder — Fake REST API (typicode.com)
// ─────────────────────────────────────────────────────────────────────────────

type jsonPlaceholderPost struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func TestAPI_Public_03_JSONPlaceholder(t *testing.T) {
	const baseURL = "https://jsonplaceholder.typicode.com"
	logDir := tests.ExecutionLogDir

	type testCase struct {
		Name        string
		Description string
		Expected    string
		RunFn       func(tc *tests.TestContext, c *client.Client)
	}

	testCases := []testCase{
		{
			Name:        "GET /posts returns a list of 100 posts",
			Description: "Verifies GET /posts returns all 100 placeholder posts.",
			Expected:    "HTTP 200 OK with exactly 100 posts",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var posts []jsonPlaceholderPost
				actions.GetAndExpectOK(tc, c, "/posts", &posts)
				actions.AssertListLength(tc, "posts", len(posts), 100)
			},
		},
		{
			Name:        "GET /posts/1 returns a single post with valid fields",
			Description: "Verifies GET /posts/1 returns a post with non-empty title.",
			Expected:    "HTTP 200 OK with id=1, non-empty title",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var post jsonPlaceholderPost
				actions.GetAndExpectOK(tc, c, "/posts/1", &post)
				actions.AssertIntEquals(tc, "id", post.ID, 1)
				actions.AssertNotEmpty(tc, "title", post.Title)
			},
		},
		{
			Name:        "GET /posts/9999 returns 404 for non-existent resource",
			Description: "Verifies that a non-existent post ID returns 404.",
			Expected:    "HTTP 404 Not Found",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				actions.GetAndExpectStatus(tc, c, "/posts/9999", 404)
			},
		},
		{
			Name:        "POST /posts creates a new post (201 Created)",
			Description: "Verifies that POST /posts with a valid body returns 201 with an assigned ID.",
			Expected:    "HTTP 201 Created with non-zero id",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				reqBody := &jsonPlaceholderPost{
					Title:  "E2E Template Test Post",
					Body:   "This post was created by the e2e-template framework test suite.",
					UserID: 1,
				}
				var created jsonPlaceholderPost
				actions.PostAndExpectCreated(tc, c, "/posts", reqBody, &created)
				actions.AssertNotZero(tc, "id", created.ID)
			},
		},
		{
			Name:        "GET /users returns a list of users",
			Description: "Verifies GET /users returns a non-empty list.",
			Expected:    "HTTP 200 OK with non-empty list of users",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var users []map[string]interface{}
				actions.GetAndExpectOK(tc, c, "/users", &users)
				actions.AssertNonEmptyList(tc, "users", len(users))
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		c := newPublicClient(baseURL, logDir)
		tests.RunAPITestWithDetails(t, testCase.Name, testCase.Description, testCase.Expected,
			func(tc *tests.TestContext) { testCase.RunFn(tc, c) },
		)
	}

	// Suppress unused import warning
	_ = logger.INFO
}
