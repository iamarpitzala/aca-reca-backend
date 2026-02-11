package domain

import (
	"time"

	"github.com/google/uuid"
)

// Normalized entry table models

// EntryHeader represents tbl_entry_header
type EntryHeader struct {
	ID                    uuid.UUID  `db:"id"`
	FormID                uuid.UUID  `db:"form_id"`
	FormName              string     `db:"form_name"`
	FormType              string     `db:"form_type"`
	CalculationMethod     string     `db:"calculation_method"`
	ClinicID              uuid.UUID  `db:"clinic_id"`
	QuarterID             *uuid.UUID `db:"quarter_id"`
	EntryDate             time.Time  `db:"entry_date"`
	Description           string     `db:"description"`
	Remarks               string     `db:"remarks"`
	PaymentResponsibility *string    `db:"payment_responsibility"`
	CreatedBy             uuid.UUID  `db:"created_by"`
	CreatedAt             time.Time  `db:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at"`
	DeletedAt             *time.Time `db:"deleted_at"`
	OriginalEntryID       *uuid.UUID `db:"original_entry_id"`
}

// EntryFieldValue represents tbl_entry_field_value
type EntryFieldValue struct {
	ID              uuid.UUID `db:"id"`
	EntryID         uuid.UUID `db:"entry_id"`
	FieldID         string    `db:"field_id"`
	FieldName       string    `db:"field_name"`
	Value           *float64  `db:"value"`
	TextValue       *string   `db:"text_value"`
	BooleanValue    *bool     `db:"boolean_value"`
	ManualGstAmount *float64  `db:"manual_gst_amount"`
	DisplayOrder    int       `db:"display_order"`
	CreatedAt       time.Time `db:"created_at"`
}

// EntryFieldCalculation represents tbl_entry_field_calculation
type EntryFieldCalculation struct {
	ID                    uuid.UUID `db:"id"`
	EntryID               uuid.UUID `db:"entry_id"`
	FieldID               string    `db:"field_id"`
	FieldName             string    `db:"field_name"`
	BaseAmount            float64   `db:"base_amount"`
	GstAmount             float64   `db:"gst_amount"`
	TotalAmount           float64   `db:"total_amount"`
	GstRate               float64   `db:"gst_rate"`
	GstType               string    `db:"gst_type"`
	Section               *string   `db:"section"`
	PaymentResponsibility *string   `db:"payment_responsibility"`
	DisplayOrder          int       `db:"display_order"`
	CreatedAt             time.Time `db:"created_at"`
}

// EntrySummary represents tbl_entry_summary
type EntrySummary struct {
	ID              uuid.UUID `db:"id"`
	EntryID         uuid.UUID `db:"entry_id"`
	TotalBaseAmount float64   `db:"total_base_amount"`
	TotalGstAmount  float64   `db:"total_gst_amount"`
	TotalAmount     float64   `db:"total_amount"`
	NetPayable      float64   `db:"net_payable"`
	NetReceivable   float64   `db:"net_receivable"`
	NetFee          *float64  `db:"net_fee"`
	BasGstOnSales1A float64   `db:"bas_gst_on_sales_1a"`
	BasGstCredit1B  float64   `db:"bas_gst_credit_1b"`
	BasTotalSalesG1 float64   `db:"bas_total_sales_g1"`
	BasExpensesG11  float64   `db:"bas_expenses_g11"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// EntryNetDetails represents tbl_entry_net_details
type EntryNetDetails struct {
	ID                     uuid.UUID `db:"id"`
	EntryID                uuid.UUID `db:"entry_id"`
	CommissionPercent      float64   `db:"commission_percent"`
	Commission             float64   `db:"commission"`
	GstOnCommission        float64   `db:"gst_on_commission"`
	TotalPaymentReceived   float64   `db:"total_payment_received"`
	SuperHoldingEnabled    bool      `db:"super_holding_enabled"`
	SuperComponentPercent  *float64  `db:"super_component_percent"`
	CommissionComponent    *float64  `db:"commission_component"`
	SuperComponent         *float64  `db:"super_component"`
	TotalForReconciliation *float64  `db:"total_for_reconciliation"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedAt              time.Time `db:"updated_at"`
}

// EntryGrossDetails represents tbl_entry_gross_details
type EntryGrossDetails struct {
	ID                        uuid.UUID `db:"id"`
	EntryID                   uuid.UUID `db:"entry_id"`
	ServiceFacilityFeePercent float64   `db:"service_facility_fee_percent"`
	ServiceFeeBase            float64   `db:"service_fee_base"`
	GstOnServiceFee           float64   `db:"gst_on_service_fee"`
	TotalServiceFee           float64   `db:"total_service_fee"`
	SubtotalAfterDeductions   *float64  `db:"subtotal_after_deductions"`
	RemittedAmount            *float64  `db:"remitted_amount"`
	CreatedAt                 time.Time `db:"created_at"`
	UpdatedAt                 time.Time `db:"updated_at"`
}

// EntryGrossReduction represents tbl_entry_gross_reduction
type EntryGrossReduction struct {
	ID                 uuid.UUID  `db:"id"`
	EntryID            uuid.UUID  `db:"entry_id"`
	FieldCalculationID *uuid.UUID `db:"field_calculation_id"`
	FieldID            string     `db:"field_id"`
	FieldName          string     `db:"field_name"`
	BaseAmount         float64    `db:"base_amount"`
	GstAmount          float64    `db:"gst_amount"`
	TotalAmount        float64    `db:"total_amount"`
	DisplayOrder       int        `db:"display_order"`
	CreatedAt          time.Time  `db:"created_at"`
}

// EntryGrossReimbursement represents tbl_entry_gross_reimbursement
type EntryGrossReimbursement struct {
	ID                 uuid.UUID  `db:"id"`
	EntryID            uuid.UUID  `db:"entry_id"`
	FieldCalculationID *uuid.UUID `db:"field_calculation_id"`
	FieldID            string     `db:"field_id"`
	FieldName          string     `db:"field_name"`
	BaseAmount         float64    `db:"base_amount"`
	GstAmount          float64    `db:"gst_amount"`
	TotalAmount        float64    `db:"total_amount"`
	DisplayOrder       int        `db:"display_order"`
	CreatedAt          time.Time  `db:"created_at"`
}

// EntryGrossAdditionalReduction represents tbl_entry_gross_additional_reduction
type EntryGrossAdditionalReduction struct {
	ID                 uuid.UUID  `db:"id"`
	EntryID            uuid.UUID  `db:"entry_id"`
	FieldCalculationID *uuid.UUID `db:"field_calculation_id"`
	FieldID            string     `db:"field_id"`
	FieldName          string     `db:"field_name"`
	BaseAmount         float64    `db:"base_amount"`
	GstAmount          float64    `db:"gst_amount"`
	TotalAmount        float64    `db:"total_amount"`
	DisplayOrder       int        `db:"display_order"`
	CreatedAt          time.Time  `db:"created_at"`
}

// EntryGrossReductionsSummary represents tbl_entry_gross_reductions_summary
type EntryGrossReductionsSummary struct {
	ID                           uuid.UUID `db:"id"`
	EntryID                      uuid.UUID `db:"entry_id"`
	TotalReductions              float64   `db:"total_reductions"`
	TotalReductionBase           float64   `db:"total_reduction_base"`
	TotalExpenseGst              float64   `db:"total_expense_gst"`
	TotalReimbursements          float64   `db:"total_reimbursements"`
	TotalAdditionalReduction     float64   `db:"total_additional_reduction"`
	TotalAdditionalReductionBase float64   `db:"total_additional_reduction_base"`
	TotalAdditionalReductionGst  float64   `db:"total_additional_reduction_gst"`
	CreatedAt                    time.Time `db:"created_at"`
	UpdatedAt                    time.Time `db:"updated_at"`
}

// EntryGrossOutwork represents tbl_entry_gross_outwork
type EntryGrossOutwork struct {
	ID                 uuid.UUID `db:"id"`
	EntryID            uuid.UUID `db:"entry_id"`
	OutworkEnabled     bool      `db:"outwork_enabled"`
	OutworkRatePercent *float64  `db:"outwork_rate_percent"`
	OutworkChargeBase  float64   `db:"outwork_charge_base"`
	OutworkChargeGst   float64   `db:"outwork_charge_gst"`
	OutworkChargeTotal float64   `db:"outwork_charge_total"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

// EntryDeductions represents tbl_entry_deductions
type EntryDeductions struct {
	ID                         uuid.UUID `db:"id"`
	EntryID                    uuid.UUID `db:"entry_id"`
	ServiceFacilityFeePercent  *float64  `db:"service_facility_fee_percent"`
	ServiceFeeOverride         *float64  `db:"service_fee_override"`
	CommissionPercent          *float64  `db:"commission_percent"`
	SuperHoldingEnabled        *bool     `db:"super_holding_enabled"`
	SuperComponentPercent      *float64  `db:"super_component_percent"`
	OutworkEnabled             *bool     `db:"outwork_enabled"`
	OutworkRatePercent         *float64  `db:"outwork_rate_percent"`
	EntryPaymentResponsibility *string   `db:"entry_payment_responsibility"`
	CreatedAt                  time.Time `db:"created_at"`
}

// NormalizedEntry represents a complete normalized entry with all related data
type NormalizedEntry struct {
	Header                    *EntryHeader
	FieldValues               []EntryFieldValue
	FieldCalculations         []EntryFieldCalculation
	Summary                   *EntrySummary
	NetDetails                *EntryNetDetails
	GrossDetails              *EntryGrossDetails
	GrossReductions           []EntryGrossReduction
	GrossReimbursements       []EntryGrossReimbursement
	GrossAdditionalReductions []EntryGrossAdditionalReduction
	GrossReductionsSummary    *EntryGrossReductionsSummary
	GrossOutwork              *EntryGrossOutwork
	Deductions                *EntryDeductions
}
