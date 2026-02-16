package calculation

import (
	"strings"

	"github.com/iamarpitzala/aca-reca-backend/util"
)

// getSection returns "expense", "income", or "reduction". REDUCTION (additional_reduction) is excluded from income/expense totals.
func getSection(s string) string {
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "expense") {
		return "expense"
	}
	if lower == "reduction" || strings.HasPrefix(lower, "additional_reduction") {
		return "reduction"
	}
	return "income"
}

const defaultServiceFeePct = 50.0

func buildBASMapping(formType string, totalBase, totalGst, totalAmount, incomeBase, incomeGst, incomeTotal, expenseBase, expenseGst, expenseTotal float64) basMapping {
	switch formType {
	case util.FormTypeIncome:
		return basMapping{GstOnSales1A: totalGst, TotalSalesG1: totalAmount}
	case util.FormTypeExpense:
		return basMapping{GstCredit1B: totalGst, ExpensesG11: totalBase}
	default:
		return basMapping{
			GstOnSales1A: round2(incomeGst),
			GstCredit1B:  round2(expenseGst),
			TotalSalesG1: round2(incomeTotal),
			ExpensesG11:  round2(expenseBase),
		}
	}
}

// resolveCalculationMethod returns whether net and/or gross method apply from form and deductions.
func resolveCalculationMethod(formCalculationMethod string, deductions *deductionsInput, formServiceFeePct *float64, hasIncome bool) (isNet, isGross bool) {
	switch formCalculationMethod {
	case util.MethodTypeNet:
		isNet = true
	case util.MethodTypeGross:
		isGross = true
	default:
		calcMethodLower := strings.ToLower(formCalculationMethod)
		switch calcMethodLower {
		case util.MethodTypeNet:
			isNet = true
		case util.MethodTypeGross:
			isGross = true
		default:
			isNet = true
		}
	}
	if deductions != nil {
		if deductions.CommissionPercent != nil && *deductions.CommissionPercent > 0 {
			isNet, isGross = true, false
		} else if deductions.ServiceFacilityFeePercent != nil && *deductions.ServiceFacilityFeePercent > 0 && !isNet {
			isGross = true
		}
	}
	if !isNet && !isGross {
		if formServiceFeePct != nil && *formServiceFeePct > 0 {
			isGross = true
		} else if hasIncome {
			isGross = true
		}
	}
	return isNet, isGross
}
