package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	constants "github.com/iamarpitzala/aca-reca-backend/constant"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

// FinancialForm represents a financial form configuration

type GST struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Type       string    `db:"type"`
	Percentage float64   `db:"percentage"`
}

type GSTResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Percentage float64   `json:"percentage"`
}

func (g *GST) ToResponse() *GSTResponse {
	return &GSTResponse{
		ID:         g.ID,
		Name:       g.Name,
		Type:       g.Type,
		Percentage: g.Percentage,
	}
}

type FinancialFormRequest struct {
	ID                uuid.UUID       `json:"id"`
	ClinicID          uuid.UUID       `json:"clinicId"`
	QuarterID         uuid.UUID       `json:"quarterId"`
	Name              string          `json:"name"`
	CalculationMethod string          `json:"calculationMethod"` // net | gross
	Configuration     json.RawMessage `json:"configuration"`     // JSONB
	IsActive          bool            `json:"isActive"`
}

type FinancialForm struct {
	ID                uuid.UUID       `db:"id"`
	ClinicID          uuid.UUID       `db:"clinic_id"`
	QuarterID         uuid.UUID       `db:"quarter_id"`
	Name              string          `db:"name"`
	CalculationMethod string          `db:"calculation_method"` // net | gross
	Configuration     json.RawMessage `db:"configuration"`      // JSONB
	IsActive          bool            `db:"is_active"`
	CreatedAt         time.Time       `db:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at"`
	DeletedAt         *time.Time      `db:"deleted_at"`
}

func (rr *FinancialFormRequest) ToRepo() (*FinancialForm, error) {
	if rr.CalculationMethod != "net" && rr.CalculationMethod != "gross" {
		return nil, fmt.Errorf("invalid calculation method: %s", rr.CalculationMethod)
	}

	id := rr.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	now := time.Now()
	return &FinancialForm{
		ID:                id,
		ClinicID:          rr.ClinicID,
		QuarterID:         rr.QuarterID,
		Name:              rr.Name,
		CalculationMethod: rr.CalculationMethod,
		Configuration:     rr.Configuration,
		IsActive:          rr.IsActive,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

type FinancialFormResponse struct {
	ID                uuid.UUID              `json:"id"`
	ClinicID          uuid.UUID              `json:"clinicId"`
	QuarterID         uuid.UUID              `json:"quarterId"`
	Name              string                 `json:"name"`
	CalculationMethod string                 `json:"calculationMethod"`
	Configuration     map[string]interface{} `json:"configuration"`
	IsActive          bool                   `json:"isActive"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
	DeletedAt         *time.Time             `json:"deletedAt,omitempty"`
}

type NetMethodConfig struct {
	CommissionType      string   `json:"commissionType"`
	ClinicCommission    float64  `json:"clinicCommission"`
	OwnerCommission     float64  `json:"ownerCommission"`
	SuperHoldingEnabled bool     `json:"superHoldingEnabled"`
	SuperPercent        *float64 `json:"superPercent,omitempty"`
	LabFees             bool     `json:"labFees"`
}

type GrossMethodConfig struct {
	CommissionType     string   `json:"commissionType"`
	ClinicCommission   float64  `json:"clinicCommission"`
	OwnerCommission    float64  `json:"ownerCommission"`
	PaidBy             string   `json:"paidBy"` // clinic | owner
	GSTOnLabFee        bool     `json:"gstOnLabFee"`
	OutworkRateEnabled bool     `json:"outworkRateEnabled"`
	OutworkRatePercent *float64 `json:"outworkRatePercent,omitempty"`
}

func (f *FinancialForm) ToFinancialFormResponse() (*FinancialFormResponse, error) {
	var cfg map[string]interface{}
	var err error

	switch f.CalculationMethod {
	case "net":
		var netCfg NetMethodConfig
		if err := json.Unmarshal(f.Configuration, &netCfg); err != nil {
			return nil, err
		}
		cfg, err = util.StructToMap(netCfg)

	case "gross":
		var grossCfg GrossMethodConfig
		if err := json.Unmarshal(f.Configuration, &grossCfg); err != nil {
			return nil, err
		}
		cfg, err = util.StructToMap(grossCfg)

	default:
		return nil, fmt.Errorf("invalid calculation method: %s", f.CalculationMethod)
	}

	if err != nil {
		return nil, err
	}

	return &FinancialFormResponse{
		ID:                f.ID,
		ClinicID:          f.ClinicID,
		QuarterID:         f.QuarterID,
		Name:              f.Name,
		CalculationMethod: f.CalculationMethod,
		Configuration:     cfg,
		IsActive:          f.IsActive,
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         f.UpdatedAt,
		DeletedAt:         f.DeletedAt,
	}, nil
}

func (n NetMethodConfig) Validate() error {
	if n.CommissionType == constants.PERCENTAGE {

		if n.ClinicCommission < 0 || n.OwnerCommission < 0 {
			return fmt.Errorf("commission values cannot be negative")
		}

		if n.ClinicCommission+n.OwnerCommission != 100 {
			return fmt.Errorf("clinic + owner commission must be 100")
		}

		if n.SuperHoldingEnabled {
			if n.SuperPercent == nil {
				return fmt.Errorf("superPercent required when superHoldingEnabled")
			}
			if *n.SuperPercent <= 0 || *n.SuperPercent > 100 {
				return fmt.Errorf("superPercent must be between 1 and 100")
			}
		}
	}

	return nil
}

func (g GrossMethodConfig) Validate() error {
	if g.CommissionType == constants.PERCENTAGE {

		if g.ClinicCommission < 0 || g.OwnerCommission < 0 {
			return fmt.Errorf("commission values cannot be negative")
		}

		if g.ClinicCommission+g.OwnerCommission != 100 {
			return fmt.Errorf("clinic + owner commission must be 100")
		}

		switch g.PaidBy {
		case constants.PAID_BY_CLINIC, constants.PAID_BY_OWNER:
			// valid
		default:
			return fmt.Errorf("paidBy must be clinic or owner")
		}
	}

	return nil
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
