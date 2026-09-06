package swagger

// ── Go text/template strings for generated test files ───────────────────────

// mainTestTmpl generates the TestMain boilerplate for the test package.
var mainTestTmpl = `package {{.ServiceName}}_test

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

// modelsTmpl generates Go struct definitions from swagger definitions.
var modelsTmpl = `package {{.ServiceName}}_test

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

// healthTestTmpl generates a simple health check test.
var healthTestTmpl = `package {{.ServiceName}}_test

import (
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .Group.Tag}}_{{pad .SeqNum}}_HealthCheck verifies the service health endpoint.
func TestAPI_{{sanitize .Group.Tag}}_{{pad .SeqNum}}_HealthCheck(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)
{{range .Group.Endpoints}}
	tests.RunAPITestWithDetails(t, "{{.Method}} {{.Path}} returns healthy status",
		"{{.Description}}",
		"{{.Summary}} — HTTP {{.SuccessCode}} OK",
		func(tc *tests.TestContext) {
			var resp HealthResponse
			actions.GetAndExpectOK(tc, c, "{{.Path}}", &resp)
			actions.AssertNotEmpty(tc, "status", resp.Status)
		},
	)
{{end}}
}
`

// crudTestTmpl generates a full CRUD lifecycle test for a resource group.
var crudTestTmpl = `package {{.ServiceName}}_test

import (
	"fmt"
	"testing"
	"time"

	"{{.ModulePath}}/pkg/api/actions"
	"{{.ModulePath}}/pkg/client"
	"{{.ModulePath}}/tests"
)

// TestAPI_{{sanitize .ServiceName}}_{{pad .SeqNum}}_{{sanitize .Group.Tag}} tests the {{.Group.Tag}} CRUD lifecycle.
func TestAPI_{{sanitize .ServiceName}}_{{pad .SeqNum}}_{{sanitize .Group.Tag}}(t *testing.T) {
	c := client.NewClient("{{.BaseURL}}", 15*time.Second, tests.ExecutionLogDir)

	// Track the created resource ID for the lifecycle
	var createdID float64
{{$group := .Group}}{{$modelName := .ModelName}}{{$goModelName := .GoModelName}}{{$defs := .Definitions}}{{$svcName := .ServiceName}}{{$seqNum := .SeqNum}}
	// ── Step 1: List {{$group.Tag}} ──────────────────────────────────────────
{{with $group.ListEndpoint}}
	tests.RunAPITestWithDetails(t, "List {{$group.Tag}} — GET {{.Path}}",
		"{{.Description}}",
		"HTTP 200 OK with array response",
		func(tc *tests.TestContext) {
			var resp []map[string]interface{}
			actions.GetAndExpectOK(tc, c, "{{.Path}}", &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — received %d items", len(resp))
		},
	)
{{end}}
	// ── Step 2: Create {{$group.Tag}} ─────────────────────────────────────────
{{with $group.CreateEndpoint}}
	tests.RunAPITestWithDetails(t, "Create {{$group.Tag}} — POST {{.Path}}",
		"{{.Description}}",
		"HTTP 201 Created with resource in response body",
		func(tc *tests.TestContext) {
			body := {{buildExampleBody $modelName $defs}}

			var resp map[string]interface{}
			actions.PostAndExpectCreated(tc, c, "{{.Path}}", &body, &resp)

			if id, ok := resp["id"]; ok {
				switch v := id.(type) {
				case float64:
					createdID = v
				}
			}
			tc.Actual = fmt.Sprintf("HTTP 201 Created — ID: %v", createdID)
		},
	)
{{end}}
	// ── Step 3: Get created {{$group.Tag}} by ID ──────────────────────────────
{{with $group.GetByIDEndpoint}}
	tests.RunAPITestWithDetails(t, "Get {{$group.Tag}} by ID — GET {{.Path}}",
		"{{.Description}} — uses the ID from the create step",
		"HTTP 200 OK with the created resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("{{pathParamReplace .Path}}", int(createdID))
			var resp map[string]interface{}
			actions.GetAndExpectOK(tc, c, path, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — retrieved resource ID %v", createdID)
		},
	)
{{end}}
	// ── Step 4: Update {{$group.Tag}} ─────────────────────────────────────────
{{with $group.UpdateEndpoint}}
	tests.RunAPITestWithDetails(t, "Update {{$group.Tag}} — {{.Method}} {{.Path}}",
		"{{.Description}} — modifies the created resource",
		"HTTP 200 OK with updated resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			body := {{buildUpdateBody $modelName $defs}}

			path := fmt.Sprintf("{{pathParamReplace .Path}}", int(createdID))
			var resp map[string]interface{}
			actions.PutAndExpectOK(tc, c, path, &body, &resp)
			tc.Actual = fmt.Sprintf("HTTP 200 OK — updated resource ID %v", createdID)
		},
	)
{{end}}
	// ── Step 5: Delete {{$group.Tag}} ─────────────────────────────────────────
{{with $group.DeleteEndpoint}}
	tests.RunAPITestWithDetails(t, "Delete {{$group.Tag}} — DELETE {{.Path}}",
		"{{.Description}}",
		"HTTP 204 No Content",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("{{pathParamReplace .Path}}", int(createdID))
			actions.DeleteAndExpectNoContent(tc, c, path)
			tc.Actual = fmt.Sprintf("HTTP 204 No Content — deleted resource ID %v", createdID)
		},
	)
{{end}}
	// ── Step 6: Verify delete — Get deleted resource ──────────────────────────
{{with $group.GetByIDEndpoint}}
	tests.RunAPITestWithDetails(t, "Verify delete {{$group.Tag}} — GET {{.Path}} after delete",
		"Verify the deleted resource returns 404 Not Found",
		"HTTP 404 Not Found for deleted resource",
		func(tc *tests.TestContext) {
			if createdID == 0 {
				tc.Skip("Skipping: no resource was created in previous step")
				return
			}
			path := fmt.Sprintf("{{pathParamReplace .Path}}", int(createdID))
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)
{{end}}
	// ── Step 7: Error — Get non-existent resource ─────────────────────────────
{{with $group.GetByIDEndpoint}}
	tests.RunAPITestWithDetails(t, "Get non-existent {{$group.Tag}} — GET {{.Path}} with invalid ID",
		"Verify 404 is returned for a non-existent resource ID",
		"HTTP 404 Not Found",
		func(tc *tests.TestContext) {
			path := fmt.Sprintf("{{pathParamReplace .Path}}", 999999)
			actions.GetAndExpectStatus(tc, c, path, 404)
		},
	)
{{end}}
	// ── Step 8: Error — Create with empty body ────────────────────────────────
{{with $group.CreateEndpoint}}
	tests.RunAPITestWithDetails(t, "Create {{$group.Tag}} with empty body — POST {{.Path}}",
		"Verify 400 Bad Request is returned for an empty/invalid request body",
		"HTTP 400 Bad Request",
		func(tc *tests.TestContext) {
			emptyBody := map[string]interface{}{}
			actions.PostAndExpectStatus(tc, c, "{{.Path}}", &emptyBody, 400)
		},
	)
{{end}}
}
`
