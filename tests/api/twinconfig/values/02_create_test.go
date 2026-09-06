package twinconfig_values

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Values_Create verifies creating a new Values resource.
func TestAPI_Values_Create(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Create Values — POST /api/v1/config/values",
		"Create a new configuration value",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"code": fmt.Sprintf("code-%d", time.Now().UnixNano()%100000),
		"description": "A temporary value for testing",
		"display_order": 1,
		"is_active": true,
		"type_code": "status",
		"value": "Test Value Content",
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/config/values", &body, &resp)

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
