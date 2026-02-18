package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type transactionRepo struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) port.TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, txn *domain.Transaction) error {
	q := `
		INSERT INTO tbl_transaction (
			id, clinic_id, source_entry_id, reference_number,
			description, transaction_date, status, created_by,
			posted_at, voided_at, void_reason,
			created_at, updated_at
		) VALUES (
			:id, :clinic_id, :source_entry_id, :reference_number,
			:description, :transaction_date, :status, :created_by,
			:posted_at, :voided_at, :void_reason,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, q, txn)
	return err
}

func (r *transactionRepo) CreateLedgerLines(ctx context.Context, lines []domain.TransactionLedger) error {
	if len(lines) == 0 {
		return nil
	}
	q := `
		INSERT INTO tbl_transaction_ledger (
			id, transaction_id, coa_id, entry_type,
			amount, gst_amount, net_amount,
			transaction_date, description, created_at
		) VALUES (
			:id, :transaction_id, :coa_id, :entry_type,
			:amount, :gst_amount, :net_amount,
			:transaction_date, :description, :created_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, q, lines)
	return err
}

func (r *transactionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	q := `
		SELECT id, clinic_id, source_entry_id, reference_number,
		       description, transaction_date, status, created_by,
		       posted_at, voided_at, void_reason,
		       created_at, updated_at, deleted_at
		FROM tbl_transaction
		WHERE id = $1 AND deleted_at IS NULL
	`
	var txn domain.Transaction
	if err := r.db.GetContext(ctx, &txn, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get transaction")
	}
	return &txn, nil
}

func (r *transactionRepo) GetLedgerByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]domain.TransactionLedger, error) {
	q := `
		SELECT id, transaction_id, coa_id, entry_type,
		       amount, gst_amount, net_amount,
		       transaction_date, description, created_at
		FROM tbl_transaction_ledger
		WHERE transaction_id = $1
		ORDER BY entry_type, created_at
	`
	var lines []domain.TransactionLedger
	if err := r.db.SelectContext(ctx, &lines, q, transactionID); err != nil {
		return nil, errors.New("failed to get transaction ledger lines")
	}
	return lines, nil
}

func (r *transactionRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.Transaction, error) {
	q := `
		SELECT id, clinic_id, source_entry_id, reference_number,
		       description, transaction_date, status, created_by,
		       posted_at, voided_at, void_reason,
		       created_at, updated_at, deleted_at
		FROM tbl_transaction
		WHERE clinic_id = $1 AND deleted_at IS NULL
		ORDER BY transaction_date DESC, created_at DESC
	`
	var txns []domain.Transaction
	if err := r.db.SelectContext(ctx, &txns, q, clinicID); err != nil {
		return nil, errors.New("failed to list transactions for clinic")
	}
	return txns, nil
}

func (r *transactionRepo) Update(ctx context.Context, txn *domain.Transaction) error {
	q := `
		UPDATE tbl_transaction SET
			reference_number = :reference_number,
			description = :description,
			transaction_date = :transaction_date,
			status = :status,
			posted_at = :posted_at,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`
	res, err := r.db.NamedExecContext(ctx, q, txn)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

func (r *transactionRepo) DeleteLedgerByTransactionID(ctx context.Context, transactionID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM tbl_transaction_ledger WHERE transaction_id = $1`,
		transactionID,
	)
	return err
}

func (r *transactionRepo) VoidTransaction(ctx context.Context, id uuid.UUID, reason string) error {
	now := time.Now()
	q := `
		UPDATE tbl_transaction SET
			status = 'VOIDED',
			voided_at = $1,
			void_reason = $2,
			updated_at = $1
		WHERE id = $3 AND deleted_at IS NULL AND status != 'VOIDED'
	`
	res, err := r.db.ExecContext(ctx, q, now, reason, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("transaction not found or already voided")
	}
	return nil
}

func (r *transactionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	q := `
		UPDATE tbl_transaction SET
			deleted_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("transaction not found")
	}
	return nil
}
