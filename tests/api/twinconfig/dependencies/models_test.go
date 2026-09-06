package twinconfig_dependencies

// ─────────────────────────────────────────────────────────────────────────────
// Auto-generated models from Swagger definitions.
// DO NOT EDIT — regenerate with: ./generate-api-tests.sh
// ─────────────────────────────────────────────────────────────────────────────

// CfgDependency represents the model.CfgDependency swagger model.
type CfgDependency struct {
	ChildValueCode string `json:"child_value_code,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	DependencyType string `json:"dependency_type,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	ParentValueCode string `json:"parent_value_code,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

// ErrorResponse represents the handler.ErrorResponse swagger model.
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

