package clinic

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ClinicFinancialSettingRequest is used for API interactions with validation.
type ClinicFinancialSettingRequest struct {
	ID                    *string    `json:"id" validate:"omitempty,required"`
	ClinicFinancialYearID int        `json:"clinicFinancialYearId" validate:"required"`
	CalculationMethod     string     `json:"calculationMethod" validate:"required,oneof=CASH ACCRUAL"`
	GSTRegistered         bool       `json:"gstRegistered" validate:"required"`
	GSTReportingFrequency string     `json:"gstReportingFrequency" validate:"required,oneof=QUARTERLY ANNUALLY"`
	DefaultAmountMode     string     `json:"defaultAmountMode" validate:"required,oneof=INCLUSIVE EXCLUSIVE"`
	LockDate              *time.Time `json:"lockDate"`
}

// Validate adds business logic validation for ClinicFinancialSettingRequest.
func (r *ClinicFinancialSettingRequest) Validate() error {
	if r.ClinicFinancialYearID <= 0 {
		return errors.New("clinicFinancialYearId must be positive")
	}
	if r.CalculationMethod != "CASH" && r.CalculationMethod != "ACCRUAL" {
		return errors.New("calculationMethod must be CASH or ACCRUAL")
	}
	if r.GSTReportingFrequency != "QUARTERLY" && r.GSTReportingFrequency != "ANNUALLY" {
		return errors.New("gstReportingFrequency must be QUARTERLY or ANNUALLY")
	}
	if r.DefaultAmountMode != "INCLUSIVE" && r.DefaultAmountMode != "EXCLUSIVE" {
		return errors.New("defaultAmountMode must be INCLUSIVE or EXCLUSIVE")
	}
	return nil
}

// ToClinicFinancialSetting maps request to DB model.
func (r *ClinicFinancialSettingRequest) ToClinicFinancialSetting(clinicID string) (*ClinicFinancialSetting, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	id := uuid.New().String()
	if r.ID != nil && *r.ID != "" {
		id = *r.ID
	}
	now := time.Now()
	return &ClinicFinancialSetting{
		ID:                    id,
		ClinicID:              clinicID,
		ClinicFinancialYearID: r.ClinicFinancialYearID,
		AccountingMethod:      r.CalculationMethod,
		GSTRegistered:         r.GSTRegistered,
		GSTReportingFrequency: r.GSTReportingFrequency,
		DefaultAmountMode:     r.DefaultAmountMode,
		LockDate:              r.LockDate,
		CreatedAt:             now,
		UpdatedAt:             now,
		DeletedAt:             nil,
	}, nil
}

// ClinicFinancialSetting is the DB representation for tbl_clinic_financial_settings.
type ClinicFinancialSetting struct {
	ID                    string     `db:"id" json:"id"`
	ClinicID              string     `db:"clinic_id" json:"clinicId"`
	ClinicFinancialYearID int        `db:"clinic_financial_year_id" json:"clinicFinancialYearId"`
	AccountingMethod      string     `db:"accounting_method" json:"calculationMethod"`
	GSTRegistered         bool       `db:"gst_registered" json:"gstRegistered"`
	GSTReportingFrequency string     `db:"gst_reporting_frequency" json:"gstReportingFrequency"`
	DefaultAmountMode     string     `db:"default_amount_mode" json:"defaultAmountMode"`
	LockDate              *time.Time `db:"lock_date" json:"lockDate"`
	CreatedAt             time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt             *time.Time `db:"deleted_at" json:"deletedAt"`
}

// ClinicFinancialSettingResponse for API response.
type ClinicFinancialSettingResponse struct {
	ID                    string     `json:"id"`
	ClinicID              string     `json:"clinicId"`
	ClinicFinancialYearID int        `json:"clinicFinancialYearId"`
	CalculationMethod     string     `json:"calculationMethod"`
	GSTRegistered         bool       `json:"gstRegistered"`
	GSTReportingFrequency string     `json:"gstReportingFrequency"`
	DefaultAmountMode     string     `json:"defaultAmountMode"`
	LockDate              *time.Time `json:"lockDate"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
	DeletedAt             *time.Time `json:"deletedAt,omitempty"`
}

func (f *ClinicFinancialSetting) ToClinicFinancialSettingResponse() *ClinicFinancialSettingResponse {
	return &ClinicFinancialSettingResponse{
		ID:                    f.ID,
		ClinicID:              f.ClinicID,
		ClinicFinancialYearID: f.ClinicFinancialYearID,
		CalculationMethod:     f.AccountingMethod,
		GSTRegistered:         f.GSTRegistered,
		GSTReportingFrequency: f.GSTReportingFrequency,
		DefaultAmountMode:     f.DefaultAmountMode,
		LockDate:              f.LockDate,
		CreatedAt:             f.CreatedAt,
		UpdatedAt:             f.UpdatedAt,
		DeletedAt:             f.DeletedAt,
	}
}
