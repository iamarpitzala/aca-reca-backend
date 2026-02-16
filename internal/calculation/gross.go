package calculation

import (
	"strings"

	"github.com/iamarpitzala/aca-reca-backend/util"
)

func applyGrossMethodCalculations(out *calculationsOutput, fields []calcField, fieldTotals []fieldCalc, deductions *deductionsInput, formServiceFeePct *float64, formOutworkEnabled bool, formOutworkRatePercent *float64) {
	pct := defaultServiceFeePct
	if deductions != nil && deductions.ServiceFacilityFeePercent != nil && *deductions.ServiceFacilityFeePercent > 0 {
		pct = *deductions.ServiceFacilityFeePercent
	} else if formServiceFeePct != nil && *formServiceFeePct > 0 {
		pct = *formServiceFeePct
	}
	if pct <= 0 {
		return
	}
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

	var totalRedGst, totalReimb float64
	var redBreak, reimbBreak []fieldCalc
	entryPayResp := ""
	if deductions != nil && deductions.EntryPaymentResponsibility != nil {
		entryPayResp = strings.ToLower(*deductions.EntryPaymentResponsibility)
	}
	for _, f := range fields {
		if getSection(f.Section) != util.SectionExpense {
			continue
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
		payResp := entryPayResp
		if payResp == "" {
			payResp = f.PaymentResp
		}
		if payResp == "" {
			payResp = util.PaymentResponsibilityOwner
		}
		if payResp == util.PaymentResponsibilityClinic {
			totalRedGst += ft.GstAmount
			redBreak = append(redBreak, *ft)
		} else {
			totalReimb += ft.TotalAmount
			reimbBreak = append(reimbBreak, *ft)
		}
	}
	totalRedGst = round2(totalRedGst)
	totalReimb = round2(totalReimb)
	out.TotalReimbursements = &totalReimb
	out.ReductionBreakdown = redBreak
	out.ReimbursementBreakdown = reimbBreak

	var totalReductionBase float64
	for _, r := range redBreak {
		totalReductionBase += r.BaseAmount
	}
	out.TotalReductionBase = ptrFloat64(round2(totalReductionBase))
	out.TotalExpenseGst = ptrFloat64(totalRedGst)

	// Additional reduction section (reduction / additional_reduction)
	var addlRedBreak []fieldCalc
	var totalAddlRed, totalAddlRedBase, totalAddlRedGst float64
	for _, f := range fields {
		if getSection(f.Section) != "reduction" {
			continue
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
		addlRedBreak = append(addlRedBreak, *ft)
		totalAddlRed += ft.TotalAmount
		totalAddlRedBase += ft.BaseAmount
		totalAddlRedGst += ft.GstAmount
	}
	if len(addlRedBreak) > 0 {
		out.AdditionalReductionBreakdown = addlRedBreak
		out.TotalAdditionalReduction = ptrFloat64(round2(totalAddlRed))
		out.TotalAdditionalReductionBase = ptrFloat64(round2(totalAddlRedBase))
		out.TotalAdditionalReductionGst = ptrFloat64(round2(totalAddlRedGst))
	}

	out.OutworkEnabled = formOutworkEnabled
	out.OutworkRatePercent = formOutworkRatePercent

	outworkRate := 0.0
	if formOutworkEnabled && formOutworkRatePercent != nil && *formOutworkRatePercent > 0 {
		outworkRate = *formOutworkRatePercent
	}
	var effectiveReductions float64
	if outworkRate > 0 {
		var totalOutworkCosts float64
		for _, r := range redBreak {
			totalOutworkCosts += r.BaseAmount
		}
		outworkChargeBase := round2(totalOutworkCosts * (outworkRate / 100.0))
		outworkChargeGst := round2(outworkChargeBase * 0.1)
		outworkChargeTotal := round2(outworkChargeBase + outworkChargeGst)
		out.OutworkChargeBase = &outworkChargeBase
		out.OutworkChargeGst = &outworkChargeGst
		out.OutworkChargeTotal = &outworkChargeTotal
		effectiveReductions = outworkChargeTotal
	} else {
		effectiveReductions = totalRedGst
	}
	out.TotalReductions = &effectiveReductions

	out.SubtotalAfterDeductions = ptrFloat64(round2(netFee - serviceBase))
	out.RemittedAmount = ptrFloat64(round2(netFee - totalSvc + totalReimb - effectiveReductions))
}
