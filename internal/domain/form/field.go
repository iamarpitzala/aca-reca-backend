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

	StartDate     string  `json:"startDate" validate:"required"` // YYYY-MM-DD
	EndDate       *string `json:"endDate" validate:"omitempty"`  // YYYY-MM-DD or empty
	Label         string  `json:"label" validate:"required,min=1,max=255"`
	SectionTypeID int     `json:"sectionTypeId" validate:"required"`
	Description   *string `json:"description" validate:"omitempty,max=500"`
	IsRequired    bool    `json:"isRequired"`
	CoaID         string  `json:"coaId" validate:"required"`

	Placeholder *string  `json:"placeholder" validate:"omitempty,max=255"`
	MinValue    *float64 `json:"minValue" validate:"omitempty"`
	MaxValue    *float64 `json:"maxValue" validate:"omitempty"`
	FieldOrder  int      `json:"fieldOrder" validate:"required,min=0"`
	TaxTypeID   *int     `json:"taxTypeId" validate:"omitempty"`
}

func (f *FieldRequest) Validate() error {
	if f.Label == "" {
		return errors.New("label is required")
	}
	if f.FormID == "" || f.FormVersionID == "" {
		return errors.New("formId and formVersionId are required")
	}
	if f.MinValue != nil && f.MaxValue != nil && *f.MinValue > *f.MaxValue {
		return errors.New("minValue cannot be greater than maxValue")
	}
	return nil
}

type Field struct {
	ID            string     `db:"id"`
	FormVersionID string     `db:"form_version_id"`
	FormID        string     `db:"form_id"`
	StartDate     time.Time  `db:"start_date"`
	EndDate       *time.Time `db:"end_date"`
	Label         string     `db:"label"`
	SectionTypeID int        `db:"section_type_id"`
	Description   *string    `db:"description"`
	IsRequired    bool       `db:"is_required"`
	CoaID         string    `db:"coa_id"`

	Placeholder *string  `db:"placeholder"`
	MinValue    *float64 `db:"min_value"`
	MaxValue    *float64 `db:"max_value"`
	FieldOrder  int      `db:"field_order"`
	TaxTypeID   *int     `db:"tax_type_id"`

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
	if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
		f.StartDate = t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			f.EndDate = &t
		}
	}
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
}

type FieldResponse struct {
	ID            string     `json:"id"`
	FormVersionID string     `json:"formVersionId"`
	FormID        string     `json:"formId"`
	StartDate     time.Time  `json:"startDate"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	Label         string     `json:"label"`
	SectionTypeID int        `json:"sectionTypeId"`
	Description   *string    `json:"description,omitempty"`
	IsRequired    bool       `json:"isRequired"`
	CoaID         string     `json:"coaId"`
	Placeholder   *string    `json:"placeholder,omitempty"`
	MinValue      *float64   `json:"minValue,omitempty"`
	MaxValue      *float64   `json:"maxValue,omitempty"`
	FieldOrder    int        `json:"fieldOrder"`
	TaxTypeID     *int       `json:"taxTypeId,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

func (f *Field) ToFieldResponse() *FieldResponse {
	return &FieldResponse{
		ID:            f.ID,
		FormVersionID: f.FormVersionID,
		FormID:        f.FormID,
		StartDate:     f.StartDate,
		EndDate:       f.EndDate,
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
