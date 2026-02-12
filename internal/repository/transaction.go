package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateTransaction(ctx context.Context, db *sqlx.DB, t *domain.Transaction) error {
	query := `INSERT INTO tbl_transaction (
		id, clinic_id, source_entry_id, source_form_id, field_id, coa_id, account_code, account_name, tax_category,
		transaction_date, reference, details, gross_amount, gst_amount, net_amount, status, created_at, updated_at
	) VALUES (
		:id, :clinic_id, :source_entry_id, :source_form_id, :field_id, :coa_id, :account_code, :account_name, :tax_category,
		:transaction_date, :reference, :details, :gross_amount, :gst_amount, :net_amount, :status, :created_at, :updated_at
	)`
	_, err := db.NamedExecContext(ctx, query, t)
	return err
}

func GetTransactionsByClinicID(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID, f *domain.ListTransactionsFilters) ([]domain.Transaction, int, error) {
	base := `FROM tbl_transaction WHERE clinic_id = $1`
	args := []interface{}{clinicID}
	argNum := 2

	if f.Search != "" {
		base += fmt.Sprintf(" AND (account_name ILIKE $%d OR reference ILIKE $%d OR details ILIKE $%d)", argNum, argNum, argNum)
		args = append(args, "%"+f.Search+"%")
		argNum++
	}
	if f.TaxCategory != "" {
		base += fmt.Sprintf(" AND tax_category = $%d", argNum)
		args = append(args, f.TaxCategory)
		argNum++
	}
	if f.Status != "" {
		base += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, f.Status)
		argNum++
	}
	if f.DateFrom != "" {
		base += fmt.Sprintf(" AND transaction_date >= $%d", argNum)
		args = append(args, f.DateFrom)
		argNum++
	}
	if f.DateTo != "" {
		base += fmt.Sprintf(" AND transaction_date <= $%d", argNum)
		args = append(args, f.DateTo)
		argNum++
	}

	sortCol := "transaction_date"
	if f.SortField != "" {
		switch f.SortField {
		case "date": sortCol = "transaction_date"
		case "account": sortCol = "account_name"
		case "reference": sortCol = "reference"
		case "gross": sortCol = "gross_amount"
		case "gst": sortCol = "gst_amount"
		case "net": sortCol = "net_amount"
		}
	}
	sortDir := "DESC"
	if f.SortDirection == "asc" {
		sortDir = "ASC"
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + base
	if err := db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	page := f.Page
	if page < 1 {
		page = 1
	}
	limit := f.Limit
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	sel := `SELECT id, clinic_id, source_entry_id, source_form_id, field_id, coa_id, account_code, account_name, tax_category,
		transaction_date, reference, details, gross_amount, gst_amount, net_amount, status, created_at, updated_at `
	listQuery := sel + base + " ORDER BY " + sortCol + " " + sortDir + fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)

	var list []domain.Transaction
	if err := db.SelectContext(ctx, &list, listQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("list transactions: %w", err)
	}
	if list == nil {
		list = []domain.Transaction{}
	}
	return list, total, nil
}

func GetTransactionsByEntryID(ctx context.Context, db *sqlx.DB, entryID uuid.UUID) ([]domain.Transaction, error) {
	query := `SELECT id, clinic_id, source_entry_id, source_form_id, field_id, coa_id, account_code, account_name, tax_category,
		transaction_date, reference, details, gross_amount, gst_amount, net_amount, status, created_at, updated_at
		FROM tbl_transaction WHERE source_entry_id = $1 ORDER BY transaction_date, account_code`
	var list []domain.Transaction
	if err := db.SelectContext(ctx, &list, query, entryID); err != nil {
		return nil, err
	}
	if list == nil {
		list = []domain.Transaction{}
	}
	return list, nil
}

func DeleteTransactionsByEntryID(ctx context.Context, db *sqlx.DB, entryID uuid.UUID) error {
	res, err := db.ExecContext(ctx, `DELETE FROM tbl_transaction WHERE source_entry_id = $1`, entryID)
	if err != nil {
		return err
	}
	_ = res
	return nil
}

func transactionToResponse(t *domain.Transaction) *domain.TransactionResponse {
	dateStr := t.TransactionDate.Format("2006-01-02")
	return &domain.TransactionResponse{
		ID:            t.ID.String(),
		ClinicID:      t.ClinicID.String(),
		SourceEntryID: t.SourceEntryID.String(),
		SourceFormID:  t.SourceFormID.String(),
		FieldID:       t.FieldID,
		AccountCode:   t.AccountCode,
		AccountName:   t.AccountName,
		TaxCategory:   t.TaxCategory,
		Date:          dateStr,
		Reference:     t.Reference,
		Details:       t.Details,
		GrossAmount:   t.GrossAmount,
		GSTAmount:     t.GSTAmount,
		NetAmount:     t.NetAmount,
		Status:        t.Status,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}

// TransactionToResponse is exported for use by service
func TransactionToResponse(t *domain.Transaction) *domain.TransactionResponse {
	return transactionToResponse(t)
}
