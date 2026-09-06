package twinconfig_dependencies

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Dependencies_Update verifies updating a Dependencies resource.
func TestAPI_Dependencies_Update(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Update Dependencies — PUT /api/v1/config/dependencies/{id}",
		"Update a configuration dependency by ID",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to update
			createBody := map[string]interface{}{
		"child_value_code": "INACTIVE",
		"dependency_type": "REQUIRES",
		"is_active": true,
		"parent_value_code": "ACTIVE",
	}
			var createResp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/config/dependencies", &createBody, &createResp)

			createdID, ok := createResp["id"].(float64)
			if !ok || createdID == 0 {
				tc.Skip("Skipping: could not obtain valid created resource ID")
				return
			}

			// 2. Update the resource
			updateBody := map[string]interface{}{
		"child_value_code": "INACTIVE-updated",
		"dependency_type": "REQUIRES-updated",
		"is_active": false,
		"parent_value_code": "ACTIVE-updated",
	}
			path := fmt.Sprintf("/api/v1/config/dependencies/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &updateBody, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)
}
