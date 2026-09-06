package twinconfig_dependencies

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Dependencies_Create verifies creating a new Dependencies resource.
func TestAPI_Dependencies_Create(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Create Dependencies — POST /api/v1/config/dependencies",
		"Create a new configuration dependency",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"child_value_code": "INACTIVE",
		"dependency_type": "REQUIRES",
		"is_active": true,
		"parent_value_code": "ACTIVE",
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/config/dependencies", &body, &resp)

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
