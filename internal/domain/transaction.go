package domain

import (
	"time"
)

// Transaction represents tbl_transaction (header)
type Transaction struct {
	ID              string     `db:"id" json:"id"`
	ClinicID        string     `db:"clinic_id" json:"clinicId"`
	SourceEntryID   *string    `db:"source_entry_id" json:"sourceEntryId,omitempty"`
	ReferenceNumber *string    `db:"reference_number" json:"referenceNumber,omitempty"`
	Description     *string    `db:"description" json:"description,omitempty"`
	TransactionDate time.Time  `db:"transaction_date" json:"transactionDate"`
	Status          string     `db:"status" json:"status"`
	CreatedBy       string     `db:"created_by" json:"createdBy"`
	PostedAt        *time.Time `db:"posted_at" json:"postedAt,omitempty"`
	VoidedAt        *time.Time `db:"voided_at" json:"voidedAt,omitempty"`
	VoidReason      *string    `db:"void_reason" json:"voidReason,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// TransactionLedger represents tbl_transaction_ledger (line item)
type TransactionLedger struct {
	ID              string    `db:"id" json:"id"`
	TransactionID   string    `db:"transaction_id" json:"transactionId"`
	COAID           string    `db:"coa_id" json:"coaId"`
	EntryType       string    `db:"entry_type" json:"entryType"`
	Amount          float64   `db:"amount" json:"amount"`
	GSTAmount       float64   `db:"gst_amount" json:"gstAmount"`
	NetAmount       float64   `db:"net_amount" json:"netAmount"`
	TransactionDate time.Time `db:"transaction_date" json:"transactionDate"`
	Description     *string   `db:"description" json:"description,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
}

// TransactionWithLedger combines header + line items for API responses
type TransactionWithLedger struct {
	Transaction Transaction         `json:"transaction"`
	Ledger      []TransactionLedger `json:"ledger"`
}

// ---- Request DTOs ----

// LedgerLineRequest represents a single ledger line in a create/update request
type LedgerLineRequest struct {
	COAID       string  `json:"coaId" validate:"required"`
	EntryType   string  `json:"entryType" validate:"required,oneof=DEBIT CREDIT"`
	Amount      float64 `json:"amount" validate:"required,gte=0"`
	GSTAmount   float64 `json:"gstAmount" validate:"gte=0"`
	NetAmount   float64 `json:"netAmount"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// CreateTransactionRequest is the request body for POST /transaction
type CreateTransactionRequest struct {
	ClinicID        string              `json:"clinicId" validate:"required"`
	SourceEntryID   *string             `json:"sourceEntryId" validate:"omitempty"`
	ReferenceNumber *string             `json:"referenceNumber" validate:"omitempty,max=50"`
	Description     *string             `json:"description" validate:"omitempty,max=500"`
	TransactionDate string              `json:"transactionDate" validate:"required"`
	Status          string              `json:"status" validate:"omitempty,oneof=DRAFT POSTED"`
	Ledger          []LedgerLineRequest `json:"ledger" validate:"required,min=1,dive"`
}

// UpdateTransactionRequest is the request body for PUT /transaction/:id
type UpdateTransactionRequest struct {
	ReferenceNumber *string             `json:"referenceNumber" validate:"omitempty,max=50"`
	Description     *string             `json:"description" validate:"omitempty,max=500"`
	TransactionDate *string             `json:"transactionDate" validate:"omitempty"`
	Ledger          []LedgerLineRequest `json:"ledger" validate:"omitempty,min=1,dive"`
}

// VoidTransactionRequest is the request body for POST /transaction/:id/void
type VoidTransactionRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

// ---- Response DTOs ----

// TransactionResponse wraps the transaction + ledger for API output
type TransactionResponse struct {
	Transaction Transaction         `json:"transaction"`
	Ledger      []TransactionLedger `json:"ledger"`
}

// TransactionListResponse wraps a list of transactions for API output
type TransactionListResponse struct {
	Transactions []TransactionWithLedger `json:"transactions"`
	Total        int                     `json:"total"`
}
