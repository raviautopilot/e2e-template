package twinconfig_values

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Values_GetByID verifies retrieving a Values resource by ID.
func TestAPI_Values_GetByID(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Get Values by ID — GET /api/v1/config/values/{id}",
		"Get a configuration value by ID",
		"HTTP 200 OK with matching resource",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to fetch
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

			// 2. Fetch resource by ID
			path := fmt.Sprintf("/api/v1/config/values/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)
}
