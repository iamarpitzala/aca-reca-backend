package postgres

import (
	"context"
	"fmt"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/jmoiron/sqlx"
)

type customFormFieldRepo struct {
	db *sqlx.DB
}

func NewCustomFormFieldRepository(db *sqlx.DB) port.CustomFormFieldRepository {
	return &customFormFieldRepo{db: db}
}

func (r *customFormFieldRepo) GetByID(ctx context.Context, id string) (*form.Field, error) {
	q := `SELECT id, section_id, label, payment_responsibility_id, created_at, updated_at, deleted_at
		FROM tbl_custom_form_field WHERE id = $1 AND deleted_at IS NULL`
	var f form.Field
	if err := r.db.GetContext(ctx, &f, q, id); err != nil {
		return nil, fmt.Errorf("failed to get custom form field: %w", err)
	}
	return &f, nil
}

func (r *customFormFieldRepo) Create(ctx context.Context, field *form.Field) error {
	q := `INSERT INTO tbl_custom_form_field (id, section_id, label, payment_responsibility_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, q,
		field.ID, field.SectionID, field.Label, field.PaymentResponsibilityID, field.CreatedAt, field.UpdatedAt,
	)
	return err
}

func (r *customFormFieldRepo) CreateBatch(ctx context.Context, fields []*form.Field) error {
	for _, f := range fields {
		if err := r.Create(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

func (r *customFormFieldRepo) Update(ctx context.Context, field *form.Field) error {
	q := `UPDATE tbl_custom_form_field SET section_id = $2, label = $3, payment_responsibility_id = $4, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q, field.ID, field.SectionID, field.Label, field.PaymentResponsibilityID)
	return err
}

func (r *customFormFieldRepo) GetByFormVersionID(ctx context.Context, formVersionID int) ([]form.Field, error) {
	q := `SELECT f.id, f.section_id, f.label, f.payment_responsibility_id, f.created_at, f.updated_at, f.deleted_at
		FROM tbl_custom_form_field f
		INNER JOIN tbl_custom_form_section s ON s.id = f.section_id
		WHERE s.form_version_id = $1 AND f.deleted_at IS NULL
		ORDER BY s.section_order, f.id`
	var list []form.Field
	if err := r.db.SelectContext(ctx, &list, q, formVersionID); err != nil {
		return nil, fmt.Errorf("failed to get custom form fields by version: %w", err)
	}
	return list, nil
}

func (r *customFormFieldRepo) GetBySectionID(ctx context.Context, sectionID int) ([]form.Field, error) {
	q := `SELECT id, section_id, label, payment_responsibility_id, created_at, updated_at, deleted_at
		FROM tbl_custom_form_field WHERE section_id = $1 AND deleted_at IS NULL ORDER BY id`
	var list []form.Field
	if err := r.db.SelectContext(ctx, &list, q, sectionID); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *customFormFieldRepo) GetByFormID(ctx context.Context, formID string) ([]form.Field, error) {
	q := `SELECT f.id, f.section_id, f.label, f.payment_responsibility_id, f.created_at, f.updated_at, f.deleted_at
		FROM tbl_custom_form_field f
		INNER JOIN tbl_custom_form_section s ON s.id = f.section_id
		INNER JOIN tbl_custom_form_version v ON v.id = s.form_version_id
		WHERE v.form_id = $1 AND f.deleted_at IS NULL ORDER BY s.section_order, f.id`
	var list []form.Field
	if err := r.db.SelectContext(ctx, &list, q, formID); err != nil {
		return nil, fmt.Errorf("failed to get custom form fields: %w", err)
	}
	return list, nil
}

func (r *customFormFieldRepo) DeleteByFormID(ctx context.Context, formID string) error {
	q := `UPDATE tbl_custom_form_field SET deleted_at = now()
		WHERE section_id IN (SELECT id FROM tbl_custom_form_section WHERE form_version_id IN (SELECT id FROM tbl_custom_form_version WHERE form_id = $1))`
	_, err := r.db.ExecContext(ctx, q, formID)
	return err
}

func (r *customFormFieldRepo) DeleteByFormVersionID(ctx context.Context, formVersionID int) error {
	q := `UPDATE tbl_custom_form_field SET deleted_at = now()
		WHERE section_id IN (SELECT id FROM tbl_custom_form_section WHERE form_version_id = $1)`
	_, err := r.db.ExecContext(ctx, q, formVersionID)
	return err
}

func (r *customFormFieldRepo) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_field SET deleted_at = now() WHERE id = $1`, id)
	return err
}
