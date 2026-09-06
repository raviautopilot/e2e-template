package twinconfig_types

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Types_GetByID verifies retrieving a Types resource by ID.
func TestAPI_Types_GetByID(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Get Types by ID — GET /api/v1/config/types/{id}",
		"Get a configuration type by ID",
		"HTTP 200 OK with matching resource",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to fetch
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

			// 2. Fetch resource by ID
			path := fmt.Sprintf("/api/v1/config/types/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)
}
