package httpbin_test

import (
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_HttpBin_05_Auth verifies Bearer and Basic authentication mechanisms.
func TestAPI_HttpBin_05_Auth(t *testing.T) {
	c := client.NewClient("https://httpbin.org", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Bearer and Basic authentication endpoints",
		"Verifies Bearer Token auth against /bearer and Basic Auth against /basic-auth.",
		"HTTP 200 OK for both Bearer and Basic authentication",
		func(tc *tests.TestContext) {
			bearerAuth := &client.BearerTokenAuth{Token: "secret-test-token-xyz"}
			var bearerResp map[string]interface{}
			actions.AuthenticatedGet(tc, c, "/bearer", bearerAuth, &bearerResp)
			if bearerResp["authenticated"] != true {
				tc.Errorf("Expected authenticated=true, got %v", bearerResp["authenticated"])
			}

			basicAuth := &client.BasicAuth{Username: "user1", Password: "pass123"}
			var basicResp map[string]interface{}
			actions.AuthenticatedGet(tc, c, "/basic-auth/user1/pass123", basicAuth, &basicResp)
			if basicResp["authenticated"] != true {
				tc.Errorf("Expected authenticated=true, got %v", basicResp["authenticated"])
			}
		},
	)
}
