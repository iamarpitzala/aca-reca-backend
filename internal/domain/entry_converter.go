package domain

import (
	"encoding/json"
	"fmt"
)

// ConvertNormalizedToJSONB converts a normalized entry structure to JSONB format (for API compatibility)
func ConvertNormalizedToJSONB(normalized *NormalizedEntry) (*CustomFormEntry, error) {
	if normalized == nil || normalized.Header == nil {
		return nil, nil
	}

	// Convert field values to JSONB
	values := make([]map[string]interface{}, 0, len(normalized.FieldValues))
	for _, fv := range normalized.FieldValues {
		val := map[string]interface{}{
			"fieldId":   fv.FieldID,
			"fieldName": fv.FieldName,
		}
		if fv.Value != nil {
			val["value"] = *fv.Value
		} else if fv.TextValue != nil {
			val["value"] = *fv.TextValue
		} else if fv.BooleanValue != nil {
			val["value"] = *fv.BooleanValue
		}
		if fv.ManualGstAmount != nil {
			val["manualGstAmount"] = *fv.ManualGstAmount
		}
		values = append(values, val)
	}
	valuesJSON, _ := json.Marshal(values)

	// Convert calculations to JSONB
	calculations := make(map[string]interface{})

	// Field totals
	fieldTotals := make([]map[string]interface{}, 0, len(normalized.FieldCalculations))
	for _, fc := range normalized.FieldCalculations {
		fieldTotals = append(fieldTotals, map[string]interface{}{
			"fieldId":     fc.FieldID,
			"fieldName":   fc.FieldName,
			"baseAmount":  fc.BaseAmount,
			"gstAmount":   fc.GstAmount,
			"totalAmount": fc.TotalAmount,
			"gstRate":     fc.GstRate,
			"gstType":     fc.GstType,
		})
	}
	calculations["fieldTotals"] = fieldTotals

	// Summary totals
	if normalized.Summary != nil {
		s := normalized.Summary
		calculations["totalBaseAmount"] = s.TotalBaseAmount
		calculations["totalGSTAmount"] = s.TotalGstAmount
		calculations["totalAmount"] = s.TotalAmount
		calculations["netPayable"] = s.NetPayable
		calculations["netReceivable"] = s.NetReceivable
		if s.NetFee != nil {
			calculations["netFee"] = *s.NetFee
		}
		calculations["basMapping"] = map[string]interface{}{
			"gstOnSales1A": s.BasGstOnSales1A,
			"gstCredit1B":  s.BasGstCredit1B,
			"totalSalesG1": s.BasTotalSalesG1,
			"expensesG11":  s.BasExpensesG11,
		}
	}

	// Net method fields
	if normalized.NetDetails != nil {
		nd := normalized.NetDetails
		calculations["commission"] = nd.Commission
		calculations["gstOnCommission"] = nd.GstOnCommission
		calculations["totalPaymentReceived"] = nd.TotalPaymentReceived
		if nd.CommissionComponent != nil {
			calculations["commissionComponent"] = *nd.CommissionComponent
		}
		if nd.SuperComponent != nil {
			calculations["superComponent"] = *nd.SuperComponent
		}
		if nd.TotalForReconciliation != nil {
			calculations["totalForReconciliation"] = *nd.TotalForReconciliation
		}
	}

	// Gross method fields
	if normalized.GrossDetails != nil {
		gd := normalized.GrossDetails
		calculations["serviceFeeBase"] = gd.ServiceFeeBase
		calculations["gstOnServiceFee"] = gd.GstOnServiceFee
		calculations["totalServiceFee"] = gd.TotalServiceFee
		if gd.SubtotalAfterDeductions != nil {
			calculations["subtotalAfterDeductions"] = *gd.SubtotalAfterDeductions
		}
		if gd.RemittedAmount != nil {
			calculations["remittedAmount"] = *gd.RemittedAmount
		}
	}

	// Gross reductions
	if len(normalized.GrossReductions) > 0 {
		reductionBreakdown := make([]map[string]interface{}, 0, len(normalized.GrossReductions))
		for _, r := range normalized.GrossReductions {
			reductionBreakdown = append(reductionBreakdown, map[string]interface{}{
				"fieldId":     r.FieldID,
				"fieldName":   r.FieldName,
				"baseAmount":  r.BaseAmount,
				"gstAmount":   r.GstAmount,
				"totalAmount": r.TotalAmount,
			})
		}
		calculations["reductionBreakdown"] = reductionBreakdown
	}

	// Gross reimbursements
	if len(normalized.GrossReimbursements) > 0 {
		reimbursementBreakdown := make([]map[string]interface{}, 0, len(normalized.GrossReimbursements))
		for _, r := range normalized.GrossReimbursements {
			reimbursementBreakdown = append(reimbursementBreakdown, map[string]interface{}{
				"fieldId":     r.FieldID,
				"fieldName":   r.FieldName,
				"baseAmount":  r.BaseAmount,
				"gstAmount":   r.GstAmount,
				"totalAmount": r.TotalAmount,
			})
		}
		calculations["reimbursementBreakdown"] = reimbursementBreakdown
	}

	// Gross additional reductions
	if len(normalized.GrossAdditionalReductions) > 0 {
		additionalReductionBreakdown := make([]map[string]interface{}, 0, len(normalized.GrossAdditionalReductions))
		for _, r := range normalized.GrossAdditionalReductions {
			additionalReductionBreakdown = append(additionalReductionBreakdown, map[string]interface{}{
				"fieldId":     r.FieldID,
				"fieldName":   r.FieldName,
				"baseAmount":  r.BaseAmount,
				"gstAmount":   r.GstAmount,
				"totalAmount": r.TotalAmount,
			})
		}
		calculations["additionalReductionBreakdown"] = additionalReductionBreakdown
	}

	// Gross reductions summary
	if normalized.GrossReductionsSummary != nil {
		rs := normalized.GrossReductionsSummary
		calculations["totalReductions"] = rs.TotalReductions
		calculations["totalReductionBase"] = rs.TotalReductionBase
		calculations["totalExpenseGst"] = rs.TotalExpenseGst
		calculations["totalReimbursements"] = rs.TotalReimbursements
		calculations["totalAdditionalReduction"] = rs.TotalAdditionalReduction
		calculations["totalAdditionalReductionBase"] = rs.TotalAdditionalReductionBase
		calculations["totalAdditionalReductionGst"] = rs.TotalAdditionalReductionGst
	}

	// Gross outwork
	if normalized.GrossOutwork != nil {
		o := normalized.GrossOutwork
		calculations["outworkEnabled"] = o.OutworkEnabled
		if o.OutworkRatePercent != nil {
			calculations["outworkRatePercent"] = *o.OutworkRatePercent
		}
		calculations["outworkChargeBase"] = o.OutworkChargeBase
		calculations["outworkChargeGst"] = o.OutworkChargeGst
		calculations["outworkChargeTotal"] = o.OutworkChargeTotal
	}

	calculationsJSON, _ := json.Marshal(calculations)

	// Convert deductions to JSONB
	var deductionsJSON json.RawMessage
	if normalized.Deductions != nil {
		deductions := make(map[string]interface{})
		if normalized.Deductions.ServiceFacilityFeePercent != nil {
			deductions["serviceFacilityFeePercent"] = *normalized.Deductions.ServiceFacilityFeePercent
		}
		if normalized.Deductions.ServiceFeeOverride != nil {
			deductions["serviceFeeOverride"] = *normalized.Deductions.ServiceFeeOverride
		}
		if normalized.Deductions.CommissionPercent != nil {
			deductions["commissionPercent"] = *normalized.Deductions.CommissionPercent
		}
		if normalized.Deductions.SuperHoldingEnabled != nil {
			deductions["superHoldingEnabled"] = *normalized.Deductions.SuperHoldingEnabled
		}
		if normalized.Deductions.SuperComponentPercent != nil {
			deductions["superComponentPercent"] = *normalized.Deductions.SuperComponentPercent
		}
		if normalized.Deductions.OutworkEnabled != nil {
			deductions["outworkEnabled"] = *normalized.Deductions.OutworkEnabled
		}
		if normalized.Deductions.OutworkRatePercent != nil {
			deductions["outworkRatePercent"] = *normalized.Deductions.OutworkRatePercent
		}
		deductionsJSON, _ = json.Marshal(deductions)
	}

	return &CustomFormEntry{
		ID:                    normalized.Header.ID,
		FormID:                normalized.Header.FormID,
		FormName:              normalized.Header.FormName,
		FormType:              normalized.Header.FormType,
		ClinicID:              normalized.Header.ClinicID,
		QuarterID:             normalized.Header.QuarterID,
		Values:                valuesJSON,
		Calculations:          calculationsJSON,
		EntryDate:             normalized.Header.EntryDate,
		Description:           normalized.Header.Description,
		Remarks:               normalized.Header.Remarks,
		PaymentResponsibility: normalized.Header.PaymentResponsibility,
		Deductions:            deductionsJSON,
		CreatedBy:             normalized.Header.CreatedBy,
		CreatedAt:             normalized.Header.CreatedAt,
		UpdatedAt:             normalized.Header.UpdatedAt,
		DeletedAt:             normalized.Header.DeletedAt,
	}, nil
}

// ConvertJSONBToNormalized converts JSONB entry to normalized structure
func ConvertJSONBToNormalized(entry *CustomFormEntry, form *CustomForm) (*NormalizedEntry, error) {
	if entry == nil {
		return nil, nil
	}

	normalized := &NormalizedEntry{
		Header: &EntryHeader{
			ID:                    entry.ID,
			FormID:                entry.FormID,
			FormName:              entry.FormName,
			FormType:              entry.FormType,
			CalculationMethod:     form.CalculationMethod,
			ClinicID:              entry.ClinicID,
			QuarterID:             entry.QuarterID,
			EntryDate:             entry.EntryDate,
			Description:           entry.Description,
			Remarks:               entry.Remarks,
			PaymentResponsibility: entry.PaymentResponsibility,
			CreatedBy:             entry.CreatedBy,
			CreatedAt:             entry.CreatedAt,
			UpdatedAt:             entry.UpdatedAt,
			DeletedAt:             entry.DeletedAt,
		},
	}

	// Parse field values
	var values []map[string]interface{}
	if len(entry.Values) > 0 {
		if err := json.Unmarshal(entry.Values, &values); err != nil {
			return nil, fmt.Errorf("failed to unmarshal field values: %w", err)
		}
		for i, v := range values {
			fv := EntryFieldValue{
				EntryID:      entry.ID,
				DisplayOrder: i,
				CreatedAt:    entry.CreatedAt,
			}
			if fieldID, ok := v["fieldId"].(string); ok {
				fv.FieldID = fieldID
			}
			if fieldName, ok := v["fieldName"].(string); ok {
				fv.FieldName = fieldName
			}
			if val, ok := v["value"]; ok {
				switch val := val.(type) {
				case float64:
					fv.Value = &val
				case string:
					fv.TextValue = &val
				case bool:
					fv.BooleanValue = &val
				}
			}
			if mgst, ok := v["manualGstAmount"].(float64); ok {
				fv.ManualGstAmount = &mgst
			}
			normalized.FieldValues = append(normalized.FieldValues, fv)
		}
	}

	// Parse calculations
	var calc map[string]interface{}
	if len(entry.Calculations) > 0 {
		if err := json.Unmarshal(entry.Calculations, &calc); err != nil {
			return nil, fmt.Errorf("failed to unmarshal calculations: %w", err)
		}
	}

	// Parse field totals
	if fieldTotals, ok := calc["fieldTotals"].([]interface{}); ok {
		for i, ft := range fieldTotals {
			ftMap, ok := ft.(map[string]interface{})
			if !ok {
				continue
			}
			fc := EntryFieldCalculation{
				EntryID:      entry.ID,
				DisplayOrder: i,
				CreatedAt:    entry.CreatedAt,
			}
			if fieldID, ok := ftMap["fieldId"].(string); ok {
				fc.FieldID = fieldID
			}
			if fieldName, ok := ftMap["fieldName"].(string); ok {
				fc.FieldName = fieldName
			}
			if baseAmount, ok := ftMap["baseAmount"].(float64); ok {
				fc.BaseAmount = baseAmount
			}
			if gstAmount, ok := ftMap["gstAmount"].(float64); ok {
				fc.GstAmount = gstAmount
			}
			if totalAmount, ok := ftMap["totalAmount"].(float64); ok {
				fc.TotalAmount = totalAmount
			}
			if gstRate, ok := ftMap["gstRate"].(float64); ok {
				fc.GstRate = gstRate
			}
			if gstType, ok := ftMap["gstType"].(string); ok {
				fc.GstType = gstType
			}
			if section, ok := ftMap["section"].(string); ok {
				fc.Section = &section
			}
			if payResp, ok := ftMap["paymentResponsibility"].(string); ok {
				fc.PaymentResponsibility = &payResp
			}
			normalized.FieldCalculations = append(normalized.FieldCalculations, fc)
		}
	}

	// Parse summary - always create summary even if calculations are empty
	summary := &EntrySummary{
		EntryID:   entry.ID,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}
	if calc != nil {
		if val, ok := calc["totalBaseAmount"].(float64); ok {
			summary.TotalBaseAmount = val
		}
		if val, ok := calc["totalGSTAmount"].(float64); ok {
			summary.TotalGstAmount = val
		}
		if val, ok := calc["totalAmount"].(float64); ok {
			summary.TotalAmount = val
		}
		if val, ok := calc["netPayable"].(float64); ok {
			summary.NetPayable = val
		}
		if val, ok := calc["netReceivable"].(float64); ok {
			summary.NetReceivable = val
		}
		if val, ok := calc["netFee"].(float64); ok {
			summary.NetFee = &val
		}
		if basMap, ok := calc["basMapping"].(map[string]interface{}); ok {
			if val, ok := basMap["gstOnSales1A"].(float64); ok {
				summary.BasGstOnSales1A = val
			}
			if val, ok := basMap["gstCredit1B"].(float64); ok {
				summary.BasGstCredit1B = val
			}
			if val, ok := basMap["totalSalesG1"].(float64); ok {
				summary.BasTotalSalesG1 = val
			}
			if val, ok := basMap["expensesG11"].(float64); ok {
				summary.BasExpensesG11 = val
			}
		}
	}
	normalized.Summary = summary

	// Parse deductions
	if len(entry.Deductions) > 0 {
		var deductions map[string]interface{}
		if err := json.Unmarshal(entry.Deductions, &deductions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal deductions: %w", err)
		}
		ded := &EntryDeductions{
			EntryID:   entry.ID,
			CreatedAt: entry.CreatedAt,
		}
		if val, ok := deductions["serviceFacilityFeePercent"].(float64); ok {
			ded.ServiceFacilityFeePercent = &val
		}
		if val, ok := deductions["serviceFeeOverride"].(float64); ok {
			ded.ServiceFeeOverride = &val
		}
		if val, ok := deductions["commissionPercent"].(float64); ok {
			ded.CommissionPercent = &val
		}
		if val, ok := deductions["superHoldingEnabled"].(bool); ok {
			ded.SuperHoldingEnabled = &val
		}
		if val, ok := deductions["superComponentPercent"].(float64); ok {
			ded.SuperComponentPercent = &val
		}
		if val, ok := deductions["outworkEnabled"].(bool); ok {
			ded.OutworkEnabled = &val
		}
		if val, ok := deductions["outworkRatePercent"].(float64); ok {
			ded.OutworkRatePercent = &val
		}
		normalized.Deductions = ded
	}

	// Parse net method details
	if form.CalculationMethod == "net" {
		nd := &EntryNetDetails{
			EntryID:    entry.ID,
			Commission: 0,
			CreatedAt:  entry.CreatedAt,
			UpdatedAt:  entry.UpdatedAt,
		}
		
		// Set commission percent from deductions (required field)
		if normalized.Deductions != nil && normalized.Deductions.CommissionPercent != nil {
			nd.CommissionPercent = *normalized.Deductions.CommissionPercent
		} else {
			// Default to 0 if not provided
			nd.CommissionPercent = 0
		}
		
		// Set super holding enabled from deductions
		if normalized.Deductions != nil && normalized.Deductions.SuperHoldingEnabled != nil {
			nd.SuperHoldingEnabled = *normalized.Deductions.SuperHoldingEnabled
		}
		
		// Set super component percent from deductions
		if normalized.Deductions != nil && normalized.Deductions.SuperComponentPercent != nil {
			nd.SuperComponentPercent = normalized.Deductions.SuperComponentPercent
		}
		
		// Populate calculated values from calculations JSONB if available
		if calc != nil {
			if commission, ok := calc["commission"].(float64); ok {
				nd.Commission = commission
			}
			if val, ok := calc["gstOnCommission"].(float64); ok {
				nd.GstOnCommission = val
			}
			if val, ok := calc["totalPaymentReceived"].(float64); ok {
				nd.TotalPaymentReceived = val
			}
			if val, ok := calc["commissionComponent"].(float64); ok {
				nd.CommissionComponent = &val
			}
			if val, ok := calc["superComponent"].(float64); ok {
				nd.SuperComponent = &val
			}
			if val, ok := calc["totalForReconciliation"].(float64); ok {
				nd.TotalForReconciliation = &val
			}
		}
		
		// If super holding is enabled but super-related fields are not set, calculate them from commission
		// This handles cases where calculations JSONB might be missing these fields
		if nd.SuperHoldingEnabled && nd.CommissionComponent == nil {
			// Calculate commission component and super component from commission
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
		
		normalized.NetDetails = nd
	}

	// Parse gross method details
	if calc != nil && form.CalculationMethod == "gross" {
		if serviceFeeBase, ok := calc["serviceFeeBase"].(float64); ok {
			gd := &EntryGrossDetails{
				EntryID:        entry.ID,
				ServiceFeeBase: serviceFeeBase,
				CreatedAt:      entry.CreatedAt,
				UpdatedAt:      entry.UpdatedAt,
			}
			if val, ok := calc["gstOnServiceFee"].(float64); ok {
				gd.GstOnServiceFee = val
			}
			if val, ok := calc["totalServiceFee"].(float64); ok {
				gd.TotalServiceFee = val
			}
			if val, ok := calc["subtotalAfterDeductions"].(float64); ok {
				gd.SubtotalAfterDeductions = &val
			}
			if val, ok := calc["remittedAmount"].(float64); ok {
				gd.RemittedAmount = &val
			}
			if normalized.Deductions != nil && normalized.Deductions.ServiceFacilityFeePercent != nil {
				gd.ServiceFacilityFeePercent = *normalized.Deductions.ServiceFacilityFeePercent
			}
			normalized.GrossDetails = gd
		}
	}

	// Parse gross reductions
	if calc != nil {
		if reductionBreakdown, ok := calc["reductionBreakdown"].([]interface{}); ok {
			for i, rb := range reductionBreakdown {
				rbMap, ok := rb.(map[string]interface{})
				if !ok {
					continue
				}
				gr := EntryGrossReduction{
					EntryID:      entry.ID,
					DisplayOrder: i,
					CreatedAt:    entry.CreatedAt,
				}
				if fieldID, ok := rbMap["fieldId"].(string); ok {
					gr.FieldID = fieldID
				}
				if fieldName, ok := rbMap["fieldName"].(string); ok {
					gr.FieldName = fieldName
				}
				if baseAmount, ok := rbMap["baseAmount"].(float64); ok {
					gr.BaseAmount = baseAmount
				}
				if gstAmount, ok := rbMap["gstAmount"].(float64); ok {
					gr.GstAmount = gstAmount
				}
				if totalAmount, ok := rbMap["totalAmount"].(float64); ok {
					gr.TotalAmount = totalAmount
				}
				normalized.GrossReductions = append(normalized.GrossReductions, gr)
			}
		}

		// Parse gross reimbursements
		if reimbursementBreakdown, ok := calc["reimbursementBreakdown"].([]interface{}); ok {
			for i, rb := range reimbursementBreakdown {
				rbMap, ok := rb.(map[string]interface{})
				if !ok {
					continue
				}
				gr := EntryGrossReimbursement{
					EntryID:      entry.ID,
					DisplayOrder: i,
					CreatedAt:    entry.CreatedAt,
				}
				if fieldID, ok := rbMap["fieldId"].(string); ok {
					gr.FieldID = fieldID
				}
				if fieldName, ok := rbMap["fieldName"].(string); ok {
					gr.FieldName = fieldName
				}
				if baseAmount, ok := rbMap["baseAmount"].(float64); ok {
					gr.BaseAmount = baseAmount
				}
				if gstAmount, ok := rbMap["gstAmount"].(float64); ok {
					gr.GstAmount = gstAmount
				}
				if totalAmount, ok := rbMap["totalAmount"].(float64); ok {
					gr.TotalAmount = totalAmount
				}
				normalized.GrossReimbursements = append(normalized.GrossReimbursements, gr)
			}
		}

		// Parse gross additional reductions
		if additionalReductionBreakdown, ok := calc["additionalReductionBreakdown"].([]interface{}); ok {
			for i, rb := range additionalReductionBreakdown {
				rbMap, ok := rb.(map[string]interface{})
				if !ok {
					continue
				}
				gar := EntryGrossAdditionalReduction{
					EntryID:      entry.ID,
					DisplayOrder: i,
					CreatedAt:    entry.CreatedAt,
				}
				if fieldID, ok := rbMap["fieldId"].(string); ok {
					gar.FieldID = fieldID
				}
				if fieldName, ok := rbMap["fieldName"].(string); ok {
					gar.FieldName = fieldName
				}
				if baseAmount, ok := rbMap["baseAmount"].(float64); ok {
					gar.BaseAmount = baseAmount
				}
				if gstAmount, ok := rbMap["gstAmount"].(float64); ok {
					gar.GstAmount = gstAmount
				}
				if totalAmount, ok := rbMap["totalAmount"].(float64); ok {
					gar.TotalAmount = totalAmount
				}
				normalized.GrossAdditionalReductions = append(normalized.GrossAdditionalReductions, gar)
			}
		}

		// Parse gross reductions summary
		if totalReductions, ok := calc["totalReductions"].(float64); ok {
			rs := &EntryGrossReductionsSummary{
				EntryID:         entry.ID,
				TotalReductions: totalReductions,
				CreatedAt:       entry.CreatedAt,
				UpdatedAt:       entry.UpdatedAt,
			}
			if val, ok := calc["totalReductionBase"].(float64); ok {
				rs.TotalReductionBase = val
			}
			if val, ok := calc["totalExpenseGst"].(float64); ok {
				rs.TotalExpenseGst = val
			}
			if val, ok := calc["totalReimbursements"].(float64); ok {
				rs.TotalReimbursements = val
			}
			if val, ok := calc["totalAdditionalReduction"].(float64); ok {
				rs.TotalAdditionalReduction = val
			}
			if val, ok := calc["totalAdditionalReductionBase"].(float64); ok {
				rs.TotalAdditionalReductionBase = val
			}
			if val, ok := calc["totalAdditionalReductionGst"].(float64); ok {
				rs.TotalAdditionalReductionGst = val
			}
			normalized.GrossReductionsSummary = rs
		}

		// Parse gross outwork
		if outworkChargeBase, ok := calc["outworkChargeBase"].(float64); ok {
			o := &EntryGrossOutwork{
				EntryID:           entry.ID,
				OutworkChargeBase: outworkChargeBase,
				CreatedAt:         entry.CreatedAt,
				UpdatedAt:         entry.UpdatedAt,
			}
			if val, ok := calc["outworkEnabled"].(bool); ok {
				o.OutworkEnabled = val
			}
			if val, ok := calc["outworkRatePercent"].(float64); ok {
				o.OutworkRatePercent = &val
			}
			if val, ok := calc["outworkChargeGst"].(float64); ok {
				o.OutworkChargeGst = val
			}
			if val, ok := calc["outworkChargeTotal"].(float64); ok {
				o.OutworkChargeTotal = val
			}
			normalized.GrossOutwork = o
		}
	}

	return normalized, nil
}
