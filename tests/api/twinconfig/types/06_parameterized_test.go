package twinconfig_types

import (
	"fmt"
	"testing"
	"time"

	"e2e-template/pkg/api/actions"
	"e2e-template/pkg/client"
	"e2e-template/tests"
)

// TestAPI_Types_Parameterized runs table-driven parameterized tests for Types.
func TestAPI_Types_Parameterized(t *testing.T) {
	c := client.NewClient("http://localhost:1705", 15*time.Second, tests.ExecutionLogDir)

	type testCase struct {
		name           string
		description    string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}

	testCases := []testCase{

		{
			name:           "Create Types with empty body returns 400",
			description:    "Ensure validation rejects an empty payload",
			method:         "POST",
			path:           "/api/v1/config/types",
			body:           map[string]interface{}{},
			expectedStatus: 400,
		},


		{
			name:           "Get non-existent Types ID returns 404",
			description:    "Ensure querying a non-existent ID returns 404 Not Found",
			method:         "GET",
			path:           fmt.Sprintf("/api/v1/config/types/%d", 999999),
			body:           nil,
			expectedStatus: 404,
		},


		{
			name:           "Delete non-existent Types ID returns 404",
			description:    "Ensure deleting a non-existent ID returns 404 Not Found",
			method:         "DELETE",
			path:           fmt.Sprintf("/api/v1/config/types/%d", 999999),
			body:           nil,
			expectedStatus: 404,
		},

	}

	for _, tcData := range testCases {
		tcData := tcData
		tests.RunAPITestWithDetails(t, tcData.name,
			tcData.description,
			fmt.Sprintf("HTTP %d", tcData.expectedStatus),
			func(tc *tests.TestContext) {
				switch tcData.method {
				case "GET":
					actions.GetAndExpectStatus(tc, c, tcData.path, tcData.expectedStatus)
				case "POST":
					actions.PostAndExpectStatus(tc, c, tcData.path, tcData.body, tcData.expectedStatus)
				case "DELETE":
					actions.DeleteAndExpectStatus(tc, c, tcData.path, tcData.expectedStatus)
				default:
					var dummy map[string]interface{}
					_ = c.SendHttpRequest(tcData.method, tcData.path, nil, tcData.body, &dummy, nil)
				}
			},
		)
	}
}
