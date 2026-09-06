package twinconfig_types

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Types_Create verifies creating a new Types resource.
func TestAPI_Types_Create(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Create Types — POST /api/v1/config/types",
		"Create a new configuration type",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"code": fmt.Sprintf("code-%d", time.Now().UnixNano()%100000),
		"description": "A temporary type for testing",
		"is_active": true,
		"module_code": "config",
		"name": "Test Type",
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/config/types", &body, &resp)

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
