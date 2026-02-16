package calculation

import (
	"encoding/json"
	"strings"

	"github.com/iamarpitzala/aca-reca-backend/util"
)

// RunGrossCalculation computes gross method calculations following the 6-phase approach:
// Phase 1: Aggregate ALL Income Fields
// Phase 2: Aggregate ALL Expense Fields
// Phase 3: Calculate Net Amount (Once)
// Phase 4: Calculate Service & Facility Fee (Once)
// Phase 5: Aggregate ALL Remitted Fields
// Phase 6: Final Payable to Dentist
func RunGrossCalculation(
	formFieldsJSON []byte,
	formServiceFeePct *float64,
	formOutworkEnabled bool,
	formOutworkRatePercent *float64,
	valuesJSON []byte,
	deductionsJSON []byte,
) ([]byte, error) {
	var fields []calcField
	if err := json.Unmarshal(formFieldsJSON, &fields); err != nil {
		return nil, err
	}
	var values []entryValue
	if len(valuesJSON) == 0 {
		values = []entryValue{}
	} else if err := json.Unmarshal(valuesJSON, &values); err != nil {
		return nil, err
	}
	var deductions *deductionsInput
	if len(deductionsJSON) > 0 {
		deductions = &deductionsInput{}
		_ = json.Unmarshal(deductionsJSON, deductions)
	}

	// Build value lookup maps
	valueByID := make(map[string]entryValue)
	valueByName := make(map[string]entryValue)
	for _, v := range values {
		if v.FieldID != "" {
			valueByID[v.FieldID] = v
		} else if v.FieldName != "" {
			key := strings.TrimSpace(strings.ToLower(v.FieldName))
			valueByName[key] = v
		}
	}

	// Phase 1: Aggregate ALL Income Fields
	incomeExclGst := AggregateIncome(fields, valueByID, valueByName)

	// Phase 2: Aggregate ALL Expense Fields
	expenseAggregates := AggregateExpenses(fields, valueByID, valueByName, deductions)

	// Phase 3: Calculate Net Amount (Once)
	// net_amount = patient_fee_excl_gst - total_net_expenses
	netAmount := round2(incomeExclGst - expenseAggregates.TotalNetExpenses)

	// Phase 4: Calculate Service & Facility Fee (Once)
	serviceFeeCalc := CalculateServiceFee(
		incomeExclGst,
		netAmount,
		expenseAggregates.TotalNetExpenses,
		expenseAggregates.TotalExpensesGST,
		formOutworkEnabled,
		formServiceFeePct,
	)

	// Phase 5: Aggregate ALL Remitted Fields
	remittedCost := AggregateRemitted(
		fields,
		valueByID,
		valueByName,
		formOutworkEnabled,
		expenseAggregates.TotalExpensesGST,
	)

	// Phase 6: Final Payable to Dentist
	amountPayableToDentist := CalculatePayableToDentist(
		netAmount,
		serviceFeeCalc,
		remittedCost,
		expenseAggregates.PayByDentistExpenses,
		formOutworkEnabled,
	)

	// Build field totals for output
	fieldTotals := buildGrossFieldTotals(fields, valueByID, valueByName)

	// Build BAS mapping
	bas := buildGrossBASMapping(incomeExclGst, expenseAggregates)

	// Build output using standard naming conventions
	out := grossCalculationOutput{
		FieldTotals:            fieldTotals,
		IncomeExclGST:          round2(incomeExclGst),
		TotalNetExpenses:       round2(expenseAggregates.TotalNetExpenses),
		TotalExpensesGST:       round2(expenseAggregates.TotalExpensesGST),
		PayByDentistExpenses:   round2(expenseAggregates.PayByDentistExpenses),
		NetAmount:              netAmount,
		ServiceAndFacilityFee:  round2(serviceFeeCalc.ServiceAndFacilityFee),
		GSTOnServiceFee:        round2(serviceFeeCalc.GSTOnServiceFee),
		RemittedCost:           round2(remittedCost),
		AmountPayableToDentist: round2(amountPayableToDentist),
		BasMapping:             bas,
	}

	if formOutworkEnabled {
		out.LabFeeWithOtherCost = ptrFloat64(round2(serviceFeeCalc.LabFeeWithOtherCost))
		out.TotalServiceAndFacility = ptrFloat64(round2(serviceFeeCalc.TotalServiceAndFacility))
		out.TotalServiceAndFacilityIncGST = ptrFloat64(round2(serviceFeeCalc.TotalServiceAndFacilityIncGST))
	} else {
		out.TotalServiceAndFacilityFee = ptrFloat64(round2(serviceFeeCalc.TotalServiceAndFacilityFee))
	}

	return json.Marshal(out)
}

// Phase 1: Aggregate ALL Income Fields
// For each field where type = income:
//
//	IF GST = manual: net = amount - gst_amount
//	ELSE: net = amount
//	income_excl_gst += net
func AggregateIncome(fields []calcField, valueByID, valueByName map[string]entryValue) float64 {
	incomeExclGst := 0.0

	for _, f := range fields {
		sec := getSection(f.Section)
		if sec != util.FormTypeIncome {
			continue
		}

		v, ok := valueByID[f.ID]
		if !ok {
			v, ok = valueByName[strings.TrimSpace(strings.ToLower(f.Name))]
		}
		if !ok {
			continue
		}

		amount := parseFloat(v.Value)
		var net float64

		if f.GstConfig != nil && f.GstConfig.Enabled && f.GstConfig.Type == "manual" {
			manualGst := 0.0
			if v.ManualGstAmount != nil {
				manualGst = *v.ManualGstAmount
			}
			net = amount - manualGst
		} else {
			net = amount
		}

		incomeExclGst += net
	}

	return incomeExclGst
}

// ExpenseAggregates holds aggregated expense calculations
// Uses consistent naming: TotalExpensesGST (capital GST) to match TotalGSTAmount pattern
type ExpenseAggregates struct {
	TotalNetExpenses     float64
	TotalExpensesGST     float64
	PayByDentistExpenses float64
}

// Phase 2: Aggregate ALL Expense Fields
// Initialize buckets: total_net_expenses, total_expenses_gst, pay_by_dentist_expenses
// For each field where type = expense:
//
//	GST Calculation (per item):
//	  IF gst_type = inclusive: net = amount - (amount / 11), gst = amount / 11
//	  ELSE IF gst_type = exclusive: net = amount, gst = amount * 0.10
//	  ELSE IF gst_type = manual: net = amount, gst = gst_amount
//	  ELSE (gst free): net = amount, gst = 0
//	Pay-by Logic:
//	  IF pay_by = clinic: total_net_expenses += net, total_expenses_gst += gst
//	  ELSE IF pay_by = dentist: pay_by_dentist_expenses += (net + gst)
func AggregateExpenses(
	fields []calcField,
	valueByID, valueByName map[string]entryValue,
	deductions *deductionsInput,
) ExpenseAggregates {
	var aggregates ExpenseAggregates

	for _, f := range fields {
		sec := getSection(f.Section)
		if sec != util.FormTypeExpense {
			continue
		}

		v, ok := valueByID[f.ID]
		if !ok {
			v, ok = valueByName[strings.TrimSpace(strings.ToLower(f.Name))]
		}
		if !ok {
			continue
		}

		amount := parseFloat(v.Value)
		gstType := util.GSTTypeExclusive
		if f.GstConfig != nil {
			gstType = f.GstConfig.Type
			if gstType == "" {
				gstType = util.GSTTypeExclusive
			}
		}

		var net, gst float64
		if f.GstConfig != nil && f.GstConfig.Enabled {
			switch gstType {
			case util.GSTTypeInclusive:
				// net = amount - (amount / 11), gst = amount / 11
				gst = amount / 11.0
				net = amount - gst
			case util.GSTTypeExclusive:
				// net = amount, gst = amount * 0.10
				net = amount
				gst = amount * 0.10
			case util.GSTTypeManual:
				// net = amount, gst = gst_amount
				net = amount
				gst = 0.0
				if v.ManualGstAmount != nil {
					gst = *v.ManualGstAmount
				}
			default:
				// GST free: net = amount, gst = 0
				net = amount
				gst = 0
			}
		} else {
			net = amount
			gst = 0
		}

		// Determine payment responsibility
		payBy := f.PaymentResp
		if deductions != nil && deductions.EntryPaymentResponsibility != nil {
			payBy = *deductions.EntryPaymentResponsibility
		}
		if payBy == "" {
			payBy = util.PaymentResponsibilityClinic
		}

		if strings.ToUpper(payBy) == util.PaymentResponsibilityClinic {
			aggregates.TotalNetExpenses += net
			aggregates.TotalExpensesGST += gst
		} else {
			aggregates.PayByDentistExpenses += (net + gst)
		}
	}

	return aggregates
}

// ServiceFeeCalculation holds service fee calculation results
// Uses consistent naming: GST (capital) for consistency with TotalGSTAmount
type ServiceFeeCalculation struct {
	ServiceAndFacilityFee         float64
	LabFeeWithOtherCost           float64
	TotalServiceAndFacility       float64
	GSTOnServiceFee               float64
	TotalServiceAndFacilityIncGST float64
	TotalServiceAndFacilityFee    float64
}

// Phase 4: Calculate Service & Facility Fee (Once)
// If outward_cost = true: service_and_facility_fee = patient_fee_excl_gst * 0.60
// If outward_cost = false: service_and_facility_fee = net_amount * 0.60
func CalculateServiceFee(
	incomeExclGst,
	netAmount,
	totalNetExpenses,
	totalExpensesGst float64,
	outwardCost bool,
	serviceFeePct *float64,
) ServiceFeeCalculation {
	var calc ServiceFeeCalculation

	// Default to 0.60 (60%) as per requirements, but allow override
	serviceFeeRate := 0.60
	if serviceFeePct != nil {
		serviceFeeRate = *serviceFeePct / 100.0
	}

	if outwardCost {
		// If outward_cost = true
		// service_and_facility_fee = patient_fee_excl_gst * 0.60
		calc.ServiceAndFacilityFee = incomeExclGst * serviceFeeRate
		// lab_fee_with_other_cost = (total_net_expenses + total_expenses_gst) * 0.40
		calc.LabFeeWithOtherCost = (totalNetExpenses + totalExpensesGst) * 0.40
		// total_service_and_facility = service_and_facility_fee + lab_fee_with_other_cost
		calc.TotalServiceAndFacility = calc.ServiceAndFacilityFee + calc.LabFeeWithOtherCost
		// gst_on_service_fee = total_service_and_facility * 0.10
		calc.GSTOnServiceFee = calc.TotalServiceAndFacility * 0.10
		// total_service_and_facility_inc_gst = total_service_and_facility + gst_on_service_fee
		calc.TotalServiceAndFacilityIncGST = calc.TotalServiceAndFacility + calc.GSTOnServiceFee
	} else {
		// If outward_cost = false
		// service_and_facility_fee = net_amount * 0.60
		calc.ServiceAndFacilityFee = incomeExclGst * serviceFeeRate
		// gst_on_service_fee = service_and_facility_fee * 0.10
		calc.GSTOnServiceFee = calc.ServiceAndFacilityFee * 0.10
		// total_service_and_facility_fee = service_and_facility_fee + gst_on_service_fee
		calc.TotalServiceAndFacilityFee = calc.ServiceAndFacilityFee + calc.GSTOnServiceFee
	}

	return calc
}

// Phase 5: Aggregate ALL Remitted Fields
// Initialize: remitted_cost = 0
// For each field where type = remitted AND outward_cost = false:
//
//	IF gst_type = inclusive: remitted = amount - (amount / 11)
//	ELSE IF gst_type = exclusive: remitted = amount + (amount * 0.10)
//	ELSE IF gst_type = manual: remitted = amount + gst_amount
//	ELSE: remitted = amount
//	remitted_cost += remitted
//
// Also add clinic-paid GST from Step 2
func AggregateRemitted(
	fields []calcField,
	valueByID, valueByName map[string]entryValue,
	outwardCost bool,
	clinicPaidGst float64,
) float64 {
	remittedCost := 0.0

	// Only process remitted fields when outward_cost = false
	if outwardCost {
		// Add clinic-paid GST from Step 2
		return clinicPaidGst
	}

	for _, f := range fields {
		sec := getSection(f.Section)
		if sec != "reduction" {
			continue
		}
		if f.Type != "number" && f.Type != "currency" {
			continue
		}
		if !f.IncludeInTotal {
			continue
		}

		v, ok := valueByID[f.ID]
		if !ok {
			v, ok = valueByName[strings.TrimSpace(strings.ToLower(f.Name))]
		}
		if !ok {
			continue
		}

		amount := parseFloat(v.Value)
		gstType := util.GSTTypeExclusive
		if f.GstConfig != nil {
			gstType = f.GstConfig.Type
			if gstType == "" {
				gstType = util.GSTTypeExclusive
			}
		}

		var remitted float64
		if f.GstConfig != nil && f.GstConfig.Enabled {
			switch gstType {
			case util.GSTTypeInclusive:
				remitted = amount - (amount / 11.0)
			case util.GSTTypeExclusive:
				remitted = amount + (amount * 0.10)
			case util.GSTTypeManual:
				manualGst := 0.0
				if v.ManualGstAmount != nil {
					manualGst = *v.ManualGstAmount
				}
				remitted = amount + manualGst
			default:
				remitted = amount
			}
		} else {
			remitted = amount
		}

		remittedCost += remitted
	}

	// Add clinic-paid GST from Step 2
	remittedCost += clinicPaidGst

	return remittedCost
}

// Phase 6: Final Payable to Dentist
// If outward_cost = true:
//
//	amount_payable_to_dentist = net_amount - total_service_and_facility_inc_gst
//
// If outward_cost = false:
//
//	amount_payable_to_dentist = net_amount - (total_service_and_facility_fee + remitted_cost) + pay_by_dentist_expenses
func CalculatePayableToDentist(
	netAmount float64,
	serviceFeeCalc ServiceFeeCalculation,
	remittedCost float64,
	payByDentistExpenses float64,
	outwardCost bool,
) float64 {
	if outwardCost {
		// If outward_cost = true
		return netAmount - serviceFeeCalc.TotalServiceAndFacilityIncGST
	}

	// If outward_cost = false
	return netAmount - (serviceFeeCalc.TotalServiceAndFacilityFee + remittedCost) + payByDentistExpenses
}

// buildGrossFieldTotals builds field totals for gross calculation output
func buildGrossFieldTotals(
	fields []calcField,
	valueByID, valueByName map[string]entryValue,
) []fieldCalc {
	fieldTotals := make([]fieldCalc, 0)

	for _, f := range fields {
		if f.Type != "number" && f.Type != "currency" {
			continue
		}
		if !f.IncludeInTotal {
			continue
		}

		v, ok := valueByID[f.ID]
		if !ok {
			v, ok = valueByName[strings.TrimSpace(strings.ToLower(f.Name))]
		}
		if !ok {
			continue
		}

		numVal := parseFloat(v.Value)
		var manual *float64
		if v.ManualGstAmount != nil {
			manual = v.ManualGstAmount
		}

		rate := 0.0
		gstType := util.GSTTypeExclusive
		if f.GstConfig != nil {
			rate = f.GstConfig.Rate
			gstType = f.GstConfig.Type
			if gstType == "" {
				gstType = util.GSTTypeExclusive
			}
		}

		var base, gst, total float64
		if f.GstConfig != nil && f.GstConfig.Enabled {
			base, gst, total = calcGST(numVal, rate, gstType, manual)
		} else {
			base, gst, total = numVal, 0, numVal
		}

		fieldTotals = append(fieldTotals, fieldCalc{
			FieldID:     f.ID,
			FieldName:   f.Name,
			BaseAmount:  base,
			GstAmount:   gst,
			TotalAmount: total,
			GstRate:     rate,
			GstType:     gstType,
		})
	}

	return fieldTotals
}

// buildGrossBASMapping builds BAS mapping for gross calculation
// Uses standard BAS field names matching calculationsOutput
func buildGrossBASMapping(
	incomeExclGst float64,
	expenseAggregates ExpenseAggregates,
) basMapping {
	return basMapping{
		GstOnSales1A: round2(expenseAggregates.TotalExpensesGST),
		GstCredit1B:  0,
		TotalSalesG1: round2(incomeExclGst),
		ExpensesG11:  round2(expenseAggregates.TotalNetExpenses),
	}
}
