package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// FinancialYearStart represents the start month of the financial year
type FinancialYearStart string

const (
	FinancialYearStartJuly    FinancialYearStart = "JULY"
	FinancialYearStartJanuary FinancialYearStart = "JANUARY"
)

// AccountingMethod represents the accounting method used
type AccountingMethod string

const (
	AccountingMethodCash    AccountingMethod = "CASH"
	AccountingMethodAccrual AccountingMethod = "ACCRUAL"
)

// GSTReportingFrequency represents how often GST is reported
type GSTReportingFrequency string

const (
	GSTReportingFrequencyQuarterly GSTReportingFrequency = "QUARTERLY"
	GSTReportingFrequencyAnnually  GSTReportingFrequency = "ANNUALLY"
)

// DefaultAmountMode represents how amounts are entered by default
type DefaultAmountMode string

const (
	DefaultAmountModeGSTInclusive DefaultAmountMode = "GST_INCLUSIVE"
	DefaultAmountModeGSTExclusive DefaultAmountMode = "GST_EXCLUSIVE"
)

// ClinicFinancialSettings represents financial settings for a clinic
type ClinicFinancialSettings struct {
	ID                    uuid.UUID           `db:"id" json:"id"`
	ClinicID              uuid.UUID           `db:"clinic_id" json:"clinicId"`
	FinancialYearStart    FinancialYearStart  `db:"financial_year_start" json:"financialYearStart"`
	AccountingMethod      AccountingMethod    `db:"accounting_method" json:"accountingMethod"`
	GSTRegistered         bool                `db:"gst_registered" json:"gstRegistered"`
	GSTReportingFrequency GSTReportingFrequency `db:"gst_reporting_frequency" json:"gstReportingFrequency"`
	DefaultAmountMode     DefaultAmountMode   `db:"default_amount_mode" json:"defaultAmountMode"`
	LockDate              *time.Time          `db:"lock_date" json:"lockDate"`
	GSTDefaults           json.RawMessage     `db:"gst_defaults" json:"gstDefaults"` // JSONB in DB
	CreatedAt             time.Time           `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time           `db:"updated_at" json:"updatedAt"`
	DeletedAt             *time.Time         `db:"deleted_at" json:"deletedAt"`
}

// ClinicFinancialSettingsRequest represents the request to create/update financial settings
type ClinicFinancialSettingsRequest struct {
	FinancialYearStart    FinancialYearStart  `json:"financialYearStart" validate:"required,oneof=JULY JANUARY"`
	AccountingMethod      AccountingMethod    `json:"accountingMethod" validate:"required,oneof=CASH ACCRUAL"`
	GSTRegistered         bool                `json:"gstRegistered" validate:"required"`
	GSTReportingFrequency GSTReportingFrequency `json:"gstReportingFrequency" validate:"required,oneof=QUARTERLY ANNUALLY"`
	DefaultAmountMode     DefaultAmountMode   `json:"defaultAmountMode" validate:"required,oneof=GST_INCLUSIVE GST_EXCLUSIVE"`
	LockDate              *time.Time          `json:"lockDate"`
	GSTDefaults           map[string]string   `json:"gstDefaults"`
}

// ToClinicFinancialSettings converts a request to a domain model
func (r *ClinicFinancialSettingsRequest) ToClinicFinancialSettings(clinicID uuid.UUID) (*ClinicFinancialSettings, error) {
	var gstDefaultsJSON json.RawMessage
	if r.GSTDefaults != nil && len(r.GSTDefaults) > 0 {
		data, err := json.Marshal(r.GSTDefaults)
		if err != nil {
			return nil, err
		}
		gstDefaultsJSON = data
	} else {
		gstDefaultsJSON = json.RawMessage("{}")
	}
	
	return &ClinicFinancialSettings{
		ClinicID:              clinicID,
		FinancialYearStart:    r.FinancialYearStart,
		AccountingMethod:      r.AccountingMethod,
		GSTRegistered:         r.GSTRegistered,
		GSTReportingFrequency: r.GSTReportingFrequency,
		DefaultAmountMode:     r.DefaultAmountMode,
		LockDate:              r.LockDate,
		GSTDefaults:           gstDefaultsJSON,
	}, nil
}

// GetGSTDefaultsMap converts JSONB GSTDefaults to a map
func (s *ClinicFinancialSettings) GetGSTDefaultsMap() (map[string]string, error) {
	if len(s.GSTDefaults) == 0 {
		return make(map[string]string), nil
	}
	var m map[string]string
	if err := json.Unmarshal(s.GSTDefaults, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// SetGSTDefaultsMap converts a map to JSONB GSTDefaults
func (s *ClinicFinancialSettings) SetGSTDefaultsMap(m map[string]string) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	s.GSTDefaults = data
	return nil
}
