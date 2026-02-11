package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type transactionRepo struct {
	db *sqlx.DB
}

// NewTransactionRepository returns a Postgres implementation of TransactionRepository.
func NewTransactionRepository(db *sqlx.DB) port.TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, t *domain.Transaction) error {
	query := `INSERT INTO tbl_transaction (
		id, clinic_id, source_entry_id, source_form_id, field_id, coa_id, account_code, account_name, tax_category,
		transaction_date, reference, details, gross_amount, gst_amount, net_amount, status, created_at, updated_at
	) VALUES (
		:id, :clinic_id, :source_entry_id, :source_form_id, :field_id, :coa_id, :account_code, :account_name, :tax_category,
		:transaction_date, :reference, :details, :gross_amount, :gst_amount, :net_amount, :status, :created_at, :updated_at
	)`
	_, err := r.db.NamedExecContext(ctx, query, t)
	return err
}

func (r *transactionRepo) ListByClinicID(ctx context.Context, clinicID uuid.UUID, f *domain.ListTransactionsFilters) ([]domain.Transaction, int, error) {
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
		case "date":
			sortCol = "transaction_date"
		case "account":
			sortCol = "account_name"
		case "reference":
			sortCol = "reference"
		case "gross":
			sortCol = "gross_amount"
		case "gst":
			sortCol = "gst_amount"
		case "net":
			sortCol = "net_amount"
		}
	}
	sortDir := "DESC"
	if f.SortDirection == "asc" {
		sortDir = "ASC"
	}

	var total int
	countQuery := "SELECT COUNT(*) " + base
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
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
	if err := r.db.SelectContext(ctx, &list, listQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("list transactions: %w", err)
	}
	if list == nil {
		list = []domain.Transaction{}
	}
	return list, total, nil
}

func (r *transactionRepo) ListByEntryID(ctx context.Context, entryID uuid.UUID) ([]domain.Transaction, error) {
	query := `SELECT id, clinic_id, source_entry_id, source_form_id, field_id, coa_id, account_code, account_name, tax_category,
		transaction_date, reference, details, gross_amount, gst_amount, net_amount, status, created_at, updated_at
		FROM tbl_transaction WHERE source_entry_id = $1 ORDER BY transaction_date, account_code`
	var list []domain.Transaction
	if err := r.db.SelectContext(ctx, &list, query, entryID); err != nil {
		return nil, err
	}
	if list == nil {
		list = []domain.Transaction{}
	}
	return list, nil
}

func (r *transactionRepo) DeleteByEntryID(ctx context.Context, entryID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tbl_transaction WHERE source_entry_id = $1`, entryID)
	return err
}
