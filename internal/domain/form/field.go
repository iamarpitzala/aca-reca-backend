package form

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// FieldRequest matches tbl_custom_form_field: section_id, label, payment_responsibility_id.
type FieldRequest struct {
	ID                     *string `json:"id" validate:"omitempty,required"`
	SectionID              int     `json:"sectionId" validate:"required"`
	Label                  string  `json:"label" validate:"required,min=1,max=255"`
	PaymentResponsibilityID *int   `json:"paymentResponsibilityId" validate:"omitempty"`
}

func (f *FieldRequest) Validate() error {
	if f.Label == "" {
		return errors.New("label is required")
	}
	if f.SectionID <= 0 {
		return errors.New("sectionId is required")
	}
	return nil
}

// Field matches tbl_custom_form_field: id, section_id, label, payment_responsibility_id.
type Field struct {
	ID                     string     `db:"id"`
	SectionID              int        `db:"section_id"`
	Label                  string     `db:"label"`
	PaymentResponsibilityID *int       `db:"payment_responsibility_id"`
	CreatedAt              time.Time  `db:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at"`
	DeletedAt              *time.Time `db:"deleted_at"`
}

func (f *Field) ToFieldDB(req *FieldRequest) {
	if req.ID != nil && *req.ID != "" {
		f.ID = *req.ID
	} else {
		f.ID = uuid.New().String()
	}
	f.SectionID = req.SectionID
	f.Label = req.Label
	f.PaymentResponsibilityID = req.PaymentResponsibilityID
}

// FieldResponse is for API return. SectionID and FormVersionID/FormID can be resolved in use case.
type FieldResponse struct {
	ID                     string     `json:"id"`
	SectionID              int        `json:"sectionId"`
	Label                  string     `json:"label"`
	PaymentResponsibilityID *int       `json:"paymentResponsibilityId,omitempty"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
	DeletedAt              *time.Time `json:"deletedAt,omitempty"`
}

// NewFieldID returns a new UUID string for field id.
func NewFieldID() string {
	return uuid.New().String()
}

func (f *Field) ToFieldResponse() *FieldResponse {
	return &FieldResponse{
		ID:                     f.ID,
		SectionID:              f.SectionID,
		Label:                  f.Label,
		PaymentResponsibilityID: f.PaymentResponsibilityID,
		CreatedAt:              f.CreatedAt,
		UpdatedAt:              f.UpdatedAt,
		DeletedAt:              f.DeletedAt,
	}
}
