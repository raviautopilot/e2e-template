package twinconfig_types

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Types_Update verifies updating a Types resource.
func TestAPI_Types_Update(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Update Types — PUT /api/v1/config/types/{id}",
		"Update a configuration type by ID",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to update
			createBody := map[string]interface{}{
		"code": fmt.Sprintf("code-%d", time.Now().UnixNano()%100000),
		"description": "A temporary type for testing",
		"is_active": true,
		"module_code": "config",
		"name": "Test Type",
	}
			var createResp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/config/types", &createBody, &createResp)

			createdID, ok := createResp["id"].(float64)
			if !ok || createdID == 0 {
				tc.Skip("Skipping: could not obtain valid created resource ID")
				return
			}

			// 2. Update the resource
			updateBody := map[string]interface{}{
		"code": fmt.Sprintf("upd-code-%d", time.Now().UnixNano()%100000),
		"description": "A temporary type for testing-updated",
		"is_active": false,
		"module_code": "config-updated",
		"name": "Test Type-updated",
	}
			path := fmt.Sprintf("/api/v1/config/types/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &updateBody, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)
}
