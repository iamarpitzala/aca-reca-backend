package form

import (
	"time"

	"github.com/google/uuid"
)

type EntryFieldRequest struct {
	ID                *string  `json:"id" validate:"omitempty,required"`
	FormID            string   `json:"formId"`
	ClinicID          string   `json:"clinicId"`
	FormVersionID     string   `json:"formVersionId"`
	CustomFormFieldID string   `json:"customFormFieldId"`
	Value             float64  `json:"value"`
	GSTAmount         *float64 `json:"gstAmount"`
}

type EntryFieldUpdateRequest struct {
	ID        string   `json:"id" validate:"omitempty,required"`
	Value     float64  `json:"value"`
	GSTAmount *float64 `json:"gstAmount"`
}

func (e *EntryFieldRequest) ToFieldEntryDB() (*FieldEntry, error) {
	id := ""
	if e.ID != nil && *e.ID != "" {
		id = uuid.NewString()
	} else {
		id = *e.ID
	}
	return &FieldEntry{
		ID:                id,
		ClinicID:          e.ClinicID,
		FormID:            e.FormID,
		FormVersionID:     e.FormVersionID,
		CustomFormFieldID: e.CustomFormFieldID,
		Value:             e.Value,
		GSTAmount:         e.GSTAmount,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		DeletedAt:         nil,
	}, nil
}

type FieldEntry struct {
	ID                string     `db:"id"`
	ClinicID          string     `db:"clinic_id"`
	FormID            string     `db:"form_id"`
	FormVersionID     string     `db:"form_version_id"`
	CustomFormFieldID string     `db:"field_id"`
	Value             float64    `db:"value"`
	GSTAmount         *float64   `db:"gst_amount"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
}

func (e *FieldEntry) ToFieldEntryResponse() *FieldEntryResponse {
	return &FieldEntryResponse{
		ID:                e.ID,
		ClinicID:          e.ClinicID,
		FormID:            e.FormID,
		FormVersionID:     e.FormVersionID,
		CustomFormFieldID: e.CustomFormFieldID,
		Value:             e.Value,
		GSTAmount:         e.GSTAmount,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
		DeletedAt:         e.DeletedAt,
	}
}

type FieldEntryResponse struct {
	ID                string     `json:"id"`
	ClinicID          string     `json:"clinicId"`
	FormID            string     `json:"formId"`
	FormVersionID     string     `json:"formVersionId"`
	CustomFormFieldID string     `json:"customFormFieldId"`
	Value             float64    `json:"value"`
	GSTAmount         *float64   `json:"gstAmount"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt"`
}
