package domain

import (
	"time"

	"github.com/google/uuid"
)

// PnlReport represents tbl_pnl_report (header)
type PnlReport struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	ClinicID          uuid.UUID  `db:"clinic_id" json:"clinicId"`
	QuarterID         *uuid.UUID `db:"quarter_id" json:"quarterId,omitempty"`
	ReportName        string     `db:"report_name" json:"reportName"`
	PeriodStart       time.Time  `db:"period_start" json:"periodStart"`
	PeriodEnd         time.Time  `db:"period_end" json:"periodEnd"`
	TotalRevenue      float64    `db:"total_revenue" json:"totalRevenue"`
	TotalCOGS         float64    `db:"total_cogs" json:"totalCogs"`
	GrossProfit       float64    `db:"gross_profit" json:"grossProfit"`
	TotalExpenses     float64    `db:"total_expenses" json:"totalExpenses"`
	NetProfit         float64    `db:"net_profit" json:"netProfit"`
	TotalGSTCollected float64    `db:"total_gst_collected" json:"totalGstCollected"`
	TotalGSTPaid      float64    `db:"total_gst_paid" json:"totalGstPaid"`
	NetGST            float64    `db:"net_gst" json:"netGst"`
	Status            string     `db:"status" json:"status"`
	GeneratedBy       uuid.UUID  `db:"generated_by" json:"generatedBy"`
	FinalizedAt       *time.Time `db:"finalized_at" json:"finalizedAt,omitempty"`
	Notes             *string    `db:"notes" json:"notes,omitempty"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// PnlReportLine represents tbl_pnl_report_line (line item)
type PnlReportLine struct {
	ID               uuid.UUID `db:"id" json:"id"`
	PnlReportID      uuid.UUID `db:"pnl_report_id" json:"pnlReportId"`
	COAID            uuid.UUID `db:"coa_id" json:"coaId"`
	AccountCode      string    `db:"account_code" json:"accountCode"`
	AccountName      string    `db:"account_name" json:"accountName"`
	LineCategory     string    `db:"line_category" json:"lineCategory"`
	DebitTotal       float64   `db:"debit_total" json:"debitTotal"`
	CreditTotal      float64   `db:"credit_total" json:"creditTotal"`
	NetAmount        float64   `db:"net_amount" json:"netAmount"`
	GSTAmount        float64   `db:"gst_amount" json:"gstAmount"`
	TransactionCount int       `db:"transaction_count" json:"transactionCount"`
	DisplayOrder     int       `db:"display_order" json:"displayOrder"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}

// PnlReportWithLines combines header + line items for API responses
type PnlReportWithLines struct {
	Report PnlReport       `json:"report"`
	Lines  []PnlReportLine `json:"lines"`
}

// ---- Request DTOs ----

// GeneratePnlRequest is the request body for POST /reports/pnl/generate
type GeneratePnlRequest struct {
	ClinicID    uuid.UUID  `json:"clinicId" validate:"required"`
	QuarterID   *uuid.UUID `json:"quarterId" validate:"omitempty"`
	ReportName  string     `json:"reportName" validate:"required,max=255"`
	PeriodStart string     `json:"periodStart" validate:"required"`
	PeriodEnd   string     `json:"periodEnd" validate:"required"`
	Notes       *string    `json:"notes" validate:"omitempty,max=1000"`
}

// ---- Response DTOs ----

// PnlReportResponse wraps the report + lines for API output
type PnlReportResponse struct {
	Report PnlReport       `json:"report"`
	Lines  []PnlReportLine `json:"lines"`
}

// PnlReportListResponse wraps a list of reports for API output
type PnlReportListResponse struct {
	Reports []PnlReport `json:"reports"`
	Total   int         `json:"total"`
}

// LedgerAggRow is used internally to aggregate ledger data during report generation
type LedgerAggRow struct {
	COAID            uuid.UUID `db:"coa_id"`
	AccountCode      string    `db:"account_code"`
	AccountName      string    `db:"account_name"`
	AccountType      string    `db:"account_type"`
	DebitTotal       float64   `db:"debit_total"`
	CreditTotal      float64   `db:"credit_total"`
	NetAmount        float64   `db:"net_amount"`
	GSTAmount        float64   `db:"gst_amount"`
	TransactionCount int       `db:"transaction_count"`
}
