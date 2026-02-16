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
