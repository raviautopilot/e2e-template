package myservice_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/tests"
)

// TestAPI_02_DeliberateFailure_Demo runs a single test scenario deliberately configured to fail.
// Used for testing and verifying failure evidence logging, assertion errors, and HTML report formatting.
func TestAPI_02_DeliberateFailure_Demo(t *testing.T) {
	adminEmail := tests.GlobalConfig.AdminCredentials.Username

	tests.RunAPITestWithClients(
		t,
		"Deliberate Failure Demo - Wrong Password Expected 200 OK",
		"Submits valid admin email with an incorrect password, but expects HTTP 200 OK to force an intentional failure.",
		"HTTP 200 OK with valid Auth Token",
		apiClient,
		client2,
		func(tc *tests.TestContext) {
			req := LoginRequest{
				Email:    adminEmail,
				Password: "invalid_wrong_password_deliberate_failure",
			}
			var resp LoginResponse

			// Submitting incorrect password triggers 401 Unauthorized from backend.
			// Expecting 200 OK forces this test to fail intentionally.
			actions.PostAndExpectOK(tc, apiClient, "/api/v1/auth/login", nil, &req, &resp, nil)
		},
	)
}
