package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type entryGrossReductionRepo struct {
	db *sqlx.DB
}

func NewEntryGrossReductionRepository(db *sqlx.DB) port.EntryGrossReductionRepository {
	return &entryGrossReductionRepo{db: db}
}

func (r *entryGrossReductionRepo) Create(ctx context.Context, reduction *domain.EntryGrossReduction) error {
	q := `INSERT INTO tbl_entry_gross_reduction (
		id, gross_details_id, source_entry_id, tbl_custom_form_field_id,
		base_amount, gst_amount, total_amount, created_at
	) VALUES (
		:id, :gross_details_id, :source_entry_id, :tbl_custom_form_field_id,
		:base_amount, :gst_amount, :total_amount, :created_at
	)`
	_, err := r.db.NamedExecContext(ctx, q, reduction)
	return err
}

func (r *entryGrossReductionRepo) CreateBatch(ctx context.Context, reductions []*domain.EntryGrossReduction) error {
	if len(reductions) == 0 {
		return nil
	}

	q := `INSERT INTO tbl_entry_gross_reduction (
		id, gross_details_id, source_entry_id, tbl_custom_form_field_id,
		base_amount, gst_amount, total_amount, created_at
	) VALUES (
		:id, :gross_details_id, :source_entry_id, :tbl_custom_form_field_id,
		:base_amount, :gst_amount, :total_amount, :created_at
	)`

	_, err := r.db.NamedExecContext(ctx, q, reductions)
	if err != nil {
		return fmt.Errorf("failed to create batch reductions: %w", err)
	}
	return nil
}

func (r *entryGrossReductionRepo) GetByGrossDetailsID(ctx context.Context, grossDetailsID uuid.UUID) ([]domain.EntryGrossReduction, error) {
	q := `
		SELECT
			id, gross_details_id, source_entry_id, tbl_custom_form_field_id,
			base_amount, gst_amount, total_amount, created_at
		FROM tbl_entry_gross_reduction
		WHERE gross_details_id = $1
		ORDER BY created_at
	`
	var reductions []domain.EntryGrossReduction
	err := r.db.SelectContext(ctx, &reductions, q, grossDetailsID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reductions: %w", err)
	}
	return reductions, nil
}

func (r *entryGrossReductionRepo) DeleteByGrossDetailsID(ctx context.Context, grossDetailsID uuid.UUID) error {
	q := `DELETE FROM tbl_entry_gross_reduction WHERE gross_details_id = $1`
	_, err := r.db.ExecContext(ctx, q, grossDetailsID)
	return err
}
