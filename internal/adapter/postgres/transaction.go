package postgres

import (
	"context"
	"fmt"
	"time"

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

func (r *transactionRepo) Create(ctx context.Context, t *domain.TransactionWithLedger) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert into tbl_transaction (header)
	headerQuery := `INSERT INTO tbl_transaction (id, clinic_id, source_entry_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, headerQuery, t.ID, t.ClinicID, t.SourceEntryID, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert into tbl_transaction_ledger (details)
	ledgerQuery := `INSERT INTO tbl_transaction_ledger (
		id, transaction_id, coa_id, account_code, account_name, tax_category,
		transaction_date, reference, details, gross_amount, gst_amount, net_amount, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err = tx.ExecContext(ctx, ledgerQuery,
		uuid.New(), t.ID, t.COAID, t.AccountCode, t.AccountName, t.TaxCategory,
		t.TransactionDate, t.Reference, t.Details, t.GrossAmount, t.GSTAmount, t.NetAmount, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *transactionRepo) ListByClinicID(ctx context.Context, clinicID uuid.UUID, f *domain.ListTransactionsFilters) ([]domain.TransactionWithLedger, int, error) {
	base := `FROM tbl_transaction t 
		INNER JOIN tbl_transaction_ledger tl ON t.id = tl.transaction_id 
		INNER JOIN tbl_custom_form_entry e ON t.source_entry_id = e.id
		WHERE t.clinic_id = $1`
	args := []interface{}{clinicID}
	argNum := 2

	if f.Search != "" {
		base += fmt.Sprintf(" AND (tl.account_name ILIKE $%d OR tl.reference ILIKE $%d OR tl.details ILIKE $%d)", argNum, argNum, argNum)
		args = append(args, "%"+f.Search+"%")
		argNum++
	}
	if f.TaxCategory != "" {
		base += fmt.Sprintf(" AND tl.tax_category = $%d", argNum)
		args = append(args, f.TaxCategory)
		argNum++
	}
	if f.Status != "" {
		// Status is not in current schema, skip this filter
		// base += fmt.Sprintf(" AND status = $%d", argNum)
		// args = append(args, f.Status)
		// argNum++
	}
	if f.COAID != "" {
		base += fmt.Sprintf(" AND tl.coa_id = $%d", argNum)
		args = append(args, f.COAID)
		argNum++
	}
	if f.DateFrom != "" {
		base += fmt.Sprintf(" AND tl.transaction_date >= $%d", argNum)
		args = append(args, f.DateFrom)
		argNum++
	}
	if f.DateTo != "" {
		base += fmt.Sprintf(" AND tl.transaction_date <= $%d", argNum)
		args = append(args, f.DateTo)
		argNum++
	}

	sortCol := "tl.transaction_date"
	if f.SortField != "" {
		switch f.SortField {
		case "date":
			sortCol = "tl.transaction_date"
		case "account":
			sortCol = "tl.account_name"
		case "reference":
			sortCol = "tl.reference"
		case "gross":
			sortCol = "tl.gross_amount"
		case "gst":
			sortCol = "tl.gst_amount"
		case "net":
			sortCol = "tl.net_amount"
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

	sel := `SELECT 
		tl.id, t.id as transaction_id, t.clinic_id, t.source_entry_id, 
		e.form_id as source_form_id,
		tl.coa_id, tl.account_code, tl.account_name, tl.tax_category,
		tl.transaction_date, tl.reference, tl.details, 
		tl.gross_amount, tl.gst_amount, tl.net_amount, 
		t.created_at, t.updated_at `
	listQuery := sel + base + " ORDER BY " + sortCol + " " + sortDir + fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)

	type dbRow struct {
		ID              uuid.UUID `db:"id"`
		TransactionID   uuid.UUID `db:"transaction_id"`
		ClinicID        uuid.UUID `db:"clinic_id"`
		SourceEntryID   uuid.UUID `db:"source_entry_id"`
		SourceFormID    uuid.UUID `db:"source_form_id"`
		COAID           uuid.UUID `db:"coa_id"`
		AccountCode     string    `db:"account_code"`
		AccountName     string    `db:"account_name"`
		TaxCategory     string    `db:"tax_category"`
		TransactionDate time.Time `db:"transaction_date"`
		Reference       string    `db:"reference"`
		Details         string    `db:"details"`
		GrossAmount     float64   `db:"gross_amount"`
		GSTAmount       float64   `db:"gst_amount"`
		NetAmount       float64   `db:"net_amount"`
		CreatedAt       time.Time `db:"created_at"`
		UpdatedAt       time.Time `db:"updated_at"`
	}

	var rows []dbRow
	if err := r.db.SelectContext(ctx, &rows, listQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("list transactions: %w", err)
	}

	list := make([]domain.TransactionWithLedger, len(rows))
	for i, row := range rows {
		list[i] = domain.TransactionWithLedger{
			Transaction: domain.Transaction{
				ID:            row.TransactionID,
				ClinicID:      row.ClinicID,
				SourceEntryID: row.SourceEntryID,
				CreatedAt:     row.CreatedAt,
				UpdatedAt:     row.UpdatedAt,
			},
			SourceFormID:    row.SourceFormID,
			COAID:           row.COAID,
			AccountCode:     row.AccountCode,
			AccountName:     row.AccountName,
			TaxCategory:     row.TaxCategory,
			TransactionDate: row.TransactionDate,
			Reference:       row.Reference,
			Details:         row.Details,
			GrossAmount:     row.GrossAmount,
			GSTAmount:       row.GSTAmount,
			NetAmount:       row.NetAmount,
			Status:          "POSTED", // Default status since not in schema
		}
	}

	return list, total, nil
}

func (r *transactionRepo) ListByEntryID(ctx context.Context, entryID uuid.UUID) ([]domain.TransactionWithLedger, error) {
	query := `SELECT 
		tl.id, t.id as transaction_id, t.clinic_id, t.source_entry_id,
		e.form_id as source_form_id,
		tl.coa_id, tl.account_code, tl.account_name, tl.tax_category,
		tl.transaction_date, tl.reference, tl.details,
		tl.gross_amount, tl.gst_amount, tl.net_amount,
		t.created_at, t.updated_at
		FROM tbl_transaction t
		INNER JOIN tbl_transaction_ledger tl ON t.id = tl.transaction_id
		INNER JOIN tbl_custom_form_entry e ON t.source_entry_id = e.id
		WHERE t.source_entry_id = $1 
		ORDER BY tl.transaction_date, tl.account_code`

	type dbRow struct {
		ID              uuid.UUID `db:"id"`
		TransactionID   uuid.UUID `db:"transaction_id"`
		ClinicID        uuid.UUID `db:"clinic_id"`
		SourceEntryID   uuid.UUID `db:"source_entry_id"`
		SourceFormID    uuid.UUID `db:"source_form_id"`
		COAID           uuid.UUID `db:"coa_id"`
		AccountCode     string    `db:"account_code"`
		AccountName     string    `db:"account_name"`
		TaxCategory     string    `db:"tax_category"`
		TransactionDate time.Time `db:"transaction_date"`
		Reference       string    `db:"reference"`
		Details         string    `db:"details"`
		GrossAmount     float64   `db:"gross_amount"`
		GSTAmount       float64   `db:"gst_amount"`
		NetAmount       float64   `db:"net_amount"`
		CreatedAt       time.Time `db:"created_at"`
		UpdatedAt       time.Time `db:"updated_at"`
	}

	var rows []dbRow
	if err := r.db.SelectContext(ctx, &rows, query, entryID); err != nil {
		return nil, err
	}

	list := make([]domain.TransactionWithLedger, len(rows))
	for i, row := range rows {
		list[i] = domain.TransactionWithLedger{
			Transaction: domain.Transaction{
				ID:            row.TransactionID,
				ClinicID:      row.ClinicID,
				SourceEntryID: row.SourceEntryID,
				CreatedAt:     row.CreatedAt,
				UpdatedAt:     row.UpdatedAt,
			},
			SourceFormID:    row.SourceFormID,
			COAID:           row.COAID,
			AccountCode:     row.AccountCode,
			AccountName:     row.AccountName,
			TaxCategory:     row.TaxCategory,
			TransactionDate: row.TransactionDate,
			Reference:       row.Reference,
			Details:         row.Details,
			GrossAmount:     row.GrossAmount,
			GSTAmount:       row.GSTAmount,
			NetAmount:       row.NetAmount,
			Status:          "POSTED",
		}
	}

	return list, nil
}

func (r *transactionRepo) DeleteByEntryID(ctx context.Context, entryID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tbl_transaction WHERE source_entry_id = $1`, entryID)
	return err
}
