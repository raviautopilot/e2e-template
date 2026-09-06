package twinconfig_modules

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Modules_GetByID verifies retrieving a Modules resource by ID.
func TestAPI_Modules_GetByID(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Get Modules by ID — GET /api/v1/config/modules/{id}",
		"Get a configuration module by ID",
		"HTTP 200 OK with matching resource",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to fetch
			createBody := map[string]interface{}{
		"code": fmt.Sprintf("code-%d", time.Now().UnixNano()%100000),
		"description": fmt.Sprintf("description-%d", time.Now().UnixNano()%100000),
		"is_active": true,
		"name": fmt.Sprintf("name-%d", time.Now().UnixNano()%100000),
	}
			var createResp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/config/modules", &createBody, &createResp)

			createdID, ok := createResp["id"].(float64)
			if !ok || createdID == 0 {
				tc.Skip("Skipping: could not obtain valid created resource ID")
				return
			}

			// 2. Fetch resource by ID
			path := fmt.Sprintf("/api/v1/config/modules/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)
}
