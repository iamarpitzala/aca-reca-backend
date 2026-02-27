package form

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// SectionCreateInput is one section to create on the form version (e.g. COLLECTION, COST, SERVICE_FACILITY, OTHER_COST).
type SectionCreateInput struct {
	SectionTypeCode string              `json:"sectionTypeCode"` // COLLECTION, COST, SERVICE_FACILITY, OTHER_COST
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Order           int                 `json:"order"`
	Fields          []FieldCreateInput `json:"fields"`
}

// FieldCreateInput is one field (row) in a section.
type FieldCreateInput struct {
	Label                 string `json:"label"`
	PaymentResponsibility string `json:"paymentResponsibility"` // OWNER, CLINIC
}

// FormRequest is used for API interactions, validation, and client communication.
// Version is created in use case; status is DRAFT/PUBLISHED/ARCHIVED.
// Sections (and their fields) are created on the initial version when provided.
type FormRequest struct {
	ID                *string              `json:"id" validate:"omitempty,required"`
	Name              string              `json:"name" validate:"required,min=3,max=255"`
	Description       *string             `json:"description" validate:"omitempty,min=3,max=225"`
	ClinicID          string              `json:"clinicId" validate:"required"`
	Status            string              `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	CalculationMethod string              `json:"calculationMethod" validate:"required,oneof=NET GROSS"`
	Sections          []SectionCreateInput `json:"sections"`
	CreatedBy         string              `json:"createdBy" validate:"omitempty"`
	CreatedAt         time.Time           `json:"createdAt" validate:"omitempty"`
	UpdatedAt         *time.Time          `json:"updatedAt" validate:"omitempty"`
	DeletedAt         *time.Time          `json:"deletedAt" validate:"omitempty"`
}

// Form is the DB representation for tbl_custom_form (no version_id; versions in tbl_custom_form_version).
type Form struct {
	ID                string     `db:"id"`
	ClinicID          string     `db:"clinic_id"`
	Name              string     `db:"name"`
	Description       *string    `db:"description"`
	Status            string     `db:"status"`
	CalculationMethod string     `db:"calculation_method"`
	CreatedBy         string     `db:"created_by"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
}

func (f *Form) ToFormDB(req *FormRequest) {
	if req.ID != nil && *req.ID != "" {
		f.ID = *req.ID
	} else {
		f.ID = uuid.New().String()
	}
	f.ClinicID = req.ClinicID
	f.Name = req.Name
	f.Description = req.Description
	f.Status = req.Status
	f.CalculationMethod = req.CalculationMethod
	f.CreatedBy = req.CreatedBy
	if !req.CreatedAt.IsZero() {
		f.CreatedAt = req.CreatedAt
	}
	if req.UpdatedAt != nil {
		f.UpdatedAt = *req.UpdatedAt
	}
	f.DeletedAt = req.DeletedAt
}

// SectionWithFieldsResponse is a section plus its fields for GET form (Build form / edit).
type SectionWithFieldsResponse struct {
	SectionTypeCode string           `json:"sectionTypeCode"` // COLLECTION, COST, SERVICE_FACILITY, OTHER_COST
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	SectionOrder    int              `json:"sectionOrder"`
	Fields          []FieldResponse  `json:"fields"`
}

// FormResponse is for API return. ActiveVersionID is populated by use case when available.
// Sections (with fields) are populated for GET by ID so the Build form UI can be restored.
type FormResponse struct {
	ID                string                      `json:"id"`
	ClinicID          string                      `json:"clinicId"`
	Name              string                      `json:"name"`
	Description       *string                     `json:"description"`
	Status            string                      `json:"status"`
	CalculationMethod string                      `json:"calculationMethod"`
	ActiveVersionID   *int                        `json:"activeVersionId,omitempty"`
	Sections          []SectionWithFieldsResponse `json:"sections,omitempty"`
	CreatedBy         string                      `json:"createdBy"`
	CreatedAt         time.Time                   `json:"createdAt"`
	UpdatedAt         *time.Time                  `json:"updatedAt"`
	DeletedAt         *time.Time                  `json:"deletedAt"`
}

func (f *Form) ToFormResponse(activeVersionID *int) *FormResponse {
	return &FormResponse{
		ID:                f.ID,
		ClinicID:          f.ClinicID,
		Name:              f.Name,
		Description:       f.Description,
		Status:            f.Status,
		CalculationMethod: f.CalculationMethod,
		ActiveVersionID:   activeVersionID,
		CreatedBy:         f.CreatedBy,
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         lo.ToPtr(f.UpdatedAt),
		DeletedAt:         f.DeletedAt,
	}
}
