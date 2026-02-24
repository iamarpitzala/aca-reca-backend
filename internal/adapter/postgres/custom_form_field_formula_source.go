package postgres

import (
	"context"
	"fmt"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/jmoiron/sqlx"
)

type customFormFieldFormulaSourceRepo struct {
	db *sqlx.DB
}

func NewCustomFormFieldFormulaSourceRepository(db *sqlx.DB) port.CustomFormFieldFormulaSourceRepository {
	return &customFormFieldFormulaSourceRepo{db: db}
}

func (r *customFormFieldFormulaSourceRepo) Create(ctx context.Context, source *form.FormulaSource) error {
	q := `INSERT INTO tbl_custom_form_field_formula_source (
		id, field_config_id, source_field_id, source_role, source_order, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, now(), now()
	)`
	_, err := r.db.ExecContext(ctx, q,
		source.ID, source.FieldConfigID, source.SourceFieldID, source.SourceRole, source.SourceOrder,
	)
	return err
}

func (r *customFormFieldFormulaSourceRepo) CreateBatch(ctx context.Context, sources []*form.FormulaSource) error {
	for _, s := range sources {
		if err := r.Create(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (r *customFormFieldFormulaSourceRepo) GetByID(ctx context.Context, id string) (*form.FormulaSource, error) {
	q := `
		SELECT id, field_config_id, source_field_id, source_role, source_order,
			created_at, updated_at, deleted_at
		FROM tbl_custom_form_field_formula_source
		WHERE id = $1 AND deleted_at IS NULL
	`
	var s form.FormulaSource
	if err := r.db.GetContext(ctx, &s, q, id); err != nil {
		return nil, fmt.Errorf("failed to get formula source: %w", err)
	}
	return &s, nil
}

func (r *customFormFieldFormulaSourceRepo) GetByFieldConfigID(ctx context.Context, fieldConfigID string) ([]form.FormulaSource, error) {
	q := `
		SELECT id, field_config_id, source_field_id, source_role, source_order,
			created_at, updated_at, deleted_at
		FROM tbl_custom_form_field_formula_source
		WHERE field_config_id = $1 AND deleted_at IS NULL
		ORDER BY source_order ASC
	`
	var list []form.FormulaSource
	if err := r.db.SelectContext(ctx, &list, q, fieldConfigID); err != nil {
		return nil, fmt.Errorf("failed to list formula sources: %w", err)
	}
	return list, nil
}

func (r *customFormFieldFormulaSourceRepo) Update(ctx context.Context, source *form.FormulaSource) error {
	q := `UPDATE tbl_custom_form_field_formula_source SET
		source_field_id = $2, source_role = $3, source_order = $4, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, source.ID, source.SourceFieldID, source.SourceRole, source.SourceOrder)
	return err
}

func (r *customFormFieldFormulaSourceRepo) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_field_formula_source SET deleted_at = now() WHERE id = $1`, id)
	return err
}

func (r *customFormFieldFormulaSourceRepo) DeleteByFieldConfigID(ctx context.Context, fieldConfigID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_field_formula_source SET deleted_at = now() WHERE field_config_id = $1`, fieldConfigID)
	return err
}
