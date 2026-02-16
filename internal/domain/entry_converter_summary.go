package domain

import (
	"time"

	"github.com/google/uuid"
)

// parseSummaryFromCalc builds EntrySummary from the calculations map (field totals, BAS mapping).
func parseSummaryFromCalc(entryID uuid.UUID, createdAt, updatedAt time.Time, calc map[string]interface{}) *EntrySummary {
	summary := &EntrySummary{
		EntryID:   entryID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if calc == nil {
		return summary
	}
	summary.TotalBaseAmount = getFloat64(calc, "totalBaseAmount")
	summary.TotalGstAmount = getFloat64(calc, "totalGSTAmount")
	summary.TotalAmount = getFloat64(calc, "totalAmount")
	summary.NetPayable = getFloat64(calc, "netPayable")
	summary.NetReceivable = getFloat64(calc, "netReceivable")
	summary.NetFee = getFloat64Ptr(calc, "netFee")
	if basMap, ok := calc["basMapping"].(map[string]interface{}); ok {
		summary.BasGstOnSales1A = getFloat64(basMap, "gstOnSales1A")
		summary.BasGstCredit1B = getFloat64(basMap, "gstCredit1B")
		summary.BasTotalSalesG1 = getFloat64(basMap, "totalSalesG1")
		summary.BasExpensesG11 = getFloat64(basMap, "expensesG11")
	}
	return summary
}
