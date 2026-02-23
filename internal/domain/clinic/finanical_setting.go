package clinic

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ClinicFinancialSettingRequest is used for API interactions with validation
type ClinicFinancialSettingRequest struct {
	ID                    *string    `json:"id" validate:"omitempty,required"`
	FinancialYearID       int        `json:"financialYearId" validate:"required"`
	FinancialQuarterID    int        `json:"financialQuarterId" validate:"required"`
	AccountingMethod      string     `json:"accountingMethod" validate:"required,oneof=NET GROSS"`
	GSTRegistered         bool       `json:"gstRegistered" validate:"required"`
	GSTReportingFrequency string     `json:"gstReportingFrequency" validate:"required,oneof=QUARTERLY ANNUALLY"`
	DefaultAmountMode     string     `json:"defaultAmountMode" validate:"required,oneof=INCLUSIVE EXCLUSIVE"`
	LockDate              *time.Time `json:"lockDate"`
	CreatedAt             *time.Time `json:"createdAt" validate:"omitempty,required"`
	UpdatedAt             *time.Time `json:"updatedAt" validate:"omitempty,required_with=CreatedAt"`
	DeletedAt             *time.Time `json:"deletedAt" validate:"omitempty"`
}

// Validate adds business logic validation for ClinicFinancialSettingRequest.
func (r *ClinicFinancialSettingRequest) Validate() error {
	if r.FinancialYearID <= 0 {
		return errors.New("financialYearId must be positive")
	}
	if r.FinancialQuarterID < 1 || r.FinancialQuarterID > 4 {
		return errors.New("financialQuarterId must be between 1 and 4")
	}

	if r.AccountingMethod != "NET" && r.AccountingMethod != "GROSS" {
		return errors.New("accountingMethod must be NET or GROSS")
	}

	if r.GSTReportingFrequency != "QUARTERLY" && r.GSTReportingFrequency != "ANNUALLY" {
		return errors.New("gstReportingFrequency must be QUARTERLY or ANNUALLY")
	}

	if r.DefaultAmountMode != "INCLUSIVE" && r.DefaultAmountMode != "EXCLUSIVE" {
		return errors.New("defaultAmountMode must be INCLUSIVE or EXCLUSIVE")
	}

	if r.CreatedAt == nil {
		now := time.Now()
		r.CreatedAt = &now
	}
	if r.UpdatedAt == nil {
		r.UpdatedAt = r.CreatedAt
	}

	return nil
}

// ToDBModel performs mapping from request to DB model, ensuring fields are properly populated.
func (r *ClinicFinancialSettingRequest) ToClinicFinancialSetting(clinicID string) (*ClinicFinancialSetting, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	id := ""
	if r.ID != nil && *r.ID != "" {
		id = *r.ID
	} else {
		id = uuid.New().String()
	}

	setting := &ClinicFinancialSetting{
		ID:                    id,
		ClinicID:              clinicID,
		FinancialYearID:       r.FinancialYearID,
		FinancialQuarterID:    r.FinancialQuarterID,
		AccountingMethod:      r.AccountingMethod,
		GSTRegistered:         r.GSTRegistered,
		GSTReportingFrequency: r.GSTReportingFrequency,
		DefaultAmountMode:     r.DefaultAmountMode,
		LockDate:              r.LockDate,
		CreatedAt:             *r.CreatedAt,
		UpdatedAt:             *r.UpdatedAt,
		DeletedAt:             r.DeletedAt,
	}

	return setting, nil
}

// ClinicFinancialSetting is the DB representation for the financial settings table
type ClinicFinancialSetting struct {
	ID                    string     `db:"id" json:"id"`
	ClinicID              string     `db:"clinic_id" json:"clinicId"`
	FinancialYearID       int        `db:"financial_year_id" json:"financialYearId"`
	FinancialQuarterID    int        `db:"financial_quarter_id" json:"financialQuarterId"`
	AccountingMethod      string     `db:"accounting_method" json:"accountingMethod"`
	GSTRegistered         bool       `db:"gst_registered" json:"gstRegistered"`
	GSTReportingFrequency string     `db:"gst_reporting_frequency" json:"gstReportingFrequency"`
	DefaultAmountMode     string     `db:"default_amount_mode" json:"defaultAmountMode"`
	LockDate              *time.Time `db:"lock_date" json:"lockDate"`
	CreatedAt             time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt             *time.Time `db:"deleted_at" json:"deletedAt"`
}

// ClinicFinancialSettingResponse is for API response payloads
type ClinicFinancialSettingResponse struct {
	ID                    string     `json:"id"`
	ClinicID              string     `json:"clinicId"`
	FinancialYearID       int        `json:"financialYearId"`
	FinancialQuarterID    int        `json:"financialQuarterId"`
	AccountingMethod      string     `json:"accountingMethod"`
	GSTRegistered         bool       `json:"gstRegistered"`
	GSTReportingFrequency string     `json:"gstReportingFrequency"`
	DefaultAmountMode     string     `json:"defaultAmountMode"`
	LockDate              *time.Time `json:"lockDate"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
	DeletedAt             *time.Time `json:"deletedAt"`
}

// ToClinicFinancialSettingResponse maps a DB model into an API response
func (f *ClinicFinancialSetting) ToClinicFinancialSettingResponse() *ClinicFinancialSettingResponse {
	return &ClinicFinancialSettingResponse{
		ID:                    f.ID,
		ClinicID:              f.ClinicID,
		FinancialYearID:       f.FinancialYearID,
		FinancialQuarterID:    f.FinancialQuarterID,
		AccountingMethod:      f.AccountingMethod,
		GSTRegistered:         f.GSTRegistered,
		GSTReportingFrequency: f.GSTReportingFrequency,
		DefaultAmountMode:     f.DefaultAmountMode,
		LockDate:              f.LockDate,
		CreatedAt:             f.CreatedAt,
		UpdatedAt:             f.UpdatedAt,
		DeletedAt:             f.DeletedAt,
	}
}
