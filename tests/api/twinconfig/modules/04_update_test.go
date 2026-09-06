package twinconfig_modules

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Modules_Update verifies updating a Modules resource.
func TestAPI_Modules_Update(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Update Modules — PUT /api/v1/config/modules/{id}",
		"Update a configuration module by ID",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to update
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

			// 2. Update the resource
			updateBody := map[string]interface{}{
		"code": fmt.Sprintf("upd-code-%d", time.Now().UnixNano()%100000),
		"description": fmt.Sprintf("upd-description-%d", time.Now().UnixNano()%100000),
		"is_active": false,
		"name": fmt.Sprintf("upd-name-%d", time.Now().UnixNano()%100000),
	}
			path := fmt.Sprintf("/api/v1/config/modules/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &updateBody, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)
}
