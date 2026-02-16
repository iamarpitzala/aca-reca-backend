package calculation

// Field definitions for parsing form.Fields JSONB (subset needed for calculation)
type calcField struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Section        string  `json:"section"`
	IncludeInTotal bool    `json:"includeInTotal"`
	GstConfig      *gstCfg `json:"gstConfig"`
	PaymentResp    string  `json:"paymentResponsibility"`
}

type gstCfg struct {
	Enabled bool    `json:"enabled"`
	Rate    float64 `json:"rate"`
	Type    string  `json:"type"` // inclusive, exclusive, manual
}

// Value from client: fieldId, value, optional manualGstAmount
type entryValue struct {
	FieldID         string      `json:"fieldId"`
	FieldName       string      `json:"fieldName"`
	Value           interface{} `json:"value"`
	ManualGstAmount *float64    `json:"manualGstAmount"`
}

// Deductions from request (entry-level payment responsibility overrides per-field when set)
type deductionsInput struct {
	EntryPaymentResponsibility *string `json:"entryPaymentResponsibility"`
}

// Output structures (match frontend EntryCalculations)
type fieldCalc struct {
	FieldID     string  `json:"fieldId"`
	FieldName   string  `json:"fieldName"`
	BaseAmount  float64 `json:"baseAmount"`
	GstAmount   float64 `json:"gstAmount"`
	TotalAmount float64 `json:"totalAmount"`
	GstRate     float64 `json:"gstRate"`
	GstType     string  `json:"gstType"`
}

type basMapping struct {
	GstOnSales1A float64 `json:"gstOnSales1A"`
	GstCredit1B  float64 `json:"gstCredit1B"`
	TotalSalesG1 float64 `json:"totalSalesG1"`
	ExpensesG11  float64 `json:"expensesG11"`
}

type calculationsOutput struct {
	FieldTotals     []fieldCalc `json:"fieldTotals"`
	TotalBaseAmount float64     `json:"totalBaseAmount"`
	TotalGSTAmount  float64     `json:"totalGSTAmount"`
	TotalAmount     float64     `json:"totalAmount"`
	NetPayable      float64     `json:"netPayable"`
	NetReceivable   float64     `json:"netReceivable"`
	BasMapping      basMapping  `json:"basMapping"`
	NetFee          *float64    `json:"netFee,omitempty"`
}

// Gross calculation output types
// Matches standard naming conventions: camelCase for JSON tags (lowercase gst to match basMapping),
// struct fields use capital GST for consistency with TotalGSTAmount pattern
type grossCalculationOutput struct {
	FieldTotals                   []fieldCalc `json:"fieldTotals"`
	IncomeExclGST                 float64     `json:"incomeExclGST"`
	TotalNetExpenses              float64     `json:"totalNetExpenses"`
	TotalExpensesGST              float64     `json:"totalExpensesGST"`
	PayByDentistExpenses          float64     `json:"payByDentistExpenses"`
	NetAmount                     float64     `json:"netAmount"`
	ServiceAndFacilityFee         float64     `json:"serviceAndFacilityFee"`
	LabFeeWithOtherCost           *float64    `json:"labFeeWithOtherCost,omitempty"`
	TotalServiceAndFacility       *float64    `json:"totalServiceAndFacility,omitempty"`
	GSTOnServiceFee               float64     `json:"GSTOnServiceFee"`
	TotalServiceAndFacilityIncGST *float64    `json:"totalServiceAndFacilityIncGST,omitempty"`
	TotalServiceAndFacilityFee    *float64    `json:"totalServiceAndFacilityFee,omitempty"`
	RemittedCost                  float64     `json:"remittedCost"`
	AmountPayableToDentist        float64     `json:"amountPayableToDentist"`
	BasMapping                    basMapping  `json:"basMapping"`
}
