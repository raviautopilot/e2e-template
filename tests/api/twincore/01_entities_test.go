package twincore_test

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_twincore_01_Entities tests the Entities CRUD lifecycle.
func TestAPI_twincore_01_Entities(t *testing.T) {
	c := client.NewClient("http://localhost:1706", 15*time.Second, tests.ExecutionLogDir)

	// Track the created resource ID for the lifecycle
	var createdID float64

	// ── Step 1: List Entities ──────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "List Entities — GET /api/v1/core/entities",
		"Get a paginated list of core entities",
		"HTTP 200 OK with array response",
		func(tc *tests.TestContext) {
			var resp []map[string]interface{}
			actions.GetAndExpectOK(tc, c, "/api/v1/core/entities", &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — received %d items", len(resp))
		},
	)

	// ── Step 2: Create Entities ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create Entities — POST /api/v1/core/entities",
		"Create a new core entity",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"entity_type": "PERSON",
		"is_active": true,
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/core/entities", &body, &resp)

			if id, ok := resp["id"]; ok {
				switch v := id.(type) {
				case float64:
					createdID = v
				}
			}
			tc.Actual = fmt.Sprintf("HTTP 201 Created — ID: %v", createdID)
		},
	)

	// ── Step 3: Get created Entities by ID ──────────────────────────────

	tests.RunAPITestWithDetails(t, "Get Entities by ID — GET /api/v1/core/entities/{id}",
		"Get a core entity by ID — uses the ID from the create step",
		"HTTP 200 OK with the created resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/entities/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)

	// ── Step 4: Update Entities ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Update Entities — PUT /api/v1/core/entities/{id}",
		"Update a core entity by ID — modifies the created resource",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			body := map[string]interface{}{
		"entity_type": "PERSON-updated",
		"is_active": false,
	}

			path := fmt.Sprintf("/api/v1/core/entities/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &body, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)

	// ── Step 5: Delete Entities ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Delete Entities — DELETE /api/v1/core/entities/{id}",
		"Soft delete a core entity by ID",
		"HTTP 204 No Content",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/entities/%d", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)
			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted resource ID %v", createdID)
		},
	)

	// ── Step 6: Verify delete — Get deleted resource ──────────────────────────

	tests.RunAPITestWithDetails(t, "Verify delete Entities — GET /api/v1/core/entities/{id} after delete",
		"Verify the deleted resource returns 404 Not Found",
		"HTTP 404 Not Found for deleted resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/entities/%d", int(createdID))
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 7: Error — Get non-existent resource ─────────────────────────────

	tests.RunAPITestWithDetails(t, "Get non-existent Entities — GET /api/v1/core/entities/{id} with invalid ID",
		"Verify 404 is returned for a non-existent resource ID",
		"HTTP 404 Not Found",
		func(tc *tests.TestContext) {
			path := fmt.Sprintf("/api/v1/core/entities/%d", 999999)
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 8: Error — Create with empty body ────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create Entities with empty body — POST /api/v1/core/entities",
		"Verify 400 Bad Request is returned for an empty/invalid request body",
		"HTTP 400 Bad Request",
		func(tc *tests.TestContext) {
			emptyBody := map[string]interface{}{}
			actions.PostAndExpectStatus(tc, c, "/api/v1/core/entities", &emptyBody, 400)
		},
	)

}
