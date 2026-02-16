package calculation

import (
	"encoding/json"
	"math"
	"strings"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

// NetCalculationInput holds input for NET method calculation.
// NetAmount is computed first (income - expenses by section); CommissionPercent is owner %.
type NetCalculationInput struct {
	NetAmount             float64  // from CalculateNetAmountBySection
	CommissionPercent     float64  // owner %
	SuperHoldingEnabled   bool
	SuperComponentPercent *float64
	GSTRate               float64
	GSTType               string
}

// NetCalculationOutput holds output from RunNetCalculation.
type NetCalculationOutput struct {
	CommissionPercent       float64
	Commission              float64
	GSTOnCommission         float64
	TotalPaymentReceived    float64
	SuperHoldingEnabled     bool
	SuperComponentPercent   *float64
	CommissionComponent     *float64
	SuperComponent          *float64
	TotalForReconciliation  *float64
}

const gstRateOnCommission = 0.1 // 10% GST on commission

// RunNetCalculation computes NET method details after net amount is known.
// If super holding enabled: commission_component = netAmount/(1+super%), super_component = commission_component*super%,
// total_for_reconciliation = commission_component + super_component, GST = commission_component*0.1, total_payment_received = commission_component + GST.
// Else: commission = net*owner%, GST = commission*0.1, total_payment_received = commission + GST; commission_component/super/total_for_reconciliation = 0.
func RunNetCalculation(input NetCalculationInput) NetCalculationOutput {
	out := NetCalculationOutput{
		CommissionPercent:     input.CommissionPercent,
		SuperHoldingEnabled:   input.SuperHoldingEnabled,
		SuperComponentPercent: input.SuperComponentPercent,
	}

	// Round net amount to 2 decimals first so 100.00 stays 100.00 (avoids float drift giving 39.97)
	netAmount := round2(input.NetAmount)
	// Commission = net * owner%; round with half-up so 100*40% = 40.00 exactly
	commission := round2HalfUp(netAmount * (input.CommissionPercent / 100))

	if input.SuperHoldingEnabled && input.SuperComponentPercent != nil && *input.SuperComponentPercent > 0 {
		// Super holding enabled: derive from commission (rounded)
		superPct := *input.SuperComponentPercent / 100
		// commission_component = commission / (1 + super_component%) — round each step
		commissionComponent := round2(commission / (1 + superPct))
		superComponent := round2(commissionComponent * superPct)
		totalForReconciliation := round2(commissionComponent + superComponent)
		gstOnCommission := round2(commissionComponent * gstRateOnCommission)
		totalPaymentReceived := round2(commissionComponent + gstOnCommission)

		out.CommissionComponent = ptrFloat64(commissionComponent)
		out.SuperComponent = ptrFloat64(superComponent)
		out.TotalForReconciliation = ptrFloat64(totalForReconciliation)
		out.GSTOnCommission = gstOnCommission
		out.TotalPaymentReceived = totalPaymentReceived
		out.Commission = totalPaymentReceived
	} else {
		// Else: commission = net*owner% (already rounded), GST = commission*0.1, total = commission + GST
		gstOnCommission := round2(commission * gstRateOnCommission)
		totalPaymentReceived := round2(commission + gstOnCommission)

		out.Commission = commission
		out.GSTOnCommission = gstOnCommission
		out.TotalPaymentReceived = totalPaymentReceived
		out.CommissionComponent = ptrFloat64(0)
		out.SuperComponent = ptrFloat64(0)
		out.SuperComponentPercent = ptrFloat64(0)
		out.TotalForReconciliation = ptrFloat64(0)
	}

	return out
}

// NetAmountResult holds IncomeExclGST and NetAmount from field values.
type NetAmountResult struct {
	IncomeExclGST float64
	NetAmount     float64
}

// CalculateNetAmountBySection computes net amount for NET method (create entry side).
// Uses section type already defined on each field: INCOME vs EXPENSE.
// Loops over field values, sums amounts by section (amount = TotalAmount or Value already received),
// then net amount = total income - total expenses.
func CalculateNetAmountBySection(
	fieldValueResponses []domain.EntryFieldValueResponse,
	fields []domain.CustomFormField,
) NetAmountResult {
	fieldByID := make(map[string]domain.CustomFormField)
	for _, f := range fields {
		fieldByID[f.ID.String()] = f
	}

	// Sum in integer cents to avoid float drift (e.g. 33.33+33.33+33.34 must equal 100.00)
	var totalIncomeCents, totalExpensesCents int64
	for _, val := range fieldValueResponses {
		f, ok := fieldByID[val.FieldID]
		if !ok {
			continue
		}
		sec := getSection(f.Section)
		amount := 0.0
		if val.TotalAmount != nil {
			amount = *val.TotalAmount
		} else {
			amount = val.Value
		}
		amountCents := int64(math.Round(amount * 100))

		switch sec {
		case "income":
			totalIncomeCents += amountCents
		case "expense":
			totalExpensesCents += amountCents
		}
	}

	netCents := totalIncomeCents - totalExpensesCents
	// Convert back to dollars with exact 2-decimal values (no float drift)
	totalIncome := float64(totalIncomeCents) / 100
	netAmount := float64(netCents) / 100
	return NetAmountResult{
		IncomeExclGST: totalIncome,
		NetAmount:     netAmount,
	}
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
