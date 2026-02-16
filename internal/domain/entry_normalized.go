package domain

import (
	"time"

	"github.com/google/uuid"
)

// Normalized entry table models

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

// EntryGrossDetails represents tbl_entry_gross_details
type EntryGrossDetails struct {
	ID                      uuid.UUID `db:"id"`
	EntryID                 uuid.UUID `db:"source_entry_id"`
	ServiceFacilityFeePercent float64  `db:"service_facility_fee_percent"`
	ServiceFeeBase          float64   `db:"service_fee_base"`
	GstOnServiceFee         float64   `db:"gst_on_service_fee"`
	TotalServiceFee         float64   `db:"total_service_fee"`
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
}

// EntryGrossReduction represents tbl_entry_gross_reduction
type EntryGrossReduction struct {
	ID              uuid.UUID `db:"id"`
	GrossDetailsID  uuid.UUID `db:"gross_details_id"`
	EntryID         uuid.UUID `db:"source_entry_id"`
	FieldID         uuid.UUID `db:"tbl_custom_form_field_id"`
	BaseAmount      float64   `db:"base_amount"`
	GstAmount       float64   `db:"gst_amount"`
	TotalAmount     float64   `db:"total_amount"`
	CreatedAt       time.Time `db:"created_at"`
}

// EntryGrossReimbursement represents tbl_entry_gross_reimbursement
type EntryGrossReimbursement struct {
	ID              uuid.UUID `db:"id"`
	GrossDetailsID  uuid.UUID `db:"gross_details_id"`
	EntryID         uuid.UUID `db:"source_entry_id"`
	FieldID         uuid.UUID `db:"tbl_custom_form_field_id"`
	BaseAmount      float64   `db:"base_amount"`
	GstAmount       float64   `db:"gst_amount"`
	TotalAmount     float64   `db:"total_amount"`
	CreatedAt       time.Time `db:"created_at"`
}

// NormalizedEntry represents a complete normalized entry with all related data
type NormalizedEntry struct {
	Header            *FieldEntry
	FieldValues       []EntryFieldValue
	FieldCalculations []EntryFieldCalculation
	Deductions        *EntryDeductions
}
