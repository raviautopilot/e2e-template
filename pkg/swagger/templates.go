package swagger

// ── Go text/template strings for generated test files ───────────────────────

// mainTestTmpl generates the TestMain boilerplate for each test package.
var mainTestTmpl = `package {{.PackageName}}

import (
	"os"
	"testing"

	"{{.ModulePath}}/tests"
)

func TestMain(m *testing.M) {
	tests.SetupSuite()
	exitCode := m.Run()
	tests.TeardownSuite()
	os.Exit(exitCode)
}
`

// modelsTmpl generates Go struct definitions from swagger definitions for the package.
var modelsTmpl = `package {{.PackageName}}

// ─────────────────────────────────────────────────────────────────────────────
// Auto-generated models from Swagger definitions.
// DO NOT EDIT — regenerate with: ./generate-api-tests.sh
// ─────────────────────────────────────────────────────────────────────────────
{{range .Models}}
// {{.GoName}} represents the {{.SwaggerName}} swagger model.
type {{.GoName}} struct {
{{- range .Fields}}
	{{.GoName}} {{.GoType}} {{.JSONTag}}
{{- end}}
}
{{end}}
`

// listTestTmpl generates a single-test file for listing resources (01_list_test.go).
var listTestTmpl = `package {{.PackageName}}

import (
	"fmt"
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .Group.Tag}}_List verifies listing {{.Group.Tag}} resources.
func TestAPI_{{sanitize .Group.Tag}}_List(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "List {{.Group.Tag}} — GET {{.Endpoint.Path}}",
		"{{.Endpoint.Description}}",
		"HTTP {{.Endpoint.SuccessCode}} OK with array response",
		func(tc *tests.TestContext) {
			var resp []map[string]interface{}
			actions.GetAndExpectOK(tc, c, "{{.Endpoint.Path}}", &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — received %d items", len(resp))
		},
	)
}
`

// createTestTmpl generates a single-test file for creating a resource (02_create_test.go).
var createTestTmpl = `package {{.PackageName}}

import (
	"fmt"
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .Group.Tag}}_Create verifies creating a new {{.Group.Tag}} resource.
func TestAPI_{{sanitize .Group.Tag}}_Create(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Create {{.Group.Tag}} — POST {{.Endpoint.Path}}",
		"{{.Endpoint.Description}}",
		"HTTP {{.Endpoint.SuccessCode}} Created with resource in response body",
		func(tc *tests.TestContext) {
			body := {{buildExampleBody .ModelName .Definitions}}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "{{.Endpoint.Path}}", &body, &resp)

			var createdID float64
			if id, ok := resp["id"]; ok {
				if v, ok := id.(float64); ok {
					createdID = v
				}
			}
			tc.Actual = fmt.Sprintf("HTTP 201 Created — ID: %v", createdID)
		},
	)
}
`

// getTestTmpl generates a self-contained single-test file for getting a resource by ID (03_get_test.go).
var getTestTmpl = `package {{.PackageName}}

import (
	"fmt"
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .Group.Tag}}_GetByID verifies retrieving a {{.Group.Tag}} resource by ID.
func TestAPI_{{sanitize .Group.Tag}}_GetByID(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Get {{.Group.Tag}} by ID — GET {{.Endpoint.Path}}",
		"{{.Endpoint.Description}}",
		"HTTP {{.Endpoint.SuccessCode}} OK with matching resource",
		func(tc *tests.TestContext) {
{{if .Group.CreateEndpoint}}
			// 1. Create a temporary resource to fetch
			createBody := {{buildExampleBody .ModelName .Definitions}}
			var createResp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "{{.Group.CreateEndpoint.Path}}", &createBody, &createResp)

			createdID, ok := createResp["id"].(float64)
			if !ok || createdID == 0 {
				tc.Skip("Skipping: could not obtain valid created resource ID")
				return
			}
{{else}}
			createdID := float64(1)
{{end}}
			// 2. Fetch resource by ID
			path := fmt.Sprintf("{{pathParamReplace .Endpoint.Path}}", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)
}
`

// updateTestTmpl generates a self-contained single-test file for updating a resource (04_update_test.go).
var updateTestTmpl = `package {{.PackageName}}

import (
	"fmt"
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .Group.Tag}}_Update verifies updating a {{.Group.Tag}} resource.
func TestAPI_{{sanitize .Group.Tag}}_Update(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Update {{.Group.Tag}} — {{.Endpoint.Method}} {{.Endpoint.Path}}",
		"{{.Endpoint.Description}}",
		"HTTP {{.Endpoint.SuccessCode}} OK with updated resource",
		func(tc *tests.TestContext) {
{{if .Group.CreateEndpoint}}
			// 1. Create a temporary resource to update
			createBody := {{buildExampleBody .ModelName .Definitions}}
			var createResp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "{{.Group.CreateEndpoint.Path}}", &createBody, &createResp)

			createdID, ok := createResp["id"].(float64)
			if !ok || createdID == 0 {
				tc.Skip("Skipping: could not obtain valid created resource ID")
				return
			}
{{else}}
			createdID := float64(1)
{{end}}
			// 2. Update the resource
			updateBody := {{buildUpdateBody .ModelName .Definitions}}
			path := fmt.Sprintf("{{pathParamReplace .Endpoint.Path}}", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &updateBody, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)
}
`

// deleteTestTmpl generates a self-contained single-test file for deleting a resource (05_delete_test.go).
var deleteTestTmpl = `package {{.PackageName}}

import (
	"fmt"
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .Group.Tag}}_Delete verifies deleting a {{.Group.Tag}} resource.
func TestAPI_{{sanitize .Group.Tag}}_Delete(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)

	tests.RunAPITestWithDetails(t, "Delete {{.Group.Tag}} — DELETE {{.Endpoint.Path}}",
		"{{.Endpoint.Description}}",
		"HTTP {{.Endpoint.SuccessCode}} No Content followed by 404 Not Found",
		func(tc *tests.TestContext) {
{{if .Group.CreateEndpoint}}
			// 1. Create a temporary resource to delete
			createBody := {{buildExampleBody .ModelName .Definitions}}
			var createResp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "{{.Group.CreateEndpoint.Path}}", &createBody, &createResp)

			createdID, ok := createResp["id"].(float64)
			if !ok || createdID == 0 {
				tc.Skip("Skipping: could not obtain valid created resource ID")
				return
			}
{{else}}
			createdID := float64(1)
{{end}}
			// 2. Delete the resource
			path := fmt.Sprintf("{{pathParamReplace .Endpoint.Path}}", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)

{{with .Group.GetByIDEndpoint}}
			// 3. Verify resource is removed (GET returns 404)
			verifyPath := fmt.Sprintf("{{pathParamReplace .Path}}", int(createdID))
			actions.GetAndExpectStatus(tc, c, verifyPath, 404)
{{end}}
			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted and verified resource ID %v", createdID)
		},
	)
}
`

// parameterizedTestTmpl generates a table-driven parameterized test file (06_parameterized_test.go).
var parameterizedTestTmpl = `package {{.PackageName}}

import (
	"fmt"
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .Group.Tag}}_Parameterized runs table-driven parameterized tests for {{.Group.Tag}}.
func TestAPI_{{sanitize .Group.Tag}}_Parameterized(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)

	type testCase struct {
		name           string
		description    string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}

	testCases := []testCase{
{{with .Group.CreateEndpoint}}
		{
			name:           "Create {{$.Group.Tag}} with empty body returns 400",
			description:    "Ensure validation rejects an empty payload",
			method:         "POST",
			path:           "{{.Path}}",
			body:           map[string]interface{}{},
			expectedStatus: 400,
		},
{{end}}
{{with .Group.GetByIDEndpoint}}
		{
			name:           "Get non-existent {{$.Group.Tag}} ID returns 404",
			description:    "Ensure querying a non-existent ID returns 404 Not Found",
			method:         "GET",
			path:           fmt.Sprintf("{{pathParamReplace .Path}}", 999999),
			body:           nil,
			expectedStatus: 404,
		},
{{end}}
{{with .Group.DeleteEndpoint}}
		{
			name:           "Delete non-existent {{$.Group.Tag}} ID returns 404",
			description:    "Ensure deleting a non-existent ID returns 404 Not Found",
			method:         "DELETE",
			path:           fmt.Sprintf("{{pathParamReplace .Path}}", 999999),
			body:           nil,
			expectedStatus: 404,
		},
{{end}}
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
`

// healthTestTmpl generates a single-test file for health/system endpoints (01_health_test.go).
var healthTestTmpl = `package {{.PackageName}}

import (
	"fmt"
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .Group.Tag}}_HealthCheck verifies the service health endpoint.
func TestAPI_{{sanitize .Group.Tag}}_HealthCheck(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)
{{range .Group.Endpoints}}
	tests.RunAPITestWithDetails(t, "{{.Method}} {{.Path}} returns healthy status",
		"{{.Description}}",
		"{{.Summary}} — HTTP {{.SuccessCode}} OK",
		func(tc *tests.TestContext) {
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, "{{.Path}}", &resp)
			if status, ok := resp["status"]; ok {
				actions.AssertNotEmpty(tc, "status", fmt.Sprintf("%v", status))
			}
		},
	)
{{end}}
}
`
