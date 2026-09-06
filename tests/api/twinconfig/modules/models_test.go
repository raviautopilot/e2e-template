package twinconfig_modules

// ─────────────────────────────────────────────────────────────────────────────
// Auto-generated models from Swagger definitions.
// DO NOT EDIT — regenerate with: ./generate-api-tests.sh
// ─────────────────────────────────────────────────────────────────────────────

// CfgModule represents the model.CfgModule swagger model.
type CfgModule struct {
	Code string `json:"code,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	Name string `json:"name,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

// ErrorResponse represents the handler.ErrorResponse swagger model.
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

