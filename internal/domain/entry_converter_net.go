package domain

// parseNetDetailsFromCalc builds EntryNetDetails when form uses NET calculation method.
func parseNetDetailsFromCalc(entry *FieldEntry, calc map[string]interface{}, deductions *EntryDeductions) *EntryNetDetails {
	if entry == nil || calc == nil {
		return nil
	}
	nd := &EntryNetDetails{
		EntryID:    entry.ID,
		Commission: 0,
		CreatedAt:  entry.CreatedAt,
		UpdatedAt:  entry.UpdatedAt,
	}
	if deductions != nil && deductions.CommissionPercent != nil {
		nd.CommissionPercent = *deductions.CommissionPercent
	}
	if deductions != nil && deductions.SuperHoldingEnabled != nil {
		nd.SuperHoldingEnabled = *deductions.SuperHoldingEnabled
	}
	if deductions != nil && deductions.SuperComponentPercent != nil {
		nd.SuperComponentPercent = deductions.SuperComponentPercent
	}
	if calc != nil {
		nd.Commission = getFloat64(calc, "commission")
		nd.GstOnCommission = getFloat64(calc, "gstOnCommission")
		nd.TotalPaymentReceived = getFloat64(calc, "totalPaymentReceived")
		nd.CommissionComponent = getFloat64Ptr(calc, "commissionComponent")
		nd.SuperComponent = getFloat64Ptr(calc, "superComponent")
		nd.TotalForReconciliation = getFloat64Ptr(calc, "totalForReconciliation")
	}
	if nd.SuperHoldingEnabled && nd.CommissionComponent == nil {
		superPercent := 12.0
		if nd.SuperComponentPercent != nil {
			superPercent = *nd.SuperComponentPercent
		}
		superMultiplier := 1.0 + (superPercent / 100.0)
		commissionComponent := nd.Commission / superMultiplier
		nd.CommissionComponent = &commissionComponent
		superComponent := commissionComponent * (superPercent / 100.0)
		nd.SuperComponent = &superComponent
		totalForReconciliation := superComponent + commissionComponent
		nd.TotalForReconciliation = &totalForReconciliation
	}
	return nd
}
