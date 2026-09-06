package configsvc_test

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

// CfgType represents the model.CfgType swagger model.
type CfgType struct {
	Code string `json:"code,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	ModuleCode string `json:"module_code,omitempty"`
	Name string `json:"name,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

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

// HealthResponse represents the handler.HealthResponse swagger model.
type HealthResponse struct {
	RtwinStatus string `json:"rtwin_status,omitempty"`
	Status string `json:"status,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

