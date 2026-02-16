package calculation

import (
	"encoding/json"
	"strings"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

// NetCalculationInput holds input for NET method calculation.
type NetCalculationInput struct {
	TotalPaymentReceived  float64
	CommissionPercent     float64
	SuperHoldingEnabled   bool
	SuperComponentPercent *float64
	GSTRate               float64
	GSTType               string
}

// NetCalculationOutput holds output from RunNetCalculation.
type NetCalculationOutput struct {
	CommissionPercent     float64
	Commission            float64
	GSTOnCommission       float64
	TotalPaymentReceived  float64
	SuperHoldingEnabled   bool
	SuperComponentPercent *float64
	CommissionComponent   *float64
	SuperComponent        *float64
	TotalForReconciliation *float64
}

// RunNetCalculation computes NET method details: commission, GST on commission, super components.
func RunNetCalculation(input NetCalculationInput) NetCalculationOutput {
	out := NetCalculationOutput{
		CommissionPercent:    input.CommissionPercent,
		TotalPaymentReceived: round2(input.TotalPaymentReceived),
		SuperHoldingEnabled:  input.SuperHoldingEnabled,
		SuperComponentPercent: input.SuperComponentPercent,
	}

	// commission = total_payment_received × (commission_percent / 100)
	commission := input.TotalPaymentReceived * (input.CommissionPercent / 100)

	// GST on commission based on type
	gstType := strings.ToLower(input.GSTType)
	rate := input.GSTRate / 100
	if gstType == util.GSTTypeInclusive {
		// gst_on_commission = commission × (gst_rate / (100 + gst_rate))
		out.GSTOnCommission = commission * (input.GSTRate / (100 + input.GSTRate))
		out.Commission = round2(commission)
	} else {
		// exclusive or default
		out.GSTOnCommission = round2(commission * rate)
		out.Commission = round2(commission)
	}

	// commission_component = commission - gst_on_commission (net to dentist before super)
	commissionComponent := out.Commission - out.GSTOnCommission
	out.CommissionComponent = ptrFloat64(round2(commissionComponent))

	if input.SuperHoldingEnabled && input.SuperComponentPercent != nil {
		superPct := *input.SuperComponentPercent / 100
		superComp := commissionComponent * superPct
		out.SuperComponent = ptrFloat64(round2(superComp))
		// total_for_reconciliation = commission_component - super_component
		recon := commissionComponent - superComp
		out.TotalForReconciliation = ptrFloat64(round2(recon))
	} else {
		out.TotalForReconciliation = ptrFloat64(round2(commissionComponent))
	}

	return out
}

// NetAmountResult holds IncomeExclGST and NetAmount from field values.
type NetAmountResult struct {
	IncomeExclGST float64
	NetAmount     float64
}

// CalculateNetAmountFromFieldValues computes net income, net expenses, and net amount from field value responses.
func CalculateNetAmountFromFieldValues(
	fieldValueResponses []domain.EntryFieldValueResponse,
	fields []domain.CustomFormField,
) NetAmountResult {
	fieldByID := make(map[string]domain.CustomFormField)
	for _, f := range fields {
		fieldByID[f.ID.String()] = f
	}

	// Parse metadata for includeInTotal and gstConfig
	fieldMeta := make(map[string]struct {
		IncludeInTotal bool
		GSTEnabled     bool
		GSTRate        float64
		GSTType        string
	})
	for _, f := range fields {
		enabled := f.GSTConfig
		rate := 0.0
		if f.GSTRate != nil {
			rate = *f.GSTRate
		}
		gstType := strings.ToLower(f.GSTType)
		if gstType == "" {
			gstType = util.GSTTypeExclusive
		}
		includeInTotal := false
		if len(f.Metadata) > 0 {
			var meta map[string]interface{}
			_ = json.Unmarshal(f.Metadata, &meta)
			if meta != nil {
				if it, ok := meta["includeInTotal"].(bool); ok {
					includeInTotal = it
				}
			}
		}
		fieldMeta[f.ID.String()] = struct {
			IncludeInTotal bool
			GSTEnabled     bool
			GSTRate        float64
			GSTType        string
		}{IncludeInTotal: includeInTotal, GSTEnabled: enabled, GSTRate: rate, GSTType: gstType}
	}

	var totalNetIncome, totalNetExpenses float64
	for _, val := range fieldValueResponses {
		f, ok := fieldByID[val.FieldID]
		if !ok {
			continue
		}
		meta := fieldMeta[val.FieldID]
		if !meta.IncludeInTotal {
			continue
		}
		sec := getSection(f.Section)

		amount := val.Value
		if val.TotalAmount != nil {
			amount = *val.TotalAmount
		}

		var net float64
		if meta.GSTEnabled && meta.GSTType == util.GSTTypeManual {
			manualGst := 0.0
			if val.ManualGSTAmount != nil {
				manualGst = *val.ManualGSTAmount
			}
			net = amount - manualGst
		} else if meta.GSTEnabled && meta.GSTType == util.GSTTypeInclusive {
			gst := amount / (1 + meta.GSTRate/100)
			net = amount - gst
		} else {
			net = amount
		}

		switch sec {
		case util.FormTypeIncome:
			totalNetIncome += net
		case util.FormTypeExpense:
			totalNetExpenses += net
		}
	}

	return NetAmountResult{
		IncomeExclGST: round2(totalNetIncome),
		NetAmount:     round2(totalNetIncome - totalNetExpenses),
	}
}
