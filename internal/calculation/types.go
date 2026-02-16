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

// Deductions from request (for service fee % and override; entry-level payment responsibility overrides per-field when set)
type deductionsInput struct {
	ServiceFacilityFeePercent  *float64 `json:"serviceFacilityFeePercent"`
	ServiceFeeOverride         *float64 `json:"serviceFeeOverride"`
	EntryPaymentResponsibility *string  `json:"entryPaymentResponsibility"`
	CommissionPercent          *float64 `json:"commissionPercent"`
	SuperHoldingEnabled        *bool    `json:"superHoldingEnabled"`
	SuperComponentPercent      *float64 `json:"superComponentPercent"`
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
	FieldTotals                  []fieldCalc `json:"fieldTotals"`
	TotalBaseAmount              float64     `json:"totalBaseAmount"`
	TotalGSTAmount              float64     `json:"totalGSTAmount"`
	TotalAmount                  float64     `json:"totalAmount"`
	NetPayable                   float64     `json:"netPayable"`
	NetReceivable                float64     `json:"netReceivable"`
	BasMapping                   basMapping  `json:"basMapping"`
	NetFee                       *float64    `json:"netFee,omitempty"`
	ServiceFeeBase               *float64    `json:"serviceFeeBase,omitempty"`
	GstOnServiceFee              *float64    `json:"gstOnServiceFee,omitempty"`
	TotalServiceFee              *float64    `json:"totalServiceFee,omitempty"`
	TotalReductions              *float64    `json:"totalReductions,omitempty"`
	TotalReimbursements          *float64    `json:"totalReimbursements,omitempty"`
	TotalReductionBase           *float64    `json:"totalReductionBase,omitempty"`
	TotalExpenseGst              *float64    `json:"totalExpenseGst,omitempty"`
	ReductionBreakdown           []fieldCalc  `json:"reductionBreakdown,omitempty"`
	ReimbursementBreakdown       []fieldCalc  `json:"reimbursementBreakdown,omitempty"`
	AdditionalReductionBreakdown []fieldCalc  `json:"additionalReductionBreakdown,omitempty"`
	TotalAdditionalReduction     *float64    `json:"totalAdditionalReduction,omitempty"`
	TotalAdditionalReductionBase *float64    `json:"totalAdditionalReductionBase,omitempty"`
	TotalAdditionalReductionGst  *float64    `json:"totalAdditionalReductionGst,omitempty"`
	SubtotalAfterDeductions      *float64    `json:"subtotalAfterDeductions,omitempty"`
	RemittedAmount               *float64    `json:"remittedAmount,omitempty"`
	OutworkEnabled               bool        `json:"outworkEnabled"`
	OutworkRatePercent           *float64    `json:"outworkRatePercent,omitempty"`
	OutworkChargeBase            *float64    `json:"outworkChargeBase,omitempty"`
	OutworkChargeGst             *float64    `json:"outworkChargeGst,omitempty"`
	OutworkChargeTotal           *float64    `json:"outworkChargeTotal,omitempty"`
	Commission                   *float64    `json:"commission,omitempty"`
	GstOnCommission              *float64    `json:"gstOnCommission,omitempty"`
	CommissionComponent          *float64    `json:"commissionComponent,omitempty"`
	SuperComponent               *float64    `json:"superComponent,omitempty"`
	TotalForReconciliation        *float64    `json:"totalForReconciliation,omitempty"`
	TotalPaymentReceived         *float64    `json:"totalPaymentReceived,omitempty"`
}
