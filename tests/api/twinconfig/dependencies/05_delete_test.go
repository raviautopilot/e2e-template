package twinconfig_dependencies

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Dependencies_Delete verifies deleting a Dependencies resource.
func TestAPI_Dependencies_Delete(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Delete Dependencies — DELETE /api/v1/config/dependencies/{id}",
		"Soft delete a configuration dependency by ID",
		"HTTP 204 No Content followed by 404 Not Found",
		func(tc *tests.TestContext) {

			// 1. Create a temporary resource to delete
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

			// 2. Delete the resource
			path := fmt.Sprintf("/api/v1/config/dependencies/%d", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)


			// 3. Verify resource is removed (GET returns 404)
			verifyPath := fmt.Sprintf("/api/v1/config/dependencies/%d", int(createdID))
			actions.GetAndExpectStatus(tc, c, verifyPath, 404)

			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted and verified resource ID %v", createdID)
		},
	)
}
