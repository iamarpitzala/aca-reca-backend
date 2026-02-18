package calculation

import (
	"encoding/json"
	"testing"

	"github.com/iamarpitzala/aca-reca-backend/util"
)

// formFieldsJSON for INCOME form with one number field
var testFormIncomeFields = []byte(`[
	{"id":"f1","name":"Fee","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"}
]`)

// formFieldsJSON for BOTH form: income + expense (clinic) + optional reduction section
var testFormBothFields = []byte(`[
	{"id":"inc1","name":"Income","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"},
	{"id":"exp1","name":"Expense","type":"number","section":"expense","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"CLINIC"},
	{"id":"red1","name":"Additional reduction","type":"number","section":"reduction","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"},"paymentResponsibility":"OWNER"}
]`)

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
		"", // calculation method no longer used
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
