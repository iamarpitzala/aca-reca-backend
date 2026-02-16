package domain

import (
	"encoding/json"
	"fmt"
)

// parseDeductionsFromEntry unmarshals entry.Deductions and builds EntryDeductions using map helpers.
func parseDeductionsFromEntry(entry *CustomFormEntry) (*EntryDeductions, error) {
	if entry == nil || len(entry.Deductions) == 0 {
		return nil, nil
	}
	var deductions map[string]interface{}
	if err := json.Unmarshal(entry.Deductions, &deductions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deductions: %w", err)
	}
	return &EntryDeductions{
		EntryID:                   entry.ID,
		CreatedAt:                 entry.CreatedAt,
		ServiceFacilityFeePercent: getFloat64Ptr(deductions, "serviceFacilityFeePercent"),
		ServiceFeeOverride:        getFloat64Ptr(deductions, "serviceFeeOverride"),
		CommissionPercent:        getFloat64Ptr(deductions, "commissionPercent"),
		SuperHoldingEnabled:      getBoolPtr(deductions, "superHoldingEnabled"),
		SuperComponentPercent:    getFloat64Ptr(deductions, "superComponentPercent"),
		OutworkEnabled:           getBoolPtr(deductions, "outworkEnabled"),
		OutworkRatePercent:       getFloat64Ptr(deductions, "outworkRatePercent"),
	}, nil
}
