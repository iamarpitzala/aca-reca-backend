package domain

import (
	"time"

	"github.com/google/uuid"
)

// parseGrossDetailsFromCalc builds EntryGrossDetails when form uses GROSS and calc has serviceFeeBase.
func parseGrossDetailsFromCalc(entry *CustomFormEntry, form *CustomForm, calc map[string]interface{}, deductions *EntryDeductions) *EntryGrossDetails {
	if entry == nil || form == nil || calc == nil || form.CalculationMethod != string(CalcMethodGross) {
		return nil
	}
	if _, hasServiceFeeBase := calc["serviceFeeBase"]; !hasServiceFeeBase {
		return nil
	}
	gd := &EntryGrossDetails{
		EntryID:                 entry.ID,
		ServiceFeeBase:          getFloat64(calc, "serviceFeeBase"),
		GstOnServiceFee:         getFloat64(calc, "gstOnServiceFee"),
		TotalServiceFee:         getFloat64(calc, "totalServiceFee"),
		SubtotalAfterDeductions: getFloat64Ptr(calc, "subtotalAfterDeductions"),
		RemittedAmount:          getFloat64Ptr(calc, "remittedAmount"),
		CreatedAt:               entry.CreatedAt,
		UpdatedAt:               entry.UpdatedAt,
	}
	if deductions != nil && deductions.ServiceFacilityFeePercent != nil {
		gd.ServiceFacilityFeePercent = *deductions.ServiceFacilityFeePercent
	} else if form.ServiceFacilityFeePercent != nil && *form.ServiceFacilityFeePercent > 0 {
		gd.ServiceFacilityFeePercent = *form.ServiceFacilityFeePercent
	} else {
		gd.ServiceFacilityFeePercent = 50
	}
	return gd
}

// parseGrossReductionsSummaryFromCalc builds EntryGrossReductionsSummary when any summary amount is present.
func parseGrossReductionsSummaryFromCalc(entry *CustomFormEntry, calc map[string]interface{}) *EntryGrossReductionsSummary {
	if entry == nil || calc == nil {
		return nil
	}
	totalReductions := getFloat64(calc, "totalReductions")
	totalReimbursements := getFloat64(calc, "totalReimbursements")
	totalAdditionalReduction := getFloat64(calc, "totalAdditionalReduction")
	if totalReductions != 0 || getFloat64(calc, "totalReductionBase") != 0 || totalReimbursements != 0 || totalAdditionalReduction != 0 {
		return &EntryGrossReductionsSummary{
			EntryID:                       entry.ID,
			TotalReductions:               totalReductions,
			TotalReductionBase:            getFloat64(calc, "totalReductionBase"),
			TotalExpenseGst:               getFloat64(calc, "totalExpenseGst"),
			TotalReimbursements:           totalReimbursements,
			TotalAdditionalReduction:      totalAdditionalReduction,
			TotalAdditionalReductionBase:  getFloat64(calc, "totalAdditionalReductionBase"),
			TotalAdditionalReductionGst:   getFloat64(calc, "totalAdditionalReductionGst"),
			CreatedAt:                     entry.CreatedAt,
			UpdatedAt:                     entry.UpdatedAt,
		}
	}
	return nil
}

// parseGrossOutworkFromCalc builds EntryGrossOutwork when outwork is enabled or any charge is non-zero.
func parseGrossOutworkFromCalc(entry *CustomFormEntry, calc map[string]interface{}) *EntryGrossOutwork {
	if entry == nil || calc == nil {
		return nil
	}
	outworkChargeBase := getFloat64(calc, "outworkChargeBase")
	outworkChargeTotal := getFloat64(calc, "outworkChargeTotal")
	outworkEnabled := getBool(calc, "outworkEnabled")
	if outworkEnabled || outworkChargeBase != 0 || outworkChargeTotal != 0 {
		return &EntryGrossOutwork{
			EntryID:            entry.ID,
			OutworkEnabled:     outworkEnabled,
			OutworkRatePercent: getFloat64Ptr(calc, "outworkRatePercent"),
			OutworkChargeBase:  outworkChargeBase,
			OutworkChargeGst:   getFloat64(calc, "outworkChargeGst"),
			OutworkChargeTotal: outworkChargeTotal,
			CreatedAt:          entry.CreatedAt,
			UpdatedAt:          entry.UpdatedAt,
		}
	}
	return nil
}

func parseBreakdownToGrossReduction(entryID uuid.UUID, displayOrder int, createdAt time.Time, m map[string]interface{}) EntryGrossReduction {
	return EntryGrossReduction{
		EntryID:      entryID,
		FieldID:      getString(m, "fieldId"),
		FieldName:    getString(m, "fieldName"),
		BaseAmount:   getFloat64(m, "baseAmount"),
		GstAmount:    getFloat64(m, "gstAmount"),
		TotalAmount:  getFloat64(m, "totalAmount"),
		DisplayOrder: displayOrder,
		CreatedAt:    createdAt,
	}
}

func parseBreakdownToGrossReimbursement(entryID uuid.UUID, displayOrder int, createdAt time.Time, m map[string]interface{}) EntryGrossReimbursement {
	return EntryGrossReimbursement{
		EntryID:      entryID,
		FieldID:      getString(m, "fieldId"),
		FieldName:    getString(m, "fieldName"),
		BaseAmount:   getFloat64(m, "baseAmount"),
		GstAmount:    getFloat64(m, "gstAmount"),
		TotalAmount:  getFloat64(m, "totalAmount"),
		DisplayOrder: displayOrder,
		CreatedAt:    createdAt,
	}
}

func parseBreakdownToGrossAdditionalReduction(entryID uuid.UUID, displayOrder int, createdAt time.Time, m map[string]interface{}) EntryGrossAdditionalReduction {
	return EntryGrossAdditionalReduction{
		EntryID:      entryID,
		FieldID:      getString(m, "fieldId"),
		FieldName:    getString(m, "fieldName"),
		BaseAmount:   getFloat64(m, "baseAmount"),
		GstAmount:    getFloat64(m, "gstAmount"),
		TotalAmount:  getFloat64(m, "totalAmount"),
		DisplayOrder: displayOrder,
		CreatedAt:    createdAt,
	}
}
