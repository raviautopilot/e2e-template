package api_test

// ─────────────────────────────────────────────────────────────────────────────
// PUBLIC SITE API TESTS
//
// These tests target real, publicly-available HTTP APIs so the template works
// out-of-the-box without any local server running.
//
// APIs used:
//   - https://httpbin.org  — echo/inspect service (great for testing HTTP mechanics)
//   - https://api.github.com — GitHub public REST API (no auth needed for public data)
//
// These tests demonstrate the full framework capabilities:
//   - Table-driven tests
//   - RunAPITestWithDetails with rich report metadata
//   - Custom client with explicit base URL (overrides config.json baseUrl)
//   - Positive + negative test cases
//   - Header validation
//   - HTTP method validation
//   - Response schema validation
//
// Run with:
//   go test -v ./tests/api/... -run TestAPI_Public
// ─────────────────────────────────────────────────────────────────────────────

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/client"
	"e2e-template/pkg/logger"
	"e2e-template/tests"
)

// newPublicClient creates a Client pointed at an external base URL.
// This bypasses the config.json baseUrl so tests work without a local server.
func newPublicClient(baseURL string, logDir string) *client.Client {
	return client.NewClient(baseURL, 15*time.Second, logDir)
}

// ─────────────────────────────────────────────────────────────────────────────
// 01 · httpbin.org — HTTP mechanics testing
//
// httpbin.org is a reliable public HTTP testing service maintained by the
// Python Requests team. It echoes back request details and supports all verbs.
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

type httpbinStatusCase struct {
	Name        string
	Path        string
	Description string
	Expected    string
	ValidateFn  func(tc *tests.TestContext, c *client.Client)
}

// TestAPI_Public_01_HttpBin validates HTTP mechanics using httpbin.org.
func TestAPI_Public_01_HttpBin(t *testing.T) {
	const baseURL = "https://httpbin.org"
	logDir := tests.ExecutionLogDir

	testCases := []httpbinStatusCase{
		{
			Name:        "GET /get returns 200 OK with request echo",
			Path:        "/get",
			Description: "httpbin.org/get echoes the request back. Verifies 200 OK and URL field is non-empty.",
			Expected:    "HTTP 200 OK with non-empty url field in response",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				var resp httpbinGetResponse
				err := c.SendHttpRequest("GET", "/get", nil, nil, &resp, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("GET /get failed: %v", err)
					tc.Fatalf("GET /get failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("HTTP 200 OK, url=%q", resp.URL)
				if resp.URL == "" {
					tc.FailureReason = "Expected non-empty url in response"
					tc.Errorf("Expected non-empty url in response")
				}
			},
		},
		{
			Name:        "GET /get with custom query params are echoed back",
			Path:        "/get?foo=bar&baz=42",
			Description: "Verifies that query parameters are correctly echoed back by httpbin.",
			Expected:    "HTTP 200 OK with args.foo=bar and args.baz=42",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				var resp httpbinGetResponse
				err := c.SendHttpRequest("GET", "/get?foo=bar&baz=42", nil, nil, &resp, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("Request failed: %v", err)
					tc.Fatalf("Request failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("args=%v", resp.Args)
				if resp.Args["foo"] != "bar" {
					tc.FailureReason = fmt.Sprintf("Expected args.foo=bar, got %q", resp.Args["foo"])
					tc.Errorf("Expected args.foo=bar, got %q", resp.Args["foo"])
				}
				if resp.Args["baz"] != "42" {
					tc.FailureReason = fmt.Sprintf("Expected args.baz=42, got %q", resp.Args["baz"])
					tc.Errorf("Expected args.baz=42, got %q", resp.Args["baz"])
				}
			},
		},
		{
			Name:        "GET /status/200 returns 200 OK",
			Path:        "/status/200",
			Description: "httpbin.org/status/:code returns the requested HTTP status code.",
			Expected:    "HTTP 200 OK",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				var dummy map[string]interface{}
				err := c.SendHttpRequest("GET", "/status/200", nil, nil, &dummy, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("Expected 200 OK, got error: %v", err)
					tc.Errorf("Expected 200 OK, got error: %v", err)
				} else {
					tc.Actual = "HTTP 200 OK as expected"
				}
			},
		},
		{
			Name:        "GET /status/404 returns 404 Not Found",
			Path:        "/status/404",
			Description: "Verifies that requesting a 404 status code returns 404 Not Found.",
			Expected:    "HTTP 404 Not Found",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				var dummy map[string]interface{}
				err := c.SendHttpRequest("GET", "/status/404", nil, nil, &dummy, nil)
				if err == nil {
					tc.FailureReason = "Expected 404, got 200 OK"
					tc.Errorf("Expected 404, got 200 OK")
				} else if err.StatusCode() != 404 {
					tc.FailureReason = fmt.Sprintf("Expected 404, got %d", err.StatusCode())
					tc.Errorf("Expected 404, got %d", err.StatusCode())
				} else {
					tc.Actual = "HTTP 404 Not Found as expected"
				}
			},
		},
		{
			Name:        "GET /status/500 returns 500 Internal Server Error",
			Path:        "/status/500",
			Description: "Verifies that a 500 status is correctly returned and handled (not swallowed).",
			Expected:    "HTTP 500 Internal Server Error returned",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				var dummy map[string]interface{}
				err := c.SendHttpRequest("GET", "/status/500", nil, nil, &dummy, nil)
				if err == nil {
					tc.FailureReason = "Expected 500 error, got 200 OK"
					tc.Errorf("Expected 500 error, got 200 OK")
				} else if err.StatusCode() != 500 {
					tc.FailureReason = fmt.Sprintf("Expected 500, got %d", err.StatusCode())
					tc.Errorf("Expected 500, got %d", err.StatusCode())
				} else {
					tc.Actual = "HTTP 500 received and handled correctly (not swallowed)"
				}
			},
		},
		{
			Name:        "POST /post with JSON body is echoed back",
			Path:        "/post",
			Description: "Verifies POST with a JSON payload is accepted (200 OK) and body is echoed.",
			Expected:    "HTTP 200 OK with json.key=value echoed in response",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				reqBody := map[string]string{"key": "value", "framework": "e2e-template"}
				var resp httpbinPostResponse
				err := c.SendHttpRequest("POST", "/post", nil, &reqBody, &resp, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("POST /post failed: %v", err)
					tc.Fatalf("POST /post failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("HTTP 200 OK, json=%v", resp.JSON)
				if resp.JSON["key"] != "value" {
					tc.FailureReason = fmt.Sprintf("Expected json.key=value, got %q", resp.JSON["key"])
					tc.Errorf("Expected json.key=value, got %q", resp.JSON["key"])
				}
			},
		},
		{
			Name:        "GET /delay/1 responds within timeout",
			Path:        "/delay/1",
			Description: "Verifies a 1-second delayed response completes within the client timeout.",
			Expected:    "HTTP 200 OK within 15-second client timeout",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				var resp httpbinGetResponse
				err := c.SendHttpRequest("GET", "/delay/1", nil, nil, &resp, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("Delayed response timed out or failed: %v", err)
					tc.Errorf("Delayed response timed out or failed: %v", err)
				} else {
					tc.Actual = "HTTP 200 OK received within timeout"
				}
			},
		},
		{
			Name:        "Bearer Token Auth via /bearer endpoint",
			Path:        "/bearer",
			Description: "Verifies that a valid Bearer token is accepted by httpbin's /bearer endpoint.",
			Expected:    "HTTP 200 OK when Authorization: Bearer <token> header is sent",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				auth := &client.BearerTokenAuth{Token: "test-framework-token"}
				var resp map[string]interface{}
				err := c.SendHttpRequest("GET", "/bearer", nil, nil, &resp, auth)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("Expected 200 OK with bearer token, got: %v", err)
					tc.Errorf("Expected 200 OK with bearer token, got: %v", err)
				} else {
					tc.Actual = fmt.Sprintf("HTTP 200 OK, authenticated=%v", resp["authenticated"])
				}
			},
		},
		{
			Name:        "Basic Auth via /basic-auth endpoint",
			Path:        "/basic-auth/user/pass",
			Description: "Verifies that correct Basic Auth credentials are accepted by httpbin.",
			Expected:    "HTTP 200 OK with correct Basic Auth user=user, password=pass",
			ValidateFn: func(tc *tests.TestContext, c *client.Client) {
				auth := &client.BasicAuth{Username: "user", Password: "pass"}
				var resp map[string]interface{}
				err := c.SendHttpRequest("GET", "/basic-auth/user/pass", nil, nil, &resp, auth)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("Basic auth failed: %v", err)
					tc.Errorf("Basic auth failed: %v", err)
				} else {
					tc.Actual = fmt.Sprintf("HTTP 200 OK, authenticated=%v, user=%v", resp["authenticated"], resp["user"])
				}
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		c := newPublicClient(baseURL, logDir)
		tests.RunAPITestWithDetails(
			t,
			testCase.Name,
			testCase.Description,
			testCase.Expected,
			func(tc *tests.TestContext) {
				testCase.ValidateFn(tc, c)
			},
		)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 02 · GitHub Public REST API
//
// GitHub's public API requires no authentication for read-only public data.
// Good for testing real-world JSON schema validation and pagination patterns.
// ─────────────────────────────────────────────────────────────────────────────

type githubUser struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
}

type githubRepo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

// TestAPI_Public_02_GitHub validates GitHub's public REST API endpoints.
func TestAPI_Public_02_GitHub(t *testing.T) {
	const baseURL = "https://api.github.com"
	logDir := tests.ExecutionLogDir

	type githubTestCase struct {
		Name        string
		Description string
		Expected    string
		RunFn       func(tc *tests.TestContext, c *client.Client)
	}

	testCases := []githubTestCase{
		{
			Name:        "GET /users/octocat returns valid GitHub user",
			Description: "Verifies the GitHub API returns octocat's public profile with login=octocat.",
			Expected:    "HTTP 200 OK with login=octocat, id>0, type=User",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				headers := map[string]string{
					"Accept":               "application/vnd.github.v3+json",
					"X-GitHub-Api-Version": "2022-11-28",
				}
				var user githubUser
				err := c.SendHttpRequest("GET", "/users/octocat", headers, nil, &user, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("GET /users/octocat failed: %v", err)
					tc.Fatalf("GET /users/octocat failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("login=%q, id=%d, type=%q", user.Login, user.ID, user.Type)
				if user.Login != "octocat" {
					tc.FailureReason = fmt.Sprintf("Expected login=octocat, got %q", user.Login)
					tc.Errorf("Expected login=octocat, got %q", user.Login)
				}
				if user.ID == 0 {
					tc.FailureReason = "Expected non-zero ID"
					tc.Errorf("Expected non-zero ID")
				}
			},
		},
		{
			Name:        "GET /users/nonexistent-user-xyz-abc-123 returns 404",
			Description: "Verifies that fetching a non-existent GitHub user returns 404 Not Found.",
			Expected:    "HTTP 404 Not Found",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				headers := map[string]string{"Accept": "application/vnd.github.v3+json"}
				var dummy map[string]interface{}
				err := c.SendHttpRequest("GET", "/users/nonexistent-user-xyz-abc-123", headers, nil, &dummy, nil)
				if err == nil {
					tc.FailureReason = "Expected 404, got 200 OK"
					tc.Errorf("Expected 404, got 200 OK")
				} else if err.StatusCode() != 404 {
					tc.FailureReason = fmt.Sprintf("Expected 404, got %d", err.StatusCode())
					tc.Errorf("Expected 404, got %d", err.StatusCode())
				} else {
					tc.Actual = "HTTP 404 Not Found as expected"
				}
			},
		},
		{
			Name:        "GET /repos/octocat/Hello-World returns valid public repo",
			Description: "Verifies the GitHub API returns the Hello-World public repo with full_name.",
			Expected:    "HTTP 200 OK with full_name=octocat/Hello-World",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				headers := map[string]string{"Accept": "application/vnd.github.v3+json"}
				var repo githubRepo
				err := c.SendHttpRequest("GET", "/repos/octocat/Hello-World", headers, nil, &repo, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("GET /repos/octocat/Hello-World failed: %v", err)
					tc.Fatalf("GET /repos/octocat/Hello-World failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("full_name=%q, id=%d", repo.FullName, repo.ID)
				if repo.FullName != "octocat/Hello-World" {
					tc.FailureReason = fmt.Sprintf("Expected full_name=octocat/Hello-World, got %q", repo.FullName)
					tc.Errorf("Expected full_name=octocat/Hello-World, got %q", repo.FullName)
				}
			},
		},
		{
			Name:        "GET /repos/nonexistent-org-xyz/nonexistent-repo-abc returns 404",
			Description: "Verifies that a non-existent repository returns 404 Not Found.",
			Expected:    "HTTP 404 Not Found",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				headers := map[string]string{"Accept": "application/vnd.github.v3+json"}
				var dummy map[string]interface{}
				err := c.SendHttpRequest("GET", "/repos/nonexistent-org-xyz/nonexistent-repo-abc", headers, nil, &dummy, nil)
				if err == nil {
					tc.FailureReason = "Expected 404, got 200 OK"
					tc.Errorf("Expected 404, got 200 OK")
				} else if err.StatusCode() != 404 {
					tc.FailureReason = fmt.Sprintf("Expected 404, got %d", err.StatusCode())
					tc.Errorf("Expected 404, got %d", err.StatusCode())
				} else {
					tc.Actual = "HTTP 404 Not Found as expected"
				}
			},
		},
		{
			Name:        "GET /rate_limit shows current rate limit status",
			Description: "Verifies the rate_limit endpoint returns valid rate limit information (no auth needed).",
			Expected:    "HTTP 200 OK with rate.limit and rate.remaining fields",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				headers := map[string]string{"Accept": "application/vnd.github.v3+json"}
				var resp map[string]interface{}
				err := c.SendHttpRequest("GET", "/rate_limit", headers, nil, &resp, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("GET /rate_limit failed: %v", err)
					tc.Fatalf("GET /rate_limit failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("HTTP 200 OK, rate_limit response received")
				if resp["rate"] == nil && resp["resources"] == nil {
					tc.FailureReason = "Expected rate or resources field in rate_limit response"
					tc.Errorf("Expected rate or resources field in rate_limit response")
				}
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		c := newPublicClient(baseURL, logDir)
		tests.RunAPITestWithDetails(
			t,
			testCase.Name,
			testCase.Description,
			testCase.Expected,
			func(tc *tests.TestContext) {
				testCase.RunFn(tc, c)
			},
		)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// 03 · JSONPlaceholder — Fake REST API (typicode.com)
//
// JSONPlaceholder is a free online REST API for quick prototyping and testing.
// Great for demonstrating CRUD test patterns.
// ─────────────────────────────────────────────────────────────────────────────

type jsonPlaceholderPost struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

// TestAPI_Public_03_JSONPlaceholder validates CRUD patterns against jsonplaceholder.typicode.com.
func TestAPI_Public_03_JSONPlaceholder(t *testing.T) {
	const baseURL = "https://jsonplaceholder.typicode.com"
	logDir := tests.ExecutionLogDir

	type placeholderCase struct {
		Name        string
		Description string
		Expected    string
		RunFn       func(tc *tests.TestContext, c *client.Client)
	}

	testCases := []placeholderCase{
		{
			Name:        "GET /posts returns a list of 100 posts",
			Description: "Verifies GET /posts returns all 100 placeholder posts.",
			Expected:    "HTTP 200 OK with exactly 100 posts",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var posts []jsonPlaceholderPost
				err := c.SendHttpRequest("GET", "/posts", nil, nil, &posts, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("GET /posts failed: %v", err)
					tc.Fatalf("GET /posts failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("HTTP 200 OK, count=%d", len(posts))
				if len(posts) != 100 {
					tc.FailureReason = fmt.Sprintf("Expected 100 posts, got %d", len(posts))
					tc.Errorf("Expected 100 posts, got %d", len(posts))
				}
			},
		},
		{
			Name:        "GET /posts/1 returns a single post with valid fields",
			Description: "Verifies GET /posts/1 returns a post with non-empty title and body.",
			Expected:    "HTTP 200 OK with id=1, non-empty title and body",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var post jsonPlaceholderPost
				err := c.SendHttpRequest("GET", "/posts/1", nil, nil, &post, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("GET /posts/1 failed: %v", err)
					tc.Fatalf("GET /posts/1 failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("id=%d, title=%q", post.ID, post.Title)
				if post.ID != 1 {
					tc.FailureReason = fmt.Sprintf("Expected id=1, got %d", post.ID)
					tc.Errorf("Expected id=1, got %d", post.ID)
				}
				if post.Title == "" {
					tc.FailureReason = "Expected non-empty title"
					tc.Errorf("Expected non-empty title")
				}
			},
		},
		{
			Name:        "GET /posts/9999 returns 404 for non-existent resource",
			Description: "Verifies that a non-existent post ID returns 404 Not Found.",
			Expected:    "HTTP 404 Not Found",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var dummy map[string]interface{}
				err := c.SendHttpRequest("GET", "/posts/9999", nil, nil, &dummy, nil)
				if err == nil {
					tc.FailureReason = "Expected 404, got 200 OK"
					tc.Errorf("Expected 404, got 200 OK")
				} else if err.StatusCode() != 404 {
					tc.FailureReason = fmt.Sprintf("Expected 404, got %d", err.StatusCode())
					tc.Errorf("Expected 404, got %d", err.StatusCode())
				} else {
					tc.Actual = "HTTP 404 Not Found as expected"
				}
			},
		},
		{
			Name:        "POST /posts creates a new post (201 Created)",
			Description: "Verifies that POST /posts with a valid body returns 201 Created with an assigned ID.",
			Expected:    "HTTP 201 Created with non-zero id in response",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				reqBody := &jsonPlaceholderPost{
					Title:  "E2E Template Test Post",
					Body:   "This post was created by the e2e-template framework test suite.",
					UserID: 1,
				}
				var created jsonPlaceholderPost
				err := c.SendHttpRequest("POST", "/posts", nil, reqBody, &created, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("POST /posts failed: %v", err)
					tc.Fatalf("POST /posts failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("id=%d, title=%q", created.ID, created.Title)
				if created.ID == 0 {
					tc.FailureReason = "Expected non-zero id in created post response"
					tc.Errorf("Expected non-zero id in created post response")
				}
			},
		},
		{
			Name:        "GET /users returns a list of users",
			Description: "Verifies GET /users returns a non-empty list of users.",
			Expected:    "HTTP 200 OK with non-empty list of users",
			RunFn: func(tc *tests.TestContext, c *client.Client) {
				var users []map[string]interface{}
				err := c.SendHttpRequest("GET", "/users", nil, nil, &users, nil)
				if err != nil {
					tc.FailureReason = fmt.Sprintf("GET /users failed: %v", err)
					tc.Fatalf("GET /users failed: %v", err)
				}
				tc.Actual = fmt.Sprintf("HTTP 200 OK, count=%d", len(users))
				if len(users) == 0 {
					tc.FailureReason = "Expected non-empty users list"
					tc.Errorf("Expected non-empty users list")
				}
			},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		c := newPublicClient(baseURL, logDir)
		tests.RunAPITestWithDetails(
			t,
			testCase.Name,
			testCase.Description,
			testCase.Expected,
			func(tc *tests.TestContext) {
				testCase.RunFn(tc, c)
			},
		)
	}

	// Suppress unused import warning
	_ = logger.INFO
}
