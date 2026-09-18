package httpbin_test

import (
	"testing"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_HttpBin_05_Auth verifies Bearer and Basic authentication mechanisms using dual service clients.
func TestAPI_HttpBin_05_Auth(t *testing.T) {
	tests.RunAPITestWithClients(t, "Bearer and Basic authentication endpoints with dual clients",
		"Verifies Bearer Token auth on apiClient and Basic Auth on client2.",
		"HTTP 200 OK for both Bearer and Basic authentication across separate clients",
		apiClient, client2,
		func(tc *tests.TestContext) {
			// Client 1 handles Bearer token authentication
			bearerAuth := &client.BearerTokenAuth{Token: "secret-test-token-xyz"}
			var bearerResp map[string]interface{}
			actions.AuthenticatedGet(tc, tc.Client, "/bearer", bearerAuth, &bearerResp)
			if bearerResp["authenticated"] != true {
				tc.Errorf("Expected authenticated=true for client 1, got %v", bearerResp["authenticated"])
			}

			// Client 2 handles Basic authentication independently
			basicAuth := &client.BasicAuth{Username: "user1", Password: "pass123"}
			var basicResp map[string]interface{}
			actions.AuthenticatedGet(tc, tc.Client2, "/basic-auth/user1/pass123", basicAuth, &basicResp)
			if basicResp["authenticated"] != true {
				tc.Errorf("Expected authenticated=true for client 2, got %v", basicResp["authenticated"])
			}
		},
	)
}
