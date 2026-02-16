package calculation

func applyNetMethodCalculations(out *calculationsOutput, deductions *deductionsInput) {
	commissionPercent := 0.0
	if deductions != nil && deductions.CommissionPercent != nil {
		commissionPercent = *deductions.CommissionPercent
	}
	netFee := out.TotalBaseAmount
	if out.NetFee != nil {
		netFee = *out.NetFee
	}
	commission := round2(netFee * (commissionPercent / 100.0))
	out.Commission = &commission

	superHoldingEnabled := false
	superComponentPercent := 12.0
	if deductions != nil {
		if deductions.SuperHoldingEnabled != nil {
			superHoldingEnabled = *deductions.SuperHoldingEnabled
		}
		if deductions.SuperComponentPercent != nil && *deductions.SuperComponentPercent > 0 {
			superComponentPercent = *deductions.SuperComponentPercent
		}
	}

	var commissionForGst float64
	if superHoldingEnabled {
		superMultiplier := 1.0 + (superComponentPercent / 100.0)
		commissionComponent := round2(commission / superMultiplier)
		out.CommissionComponent = &commissionComponent
		superComponent := round2(commissionComponent * (superComponentPercent / 100.0))
		out.SuperComponent = &superComponent
		out.TotalForReconciliation = ptrFloat64(round2(superComponent + commissionComponent))
		commissionForGst = commissionComponent
	} else {
		commissionForGst = commission
	}
	gstOnCommission := round2(commissionForGst * 0.1)
	out.GstOnCommission = &gstOnCommission
	out.TotalPaymentReceived = ptrFloat64(round2(commissionForGst + gstOnCommission))
}
