# E2E Test Execution Report{{if .Summary.SuiteName}} - {{.Summary.SuiteName}}{{end}}{{if .Summary.TargetName}} ({{.Summary.TargetName}}){{end}}

## Summary

| Metric | Value |
| :--- | :--- |
| **Suite** | {{if .Summary.SuiteName}}{{.Summary.SuiteName}}{{else}}E2E Tests{{end}} |
| **Target Scope** | {{if .Summary.TargetName}}{{.Summary.TargetName}}{{else}}All{{end}} |
| **Run ID** | `{{if .Summary.RunID}}{{.Summary.RunID}}{{else}}N/A{{end}}` |
| **Total Tests** | {{.Summary.Total}} |
| **Passed** | {{.Summary.Passed}} |
| **Failed** | {{.Summary.Failed}} |
| **Skipped** | {{.Summary.Skipped}} |
| **Start Time** | {{.Summary.StartTime.Format "2006-01-02 15:04:05"}} |
| **End Time** | {{.Summary.EndTime.Format "2006-01-02 15:04:05"}} |
| **Total Duration** | {{.Summary.TotalDurationStr}} |

---

## Detailed Test Results Categorized by Test File

{{range .GroupedResults}}
### 📄 `{{.Category}}` (Passed: {{.Passed}}, Failed: {{.Failed}}, Total: {{.Total}})

| Test Case Name & Purpose | Expected Result | Actual Result (Got) | Status | Duration | Evidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
{{range .Results}}| **{{.Name}}**<br>_{{.Description}}_ | `{{if .Expected}}{{.Expected}}{{else}}200 OK Response{{end}}` | `{{if .Actual}}{{.Actual}}{{else}}{{.Status}}{{end}}` | {{if eq .Status "passed"}}🟢 Passed{{else if eq .Status "failed"}}🔴 Failed{{else}}🟡 Skipped{{end}} | {{.DurationStr}} | {{if .Screenshots}}📸 {{len .Screenshots}} screenshots{{end}}{{if and .Screenshots .RequestLogs}}<br>{{end}}{{if .RequestLogs}}🌐 {{len .RequestLogs}} requests{{end}}{{if and (not .Screenshots) (not .RequestLogs)}}-{{end}} |
{{end}}

{{end}}

{{if gt .Summary.Failed 0}}
---

## Failed Tests & Diagnoses

{{range .GroupedResults}}{{range .Results}}{{if eq .Status "failed"}}
### ❌ [{{.Category}}] {{.Name}}

- **Description / Purpose**: {{.Description}}
- **Expected Result**: `{{.Expected}}`
- **Actual Result (Got)**: `{{.Actual}}`
- **Duration**: {{.DurationStr}}
- **Failure Reason**:
```
{{if .FailureReason}}{{.FailureReason}}{{else}}{{.Error}}{{end}}
```
{{if .Screenshots}}
- **Step Screenshots ({{len .Screenshots}})**:
{{range .Screenshots}}  - `{{.}}`
{{end}}
{{else if .Screenshot}}
- **Failure Screenshot**:
  ![Screenshot]({{.Screenshot}})
{{end}}
{{if .RequestLogs}}
- **Request / Response Logs ({{len .RequestLogs}})**:
{{range .RequestLogs}}  - `{{.}}`
{{end}}
{{end}}

---
{{end}}{{end}}{{end}}
{{end}}

