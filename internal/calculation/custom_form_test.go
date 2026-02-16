package calculation

import (
	"encoding/json"
	"testing"

	"github.com/iamarpitzala/aca-reca-backend/util"
)

// formFieldsJSON for INCOME form with one number field (used in gross tests)
var testFormIncomeFields = []byte(`[
	{"id":"f1","name":"Fee","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"}
]`)

// formFieldsJSON for BOTH form: income + expense (clinic) + optional reduction section
var testFormBothFields = []byte(`[
	{"id":"inc1","name":"Income","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"},
	{"id":"exp1","name":"Expense","type":"number","section":"expense","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"CLINIC"},
	{"id":"red1","name":"Additional reduction","type":"number","section":"reduction","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"}
]`)

func TestRunEntryCalculation_GrossMethod_Income(t *testing.T) {
	values := []byte(`[{"fieldId":"f1","value":1000}]`)
	out, err := RunEntryCalculation(
		testFormIncomeFields,
		util.FormTypeIncome,
		util.MethodTypeGross,
		nil,
		false,
		nil,
		values,
		nil,
	)
	if err != nil {
		t.Fatalf("RunEntryCalculation: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// Default service fee 50% of net fee (1000 base) -> 500 base, 50 GST, 550 total
	if m["serviceFeeBase"] == nil {
		t.Error("expected serviceFeeBase for gross method")
	}
	if m["totalServiceFee"] == nil {
		t.Error("expected totalServiceFee")
	}
	if m["subtotalAfterDeductions"] == nil {
		t.Error("expected subtotalAfterDeductions")
	}
	if m["remittedAmount"] == nil {
		t.Error("expected remittedAmount")
	}
	if m["outworkEnabled"] != false {
		t.Errorf("outworkEnabled want false, got %v", m["outworkEnabled"])
	}
	// No expense fields -> no reduction/reimbursement breakdown
	if m["totalReductionBase"] == nil {
		t.Error("expected totalReductionBase (can be 0)")
	}
}

func TestRunEntryCalculation_GrossMethod_WithDeductions(t *testing.T) {
	values := []byte(`[{"fieldId":"f1","value":1000}]`)
	deductions := []byte(`{"serviceFacilityFeePercent":40,"serviceFeeOverride":300}`)
	out, err := RunEntryCalculation(
		testFormIncomeFields,
		util.FormTypeIncome,
		util.MethodTypeGross,
		nil,
		false,
		nil,
		values,
		deductions,
	)
	if err != nil {
		t.Fatalf("RunEntryCalculation: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// Override 300 -> service fee base 300
	if sb, ok := m["serviceFeeBase"].(float64); !ok || sb != 300 {
		t.Errorf("serviceFeeBase want 300 (override), got %v", m["serviceFeeBase"])
	}
}

func TestRunEntryCalculation_GrossMethod_OutworkEnabled(t *testing.T) {
	// BOTH form: income 1100 (1000+100 GST), expense 220 (200+20 GST) clinic-paid -> net fee 880, then outwork on expense base 200
	formBothOneExpense := []byte(`[
		{"id":"inc1","name":"Income","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"},
		{"id":"exp1","name":"Expense","type":"number","section":"expense","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"CLINIC"}
	]`)
	values := []byte(`[{"fieldId":"inc1","value":1000},{"fieldId":"exp1","value":200}]`)
	svcFeePct := 50.0
	outworkPct := 5.0
	out, err := RunEntryCalculation(
		formBothOneExpense,
		util.FormTypeBoth,
		util.MethodTypeGross,
		&svcFeePct,
		true,
		&outworkPct,
		values,
		nil,
	)
	if err != nil {
		t.Fatalf("RunEntryCalculation: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["outworkEnabled"] != true {
		t.Errorf("outworkEnabled want true, got %v", m["outworkEnabled"])
	}
	if m["outworkChargeBase"] == nil {
		t.Error("expected outworkChargeBase when outwork enabled with expense")
	}
	if m["outworkChargeTotal"] == nil {
		t.Error("expected outworkChargeTotal")
	}
	// When expense is clinic-paid, outwork 5% of expense base gives non-zero charge
	if ob, ok := m["outworkChargeBase"].(float64); ok && ob > 0 {
		if m["outworkChargeTotal"].(float64) <= 0 {
			t.Error("outworkChargeTotal should be positive when outworkChargeBase is")
		}
	}
}

func TestRunEntryCalculation_NetMethod_Income(t *testing.T) {
	values := []byte(`[{"fieldId":"f1","value":1000}]`)
	deductions := []byte(`{"commissionPercent":10}`)
	out, err := RunEntryCalculation(
		testFormIncomeFields,
		util.FormTypeIncome,
		util.MethodTypeNet,
		nil,
		false,
		nil,
		values,
		deductions,
	)
	if err != nil {
		t.Fatalf("RunEntryCalculation: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["commission"] == nil {
		t.Error("expected commission for net method")
	}
	if m["gstOnCommission"] == nil {
		t.Error("expected gstOnCommission")
	}
	if m["totalPaymentReceived"] == nil {
		t.Error("expected totalPaymentReceived")
	}
	// Gross method fields should be absent
	if m["serviceFeeBase"] != nil {
		t.Error("net method should not set serviceFeeBase")
	}
}

func TestRunEntryCalculation_GrossMethod_AdditionalReduction(t *testing.T) {
	// Income 1000, expense 200 (clinic), reduction 50
	values := []byte(`[
		{"fieldId":"inc1","value":1000},
		{"fieldId":"exp1","value":200},
		{"fieldId":"red1","value":50}
	]`)
	svcFeePct := 50.0
	out, err := RunEntryCalculation(
		testFormBothFields,
		util.FormTypeBoth,
		util.MethodTypeGross,
		&svcFeePct,
		false,
		nil,
		values,
		nil,
	)
	if err != nil {
		t.Fatalf("RunEntryCalculation: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["additionalReductionBreakdown"] == nil {
		t.Error("expected additionalReductionBreakdown when form has reduction section")
	}
	if m["totalAdditionalReduction"] == nil {
		t.Error("expected totalAdditionalReduction")
	}
	if m["totalAdditionalReductionBase"] == nil {
		t.Error("expected totalAdditionalReductionBase")
	}
}

func TestRunEntryCalculation_GrossMethod_SummaryFields(t *testing.T) {
	// Gross run produces totalReductionBase, totalExpenseGst for summary
	formBoth := []byte(`[
		{"id":"inc1","name":"Income","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"},
		{"id":"exp1","name":"Expense","type":"number","section":"expense","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"CLINIC"}
	]`)
	values := []byte(`[{"fieldId":"inc1","value":1000},{"fieldId":"exp1","value":200}]`)
	out, err := RunEntryCalculation(
		formBoth,
		util.FormTypeBoth,
		util.MethodTypeGross,
		nil,
		false,
		nil,
		values,
		nil,
	)
	if err != nil {
		t.Fatalf("RunEntryCalculation: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["totalReductionBase"] == nil {
		t.Error("expected totalReductionBase for gross summary")
	}
	if m["totalExpenseGst"] == nil {
		t.Error("expected totalExpenseGst for gross summary")
	}
}

func TestRunEntryCalculation_EmptyFieldId_MatchesByName(t *testing.T) {
	// Values with empty fieldId but fieldName - should match form field by name
	formFields := []byte(`[
		{"id":"f1","name":"Patient Fee","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"},
		{"id":"f2","name":"Lab Fee","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"}
	]`)
	values := []byte(`[{"fieldId":"","fieldName":"Patient Fee","value":"21211"},{"fieldId":"","fieldName":"Lab Fee","value":"212"}]`)
	out, err := RunEntryCalculation(
		formFields,
		util.FormTypeIncome,
		util.MethodTypeGross,
		nil,
		false,
		nil,
		values,
		nil,
	)
	if err != nil {
		t.Fatalf("RunEntryCalculation: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	ft, ok := m["fieldTotals"].([]interface{})
	if !ok || len(ft) != 2 {
		t.Fatalf("fieldTotals want 2 entries, got %v", m["fieldTotals"])
	}
	totalAmount, _ := m["totalAmount"].(float64)
	// 21211+10% + 212+10% = 23332.1 + 233.2 = 23565.3 (approx)
	if totalAmount < 23000 || totalAmount > 24000 {
		t.Errorf("totalAmount want ~23565, got %v", totalAmount)
	}
}
