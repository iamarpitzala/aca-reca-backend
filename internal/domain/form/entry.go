package form

import (
	"time"

	"github.com/google/uuid"
)

// Entry is one submission of a form version (tbl_custom_form_entry).
type Entry struct {
	ID            string     `db:"id"`
	FormVersionID int        `db:"form_version_id"`
	SubmittedBy   *string    `db:"submitted_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

// EntryValue is one field value for an entry (tbl_custom_form_entry_value).
type EntryValue struct {
	ID        int        `db:"id"`
	EntryID   string     `db:"entry_id"`
	FieldID   string     `db:"field_id"`
	Value     float64    `db:"value"`
	GSTAmount *float64   `db:"gst_amount"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// EntryCreateRequest is the API request to create an entry with values.
type EntryCreateRequest struct {
	FormVersionID string              `json:"formVersionId" binding:"required"`
	Values        []EntryValueRequest `json:"values" binding:"required,min=1"`
}

// EntryValueRequest is one field value in create/update.
type EntryValueRequest struct {
	FieldID   string   `json:"fieldId" binding:"required"`
	Value     float64  `json:"value"`
	GSTAmount *float64 `json:"gstAmount,omitempty"`
}

// EntryUpdateRequest is the API request to update an entry's values.
type EntryUpdateRequest struct {
	Values []EntryValueRequest `json:"values" binding:"required,min=1"`
}

// EntryResponse is the API response for an entry with its values.
type EntryResponse struct {
	ID            string               `json:"id"`
	FormVersionID int                 `json:"formVersionId"`
	SubmittedBy   *string             `json:"submittedBy,omitempty"`
	Values        []EntryValueResponse `json:"values"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
	DeletedAt     *time.Time           `json:"deletedAt,omitempty"`
}

// EntryValueResponse is one value in the entry response.
type EntryValueResponse struct {
	ID        int      `json:"id"`
	EntryID   string   `json:"entryId"`
	FieldID   string   `json:"fieldId"`
	Value     float64  `json:"value"`
	GSTAmount *float64 `json:"gstAmount,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (e *Entry) ToEntryResponse(values []EntryValue) *EntryResponse {
	valResp := make([]EntryValueResponse, len(values))
	for i := range values {
		valResp[i] = EntryValueResponse{
			ID:        values[i].ID,
			EntryID:   values[i].EntryID,
			FieldID:   values[i].FieldID,
			Value:     values[i].Value,
			GSTAmount: values[i].GSTAmount,
			CreatedAt: values[i].CreatedAt,
			UpdatedAt: values[i].UpdatedAt,
		}
	}
	return &EntryResponse{
		ID:            e.ID,
		FormVersionID: e.FormVersionID,
		SubmittedBy:   e.SubmittedBy,
		Values:        valResp,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
		DeletedAt:     e.DeletedAt,
	}
}

// NewEntry builds an Entry for create (ID generated).
func NewEntry(formVersionID int, submittedBy *string) *Entry {
	now := time.Now()
	return &Entry{
		ID:            uuid.New().String(),
		FormVersionID: formVersionID,
		SubmittedBy:   submittedBy,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
	}
}

// EntryValueForCreate builds an EntryValue for create (no ID; DB assigns).
func EntryValueForCreate(entryID, fieldID string, value float64, gstAmount *float64) EntryValue {
	now := time.Now()
	return EntryValue{
		EntryID:   entryID,
		FieldID:   fieldID,
		Value:     value,
		GSTAmount: gstAmount,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}
}
