package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ConvertNormalizedToJSONB converts a normalized entry structure to JSONB format (for API compatibility)
func ConvertNormalizedToJSONB(normalized *NormalizedEntry) (*CustomFormEntry, error) {
	if normalized == nil || normalized.Header == nil {
		return nil, nil
	}

	// Convert field values to JSONB
	values := make([]map[string]interface{}, 0, len(normalized.FieldValues))
	if len(normalized.FieldValues) > 0 {
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
	} else if len(normalized.FieldCalculations) > 0 {
		for _, fc := range normalized.FieldCalculations {
			val := map[string]interface{}{
				"fieldId":   fc.FieldID,
				"fieldName": fc.FieldName,
			}
			gstType := strings.ToLower(fc.GstType)
			if gstType == "inclusive" {
				val["value"] = fc.TotalAmount
			} else {
				val["value"] = fc.BaseAmount
			}
			values = append(values, val)
		}
	} else {
		seen := make(map[string]bool)
		addFromBreakdown := func(fieldID, fieldName string, totalAmount float64) {
			if fieldID == "" || seen[fieldID] {
				return
			}
			seen[fieldID] = true
			values = append(values, map[string]interface{}{
				"fieldId":   fieldID,
				"fieldName": fieldName,
				"value":     totalAmount,
			})
		}
		for _, r := range normalized.GrossReductions {
			addFromBreakdown(r.FieldID, r.FieldName, r.TotalAmount)
		}
		for _, r := range normalized.GrossReimbursements {
			addFromBreakdown(r.FieldID, r.FieldName, r.TotalAmount)
		}
		for _, r := range normalized.GrossAdditionalReductions {
			addFromBreakdown(r.FieldID, r.FieldName, r.TotalAmount)
		}
	}
	valuesJSON, _ := json.Marshal(values)

	// Convert calculations to JSONB
	calculations := make(map[string]interface{})
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

	if len(normalized.GrossReductions) > 0 {
		reductionBreakdown := make([]map[string]interface{}, 0, len(normalized.GrossReductions))
		for _, r := range normalized.GrossReductions {
			reductionBreakdown = append(reductionBreakdown, map[string]interface{}{
				"fieldId": r.FieldID, "fieldName": r.FieldName,
				"baseAmount": r.BaseAmount, "gstAmount": r.GstAmount, "totalAmount": r.TotalAmount,
			})
		}
		calculations["reductionBreakdown"] = reductionBreakdown
	}
	if len(normalized.GrossReimbursements) > 0 {
		reimbursementBreakdown := make([]map[string]interface{}, 0, len(normalized.GrossReimbursements))
		for _, r := range normalized.GrossReimbursements {
			reimbursementBreakdown = append(reimbursementBreakdown, map[string]interface{}{
				"fieldId": r.FieldID, "fieldName": r.FieldName,
				"baseAmount": r.BaseAmount, "gstAmount": r.GstAmount, "totalAmount": r.TotalAmount,
			})
		}
		calculations["reimbursementBreakdown"] = reimbursementBreakdown
	}
	if len(normalized.GrossAdditionalReductions) > 0 {
		additionalReductionBreakdown := make([]map[string]interface{}, 0, len(normalized.GrossAdditionalReductions))
		for _, r := range normalized.GrossAdditionalReductions {
			additionalReductionBreakdown = append(additionalReductionBreakdown, map[string]interface{}{
				"fieldId": r.FieldID, "fieldName": r.FieldName,
				"baseAmount": r.BaseAmount, "gstAmount": r.GstAmount, "totalAmount": r.TotalAmount,
			})
		}
		calculations["additionalReductionBreakdown"] = additionalReductionBreakdown
	}
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

	var calc map[string]interface{}
	if len(entry.Calculations) > 0 {
		if err := json.Unmarshal(entry.Calculations, &calc); err != nil {
			return nil, fmt.Errorf("failed to unmarshal calculations: %w", err)
		}
	}

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

	normalized.Summary = parseSummaryFromCalc(entry.ID, entry.CreatedAt, entry.UpdatedAt, calc)

	ded, err := parseDeductionsFromEntry(entry)
	if err != nil {
		return nil, err
	}
	normalized.Deductions = ded

	normalized.NetDetails = parseNetDetailsFromCalc(entry, form, calc, ded)
	normalized.GrossDetails = parseGrossDetailsFromCalc(entry, form, calc, ded)

	if calc != nil {
		if reductionBreakdown, ok := calc["reductionBreakdown"].([]interface{}); ok {
			for i, rb := range reductionBreakdown {
				if rbMap, ok := rb.(map[string]interface{}); ok {
					normalized.GrossReductions = append(normalized.GrossReductions, parseBreakdownToGrossReduction(entry.ID, i, entry.CreatedAt, rbMap))
				}
			}
		}
		if reimbursementBreakdown, ok := calc["reimbursementBreakdown"].([]interface{}); ok {
			for i, rb := range reimbursementBreakdown {
				if rbMap, ok := rb.(map[string]interface{}); ok {
					normalized.GrossReimbursements = append(normalized.GrossReimbursements, parseBreakdownToGrossReimbursement(entry.ID, i, entry.CreatedAt, rbMap))
				}
			}
		}
		if additionalReductionBreakdown, ok := calc["additionalReductionBreakdown"].([]interface{}); ok {
			for i, rb := range additionalReductionBreakdown {
				if rbMap, ok := rb.(map[string]interface{}); ok {
					normalized.GrossAdditionalReductions = append(normalized.GrossAdditionalReductions, parseBreakdownToGrossAdditionalReduction(entry.ID, i, entry.CreatedAt, rbMap))
				}
			}
		}
		normalized.GrossReductionsSummary = parseGrossReductionsSummaryFromCalc(entry, calc)
		normalized.GrossOutwork = parseGrossOutworkFromCalc(entry, calc)
	}

	return normalized, nil
}
