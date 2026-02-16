package calculation

import (
	"encoding/json"
	"fmt"
)

// ParseNetDeductions parses deductions JSON and returns commission percent, super holding enabled, and super component percent.
func ParseNetDeductions(deductionsJSON json.RawMessage) (commissionPercent float64, superHoldingEnabled bool, superComponentPercent *float64, err error) {
	if len(deductionsJSON) == 0 {
		return 0, false, nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(deductionsJSON, &m); err != nil {
		return 0, false, nil, fmt.Errorf("parse deductions: %w", err)
	}
	commissionPercent = getFloat64(m, "commissionPercent")
	superHoldingEnabled = getBool(m, "superHoldingEnabled")
	if p := getFloat64Ptr(m, "superComponentPercent"); p != nil {
		superComponentPercent = p
	}
	return commissionPercent, superHoldingEnabled, superComponentPercent, nil
}

// ParseGrossDeductions parses deductions JSON and returns service facility fee percent, outwork enabled, and outwork rate percent.
func ParseGrossDeductions(deductionsJSON json.RawMessage) (serviceFacilityFeePercent *float64, outworkEnabled bool, outworkRatePercent *float64, err error) {
	if len(deductionsJSON) == 0 {
		return nil, false, nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(deductionsJSON, &m); err != nil {
		return nil, false, nil, fmt.Errorf("parse deductions: %w", err)
	}
	serviceFacilityFeePercent = getFloat64Ptr(m, "serviceFacilityFeePercent")
	outworkEnabled = getBool(m, "outworkEnabled")
	outworkRatePercent = getFloat64Ptr(m, "outworkRatePercent")
	return serviceFacilityFeePercent, outworkEnabled, outworkRatePercent, nil
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch x := v.(type) {
		case float64:
			return x
		case int:
			return float64(x)
		case int64:
			return float64(x)
		}
	}
	return 0
}

func getFloat64Ptr(m map[string]interface{}, key string) *float64 {
	if v, ok := m[key]; ok {
		switch x := v.(type) {
		case float64:
			return &x
		case int:
			f := float64(x)
			return &f
		case int64:
			f := float64(x)
			return &f
		}
	}
	return nil
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}
