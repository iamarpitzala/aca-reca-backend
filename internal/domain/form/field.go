package form

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type FieldRequest struct {
	ID          *string `json:"id" validate:"omitempty,required"`
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	DisplayName string  `json:"displayName" validate:"omitempty,required,min=3,max=255"`
	Placeholder string  `json:"placeholder" validate:"omitempty,required,min=3,max=255"`
	HelpText    string  `json:"helpText" validate:"omitempty,required,min=3,max=255"`

	Required        bool             `json:"required" validate:"required"`
	ValidationRules []ValidationRule `json:"validationRules" validate:"omitempty,required,oneof=REQUIRED MIN MAX PATTERN CUSTOM"`

	Value json.RawMessage `json:"value" validate:"omitempty,required"`

	Section     Section `json:"section" validate:"required,oneof=INCOME EXPENSE"`
	CoaID       string  `json:"coaId" validate:"required"`
	TaxConfigID string  `json:"taxConfigId" validate:"required"`

	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt" validate:"required_with=CreatedAt"`
	DeletedAt *time.Time `json:"deletedAt" validate:"required_with=UpdatedAt"`
}

func (f *FieldRequest) Validate() error {
	// If ValidationRules contains REQUIRED, value must be present (not nil or empty).
	if f.ValidationRules != nil {
		for _, rule := range f.ValidationRules {
			if rule == ValidationRuleRequired {
				if f.Value == nil || len(f.Value) == 0 {
					return errors.New("value is required due to REQUIRED validation rule")
				}
				// If REQUIRED is found, no need to check f.Required - we already require value.
				return nil
			}
		}
	}
	// Otherwise, if the Required field is set, also require value.
	if f.Required {
		if f.Value == nil || len(f.Value) == 0 {
			return errors.New("value is required")
		}
	}
	return nil
}

type Field struct {
	ID uuid.UUID `db:"id"`

	Name        string `db:"name"`
	DisplayName string `db:"displayName"`
	Placeholder string `db:"placeholder"`

	Required        bool             `db:"required"`
	ValidationRules []ValidationRule `db:"validationRules"`

	Value json.RawMessage `db:"value"`

	Section     Section    `db:"section"`
	CoaID       uuid.UUID  `db:"coaId"`
	TaxConfigID uuid.UUID  `db:"taxConfigId"`
	CreatedAt   time.Time  `db:"createdAt"`
	UpdatedAt   *time.Time `db:"updatedAt"`
	DeletedAt   *time.Time `db:"deletedAt"`
}

func (f *Field) ToFieldDB(field *FieldRequest) {
	if field.ID != nil && *field.ID != "" {
		f.ID = uuid.MustParse(*field.ID)
	}
	f.Name = field.Name
	f.DisplayName = field.DisplayName
	f.Placeholder = field.Placeholder
	f.Required = field.Required
	f.ValidationRules = field.ValidationRules
	f.Value = field.Value
	f.Section = field.Section
	if field.CoaID != "" {
		f.CoaID = uuid.MustParse(field.CoaID)
	}
	if field.TaxConfigID != "" {
		f.TaxConfigID = uuid.MustParse(field.TaxConfigID)
	}
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
	Required        bool             `json:"required"`
	ValidationRules []ValidationRule `json:"validationRules"`

	Value json.RawMessage `json:"value"`

	Section   Section            `json:"section"`
	CoaID     string             `json:"coaId"`
	TaxConfig *TaxConfigResponse `json:"taxConfig"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt *time.Time         `json:"updatedAt"`
	DeletedAt *time.Time         `json:"deletedAt"`
}

func (f *Field) ToFieldResponse(taxConfig *TaxConfig) *FieldResponse {
	var taxConfigResp *TaxConfigResponse
	if taxConfig != nil {
		taxConfigResp = taxConfig.ToTaxConfigResponse()
	}
	return &FieldResponse{
		ID:              f.ID.String(),
		Name:            f.Name,
		DisplayName:     f.DisplayName,
		Placeholder:     f.Placeholder,
		HelpText:        "", // HelpText is not present on Field model; adjust as needed if logic changes
		Required:        f.Required,
		ValidationRules: f.ValidationRules,
		Value:           f.Value,
		Section:         f.Section,
		CoaID:           f.CoaID.String(),
		TaxConfig:       taxConfigResp,
		CreatedAt:       f.CreatedAt,
		UpdatedAt:       f.UpdatedAt,
		DeletedAt:       f.DeletedAt,
	}
}
