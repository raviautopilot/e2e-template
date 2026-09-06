package twincore_test

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_twincore_02_OrganizationAddresses tests the OrganizationAddresses CRUD lifecycle.
func TestAPI_twincore_02_OrganizationAddresses(t *testing.T) {
	c := client.NewClient("http://localhost:1706", 15*time.Second, tests.ExecutionLogDir)

	// Track the created resource ID for the lifecycle
	var createdID float64

	// ── Step 1: List OrganizationAddresses ──────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "List OrganizationAddresses — GET /api/v1/core/organization-addresses",
		"Get a paginated list of organization addresses",
		"HTTP 200 OK with array response",
		func(tc *tests.TestContext) {
			var resp []map[string]interface{}
			actions.GetAndExpectOK(tc, c, "/api/v1/core/organization-addresses", &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — received %d items", len(resp))
		},
	)

	// ── Step 2: Create OrganizationAddresses ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create OrganizationAddresses — POST /api/v1/core/organization-addresses",
		"Create a new organization address",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := map[string]interface{}{
		"address_line1": "test-address_line1",
		"address_line2": "test-address_line2",
		"address_type": "WORK",
		"city": "test-city",
		"country": "USA",
		"is_active": true,
		"is_primary": true,
		"organization_id": 2,
		"postal_code": "test-postal_code",
		"state": "test-state",
	}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "/api/v1/core/organization-addresses", &body, &resp)

			if id, ok := resp["id"]; ok {
				switch v := id.(type) {
				case float64:
					createdID = v
				}
			}
			tc.Actual = fmt.Sprintf("HTTP 201 Created — ID: %v", createdID)
		},
	)

	// ── Step 3: Get created OrganizationAddresses by ID ──────────────────────────────

	tests.RunAPITestWithDetails(t, "Get OrganizationAddresses by ID — GET /api/v1/core/organization-addresses/{id}",
		"Get an organization address by ID — uses the ID from the create step",
		"HTTP 200 OK with the created resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/organization-addresses/%d", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)

	// ── Step 4: Update OrganizationAddresses ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Update OrganizationAddresses — PUT /api/v1/core/organization-addresses/{id}",
		"Update an organization address by ID — modifies the created resource",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			body := map[string]interface{}{
		"address_line1": "updated-address_line1",
		"address_line2": "updated-address_line2",
		"address_type": "WORK-updated",
		"city": "updated-city",
		"country": "USA-updated",
		"is_active": false,
		"is_primary": false,
		"organization_id": 3,
		"postal_code": "updated-postal_code",
		"state": "updated-state",
	}

			path := fmt.Sprintf("/api/v1/core/organization-addresses/%d", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &body, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)

	// ── Step 5: Delete OrganizationAddresses ─────────────────────────────────────────

	tests.RunAPITestWithDetails(t, "Delete OrganizationAddresses — DELETE /api/v1/core/organization-addresses/{id}",
		"Soft delete an organization address by ID",
		"HTTP 204 No Content",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/organization-addresses/%d", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)
			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted resource ID %v", createdID)
		},
	)

	// ── Step 6: Verify delete — Get deleted resource ──────────────────────────

	tests.RunAPITestWithDetails(t, "Verify delete OrganizationAddresses — GET /api/v1/core/organization-addresses/{id} after delete",
		"Verify the deleted resource returns 404 Not Found",
		"HTTP 404 Not Found for deleted resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("/api/v1/core/organization-addresses/%d", int(createdID))
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 7: Error — Get non-existent resource ─────────────────────────────

	tests.RunAPITestWithDetails(t, "Get non-existent OrganizationAddresses — GET /api/v1/core/organization-addresses/{id} with invalid ID",
		"Verify 404 is returned for a non-existent resource ID",
		"HTTP 404 Not Found",
		func(tc *tests.TestContext) {
			path := fmt.Sprintf("/api/v1/core/organization-addresses/%d", 999999)
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)

	// ── Step 8: Error — Create with empty body ────────────────────────────────

	tests.RunAPITestWithDetails(t, "Create OrganizationAddresses with empty body — POST /api/v1/core/organization-addresses",
		"Verify 400 Bad Request is returned for an empty/invalid request body",
		"HTTP 400 Bad Request",
		func(tc *tests.TestContext) {
			emptyBody := map[string]interface{}{}
			actions.PostAndExpectStatus(tc, c, "/api/v1/core/organization-addresses", &emptyBody, 400)
		},
	)

}
