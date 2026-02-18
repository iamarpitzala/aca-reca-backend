package util

// Transaction status (tbl_transaction.status)
const (
	TransactionStatusPosted = "POSTED"
	TransactionStatusDraft  = "DRAFT"
	TransactionStatusVoided = "VOIDED"
)

// Form status (tbl_custom_form.status)
const (
	FormStatusDraft     = "DRAFT"
	FormStatusPublished = "PUBLISHED"
	FormStatusArchived  = "ARCHIVED"
)

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

// BAS snapshot status (tbl_bas_snapshot.status)
const (
	BASStatusDraft     = "DRAFT"
	BASStatusFinalised = "FINALISED"
	BASStatusLocked    = "LOCKED"
)

// BAS period types (tbl_bas_snapshot.period_type)
const (
	BASPeriodTypeQuarterly = "QUARTERLY"
	BASPeriodTypeAnnually  = "ANNUALLY"
)
