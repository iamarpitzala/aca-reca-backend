package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

// Custom form field and form types (match frontend customForm.ts)

type CustomFieldType string

const (
	FieldTypeText   CustomFieldType = "text"
	FieldTypeNumber CustomFieldType = "number"
	// FieldTypeDate     CustomFieldType = "date"
	// FieldTypeDropdown CustomFieldType = "dropdown"
	// FieldTypeCheckbox CustomFieldType = "checkbox"
	// FieldTypeTextarea CustomFieldType = "textarea"
	// FieldTypeCurrency CustomFieldType = "currency"
)

type FormStatus string

const (
	FormStatusDraft     FormStatus = FormStatus(util.FormStatusDraft)
	FormStatusPublished FormStatus = FormStatus(util.FormStatusPublished)
	FormStatusArchived  FormStatus = FormStatus(util.FormStatusArchived)
)

type FormType string

const (
	FormTypeIncome  FormType = FormType(util.FormTypeIncome)
	FormTypeExpense FormType = FormType(util.FormTypeExpense)
	FormTypeBoth    FormType = FormType(util.FormTypeBoth)
)

type CalculationMethod string

const (
	CalcMethodNet   CalculationMethod = CalculationMethod(util.MethodTypeNet)
	CalcMethodGross CalculationMethod = CalculationMethod(util.MethodTypeGross)
)

type PaymentResponsibility string

const (
	PaymentOwner  PaymentResponsibility = PaymentResponsibility(util.PaymentResponsibilityOwner)
	PaymentClinic PaymentResponsibility = PaymentResponsibility(util.PaymentResponsibilityClinic)
)

//DB Models

type CustomForm struct {
	ID                           uuid.UUID  `db:"id"`
	ClinicID                     uuid.UUID  `db:"clinic_id"`
	Name                         string     `db:"name"`
	Description                  string     `db:"description"`
	CalculationMethod            string     `db:"calculation_method"`
	FormType                     string     `db:"form_type"`
	Status                       string     `db:"status"`
	DefaultPaymentResponsibility *string    `db:"default_payment_responsibility"`
	CreatedBy                    uuid.UUID  `db:"created_by"`
	CreatedAt                    time.Time  `db:"created_at"`
	UpdatedAt                    time.Time  `db:"updated_at"`
	DeletedAt                    *time.Time `db:"deleted_at"`
}

type CustomFormVersion struct {
	ID        uuid.UUID `db:"id"`
	FormID    uuid.UUID `db:"form_id"`
	Version   int       `db:"version"`
	IsActive  bool      `db:"is_active"`
	CreatedBy uuid.UUID `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}

type CustomFormField struct {
	ID            uuid.UUID       `db:"id"`
	FormVersionID uuid.UUID       `db:"form_version_id"`
	FormID        uuid.UUID       `db:"form_id"`
	FieldKey      string          `db:"field_key"`
	Label         string          `db:"label"`
	FieldType     string          `db:"field_type"`
	Section       string          `db:"section"`
	IsRequired    bool            `db:"is_required"`
	CoaID         uuid.UUID       `db:"coa_id"`
	Placeholder   string          `db:"placeholder"`
	MinValue      *float64        `db:"min_value"`
	MaxValue      *float64        `db:"max_value"`
	FieldOrder    int             `db:"field_order"`
	GSTConfig     bool            `db:"gst_config"`
	GSTRate       *float64        `db:"gst_rate"`
	GSTType       string          `db:"gst_type"`
	Metadata      json.RawMessage `db:"metadata"`
}

type CustomFormCalculation struct {
	ID                           uuid.UUID `db:"id"`
	FormVersionID                uuid.UUID `db:"form_version_id"`
	FormID                       uuid.UUID `db:"form_id"`
	CalculationMethod            string    `db:"calculation_method"`
	DefaultPaymentResponsibility *string   `db:"default_payment_responsibility"`
	ServiceFacilityFeePercent    *float64  `db:"service_facility_fee_percent"`
	OutworkEnabled               bool      `db:"outwork_enabled"`
	OutworkRatePercent           *float64  `db:"outwork_rate_percent"`
	CreatedAt                    time.Time `db:"created_at"`
}

type CustomFormPublish struct {
	ID            uuid.UUID `db:"id"`
	FormVersionID uuid.UUID `db:"form_version_id"`
	FormID        uuid.UUID `db:"form_id"`
	PublishedBy   uuid.UUID `db:"published_by"`
	PublishedAt   time.Time `db:"published_at"`
}

// Request/Response DTOs

type CreateCustomFormRequest struct {
	ClinicID                     string                 `json:"clinicId"`
	Name                         string                 `json:"name"`
	Description                  string                 `json:"description"`
	CalculationMethod            string                 `json:"calculationMethod"`
	FormType                     string                 `json:"formType"`
	DefaultPaymentResponsibility *string                `json:"defaultPaymentResponsibility,omitempty"`
	Fields                       []CustomFormFieldInput `json:"fields,omitempty"`
	ServiceFacilityFeePercent    *float64               `json:"serviceFacilityFeePercent,omitempty"`
	OutworkEnabled               *bool                  `json:"outworkEnabled,omitempty"`
	OutworkRatePercent           *float64               `json:"outworkRatePercent,omitempty"`
}

type CustomFormFieldInput struct {
	Name                  string          `json:"name"`
	Label                 string          `json:"label"`
	Type                  string          `json:"type"`
	Required              bool            `json:"required"`
	DefaultValue          interface{}     `json:"defaultValue,omitempty"`
	Placeholder           string          `json:"placeholder,omitempty"`
	Description           string          `json:"description,omitempty"`
	GSTConfig             json.RawMessage `json:"gstConfig,omitempty"`
	DropdownOptions       json.RawMessage `json:"dropdownOptions,omitempty"`
	Validation            json.RawMessage `json:"validation,omitempty"`
	Order                 int             `json:"order"`
	IncludeInTotal        bool            `json:"includeInTotal"`
	Section               string          `json:"section,omitempty"`
	AccountID             *string         `json:"accountId,omitempty"`
	PaymentResponsibility *string         `json:"paymentResponsibility,omitempty"`
}

// Request to DB mapping for CreateCustomFormRequest
func (req *CreateCustomFormRequest) ToDBModel(createdBy uuid.UUID) (*CustomForm, error) {
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, err
	}
	customForm := &CustomForm{
		ID:                           uuid.New(),
		ClinicID:                     clinicID,
		Name:                         req.Name,
		Description:                  req.Description,
		CalculationMethod:            req.CalculationMethod,
		FormType:                     req.FormType,
		Status:                       string(FormStatusDraft),
		DefaultPaymentResponsibility: req.DefaultPaymentResponsibility,
		CreatedBy:                    createdBy,
		CreatedAt:                    time.Now(),
		UpdatedAt:                    time.Now(),
		DeletedAt:                    nil,
	}
	return customForm, nil
}

type UpdateCustomFormRequest struct {
	Name                         *string                `json:"name,omitempty"`
	Description                  *string                `json:"description,omitempty"`
	CalculationMethod            *string                `json:"calculationMethod,omitempty"`
	FormType                     *string                `json:"formType,omitempty"`
	DefaultPaymentResponsibility *string                `json:"defaultPaymentResponsibility,omitempty"`
	Fields                       []CustomFormFieldInput `json:"fields,omitempty"`
	ServiceFacilityFeePercent    *float64               `json:"serviceFacilityFeePercent,omitempty"`
	OutworkEnabled               *bool                  `json:"outworkEnabled,omitempty"`
	OutworkRatePercent           *float64               `json:"outworkRatePercent,omitempty"`
}

type CustomFormResponse struct {
	ID                           string                    `json:"id"`
	ClinicID                     string                    `json:"clinicId"`
	Name                         string                    `json:"name"`
	Description                  string                    `json:"description"`
	CalculationMethod            string                    `json:"calculationMethod"`
	FormType                     string                    `json:"formType"`
	Status                       string                    `json:"status"`
	DefaultPaymentResponsibility *string                   `json:"defaultPaymentResponsibility,omitempty"`
	Fields                       []CustomFormFieldResponse `json:"fields,omitempty"`
	ServiceFacilityFeePercent    *float64                  `json:"serviceFacilityFeePercent,omitempty"`
	OutworkEnabled               *bool                     `json:"outworkEnabled,omitempty"`
	OutworkRatePercent           *float64                  `json:"outworkRatePercent,omitempty"`
	CreatedBy                    string                    `json:"createdBy"`
	CreatedAt                    time.Time                 `json:"createdAt"`
	UpdatedAt                    time.Time                 `json:"updatedAt"`
}

// Map frontend field input to backend CustomFormField
func (input *CustomFormFieldInput) ToDBModel(formID uuid.UUID, formVersionID uuid.UUID, userID uuid.UUID) (*CustomFormField, error) {
	// Parse GST config from JSON
	var gstConfig struct {
		Enabled bool     `json:"enabled"`
		Rate    *float64 `json:"rate"`
		Type    string   `json:"type"`
	}
	if len(input.GSTConfig) > 0 {
		if err := json.Unmarshal(input.GSTConfig, &gstConfig); err != nil {
			return nil, err
		}
	}

	// Build metadata JSON with all extra field info
	metadata := map[string]interface{}{
		"name":                  input.Name,
		"description":           input.Description,
		"defaultValue":          input.DefaultValue,
		"includeInTotal":        input.IncludeInTotal,
		"dropdownOptions":       nil,
		"validation":            nil,
		"paymentResponsibility": input.PaymentResponsibility,
	}
	if len(input.DropdownOptions) > 0 {
		var dropdownOptions interface{}
		if err := json.Unmarshal(input.DropdownOptions, &dropdownOptions); err == nil {
			metadata["dropdownOptions"] = dropdownOptions
		}
	}
	if len(input.Validation) > 0 {
		var validation interface{}
		if err := json.Unmarshal(input.Validation, &validation); err == nil {
			metadata["validation"] = validation
		}
	}
	metadataJSON, _ := json.Marshal(metadata)

	// Parse account ID (COA ID) - required by DB constraint
	var coaID uuid.UUID
	if input.AccountID != nil && *input.AccountID != "" {
		var err error
		coaID, err = uuid.Parse(*input.AccountID)
		if err != nil {
			return nil, fmt.Errorf("invalid account ID: %w", err)
		}
	} else {
		// TODO: DB requires coa_id to be NOT NULL, but frontend may not always provide accountId
		// Consider making coa_id nullable in DB or providing a default account
		// For now, return error if accountId is not provided
		return nil, errors.New("accountId is required for custom form fields")
	}

	// Determine section (default to INCOME if not provided)
	section := input.Section
	if section == "" {
		section = "INCOME"
	}
	section = strings.ToUpper(section)

	// Determine GST type
	gstType := "EXCLUSIVE"
	if gstConfig.Type != "" {
		gstType = strings.ToUpper(gstConfig.Type)
	}

	// Determine field type
	fieldType := strings.ToUpper(input.Type)

	return &CustomFormField{
		ID:            uuid.New(),
		FormVersionID: formVersionID,
		FormID:        formID,
		FieldKey:      input.Name, // Use name as field_key
		Label:         input.Label,
		FieldType:     fieldType,
		Section:       section,
		IsRequired:    input.Required,
		CoaID:         coaID,
		Placeholder:   input.Placeholder,
		MinValue:      nil, // Can be extracted from validation if needed
		MaxValue:      nil, // Can be extracted from validation if needed
		FieldOrder:    input.Order,
		GSTConfig:     gstConfig.Enabled,
		GSTRate:       gstConfig.Rate,
		GSTType:       gstType,
		Metadata:      metadataJSON,
	}, nil
}

// DB to Response mapping for CustomForm
func (f *CustomForm) ToResponse(fields []CustomFormFieldResponse) *CustomFormResponse {
	return &CustomFormResponse{
		ID:                           f.ID.String(),
		ClinicID:                     f.ClinicID.String(),
		Name:                         f.Name,
		Description:                  f.Description,
		CalculationMethod:            f.CalculationMethod,
		FormType:                     f.FormType,
		Status:                       f.Status,
		DefaultPaymentResponsibility: f.DefaultPaymentResponsibility,
		Fields:                       fields,
		CreatedBy:                    f.CreatedBy.String(),
		CreatedAt:                    f.CreatedAt,
		UpdatedAt:                    f.UpdatedAt,
	}
}

type CustomFormVersionResponse struct {
	ID        string    `json:"id"`
	FormID    string    `json:"formId"`
	Version   int       `json:"version"`
	IsActive  bool      `json:"isActive"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

// DB to Response mapping for CustomFormVersion
func (v *CustomFormVersion) ToResponse() *CustomFormVersionResponse {
	return &CustomFormVersionResponse{
		ID:        v.ID.String(),
		FormID:    v.FormID.String(),
		Version:   v.Version,
		IsActive:  v.IsActive,
		CreatedBy: v.CreatedBy.String(),
		CreatedAt: v.CreatedAt,
	}
}

type CustomFormFieldRequest struct {
	FormVersionID string          `json:"formVersionId"`
	FormID        string          `json:"formId"`
	FieldKey      string          `json:"fieldKey"`
	Label         string          `json:"label"`
	FieldType     string          `json:"fieldType"`
	Section       string          `json:"section"`
	IsRequired    bool            `json:"isRequired"`
	CoaID         string          `json:"coaId"`
	Placeholder   string          `json:"placeholder"`
	MinValue      *float64        `json:"minValue,omitempty"`
	MaxValue      *float64        `json:"maxValue,omitempty"`
	FieldOrder    int             `json:"fieldOrder"`
	GSTConfig     bool            `json:"gstConfig"`
	GSTRate       *float64        `json:"gstRate,omitempty"`
	GSTType       string          `json:"gstType"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
}

// Request to DB mapping for CustomFormFieldRequest
func (req *CustomFormFieldRequest) ToDBModel() (*CustomFormField, error) {
	formVersionID, err := uuid.Parse(req.FormVersionID)
	if err != nil {
		return nil, err
	}
	formID, err := uuid.Parse(req.FormID)
	if err != nil {
		return nil, err
	}
	coaID, err := uuid.Parse(req.CoaID)
	if err != nil {
		return nil, err
	}
	return &CustomFormField{
		ID:            uuid.New(),
		FormVersionID: formVersionID,
		FormID:        formID,
		FieldKey:      req.FieldKey,
		Label:         req.Label,
		FieldType:     req.FieldType,
		Section:       req.Section,
		IsRequired:    req.IsRequired,
		CoaID:         coaID,
		Placeholder:   req.Placeholder,
		MinValue:      req.MinValue,
		MaxValue:      req.MaxValue,
		FieldOrder:    req.FieldOrder,
		GSTConfig:     req.GSTConfig,
		GSTRate:       req.GSTRate,
		GSTType:       req.GSTType,
		Metadata:      req.Metadata,
	}, nil
}

type CustomFormFieldResponse struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	Label                 string          `json:"label"`
	Type                  string          `json:"type"`
	Required              bool            `json:"required"`
	DefaultValue          interface{}     `json:"defaultValue,omitempty"`
	Placeholder           string          `json:"placeholder,omitempty"`
	Description           string          `json:"description,omitempty"`
	GSTConfig             json.RawMessage `json:"gstConfig"`
	DropdownOptions       json.RawMessage `json:"dropdownOptions,omitempty"`
	Validation            json.RawMessage `json:"validation,omitempty"`
	Order                 int             `json:"order"`
	IncludeInTotal        bool            `json:"includeInTotal"`
	Section               string          `json:"section,omitempty"`
	AccountID             string          `json:"accountId,omitempty"`
	PaymentResponsibility *string         `json:"paymentResponsibility,omitempty"`
}

// DB to Response mapping for CustomFormField
// Maps backend structure to frontend CustomFormField format
func (f *CustomFormField) ToResponse() *CustomFormFieldResponse {
	// Parse metadata to extract frontend fields
	var metadata map[string]interface{}
	if len(f.Metadata) > 0 {
		json.Unmarshal(f.Metadata, &metadata)
	}

	// Extract values from metadata with defaults
	name := f.FieldKey
	if metadata != nil {
		if n, ok := metadata["name"].(string); ok && n != "" {
			name = n
		}
	}

	description := ""
	if metadata != nil {
		if d, ok := metadata["description"].(string); ok {
			description = d
		}
	}

	defaultValue := interface{}(nil)
	if metadata != nil {
		if dv, ok := metadata["defaultValue"]; ok {
			defaultValue = dv
		}
	}

	includeInTotal := false
	if metadata != nil {
		if it, ok := metadata["includeInTotal"].(bool); ok {
			includeInTotal = it
		}
	}

	// Build GST config object
	gstConfigJSON := json.RawMessage(`{"enabled":false,"rate":0,"type":"exclusive"}`)
	if f.GSTConfig {
		gstType := strings.ToLower(f.GSTType)
		if gstType == "" {
			gstType = "exclusive"
		}
		rate := 0.0
		if f.GSTRate != nil {
			rate = *f.GSTRate
		}
		gstConfigMap := map[string]interface{}{
			"enabled": true,
			"rate":    rate,
			"type":    gstType,
		}
		if gstConfigBytes, err := json.Marshal(gstConfigMap); err == nil {
			gstConfigJSON = gstConfigBytes
		}
	}

	// Extract dropdown options and validation from metadata
	var dropdownOptions json.RawMessage
	var validation json.RawMessage
	var paymentResponsibility *string
	if metadata != nil {
		if do, ok := metadata["dropdownOptions"]; ok && do != nil {
			if doBytes, err := json.Marshal(do); err == nil {
				dropdownOptions = doBytes
			}
		}
		if v, ok := metadata["validation"]; ok && v != nil {
			if vBytes, err := json.Marshal(v); err == nil {
				validation = vBytes
			}
		}
		if pr, ok := metadata["paymentResponsibility"].(string); ok {
			paymentResponsibility = &pr
		}
	}

	// Map section to lowercase
	section := strings.ToLower(f.Section)
	if section == "" {
		section = "income"
	}

	// Map field type to lowercase
	fieldType := strings.ToLower(f.FieldType)

	// Account ID
	accountID := ""
	if f.CoaID != uuid.Nil {
		accountID = f.CoaID.String()
	}

	return &CustomFormFieldResponse{
		ID:                    f.ID.String(),
		Name:                  name,
		Label:                 f.Label,
		Type:                  fieldType,
		Required:              f.IsRequired,
		DefaultValue:          defaultValue,
		Placeholder:           f.Placeholder,
		Description:           description,
		GSTConfig:             gstConfigJSON,
		DropdownOptions:       dropdownOptions,
		Validation:            validation,
		Order:                 f.FieldOrder,
		IncludeInTotal:        includeInTotal,
		Section:               section,
		AccountID:             accountID,
		PaymentResponsibility: paymentResponsibility,
	}
}

type CustomFormCalculationResponse struct {
	ID                           string    `json:"id"`
	FormVersionID                string    `json:"formVersionId"`
	FormID                       string    `json:"formId"`
	CalculationMethod            string    `json:"calculationMethod"`
	DefaultPaymentResponsibility *string   `json:"defaultPaymentResponsibility,omitempty"`
	ServiceFacilityFeePercent    *float64  `json:"serviceFacilityFeePercent,omitempty"`
	OutworkEnabled               bool      `json:"outworkEnabled"`
	OutworkRatePercent           *float64  `json:"outworkRatePercent,omitempty"`
	CreatedAt                    time.Time `json:"createdAt"`
}

// DB to Response mapping for CustomFormCalculation
func (c *CustomFormCalculation) ToResponse() *CustomFormCalculationResponse {
	return &CustomFormCalculationResponse{
		ID:                           c.ID.String(),
		FormVersionID:                c.FormVersionID.String(),
		FormID:                       c.FormID.String(),
		CalculationMethod:            c.CalculationMethod,
		DefaultPaymentResponsibility: c.DefaultPaymentResponsibility,
		ServiceFacilityFeePercent:    c.ServiceFacilityFeePercent,
		OutworkEnabled:               c.OutworkEnabled,
		OutworkRatePercent:           c.OutworkRatePercent,
		CreatedAt:                    c.CreatedAt,
	}
}

type CustomFormPublishResponse struct {
	ID            string    `json:"id"`
	FormVersionID string    `json:"formVersionId"`
	FormID        string    `json:"formId"`
	PublishedBy   string    `json:"publishedBy"`
	PublishedAt   time.Time `json:"publishedAt"`
}

// DB to Response mapping for CustomFormPublish
func (p *CustomFormPublish) ToResponse() *CustomFormPublishResponse {
	return &CustomFormPublishResponse{
		ID:            p.ID.String(),
		FormVersionID: p.FormVersionID.String(),
		FormID:        p.FormID.String(),
		PublishedBy:   p.PublishedBy.String(),
		PublishedAt:   p.PublishedAt,
	}
}
