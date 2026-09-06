package example_test

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
type ExampleListResponse struct {
	Items []map[string]interface{} `json:"items,omitempty"`
	Total int                      `json:"total,omitempty"`
	Page  int                      `json:"page,omitempty"`
	Limit int                      `json:"limit,omitempty"`
}
