package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type entryGrossReimbursementRepo struct {
	db *sqlx.DB
}

func NewEntryGrossReimbursementRepository(db *sqlx.DB) port.EntryGrossReimbursementRepository {
	return &entryGrossReimbursementRepo{db: db}
}

func (r *entryGrossReimbursementRepo) Create(ctx context.Context, reimbursement *domain.EntryGrossReimbursement) error {
	q := `INSERT INTO tbl_entry_gross_reimbursement (
		id, gross_details_id, source_entry_id, tbl_custom_form_field_id,
		base_amount, gst_amount, total_amount, created_at
	) VALUES (
		:id, :gross_details_id, :source_entry_id, :tbl_custom_form_field_id,
		:base_amount, :gst_amount, :total_amount, :created_at
	)`
	_, err := r.db.NamedExecContext(ctx, q, reimbursement)
	return err
}

func (r *entryGrossReimbursementRepo) CreateBatch(ctx context.Context, reimbursements []*domain.EntryGrossReimbursement) error {
	if len(reimbursements) == 0 {
		return nil
	}

	// Use transaction for batch insert to ensure all or nothing
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	q := `INSERT INTO tbl_entry_gross_reimbursement (
		id, gross_details_id, source_entry_id, tbl_custom_form_field_id,
		base_amount, gst_amount, total_amount, created_at
	) VALUES (
		:id, :gross_details_id, :source_entry_id, :tbl_custom_form_field_id,
		:base_amount, :gst_amount, :total_amount, :created_at
	)`

	// Insert each reimbursement individually within the transaction
	for _, reimbursement := range reimbursements {
		_, err := tx.NamedExecContext(ctx, q, reimbursement)
		if err != nil {
			return fmt.Errorf("failed to create reimbursement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *entryGrossReimbursementRepo) GetByGrossDetailsID(ctx context.Context, grossDetailsID uuid.UUID) ([]domain.EntryGrossReimbursement, error) {
	q := `
		SELECT
			id, gross_details_id, source_entry_id, tbl_custom_form_field_id,
			base_amount, gst_amount, total_amount, created_at
		FROM tbl_entry_gross_reimbursement
		WHERE gross_details_id = $1
		ORDER BY created_at
	`
	var reimbursements []domain.EntryGrossReimbursement
	err := r.db.SelectContext(ctx, &reimbursements, q, grossDetailsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reimbursements: %w", err)
	}
	return reimbursements, nil
}

func (r *entryGrossReimbursementRepo) DeleteByGrossDetailsID(ctx context.Context, grossDetailsID uuid.UUID) error {
	q := `DELETE FROM tbl_entry_gross_reimbursement WHERE gross_details_id = $1`
	_, err := r.db.ExecContext(ctx, q, grossDetailsID)
	return err
}
