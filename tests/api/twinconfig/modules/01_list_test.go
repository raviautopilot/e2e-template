package twinconfig_modules

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Modules_List verifies listing Modules resources.
func TestAPI_Modules_List(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "List Modules — GET /api/v1/config/modules",
		"Get a paginated list of configuration modules",
		"HTTP 200 OK with array response",
		func(tc *tests.TestContext) {
			var resp []map[string]interface{}
			actions.GetAndExpectOK(tc, c, "/api/v1/config/modules", &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — received %d items", len(resp))
		},
	)
}
