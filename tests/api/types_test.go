package api_test

// ─────────────────────────────────────────────────────────────────────────────
// TEMPLATE: API Response Types
//
// Define the Go structs that map to your API's response payloads here.
// These types are used in the example test file and across the api_test package.
//
// Replace/extend these with your own application's response models.
// ─────────────────────────────────────────────────────────────────────────────

// HealthResponse represents a generic health-check response payload.
// Adapt the fields to match your application's /health endpoint.
type HealthResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// RootResponse represents the root welcome/info response.
type RootResponse struct {
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}

// GenericErrorResponse represents a standard API error payload.
type GenericErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// ExampleListResponse is a placeholder for paginated list responses.
// TODO: Replace with your actual list response model.
type ExampleListResponse struct {
	Items []map[string]interface{} `json:"items,omitempty"`
	Total int                      `json:"total,omitempty"`
	Page  int                      `json:"page,omitempty"`
	Limit int                      `json:"limit,omitempty"`
}
