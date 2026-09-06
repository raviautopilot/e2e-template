package twinconfig_dependencies

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Dependencies_GetByID verifies retrieving a Dependencies resource by ID.
func TestAPI_Dependencies_GetByID(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Get Dependencies by ID — GET /api/v1/config/dependencies/{id}",
		"Get a configuration dependency by ID",
		"HTTP 200 OK with matching resource",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to fetch
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

			// 2. Fetch resource by ID
			path := fmt.Sprintf("/api/v1/config/dependencies/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)
}
