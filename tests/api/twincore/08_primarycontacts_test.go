package twincore_test

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_twincore_08_PrimaryContacts tests the PrimaryContacts CRUD lifecycle.
func TestAPI_twincore_08_PrimaryContacts(t *testing.T) {
	c := client.NewClient("http://localhost:1706", 15*time.Second, tests.ExecutionLogDir)

	// Track the created resource ID for the lifecycle
	var createdID float64

	// ── Step 1: List PrimaryContacts ──────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "List PrimaryContacts — GET /api/v1/core/primary-contacts",
		"Get a paginated list of primary contacts",
		"HTTP 200 OK with array response",
		func(tc *tests.TestContext) {
			var resp []map[string]interface{}
			actions.GetAndExpectOK(tc, c, "/api/v1/core/primary-contacts", &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — received %d items", len(resp))
		},
	)

	// ── Step 2: Create PrimaryContacts ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create PrimaryContacts — POST /api/v1/core/primary-contacts",
		"Create a new primary contact",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"is_active": true,
		"notes": "test-notes",
		"organization_id": 2,
		"person_id": 1,
		"role": "Manager",
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/core/primary-contacts", &body, &resp)

			if id, ok := resp["id"]; ok {
				switch v := id.(type) {
				case float64:
					createdID = v
				}
			}
			tc.Actual = fmt.Sprintf("HTTP 201 Created — ID: %v", createdID)
		},
	)

	// ── Step 3: Get created PrimaryContacts by ID ──────────────────────────────

	tests.RunAPITestWithDetails(t, "Get PrimaryContacts by ID — GET /api/v1/core/primary-contacts/{id}",
		"Get a primary contact by ID — uses the ID from the create step",
		"HTTP 200 OK with the created resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/primary-contacts/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)

	// ── Step 4: Update PrimaryContacts ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Update PrimaryContacts — PUT /api/v1/core/primary-contacts/{id}",
		"Update a primary contact by ID — modifies the created resource",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			body := map[string]interface{}{
		"is_active": false,
		"notes": "updated-notes",
		"organization_id": 3,
		"person_id": 2,
		"role": "Manager-updated",
	}

			path := fmt.Sprintf("/api/v1/core/primary-contacts/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &body, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)

	// ── Step 5: Delete PrimaryContacts ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Delete PrimaryContacts — DELETE /api/v1/core/primary-contacts/{id}",
		"Soft delete a primary contact by ID",
		"HTTP 204 No Content",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/primary-contacts/%d", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)
			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted resource ID %v", createdID)
		},
	)

	// ── Step 6: Verify delete — Get deleted resource ──────────────────────────

	tests.RunAPITestWithDetails(t, "Verify delete PrimaryContacts — GET /api/v1/core/primary-contacts/{id} after delete",
		"Verify the deleted resource returns 404 Not Found",
		"HTTP 404 Not Found for deleted resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/primary-contacts/%d", int(createdID))
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 7: Error — Get non-existent resource ─────────────────────────────

	tests.RunAPITestWithDetails(t, "Get non-existent PrimaryContacts — GET /api/v1/core/primary-contacts/{id} with invalid ID",
		"Verify 404 is returned for a non-existent resource ID",
		"HTTP 404 Not Found",
		func(tc *tests.TestContext) {
			path := fmt.Sprintf("/api/v1/core/primary-contacts/%d", 999999)
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 8: Error — Create with empty body ────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create PrimaryContacts with empty body — POST /api/v1/core/primary-contacts",
		"Verify 400 Bad Request is returned for an empty/invalid request body",
		"HTTP 400 Bad Request",
		func(tc *tests.TestContext) {
			emptyBody := map[string]interface{}{}
			actions.PostAndExpectStatus(tc, c, "/api/v1/core/primary-contacts", &emptyBody, 400)
		},
	)

}
