package twinconfig_system

// ─────────────────────────────────────────────────────────────────────────────
// Auto-generated models from Swagger definitions.
// DO NOT EDIT — regenerate with: ./generate-api-tests.sh
// ─────────────────────────────────────────────────────────────────────────────

// ErrorResponse represents the handler.ErrorResponse swagger model.
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// HealthResponse represents the handler.HealthResponse swagger model.
type HealthResponse struct {
	RtwinStatus string `json:"rtwin_status,omitempty"`
	Status string `json:"status,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

