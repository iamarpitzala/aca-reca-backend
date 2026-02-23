package util

// Transaction status (tbl_transaction.status)
const (
	TransactionStatusPosted = "POSTED"
	TransactionStatusDraft  = "DRAFT"
	TransactionStatusVoided = "VOIDED"
)

// Form status (tbl_custom_form.status)

// Quarter status (calculated quarter status in APIs)
const (
	QuarterStatusOpen   = "OPEN"
	QuarterStatusLocked = "LOCKED"
	QuarterStatusDraft  = "DRAFT"
)

// P&L report status (tbl_pnl_report.status)
const (
	PnlReportStatusDraft    = "DRAFT"
	PnlReportStatusFinal    = "FINAL"
	PnlReportStatusArchived = "ARCHIVED"
)
