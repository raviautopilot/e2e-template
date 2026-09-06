package twinconfig_modules

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Modules_Create verifies creating a new Modules resource.
func TestAPI_Modules_Create(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Create Modules — POST /api/v1/config/modules",
		"Create a new configuration module",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"code": fmt.Sprintf("code-%d", time.Now().UnixNano()%100000),
		"description": fmt.Sprintf("description-%d", time.Now().UnixNano()%100000),
		"is_active": true,
		"name": fmt.Sprintf("name-%d", time.Now().UnixNano()%100000),
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/config/modules", &body, &resp)

			var createdID float64
			if id, ok := resp["id"]; ok {
				if v, ok := id.(float64); ok {
					createdID = v
				}
			}
			tc.Actual = fmt.Sprintf("HTTP 201 Created — ID: %v", createdID)
		},
	)
}
