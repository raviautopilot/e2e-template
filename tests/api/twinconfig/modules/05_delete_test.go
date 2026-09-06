package twinconfig_modules

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Modules_Delete verifies deleting a Modules resource.
func TestAPI_Modules_Delete(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Delete Modules — DELETE /api/v1/config/modules/{id}",
		"Soft delete a configuration module by ID",
		"HTTP 204 No Content followed by 404 Not Found",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to delete
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

			// 2. Delete the resource
			path := fmt.Sprintf("/api/v1/config/modules/%d", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)


			// 3. Verify resource is removed (GET returns 404)
			verifyPath := fmt.Sprintf("/api/v1/config/modules/%d", int(createdID))
			actions.GetAndExpectStatus(tc, c, verifyPath, 404)

			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted and verified resource ID %v", createdID)
		},
	)
}
