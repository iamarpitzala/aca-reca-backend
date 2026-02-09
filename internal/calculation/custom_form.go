package calculation

import (
	"encoding/json"
	"math"
	"strconv"
)

// Field definitions for parsing form.Fields JSONB (subset needed for calculation)
type calcField struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Section        string  `json:"section"`
	IncludeInTotal bool    `json:"includeInTotal"`
	GstConfig      *gstCfg `json:"gstConfig"`
	PaymentResp    string  `json:"paymentResponsibility"`
}
type gstCfg struct {
	Enabled bool    `json:"enabled"`
	Rate    float64 `json:"rate"`
	Type    string  `json:"type"` // inclusive, exclusive, manual
}

// Value from client: fieldId, value, optional manualGstAmount
type entryValue struct {
	FieldID        string      `json:"fieldId"`
	FieldName      string      `json:"fieldName"`
	Value          interface{} `json:"value"`
	ManualGstAmount *float64   `json:"manualGstAmount"`
}

// Deductions from request (for service fee % and override)
type deductionsInput struct {
	ServiceFacilityFeePercent *float64 `json:"serviceFacilityFeePercent"`
	ServiceFeeOverride        *float64 `json:"serviceFeeOverride"`
}

// Output structures (match frontend EntryCalculations)
type fieldCalc struct {
	FieldID     string  `json:"fieldId"`
	FieldName   string  `json:"fieldName"`
	BaseAmount  float64 `json:"baseAmount"`
	GstAmount   float64 `json:"gstAmount"`
	TotalAmount float64 `json:"totalAmount"`
	GstRate     float64 `json:"gstRate"`
	GstType     string  `json:"gstType"`
}
type basMapping struct {
	GstOnSales1A float64 `json:"gstOnSales1A"`
	GstCredit1B  float64 `json:"gstCredit1B"`
	TotalSalesG1 float64 `json:"totalSalesG1"`
	ExpensesG11  float64 `json:"expensesG11"`
}
type calculationsOutput struct {
	FieldTotals               []fieldCalc   `json:"fieldTotals"`
	TotalBaseAmount           float64       `json:"totalBaseAmount"`
	TotalGSTAmount            float64       `json:"totalGSTAmount"`
	TotalAmount               float64       `json:"totalAmount"`
	NetPayable                float64       `json:"netPayable"`
	NetReceivable             float64       `json:"netReceivable"`
	BasMapping                basMapping    `json:"basMapping"`
	NetFee                    *float64      `json:"netFee,omitempty"`
	ServiceFeeBase            *float64      `json:"serviceFeeBase,omitempty"`
	GstOnServiceFee           *float64      `json:"gstOnServiceFee,omitempty"`
	TotalServiceFee           *float64      `json:"totalServiceFee,omitempty"`
	TotalReductions           *float64      `json:"totalReductions,omitempty"`
	TotalReimbursements       *float64      `json:"totalReimbursements,omitempty"`
	ReductionBreakdown        []fieldCalc   `json:"reductionBreakdown,omitempty"`
	ReimbursementBreakdown    []fieldCalc   `json:"reimbursementBreakdown,omitempty"`
	SubtotalAfterDeductions   *float64      `json:"subtotalAfterDeductions,omitempty"`
	RemittedAmount            *float64      `json:"remittedAmount,omitempty"`
}

func round2(n float64) float64 { return math.Round(n*100) / 100 }

// parseFloat parses a numeric value from JSON (number or string).
func parseFloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	}
	return 0
}

func calcGST(amount, rate float64, gstType string, manualGst *float64) (base, gst, total float64) {
	if gstType == "manual" {
		m := 0.0
		if manualGst != nil {
			m = *manualGst
		}
		return round2(amount), round2(m), round2(amount+m)
	}
	if rate == 0 {
		return amount, 0, amount
	}
	rateDec := rate / 100
	if gstType == "inclusive" {
		gst = amount - (amount / (1 + rateDec))
		base = amount - gst
		return round2(base), round2(gst), amount
	}
	gst = amount * rateDec
	total = amount + gst
	return amount, round2(gst), round2(total)
}

func getSection(s string) string {
	if len(s) >= 7 && (s[0]|32) == 'e' && (s[1]|32) == 'x' && (s[2]|32) == 'p' {
		return "expense"
	}
	return "income"
}

const defaultServiceFeePct = 50.0

// RunEntryCalculation computes field totals, NET FEE, and deductions from form definition and raw values.
// formFieldsJSON and valuesJSON are the raw JSONB from DB; formType, formCalculationMethod, formServiceFeePct from form row.
// deductionsJSON can be nil; if present it may contain serviceFacilityFeePercent and serviceFeeOverride.
func RunEntryCalculation(
	formFieldsJSON []byte,
	formType string,
	formServiceFeePct *float64,
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

	fieldTotals := make([]fieldCalc, 0)
	var totalBase, totalGst, totalAmount float64
	var incomeBase, incomeGst, incomeTotal float64
	var expenseBase, expenseGst, expenseTotal float64

	valueByID := make(map[string]entryValue)
	for _, v := range values {
		if v.FieldID != "" {
			valueByID[v.FieldID] = v
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
		if formType == "both" {
			sec := getSection(f.Section)
			if sec == "expense" {
				expenseBase += base
				expenseGst += gst
				expenseTotal += total
			} else {
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

	if formType == "both" {
		totalBase = incomeBase - expenseBase
		totalGst = incomeGst - expenseGst
		totalAmount = incomeTotal - expenseTotal
	}

	var bas basMapping
	switch formType {
	case "income":
		bas = basMapping{GstOnSales1A: totalGst, TotalSalesG1: totalAmount}
	case "expense":
		bas = basMapping{GstCredit1B: totalGst, ExpensesG11: totalBase}
	default:
		bas = basMapping{
			GstOnSales1A: round2(incomeGst),
			GstCredit1B:  round2(expenseGst),
			TotalSalesG1: round2(incomeTotal),
			ExpensesG11:  round2(expenseBase),
		}
	}

	out := calculationsOutput{
		FieldTotals:     fieldTotals,
		TotalBaseAmount: round2(totalBase),
		TotalGSTAmount:  round2(totalGst),
		TotalAmount:     round2(totalAmount),
		NetPayable:      0,
		NetReceivable:   0,
		BasMapping:      bas,
	}
	if formType == "expense" {
		out.NetPayable = round2(totalAmount)
	}
	if formType == "income" {
		out.NetReceivable = round2(totalAmount)
	}
	if formType == "income" || formType == "both" {
		nf := round2(totalBase)
		out.NetFee = &nf
	}

	// Deductions (income/both)
	hasIncome := formType == "income" || formType == "both"
	pct := defaultServiceFeePct
	if deductions != nil && deductions.ServiceFacilityFeePercent != nil && *deductions.ServiceFacilityFeePercent > 0 {
		pct = *deductions.ServiceFacilityFeePercent
	} else if formServiceFeePct != nil && *formServiceFeePct > 0 {
		pct = *formServiceFeePct
	}
	if hasIncome && pct > 0 {
		netFee := out.TotalBaseAmount
		serviceBase := netFee * (pct / 100.0)
		if deductions != nil && deductions.ServiceFeeOverride != nil {
			serviceBase = *deductions.ServiceFeeOverride
		}
		serviceBase = round2(serviceBase)
		gstOnSvc := round2(serviceBase * 0.1)
		totalSvc := round2(serviceBase + gstOnSvc)
		out.ServiceFeeBase = &serviceBase
		out.GstOnServiceFee = &gstOnSvc
		out.TotalServiceFee = &totalSvc

		// Additional Reductions / Reimbursements: expense-section fields when paid by clinic/owner.
		// Include each expense's GST (manual, inclusive, or exclusive) in the breakdown so Additional Reductions shows base, GST, and total per expense.
		var totalRed, totalReimb float64
		var redBreak, reimbBreak []fieldCalc
		for _, f := range fields {
			sec := getSection(f.Section)
			if sec != "expense" {
				continue // only take expense-section fields (and their GST)
			}
			var ft *fieldCalc
			for i := range fieldTotals {
				if fieldTotals[i].FieldID == f.ID {
					ft = &fieldTotals[i]
					break
				}
			}
			if ft == nil {
				continue
			}
			payResp := f.PaymentResp
			if payResp == "" {
				payResp = "owner"
			}
			if payResp == "clinic" {
				totalRed += ft.TotalAmount
				redBreak = append(redBreak, *ft)
			} else {
				totalReimb += ft.TotalAmount
				reimbBreak = append(reimbBreak, *ft)
			}
		}
		totalRed = round2(totalRed)
		totalReimb = round2(totalReimb)
		out.TotalReductions = &totalRed
		out.TotalReimbursements = &totalReimb
		out.ReductionBreakdown = redBreak
		out.ReimbursementBreakdown = reimbBreak
		sub := round2(netFee - serviceBase)
		var totalRedGst float64
		for i := range redBreak {
			totalRedGst += redBreak[i].GstAmount
		}
		totalRedGst = round2(totalRedGst)
		rem := round2(netFee - totalSvc + totalReimb - totalRedGst)
		out.SubtotalAfterDeductions = &sub
		out.RemittedAmount = &rem
	}

	return json.Marshal(out)
}
