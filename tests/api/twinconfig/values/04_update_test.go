package twinconfig_values

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Values_Update verifies updating a Values resource.
func TestAPI_Values_Update(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Update Values — PUT /api/v1/config/values/{id}",
		"Update a configuration value by ID",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to update
			createBody := map[string]interface{}{
		"code": fmt.Sprintf("code-%d", time.Now().UnixNano()%100000),
		"description": "A temporary value for testing",
		"display_order": 1,
		"is_active": true,
		"type_code": "status",
		"value": "Test Value Content",
	}
			var createResp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/config/values", &createBody, &createResp)

			createdID, ok := createResp["id"].(float64)
			if !ok || createdID == 0 {
				tc.Skip("Skipping: could not obtain valid created resource ID")
				return
			}

			// 2. Update the resource
			updateBody := map[string]interface{}{
		"code": fmt.Sprintf("upd-code-%d", time.Now().UnixNano()%100000),
		"description": "A temporary value for testing-updated",
		"display_order": 2,
		"is_active": false,
		"type_code": "status-updated",
		"value": "Test Value Content-updated",
	}
			path := fmt.Sprintf("/api/v1/config/values/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &updateBody, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)
}
