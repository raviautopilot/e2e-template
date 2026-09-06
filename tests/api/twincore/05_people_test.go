package twincore_test

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_twincore_05_People tests the People CRUD lifecycle.
func TestAPI_twincore_05_People(t *testing.T) {
	c := client.NewClient("http://localhost:1706", 15*time.Second, tests.ExecutionLogDir)

	// Track the created resource ID for the lifecycle
	var createdID float64

	// ── Step 1: List People ──────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "List People — GET /api/v1/core/people",
		"Get a paginated list of core people",
		"HTTP 200 OK with array response",
		func(tc *tests.TestContext) {
			var resp []map[string]interface{}
			actions.GetAndExpectOK(tc, c, "/api/v1/core/people", &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — received %d items", len(resp))
		},
	)

	// ── Step 2: Create People ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create People — POST /api/v1/core/people",
		"Create a new core person",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"date_of_birth": "1990-01-01T00:00:00Z",
		"entity_id": 3,
		"first_name": "test-first_name",
		"gender": "MALE",
		"is_active": true,
		"last_name": "test-last_name",
		"middle_name": "test-middle_name",
		"national_id": "test-national_id",
		"notes": "test-notes",
		"preferred_name": "test-preferred_name",
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/core/people", &body, &resp)

			if id, ok := resp["id"]; ok {
				switch v := id.(type) {
				case float64:
					createdID = v
				}
			}
			tc.Actual = fmt.Sprintf("HTTP 201 Created — ID: %v", createdID)
		},
	)

	// ── Step 3: Get created People by ID ──────────────────────────────

	tests.RunAPITestWithDetails(t, "Get People by ID — GET /api/v1/core/people/{id}",
		"Get a core person by ID — uses the ID from the create step",
		"HTTP 200 OK with the created resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/people/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)

	// ── Step 4: Update People ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Update People — PUT /api/v1/core/people/{id}",
		"Update a core person by ID — modifies the created resource",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			body := map[string]interface{}{
		"date_of_birth": "1990-01-01T00:00:00Z-updated",
		"entity_id": 4,
		"first_name": "updated-first_name",
		"gender": "MALE-updated",
		"is_active": false,
		"last_name": "updated-last_name",
		"middle_name": "updated-middle_name",
		"national_id": "updated-national_id",
		"notes": "updated-notes",
		"preferred_name": "updated-preferred_name",
	}

			path := fmt.Sprintf("/api/v1/core/people/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &body, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)

	// ── Step 5: Delete People ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Delete People — DELETE /api/v1/core/people/{id}",
		"Soft delete a core person by ID",
		"HTTP 204 No Content",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/people/%d", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)
			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted resource ID %v", createdID)
		},
	)

	// ── Step 6: Verify delete — Get deleted resource ──────────────────────────

	tests.RunAPITestWithDetails(t, "Verify delete People — GET /api/v1/core/people/{id} after delete",
		"Verify the deleted resource returns 404 Not Found",
		"HTTP 404 Not Found for deleted resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/people/%d", int(createdID))
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 7: Error — Get non-existent resource ─────────────────────────────

	tests.RunAPITestWithDetails(t, "Get non-existent People — GET /api/v1/core/people/{id} with invalid ID",
		"Verify 404 is returned for a non-existent resource ID",
		"HTTP 404 Not Found",
		func(tc *tests.TestContext) {
			path := fmt.Sprintf("/api/v1/core/people/%d", 999999)
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 8: Error — Create with empty body ────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create People with empty body — POST /api/v1/core/people",
		"Verify 400 Bad Request is returned for an empty/invalid request body",
		"HTTP 400 Bad Request",
		func(tc *tests.TestContext) {
			emptyBody := map[string]interface{}{}
			actions.PostAndExpectStatus(tc, c, "/api/v1/core/people", &emptyBody, 400)
		},
	)

}
