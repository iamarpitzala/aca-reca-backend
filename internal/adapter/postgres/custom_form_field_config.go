package postgres

import (
	"context"
	"fmt"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/jmoiron/sqlx"
)

type customFormFieldConfigRepo struct {
	db *sqlx.DB
}

func NewCustomFormFieldConfigRepository(db *sqlx.DB) port.CustomFormFieldConfigRepository {
	return &customFormFieldConfigRepo{db: db}
}

func (r *customFormFieldConfigRepo) Create(ctx context.Context, config *form.FieldConfig) error {
	q := `INSERT INTO tbl_custom_form_field_config (
		id, form_field_id, tax_type_id, arrangement_id, is_formula, operator, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, now(), now()
	)`
	_, err := r.db.ExecContext(ctx, q,
		config.ID, config.FormFieldID, config.TaxTypeID, config.ArrangementID,
		config.IsFormula, config.Operator,
	)
	return err
}

func (r *customFormFieldConfigRepo) GetByID(ctx context.Context, id string) (*form.FieldConfig, error) {
	q := `
		SELECT id, form_field_id, tax_type_id, arrangement_id, is_formula, operator,
			created_at, updated_at, deleted_at
		FROM tbl_custom_form_field_config
		WHERE id = $1 AND deleted_at IS NULL
	`
	var c form.FieldConfig
	if err := r.db.GetContext(ctx, &c, q, id); err != nil {
		return nil, fmt.Errorf("failed to get field config: %w", err)
	}
	return &c, nil
}

func (r *customFormFieldConfigRepo) GetByFormFieldID(ctx context.Context, formFieldID string) ([]form.FieldConfig, error) {
	q := `
		SELECT id, form_field_id, tax_type_id, arrangement_id, is_formula, operator,
			created_at, updated_at, deleted_at
		FROM tbl_custom_form_field_config
		WHERE form_field_id = $1 AND deleted_at IS NULL
	`
	var list []form.FieldConfig
	if err := r.db.SelectContext(ctx, &list, q, formFieldID); err != nil {
		return nil, fmt.Errorf("failed to list field configs: %w", err)
	}
	return list, nil
}

func (r *customFormFieldConfigRepo) Update(ctx context.Context, config *form.FieldConfig) error {
	q := `UPDATE tbl_custom_form_field_config SET
		tax_type_id = $2, arrangement_id = $3, is_formula = $4, operator = $5, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q,
		config.ID, config.TaxTypeID, config.ArrangementID, config.IsFormula, config.Operator,
	)
	return err
}

func (r *customFormFieldConfigRepo) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_field_config SET deleted_at = now() WHERE id = $1`, id)
	return err
}

func (r *customFormFieldConfigRepo) DeleteByFormFieldID(ctx context.Context, formFieldID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_field_config SET deleted_at = now() WHERE form_field_id = $1`, formFieldID)
	return err
}
