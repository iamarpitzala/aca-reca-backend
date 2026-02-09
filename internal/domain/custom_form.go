package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Custom form field and form types (match frontend customForm.ts)

type CustomFieldType string

const (
	FieldTypeText     CustomFieldType = "text"
	FieldTypeNumber   CustomFieldType = "number"
	FieldTypeDate     CustomFieldType = "date"
	FieldTypeDropdown CustomFieldType = "dropdown"
	FieldTypeCheckbox CustomFieldType = "checkbox"
	FieldTypeTextarea CustomFieldType = "textarea"
	FieldTypeCurrency CustomFieldType = "currency"
)

type FormStatus string

const (
	FormStatusDraft    FormStatus = "draft"
	FormStatusPublished FormStatus = "published"
	FormStatusArchived FormStatus = "archived"
)

type FormType string

const (
	FormTypeIncome   FormType = "income"
	FormTypeExpense  FormType = "expense"
	FormTypeBoth     FormType = "both"
)

type CalculationMethod string

const (
	CalcMethodNet   CalculationMethod = "net"
	CalcMethodGross CalculationMethod = "gross"
)

type PaymentResponsibility string

const (
	PaymentOwner  PaymentResponsibility = "owner"
	PaymentClinic PaymentResponsibility = "clinic"
)

// FieldGSTConfig, DropdownOption, ConditionalLogicRule, FieldValidation, CustomFormField
// are stored inside fields JSONB; we use json.RawMessage for flexibility.

// CustomForm DB model
type CustomForm struct {
	ID                          uuid.UUID       `db:"id"`
	ClinicID                    uuid.UUID       `db:"clinic_id"`
	Name                        string          `db:"name"`
	Description                 string          `db:"description"`
	CalculationMethod          string          `db:"calculation_method"`
	FormType                   string          `db:"form_type"`
	Status                     string          `db:"status"`
	Fields                     json.RawMessage `db:"fields"`
	DefaultPaymentResponsibility *string        `db:"default_payment_responsibility"`
	ServiceFacilityFeePercent   *float64        `db:"service_facility_fee_percent"`
	Version                     int             `db:"version"`
	CreatedBy                   uuid.UUID       `db:"created_by"`
	CreatedAt                   time.Time       `db:"created_at"`
	UpdatedAt                   time.Time       `db:"updated_at"`
	PublishedAt                 *time.Time      `db:"published_at"`
	DeletedAt                   *time.Time      `db:"deleted_at"`
}

// CustomFormEntry DB model
type CustomFormEntry struct {
	ID                    uuid.UUID       `db:"id"`
	FormID                uuid.UUID       `db:"form_id"`
	FormName              string          `db:"form_name"`
	FormType              string          `db:"form_type"`
	ClinicID              uuid.UUID       `db:"clinic_id"`
	QuarterID             *uuid.UUID      `db:"quarter_id"`
	Values                json.RawMessage `db:"values"`
	Calculations          json.RawMessage `db:"calculations"`
	EntryDate             time.Time       `db:"entry_date"`
	Description           string          `db:"description"`
	Remarks               string          `db:"remarks"`
	PaymentResponsibility *string         `db:"payment_responsibility"`
	Deductions            json.RawMessage `db:"deductions"`
	CreatedBy             uuid.UUID       `db:"created_by"`
	CreatedAt             time.Time       `db:"created_at"`
	UpdatedAt             time.Time       `db:"updated_at"`
	DeletedAt             *time.Time      `db:"deleted_at"`
}

// Request/Response DTOs (JSON camelCase for API)

type CreateCustomFormRequest struct {
	ClinicID                     string          `json:"clinicId"`
	Name                         string          `json:"name"`
	Description                  string          `json:"description"`
	CalculationMethod           string          `json:"calculationMethod"`
	FormType                     string          `json:"formType"`
	Fields                       json.RawMessage `json:"fields"`
	DefaultPaymentResponsibility *string         `json:"defaultPaymentResponsibility,omitempty"`
	ServiceFacilityFeePercent    *float64        `json:"serviceFacilityFeePercent,omitempty"`
}

type UpdateCustomFormRequest struct {
	Name                         *string         `json:"name,omitempty"`
	Description                  *string         `json:"description,omitempty"`
	Fields                       json.RawMessage `json:"fields,omitempty"`
	DefaultPaymentResponsibility *string         `json:"defaultPaymentResponsibility,omitempty"`
	ServiceFacilityFeePercent    *float64        `json:"serviceFacilityFeePercent,omitempty"`
}

type CustomFormResponse struct {
	ID                          string          `json:"id"`
	ClinicID                    string          `json:"clinicId"`
	Name                        string          `json:"name"`
	Description                 string          `json:"description"`
	CalculationMethod          string          `json:"calculationMethod"`
	FormType                   string          `json:"formType"`
	Status                     string          `json:"status"`
	Fields                     json.RawMessage `json:"fields"`
	DefaultPaymentResponsibility *string        `json:"defaultPaymentResponsibility,omitempty"`
	ServiceFacilityFeePercent   *float64        `json:"serviceFacilityFeePercent,omitempty"`
	Version                    int             `json:"version"`
	CreatedBy                  string          `json:"createdBy"`
	CreatedAt                  time.Time       `json:"createdAt"`
	UpdatedAt                  time.Time       `json:"updatedAt"`
	PublishedAt                *time.Time      `json:"publishedAt,omitempty"`
}

type CreateEntryRequest struct {
	FormID                string          `json:"formId"`
	ClinicID              string          `json:"clinicId"`
	QuarterID             *string         `json:"quarterId,omitempty"`
	Values                json.RawMessage `json:"values"`
	EntryDate             string          `json:"entryDate"` // ISO date
	Description           string          `json:"description,omitempty"`
	Remarks               string          `json:"remarks,omitempty"`
	PaymentResponsibility *string         `json:"paymentResponsibility,omitempty"`
	Deductions            json.RawMessage `json:"deductions,omitempty"`
	// Pre-calculated totals (frontend sends after running calculateEntryTotals / applyDeductionsToCalculations)
	Calculations          json.RawMessage `json:"calculations,omitempty"`
}

type UpdateEntryRequest struct {
	Values       json.RawMessage `json:"values"`
	Calculations json.RawMessage `json:"calculations,omitempty"`
}

// PreviewCalculationsRequest for live calculation (no save)
type PreviewCalculationsRequest struct {
	FormID     string          `json:"formId"`
	Values     json.RawMessage `json:"values"`
	Deductions json.RawMessage `json:"deductions,omitempty"`
}

type CustomFormEntryResponse struct {
	ID                    string          `json:"id"`
	FormID                string          `json:"formId"`
	FormName              string          `json:"formName"`
	FormType              string          `json:"formType"`
	ClinicID              string          `json:"clinicId"`
	QuarterID             *string         `json:"quarterId,omitempty"`
	Values                json.RawMessage `json:"values"`
	Calculations          json.RawMessage `json:"calculations"`
	EntryDate             time.Time       `json:"entryDate"`
	Description           string          `json:"description,omitempty"`
	Remarks               string          `json:"remarks,omitempty"`
	PaymentResponsibility *string         `json:"paymentResponsibility,omitempty"`
	Deductions            json.RawMessage `json:"deductions,omitempty"`
	CreatedBy             string          `json:"createdBy"`
	CreatedAt             time.Time       `json:"createdAt"`
	UpdatedAt             time.Time       `json:"updatedAt"`
}
