package form

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FieldType string

const (
	FieldTypeText   FieldType = "TEXT"
	FieldTypeNumber FieldType = "NUMBER"
)

type FieldWidth string

const (
	FieldWidthSmall  FieldWidth = "SMALL"
	FieldWidthMedium FieldWidth = "MEDIUM"
	FieldWidthLarge  FieldWidth = "LARGE"
	FieldWidthFull   FieldWidth = "FULL"
)

type FieldRequest struct {
	ID          *string    `json:"id" validate:"omitempty,required"`
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	DisplayName string     `json:"displayName" validate:"omitempty,required,min=3,max=255"`
	Placeholder string     `json:"placeholder" validate:"omitempty,required,min=3,max=255"`
	HelpText    string     `json:"helpText" validate:"omitempty,required,min=3,max=255"`
	Width       FieldWidth `json:"width" validate:"omitempty,required,oneof=SMALL MEDIUM LARGE FULL"`

	Required        bool             `json:"required" validate:"required"`
	ValidationRules []ValidationRule `json:"validationRules" validate:"omitempty,required,oneof=REQUIRED MIN MAX PATTERN CUSTOM"`

	Value json.RawMessage `json:"value" validate:"omitempty,required"`

	Type                  FieldType              `json:"type" validate:"required,oneof=TEXT NUMBER"`
	Section               Section                `json:"section" validate:"required,oneof=INCOME EXPENSE REDUCTION"`
	AocID                 int                    `json:"aocId" validate:"required,min=1"`
	IncludeInTotal        bool                   `json:"includeInTotal" validate:"required"`
	GstConfig             *GstConfigRequest      `json:"gstConfig" validate:"omitempty,required"`
	PaymentResponsibility *PaymentResponsibility `json:"paymentResponsibility" validate:"required,oneof=OWNER CLINIC"`
	CreatedAt             time.Time              `json:"createdAt" validate:"required"`
	UpdatedAt             *time.Time             `json:"updatedAt" validate:"required_with=CreatedAt"`
	DeletedAt             *time.Time             `json:"deletedAt" validate:"required_with=UpdatedAt"`
}

func (f *FieldRequest) Validate() error {
	if f.Required && f.Value == nil {
		return fmt.Errorf("value is required")
	}
	if f.ValidationRules != nil {
		for _, rule := range f.ValidationRules {
			if rule == ValidationRuleRequired && f.Value == nil {
				return fmt.Errorf("value is required")
			}
		}
	}
	return nil
}

type Field struct {
	ID uuid.UUID `db:"id"`

	Name        string     `db:"name"`
	DisplayName string     `db:"displayName"`
	Placeholder string     `db:"placeholder"`
	HelpText    string     `db:"helpText"`
	Width       FieldWidth `db:"width"`

	Required        bool             `db:"required"`
	ValidationRules []ValidationRule `db:"validationRules"`

	Value json.RawMessage `db:"value"`

	Type                  FieldType              `db:"type"`
	Section               Section                `db:"section"`
	AocID                 int                    `db:"aocId"`
	IncludeInTotal        bool                   `db:"includeInTotal"`
	GstConfig             *GstConfig             `db:"gstConfig"`
	PaymentResponsibility *PaymentResponsibility `db:"paymentResponsibility"`
	CreatedAt             time.Time              `db:"createdAt"`
	UpdatedAt             *time.Time             `db:"updatedAt"`
	DeletedAt             *time.Time             `db:"deletedAt"`
}

func (f *Field) ToFieldDB(field *FieldRequest) {

	gstConfig := &GstConfig{}
	gstConfig.ToGstDB(field.GstConfig)

	f.ID = uuid.MustParse(*field.ID)
	f.Name = field.Name
	f.DisplayName = field.DisplayName
	f.Placeholder = field.Placeholder
	f.HelpText = field.HelpText
	f.Width = field.Width
	f.Required = field.Required
	f.ValidationRules = field.ValidationRules
	f.Value = field.Value
	f.Type = field.Type
	f.Section = field.Section
	f.AocID = field.AocID
	f.IncludeInTotal = field.IncludeInTotal
	f.GstConfig = gstConfig
	f.PaymentResponsibility = field.PaymentResponsibility
	f.CreatedAt = field.CreatedAt
	f.UpdatedAt = field.UpdatedAt
	f.DeletedAt = field.DeletedAt
}

type FieldResponse struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	DisplayName     string           `json:"displayName"`
	Placeholder     string           `json:"placeholder"`
	HelpText        string           `json:"helpText"`
	Width           FieldWidth       `json:"width"`
	Required        bool             `json:"required"`
	ValidationRules []ValidationRule `json:"validationRules"`

	Value json.RawMessage `json:"value"`

	Type                  FieldType              `json:"type"`
	Section               Section                `json:"section"`
	AocID                 int                    `json:"aocId"`
	IncludeInTotal        bool                   `json:"includeInTotal"`
	GstConfig             *GstResponse           `json:"gstConfig"`
	PaymentResponsibility *PaymentResponsibility `json:"paymentResponsibility"`
	CreatedAt             time.Time              `json:"createdAt"`
	UpdatedAt             *time.Time             `json:"updatedAt"`
	DeletedAt             *time.Time             `json:"deletedAt"`
}

func (f *Field) ToFieldResponse() *FieldResponse {
	return &FieldResponse{
		ID:                    f.ID.String(),
		Name:                  f.Name,
		DisplayName:           f.DisplayName,
		Placeholder:           f.Placeholder,
		HelpText:              f.HelpText,
		Width:                 f.Width,
		Required:              f.Required,
		ValidationRules:       f.ValidationRules,
		Value:                 f.Value,
		Type:                  f.Type,
		Section:               f.Section,
		AocID:                 f.AocID,
		IncludeInTotal:        f.IncludeInTotal,
		GstConfig:             f.GstConfig.ToGstResponse(),
		PaymentResponsibility: f.PaymentResponsibility,
		CreatedAt:             f.CreatedAt,
		UpdatedAt:             f.UpdatedAt,
		DeletedAt:             f.DeletedAt,
	}
}
