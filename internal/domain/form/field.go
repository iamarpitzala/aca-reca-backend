package form

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type FieldRequest struct {
	ID            *string `json:"id" validate:"omitempty,required"`
	FormVersionID string  `json:"formVersionId" validate:"required"`
	FormID        string  `json:"formId" validate:"required"`

	Label         string  `json:"label" validate:"required,min=3,max=255"`
	SectionTypeID string  `json:"sectionTypeId" validate:"required"`
	Description   *string `json:"description" validate:"omitempty,min=3,max=500"`
	IsRequired    bool    `json:"isRequired" validate:"required"`
	CoaID         string  `json:"coaId" validate:"required"`

	Placeholder *string  `json:"placeholder" validate:"omitempty,min=1,max=255"`
	MinValue    *float64 `json:"minValue" validate:"omitempty"`
	MaxValue    *float64 `json:"maxValue" validate:"omitempty"`
	FieldOrder  int      `json:"fieldOrder" validate:"required,min=1"`
	TaxTypeID   *string  `json:"taxTypeId" validate:"omitempty"`

	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt" validate:"required_with=CreatedAt"`
	DeletedAt *time.Time `json:"deletedAt" validate:"omitempty"`
}

func (f *FieldRequest) Validate() error {
	if f.IsRequired {
		if f.MinValue != nil && f.MaxValue != nil && *f.MinValue > *f.MaxValue {
			return errors.New("minValue cannot be greater than maxValue")
		}
	}
	if f.Label == "" {
		return errors.New("label is required")
	}
	if f.FormID == "" || f.FormVersionID == "" {
		return errors.New("formId and formVersionId are required")
	}
	return nil
}

type Field struct {
	ID            string  `db:"id"`
	FormVersionID string  `db:"form_version_id"`
	FormID        string  `db:"form_id"`
	Label         string  `db:"label"`
	SectionTypeID string  `db:"section_type_id"`
	Description   *string `db:"description"`
	IsRequired    bool    `db:"is_required"`
	CoaID         string  `db:"coa_id"`

	Placeholder *string  `db:"placeholder"`
	MinValue    *float64 `db:"min_value"`
	MaxValue    *float64 `db:"max_value"`
	FieldOrder  int      `db:"field_order"`
	TaxTypeID   *string  `db:"tax_type_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (f *Field) ToFieldDB(req *FieldRequest) {
	if req.ID != nil && *req.ID != "" {
		f.ID = *req.ID
	} else {
		f.ID = uuid.New().String()
	}
	f.FormVersionID = req.FormVersionID
	f.FormID = req.FormID
	f.Label = req.Label
	f.SectionTypeID = req.SectionTypeID
	f.Description = req.Description
	f.IsRequired = req.IsRequired
	f.CoaID = req.CoaID
	f.Placeholder = req.Placeholder
	f.MinValue = req.MinValue
	f.MaxValue = req.MaxValue
	f.FieldOrder = req.FieldOrder
	f.TaxTypeID = req.TaxTypeID
	f.CreatedAt = req.CreatedAt
	if req.UpdatedAt != nil {
		f.UpdatedAt = *req.UpdatedAt
	} else {
		f.UpdatedAt = req.CreatedAt
	}
	f.DeletedAt = req.DeletedAt
}

type FieldResponse struct {
	ID            string `json:"id"`
	FormVersionID string `json:"formVersionId"`
	FormID        string `json:"formId"`

	Label         string  `json:"label"`
	SectionTypeID string  `json:"sectionTypeId"`
	Description   *string `json:"description"`
	IsRequired    bool    `json:"isRequired"`
	CoaID         string  `json:"coaId"`

	Placeholder *string  `json:"placeholder"`
	MinValue    *float64 `json:"minValue"`
	MaxValue    *float64 `json:"maxValue"`
	FieldOrder  int      `json:"fieldOrder"`
	TaxTypeID   *string  `json:"taxTypeId"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (f *Field) ToFieldResponse() *FieldResponse {
	return &FieldResponse{
		ID:            f.ID,
		FormVersionID: f.FormVersionID,
		FormID:        f.FormID,
		Label:         f.Label,
		SectionTypeID: f.SectionTypeID,
		Description:   f.Description,
		IsRequired:    f.IsRequired,
		CoaID:         f.CoaID,
		Placeholder:   f.Placeholder,
		MinValue:      f.MinValue,
		MaxValue:      f.MaxValue,
		FieldOrder:    f.FieldOrder,
		TaxTypeID:     f.TaxTypeID,
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
		DeletedAt:     f.DeletedAt,
	}
}
