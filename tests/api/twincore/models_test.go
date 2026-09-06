package twincore_test

// ─────────────────────────────────────────────────────────────────────────────
// Auto-generated models from Swagger definitions.
// DO NOT EDIT — regenerate with: ./generate-api-tests.sh
// ─────────────────────────────────────────────────────────────────────────────

// CoreEntity represents the model.CoreEntity swagger model.
type CoreEntity struct {
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

// CoreOrganization represents the model.CoreOrganization swagger model.
type CoreOrganization struct {
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	EntityId int `json:"entity_id,omitempty"`
	Id int `json:"id,omitempty"`
	Industry string `json:"industry,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	LegalName string `json:"legal_name,omitempty"`
	Notes string `json:"notes,omitempty"`
	OrganizationType string `json:"organization_type,omitempty"`
	RegistrationNumber string `json:"registration_number,omitempty"`
	TaxIdentifier string `json:"tax_identifier,omitempty"`
	TradeName string `json:"trade_name,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
	Website string `json:"website,omitempty"`
}

// CorePerson represents the model.CorePerson swagger model.
type CorePerson struct {
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	EntityId int `json:"entity_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Gender string `json:"gender,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	LastName string `json:"last_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	NationalId string `json:"national_id,omitempty"`
	Notes string `json:"notes,omitempty"`
	PreferredName string `json:"preferred_name,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

// CoreRelationship represents the model.CoreRelationship swagger model.
type CoreRelationship struct {
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	EndDate string `json:"end_date,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	RelationshipType string `json:"relationship_type,omitempty"`
	SourceEntityId int `json:"source_entity_id,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	TargetEntityId int `json:"target_entity_id,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
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

// OrganizationAddress represents the model.OrganizationAddress swagger model.
type OrganizationAddress struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	AddressType string `json:"address_type,omitempty"`
	City string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	IsPrimary bool `json:"is_primary,omitempty"`
	OrganizationId int `json:"organization_id,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	State string `json:"state,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

// OrganizationContact represents the model.OrganizationContact swagger model.
type OrganizationContact struct {
	ContactType string `json:"contact_type,omitempty"`
	ContactValue string `json:"contact_value,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	IsPrimary bool `json:"is_primary,omitempty"`
	Notes string `json:"notes,omitempty"`
	OrganizationId int `json:"organization_id,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

// PersonAddress represents the model.PersonAddress swagger model.
type PersonAddress struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	AddressType string `json:"address_type,omitempty"`
	City string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	IsPrimary bool `json:"is_primary,omitempty"`
	PersonId int `json:"person_id,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	State string `json:"state,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

// PersonContact represents the model.PersonContact swagger model.
type PersonContact struct {
	ContactType string `json:"contact_type,omitempty"`
	ContactValue string `json:"contact_value,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	IsPrimary bool `json:"is_primary,omitempty"`
	Notes string `json:"notes,omitempty"`
	PersonId int `json:"person_id,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

// PrimaryContact represents the model.PrimaryContact swagger model.
type PrimaryContact struct {
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy int `json:"created_by,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Id int `json:"id,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
	Notes string `json:"notes,omitempty"`
	OrganizationId int `json:"organization_id,omitempty"`
	PersonId int `json:"person_id,omitempty"`
	Role string `json:"role,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UpdatedBy int `json:"updated_by,omitempty"`
}

