package domain

import (
	"time"

	"github.com/google/uuid"
)

// FinancialForm represents a financial form configuration

type GST struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Type       string    `db:"type" json:"type"`
	Percentage float64   `db:"percentage" json:"percentage"`
}

type FinancialForm struct {
	ID                uuid.UUID              `db:"id" json:"id"`
	ClinicID          uuid.UUID              `db:"clinic_id" json:"clinicId"`
	QuarterID         uuid.UUID              `db:"quarter_id" json:"quarterId"`
	GSTID             *int                   `db:"gst_id" json:"gstId"`
	Name              string                 `db:"name" json:"name"`
	CalculationMethod string                 `db:"calculation_method" json:"calculationMethod"` // "net" or "gross"
	Configuration     map[string]interface{} `db:"configuration" json:"configuration"`          // JSONB stored as map
	IsActive          bool                   `db:"is_active" json:"isActive"`
	CreatedAt         time.Time              `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time              `db:"updated_at" json:"updatedAt"`
	DeletedAt         *time.Time             `db:"deleted_at" json:"deletedAt"`
	GST               *GST                   `db:"gst" json:"gst,omitempty"`
	Quarter           *Quarter               `db:"quarter" json:"quarter,omitempty"`
}

// NetMethodConfig represents configuration for Net Method
type NetMethodConfig struct {
	CommissionPercent      float64 `json:"commissionPercent"`      // Default: 40%
	GSTOnCommissionPercent float64 `json:"gstOnCommissionPercent"` // Default: 10%
	LabFeeEnabled          bool    `json:"labFeeEnabled"`          // Default: true
	SuperHoldingEnabled    bool    `json:"superHoldingEnabled"`    // Default: false
	SuperPercent           float64 `json:"superPercent"`           // Fixed: 12% (non-editable when enabled)
}

// GrossMethodConfig represents configuration for Gross Method
type GrossMethodConfig struct {
	ServiceFeePercent        float64 `json:"serviceFeePercent"`        // Default: 60%
	GSTOnServiceFeePercent   float64 `json:"gstOnServiceFeePercent"`   // Default: 10%
	LabFeeEnabled            bool    `json:"labFeeEnabled"`            // Default: true
	LabFeePaidBy             string  `json:"labFeePaidBy"`             // "clinic" or "dentist"
	GSTOnLabFee              bool    `json:"gstOnLabFee"`              // Default: false
	GSTOnPatientFee          bool    `json:"gstOnPatientFee"`          // Default: false
	OutworkChargeRateEnabled bool    `json:"outworkChargeRateEnabled"` // Default: false
	OutworkRatePercent       float64 `json:"outworkRatePercent"`
	// Default: 40%
}

// CalculationInput represents input values for calculations
type CalculationInput struct {
	// Common fields
	GrossPatientFee float64 `json:"grossPatientFee"` // A

	// Net Method fields
	LabFee float64 `json:"labFee,omitempty"` // B (shown if LabFeeEnabled)

	// Gross Method fields
	GSTOnPatientFee     float64 `json:"gstOnPatientFee,omitempty"`     // B (for B4)
	LabFeePaidByDentist float64 `json:"labFeePaidByDentist,omitempty"` // I (for B4)
	MerchantFeeInclGST  float64 `json:"merchantFeeInclGST,omitempty"`  // G (for B3)
	BankFee             float64 `json:"bankFee,omitempty"`             // K (for B3)
	GSTOnLabFee         float64 `json:"gstOnLabFee,omitempty"`         // G (for B2)

	// Outwork Charge Rate fields (B5)
	OutworkLabFee           float64 `json:"outworkLabFee,omitempty"`           // w
	OutworkMerchantFee      float64 `json:"outworkMerchantFee,omitempty"`      // x
	OutworkGSTOnLabFee      float64 `json:"outworkGSTOnLabFee,omitempty"`      // y
	OutworkGSTOnMerchantFee float64 `json:"outworkGSTOnMerchantFee,omitempty"` // z
}

// CalculationResult represents calculated results
type CalculationResult struct {
	// Net Method results
	NetPatientFee          float64 `json:"netPatientFee,omitempty"`          // C
	CommissionForDentist   float64 `json:"commissionForDentist,omitempty"`   // D
	CommissionComponent    float64 `json:"commissionComponent,omitempty"`    // F (with super)
	SuperComponent         float64 `json:"superComponent,omitempty"`         // E (with super)
	TotalForReconciliation float64 `json:"totalForReconciliation,omitempty"` // G (with super)
	GSTOnCommission        float64 `json:"gstOnCommission,omitempty"`        // E (without super) or H (with super)
	TotalCommission        float64 `json:"totalCommission,omitempty"`        // F (without super) or I (with super)

	// Gross Method results
	ServiceFacilityFee      float64 `json:"serviceFacilityFee,omitempty"`      // D
	GSTOnServiceFee         float64 `json:"gstOnServiceFee,omitempty"`         // E
	TotalServiceFee         float64 `json:"totalServiceFee,omitempty"`         // F
	AmountRemittedToDentist float64 `json:"amountRemittedToDentist,omitempty"` // G, I, H, or J (varies by variant)

	// B3 specific
	MerchantFeeGSTComponent float64 `json:"merchantFeeGSTComponent,omitempty"` // I
	NetMerchantFee          float64 `json:"netMerchantFee,omitempty"`          // J

	// B5 specific
	TotalOutworkCharge          float64 `json:"totalOutworkCharge,omitempty"`          // B
	LabOtherCostCharge          float64 `json:"labOtherCostCharge,omitempty"`          // E
	TotalServiceFeeOtherCharges float64 `json:"totalServiceFeeOtherCharges,omitempty"` // F
	TotalChargesInclGST         float64 `json:"totalChargesInclGST,omitempty"`         // H
}

// BASMapping represents BAS field mappings
type BASMapping struct {
	Field1A  float64 `json:"field1A,omitempty"`  // GST on Sales
	Field1B  float64 `json:"field1B,omitempty"`  // GST Credit
	FieldG1  float64 `json:"fieldG1,omitempty"`  // Total Sales incl GST
	FieldG11 float64 `json:"fieldG11,omitempty"` // Clinic Expenses
}

// FinancialCalculation represents a calculation record
type FinancialCalculation struct {
	ID              uuid.UUID              `db:"id" json:"id"`
	FinancialFormID uuid.UUID              `db:"financial_form_id" json:"financialFormId"`
	InputData       map[string]interface{} `db:"input_data" json:"inputData"`
	CalculatedData  map[string]interface{} `db:"calculated_data" json:"calculatedData"`
	BASMapping      map[string]interface{} `db:"bas_mapping" json:"basMapping"`
	CreatedAt       time.Time              `db:"created_at" json:"createdAt"`
	CreatedBy       *uuid.UUID             `db:"created_by" json:"createdBy"`
}
