package calculation

import (
	"encoding/json"
)

// NetCalculationInput contains input parameters for NET method calculation
type NetCalculationInput struct {
	TotalPaymentReceived   float64  // Total payment received (gross amount)
	CommissionPercent      float64  // Commission percentage (e.g., 30.0 for 30%)
	SuperHoldingEnabled    bool     // Whether super holding is enabled
	SuperComponentPercent   *float64 // Superannuation percentage (e.g., 10.5 for 10.5%)
	GSTRate                float64  // GST rate (e.g., 10.0 for 10%)
	GSTType                string   // "inclusive" or "exclusive"
}

// NetCalculationOutput contains calculated NET method values
type NetCalculationOutput struct {
	CommissionPercent       float64  `json:"commissionPercent"`
	Commission              float64  `json:"commission"`
	GSTOnCommission         float64  `json:"gstOnCommission"`
	TotalPaymentReceived     float64  `json:"totalPaymentReceived"`
	SuperHoldingEnabled      bool     `json:"superHoldingEnabled"`
	SuperComponentPercent    *float64 `json:"superComponentPercent,omitempty"`
	CommissionComponent      *float64 `json:"commissionComponent,omitempty"`
	SuperComponent           *float64 `json:"superComponent,omitempty"`
	TotalForReconciliation   *float64 `json:"totalForReconciliation,omitempty"`
}

// RunNetCalculation calculates NET method values based on commission and super settings
// This implements the NET calculation method for dentist commission scenarios
func RunNetCalculation(input NetCalculationInput) NetCalculationOutput {
	output := NetCalculationOutput{
		CommissionPercent:      input.CommissionPercent,
		TotalPaymentReceived:    input.TotalPaymentReceived,
		SuperHoldingEnabled:     input.SuperHoldingEnabled,
		SuperComponentPercent:   input.SuperComponentPercent,
	}

	// Step 1: Calculate commission from total payment received
	// commission = total_payment_received * (commission_percent / 100)
	commission := round2(input.TotalPaymentReceived * (input.CommissionPercent / 100))
	output.Commission = commission

	// Step 2: Calculate GST on commission
	// GST calculation depends on whether GST is inclusive or exclusive
	if input.GSTType == "inclusive" {
		// For inclusive GST: commission already includes GST
		// Extract GST: gst = commission - (commission / (1 + gst_rate/100))
		if input.GSTRate > 0 {
			gstRateDec := input.GSTRate / 100
			gstOnCommission := commission - (commission / (1 + gstRateDec))
			output.GSTOnCommission = round2(gstOnCommission)
		} else {
			output.GSTOnCommission = 0
		}
	} else {
		// For exclusive GST: GST is added on top
		// gst_on_commission = commission * (gst_rate / 100)
		if input.GSTRate > 0 {
			gstOnCommission := commission * (input.GSTRate / 100)
			output.GSTOnCommission = round2(gstOnCommission)
		} else {
			output.GSTOnCommission = 0
		}
	}

	// Step 3: Calculate commission component (commission after GST consideration)
	// If GST is inclusive, commission_component = commission (already includes GST)
	// If GST is exclusive, commission_component = commission + gst_on_commission
	var commissionComponent float64
	if input.GSTType == "inclusive" {
		commissionComponent = commission
	} else {
		commissionComponent = commission + output.GSTOnCommission
	}
	output.CommissionComponent = &commissionComponent

	// Step 4: Calculate super component if super holding is enabled
	if input.SuperHoldingEnabled && input.SuperComponentPercent != nil && *input.SuperComponentPercent > 0 {
		// super_component = commission_component * (super_component_percent / 100)
		superComponent := round2(commissionComponent * (*input.SuperComponentPercent / 100))
		output.SuperComponent = &superComponent

		// Step 5: Calculate total for reconciliation
		// total_for_reconciliation = commission_component - super_component
		totalForReconciliation := round2(commissionComponent - superComponent)
		output.TotalForReconciliation = &totalForReconciliation
	} else {
		// If super holding is not enabled, total_for_reconciliation = commission_component
		totalForReconciliation := commissionComponent
		output.TotalForReconciliation = &totalForReconciliation
	}

	return output
}

// ParseNetDeductions extracts NET method parameters from deductions JSON
func ParseNetDeductions(deductionsJSON json.RawMessage) (commissionPercent float64, superHoldingEnabled bool, superComponentPercent *float64, err error) {
	if len(deductionsJSON) == 0 {
		return 0, false, nil, nil
	}

	var deductions map[string]interface{}
	if err := json.Unmarshal(deductionsJSON, &deductions); err != nil {
		return 0, false, nil, err
	}

	// Extract commission_percent
	if cp, ok := deductions["commissionPercent"].(float64); ok {
		commissionPercent = cp
	} else if cp, ok := deductions["commission_percent"].(float64); ok {
		commissionPercent = cp
	}

	// Extract super_holding_enabled
	if she, ok := deductions["superHoldingEnabled"].(bool); ok {
		superHoldingEnabled = she
	} else if she, ok := deductions["super_holding_enabled"].(bool); ok {
		superHoldingEnabled = she
	}

	// Extract super_component_percent
	if scp, ok := deductions["superComponentPercent"].(float64); ok {
		superComponentPercent = &scp
	} else if scp, ok := deductions["super_component_percent"].(float64); ok {
		superComponentPercent = &scp
	}

	return commissionPercent, superHoldingEnabled, superComponentPercent, nil
}
