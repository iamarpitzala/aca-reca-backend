package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type customFormFieldRepo struct {
	db *sqlx.DB
}

func NewCustomFormFieldRepository(db *sqlx.DB) port.CustomFormFieldRepository {
	return &customFormFieldRepo{db: db}
}

func (r *customFormFieldRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormField, error) {
	q := `
		SELECT
			id, form_version_id, form_id, field_key, label, field_type, section,
			is_required, coa_id, placeholder, min_value, max_value, field_order,
			gst_config, gst_rate, gst_type, metadata
		FROM tbl_custom_form_field
		WHERE id = $1
	`
	var field domain.CustomFormField
	if err := r.db.GetContext(ctx, &field, q, id); err != nil {
		return nil, fmt.Errorf("failed to get custom form field: %w", err)
	}
	return &field, nil
}

func (r *customFormFieldRepo) Create(ctx context.Context, field *domain.CustomFormField) error {
	q := `INSERT INTO tbl_custom_form_field (
		id, form_version_id, form_id, field_key, label, field_type, section,
		is_required, coa_id, placeholder, min_value, max_value, field_order,
		gst_config, gst_rate, gst_type, metadata
	) VALUES (
		:id, :form_version_id, :form_id, :field_key, :label, :field_type, :section,
		:is_required, :coa_id, :placeholder, :min_value, :max_value, :field_order,
		:gst_config, :gst_rate, :gst_type, :metadata
	)`
	_, err := r.db.NamedExecContext(ctx, q, field)
	return err
}

func (r *customFormFieldRepo) CreateBatch(ctx context.Context, fields []*domain.CustomFormField) error {
	if len(fields) == 0 {
		return nil
	}
	q := `INSERT INTO tbl_custom_form_field (
		id, form_version_id, form_id, field_key, label, field_type, section,
		is_required, coa_id, placeholder, min_value, max_value, field_order,
		gst_config, gst_rate, gst_type, metadata
	) VALUES (
		:id, :form_version_id, :form_id, :field_key, :label, :field_type, :section,
		:is_required, :coa_id, :placeholder, :min_value, :max_value, :field_order,
		:gst_config, :gst_rate, :gst_type, :metadata
	)`
	_, err := r.db.NamedExecContext(ctx, q, fields)
	return err
}

func (r *customFormFieldRepo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormField, error) {
	q := `
		SELECT
			id, form_version_id, form_id, field_key, label, field_type, section,
			is_required, coa_id, placeholder, min_value, max_value, field_order,
			gst_config, gst_rate, gst_type, metadata
		FROM tbl_custom_form_field
		WHERE form_id = $1
		ORDER BY field_order ASC
	`
	var fields []domain.CustomFormField
	if err := r.db.SelectContext(ctx, &fields, q, formID); err != nil {
		return nil, fmt.Errorf("failed to get custom form fields: %w", err)
	}
	return fields, nil
}

func (r *customFormFieldRepo) GetByFormVersionID(ctx context.Context, formVersionID uuid.UUID) ([]domain.CustomFormField, error) {
	q := `
		SELECT
			id, form_version_id, form_id, field_key, label, field_type, section,
			is_required, coa_id, placeholder, min_value, max_value, field_order,
			gst_config, gst_rate, gst_type, metadata
		FROM tbl_custom_form_field
		WHERE form_version_id = $1
		ORDER BY field_order ASC
	`
	var fields []domain.CustomFormField
	if err := r.db.SelectContext(ctx, &fields, q, formVersionID); err != nil {
		return nil, fmt.Errorf("failed to get custom form fields by version: %w", err)
	}
	return fields, nil
}

func (r *customFormFieldRepo) DeleteByFormID(ctx context.Context, formID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tbl_custom_form_field WHERE form_id = $1`, formID)
	return err
}

func (r *customFormFieldRepo) DeleteByFormVersionID(ctx context.Context, formVersionID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tbl_custom_form_field WHERE form_version_id = $1`, formVersionID)
	return err
}
