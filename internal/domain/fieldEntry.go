package domain

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type GSTConfig struct {
	Enabled bool    `json:"enabled"`
	Rate    float64 `json:"rate"`
	Type    string  `json:"type"`
}

// DB model for field entry -- matches tbl_custom_form_entry schema
type FieldEntry struct {
	ID                uuid.UUID  `db:"id"`
	FormID            uuid.UUID  `db:"form_id"`
	FormVersionID     uuid.UUID  `db:"form_version_id"`
	CustomFormFieldID uuid.UUID  `db:"tbl_custom_form_field_id"`
	Value             float64    `db:"value"`
	CreatedBy         uuid.UUID  `db:"created_by"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
}

// Entry field value from frontend
type EntryFieldValueInput struct {
	FieldID         string   `json:"fieldId"`
	FieldName       string   `json:"fieldName"`
	Value           float64  `json:"value"`
	ManualGSTAmount *float64 `json:"manualGstAmount,omitempty"`
}

// UnmarshalJSON custom unmarshaler to handle string/number/boolean values
func (e *EntryFieldValueInput) UnmarshalJSON(data []byte) error {
	type Alias EntryFieldValueInput
	aux := &struct {
		Value interface{} `json:"value"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert value to float64
	switch v := aux.Value.(type) {
	case float64:
		e.Value = v
	case float32:
		e.Value = float64(v)
	case int:
		e.Value = float64(v)
	case int64:
		e.Value = float64(v)
	case string:
		// Try to parse string as float64
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			e.Value = parsed
		} else {
			// If parsing fails, default to 0
			e.Value = 0
		}
	case bool:
		if v {
			e.Value = 1
		} else {
			e.Value = 0
		}
	case nil:
		e.Value = 0
	default:
		e.Value = 0
	}

	return nil
}

// Full entry request structure (matches frontend CreateEntryRequest)
type CreateEntryRequest struct {
	FormID                string                 `json:"formId"`
	ClinicID              string                 `json:"clinicId"`
	QuarterID             *string                `json:"quarterId,omitempty"`
	Values                []EntryFieldValueInput `json:"values"`
	EntryDate             string                 `json:"entryDate"`
	Description           string                 `json:"description,omitempty"`
	Remarks               string                 `json:"remarks,omitempty"`
	PaymentResponsibility *string                `json:"paymentResponsibility,omitempty"`
	Deductions            json.RawMessage        `json:"deductions,omitempty"`
}

// API request structure (legacy - for single field entries)
type FieldEntryRequest struct {
	FormID            string     `json:"formId"`
	CustomFormFieldID string     `json:"customFormFieldId"`
	Value             float64    `json:"value"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt"`
}

// API response structure
type FieldEntryResponse struct {
	ID                string     `json:"id"`
	FormID            string     `json:"formId"`
	CustomFormFieldID string     `json:"customFormFieldId"`
	Value             float64    `json:"value"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt"`
}

// Entry field value response (matches frontend EntryFieldValue)
type EntryFieldValueResponse struct {
	FieldID         string     `json:"fieldId"`
	FieldName       string     `json:"fieldName"`
	Value           float64    `json:"value"`
	BaseAmount      *float64   `json:"baseAmount,omitempty"`
	GSTAmount       *float64   `json:"gstAmount,omitempty"`
	TotalAmount     *float64   `json:"totalAmount,omitempty"`
	ManualGSTAmount *float64   `json:"manualGstAmount,omitempty"`
	GSTConfig       *GSTConfig `json:"gstConfig,omitempty"`
}

// Full entry response (matches frontend ApiEntry/CustomFormEntry)
type EntryResponse struct {
	ID                    string                    `json:"id"`
	FormID                string                    `json:"formId"`
	FormName              string                    `json:"formName"`
	FormType              string                    `json:"formType"`
	ClinicID              string                    `json:"clinicId"`
	QuarterID             *string                   `json:"quarterId,omitempty"`
	Values                []EntryFieldValueResponse `json:"values"`
	Calculations          json.RawMessage           `json:"calculations"`
	EntryDate             string                    `json:"entryDate"`
	Description           string                    `json:"description,omitempty"`
	Remarks               string                    `json:"remarks,omitempty"`
	PaymentResponsibility *string                   `json:"paymentResponsibility,omitempty"`
	Deductions            json.RawMessage           `json:"deductions,omitempty"`
	CreatedBy             string                    `json:"createdBy"`
	CreatedAt             time.Time                 `json:"createdAt"`
	UpdatedAt             time.Time                 `json:"updatedAt"`
}

// Request to DB mapping
func (req *FieldEntryRequest) ToDBModel(createdBy uuid.UUID) (*FieldEntry, error) {
	formID, err := uuid.Parse(req.FormID)
	if err != nil {
		return nil, err
	}
	customFormFieldID, err := uuid.Parse(req.CustomFormFieldID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &FieldEntry{
		ID:                uuid.New(),
		FormID:            formID,
		CustomFormFieldID: customFormFieldID,
		Value:             req.Value,
		CreatedBy:         createdBy,
		CreatedAt:         now,
		UpdatedAt:         now,
		DeletedAt:         nil,
	}, nil
}

// DB to Response mapping
func (f *FieldEntry) ToResponse() *FieldEntryResponse {
	return &FieldEntryResponse{
		ID:                f.ID.String(),
		FormID:            f.FormID.String(),
		CustomFormFieldID: f.CustomFormFieldID.String(),
		Value:             f.Value,
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         f.UpdatedAt,
		DeletedAt:         f.DeletedAt,
	}
}

func GetGSTConfig(gstConfig GSTConfig) *GSTConfig {
	if !gstConfig.Enabled {
		return &GSTConfig{
			Enabled: false,
			Rate:    0,
			Type:    "",
		}
	}

	return &gstConfig
}
