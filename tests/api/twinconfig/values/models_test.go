package twinconfig_values

// ─────────────────────────────────────────────────────────────────────────────
// Auto-generated models from Swagger definitions.
// DO NOT EDIT — regenerate with: ./generate-api-tests.sh
// ─────────────────────────────────────────────────────────────────────────────

// CfgValue represents the model.CfgValue swagger model.
type CfgValue struct {
	Code string `json:"code,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	DisplayOrder int `json:"display_order,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	TypeCode string `json:"type_code,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
	Value string `json:"value,omitempty"`
}

// ErrorResponse represents the handler.ErrorResponse swagger model.
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

