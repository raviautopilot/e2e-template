package jsonplaceholder_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

type jsonPlaceholderUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// TestAPI_JSONPlaceholder_03_Users verifies GET /users returns list of users.
func TestAPI_JSONPlaceholder_03_Users(t *testing.T) {
	c := client.NewClient("https://jsonplaceholder.typicode.com", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "GET /users returns list of users",
		"Verifies that GET /users returns 10 users with valid profiles.",
		"HTTP 200 OK with 10 user objects",
		func(tc *tests.TestContext) {
			var users []jsonPlaceholderUser
			actions.GetAndExpectOK(tc, c, "/users", &users)
			actions.AssertListLength(tc, "users", len(users), 10)
			actions.AssertNotEmpty(tc, "username", users[0].Username)
		},
	)
}
