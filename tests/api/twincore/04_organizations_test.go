package twincore_test

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_twincore_04_Organizations tests the Organizations CRUD lifecycle.
func TestAPI_twincore_04_Organizations(t *testing.T) {
	c := client.NewClient("http://localhost:1706", 15*time.Second, tests.ExecutionLogDir)

	// Track the created resource ID for the lifecycle
	var createdID float64

	// ── Step 1: List Organizations ──────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "List Organizations — GET /api/v1/core/organizations",
		"Get a paginated list of core organizations",
		"HTTP 200 OK with array response",
		func(tc *tests.TestContext) {
			var resp []map[string]interface{}
			actions.GetAndExpectOK(tc, c, "/api/v1/core/organizations", &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — received %d items", len(resp))
		},
	)

	// ── Step 2: Create Organizations ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create Organizations — POST /api/v1/core/organizations",
		"Create a new core organization",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"entity_id": 4,
		"industry": "test-industry",
		"is_active": true,
		"legal_name": "test-legal_name",
		"notes": "test-notes",
		"organization_type": "BANK",
		"registration_number": "test-registration_number",
		"tax_identifier": "test-tax_identifier",
		"trade_name": "test-trade_name",
		"website": "test-website",
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/core/organizations", &body, &resp)

			if id, ok := resp["id"]; ok {
				switch v := id.(type) {
				case float64:
					createdID = v
				}
			}
			tc.Actual = fmt.Sprintf("HTTP 201 Created — ID: %v", createdID)
		},
	)

	// ── Step 3: Get created Organizations by ID ──────────────────────────────

	tests.RunAPITestWithDetails(t, "Get Organizations by ID — GET /api/v1/core/organizations/{id}",
		"Get a core organization by ID — uses the ID from the create step",
		"HTTP 200 OK with the created resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/organizations/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)

	// ── Step 4: Update Organizations ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Update Organizations — PUT /api/v1/core/organizations/{id}",
		"Update a core organization by ID — modifies the created resource",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			body := map[string]interface{}{
		"entity_id": 5,
		"industry": "updated-industry",
		"is_active": false,
		"legal_name": "updated-legal_name",
		"notes": "updated-notes",
		"organization_type": "BANK-updated",
		"registration_number": "updated-registration_number",
		"tax_identifier": "updated-tax_identifier",
		"trade_name": "updated-trade_name",
		"website": "updated-website",
	}

			path := fmt.Sprintf("/api/v1/core/organizations/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &body, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)

	// ── Step 5: Delete Organizations ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Delete Organizations — DELETE /api/v1/core/organizations/{id}",
		"Soft delete a core organization by ID",
		"HTTP 204 No Content",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/organizations/%d", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)
			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted resource ID %v", createdID)
		},
	)

	// ── Step 6: Verify delete — Get deleted resource ──────────────────────────

	tests.RunAPITestWithDetails(t, "Verify delete Organizations — GET /api/v1/core/organizations/{id} after delete",
		"Verify the deleted resource returns 404 Not Found",
		"HTTP 404 Not Found for deleted resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/organizations/%d", int(createdID))
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 7: Error — Get non-existent resource ─────────────────────────────

	tests.RunAPITestWithDetails(t, "Get non-existent Organizations — GET /api/v1/core/organizations/{id} with invalid ID",
		"Verify 404 is returned for a non-existent resource ID",
		"HTTP 404 Not Found",
		func(tc *tests.TestContext) {
			path := fmt.Sprintf("/api/v1/core/organizations/%d", 999999)
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 8: Error — Create with empty body ────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create Organizations with empty body — POST /api/v1/core/organizations",
		"Verify 400 Bad Request is returned for an empty/invalid request body",
		"HTTP 400 Bad Request",
		func(tc *tests.TestContext) {
			emptyBody := map[string]interface{}{}
			actions.PostAndExpectStatus(tc, c, "/api/v1/core/organizations", &emptyBody, 400)
		},
	)

}
