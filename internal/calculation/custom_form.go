package calculation

import (
	"encoding/json"
	"strings"

	"github.com/iamarpitzala/aca-reca-backend/util"
)

// RunEntryCalculation computes field totals, NET FEE, and deductions from form definition and raw values.
// formFieldsJSON and valuesJSON are the raw JSONB from DB; formType, formServiceFeePct from form row.
// deductionsJSON can be nil; if present it may contain serviceFacilityFeePercent and serviceFeeOverride.
// When formOutworkEnabled is true and formOutworkRatePercent > 0, expense GST is consolidated into a single outwork charge.
func RunEntryCalculation(
	formFieldsJSON []byte,
	formType string,
	formCalculationMethod string,
	formServiceFeePct *float64,
	formOutworkEnabled bool,
	formOutworkRatePercent *float64,
	valuesJSON []byte,
	deductionsJSON []byte,
) ([]byte, error) {
	// Route to gross calculation if method is GROSS
	if formCalculationMethod == util.MethodTypeGross {
		return RunGrossCalculation(
			formFieldsJSON,
			formServiceFeePct,
			formOutworkEnabled,
			formOutworkRatePercent,
			valuesJSON,
			deductionsJSON,
		)
	}
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

	fieldTotals := make([]fieldCalc, 0)
	var totalBase, totalGst, totalAmount float64
	var incomeBase, incomeGst, incomeTotal float64
	var expenseBase, expenseGst, expenseTotal float64

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
		gstType := "exclusive"
		if f.GstConfig != nil {
			rate = f.GstConfig.Rate
			gstType = f.GstConfig.Type
			if gstType == "" {
				gstType = "exclusive"
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
		if formType == util.FormTypeBoth {
			sec := getSection(f.Section)
			switch sec {
			case "expense":
				expenseBase += base
				expenseGst += gst
				expenseTotal += total
			case "reduction":
				// additional_reduction / REDUCTION: exclude from NET FEE (income − expense)
			default:
				incomeBase += base
				incomeGst += gst
				incomeTotal += total
			}
		} else {
			totalBase += base
			totalGst += gst
			totalAmount += total
		}
	}

	if formType == util.FormTypeBoth {
		totalBase = incomeBase - expenseBase
		totalGst = incomeGst - expenseGst
		totalAmount = incomeTotal - expenseTotal
	}

	bas := buildBASMapping(formType, totalBase, totalGst, totalAmount, incomeBase, incomeGst, incomeTotal, expenseBase, expenseGst, expenseTotal)

	out := calculationsOutput{
		FieldTotals:     fieldTotals,
		TotalBaseAmount: round2(totalBase),
		TotalGSTAmount:  round2(totalGst),
		TotalAmount:     round2(totalAmount),
		NetPayable:      0,
		NetReceivable:   0,
		BasMapping:      bas,
	}
	if formType == util.FormTypeExpense {
		out.NetPayable = round2(totalAmount)
	}
	if formType == util.FormTypeIncome {
		out.NetReceivable = round2(totalAmount)
	}
	if formType == util.FormTypeIncome || formType == util.FormTypeBoth {
		nf := round2(totalBase)
		out.NetFee = &nf
	}

	return json.Marshal(out)
}
