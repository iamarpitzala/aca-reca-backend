package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
	"github.com/jmoiron/sqlx"
)

// customFormRepo implements port.CustomFormRepository with queries based on the custom form tables.
type customFormRepo struct {
	db *sqlx.DB
}

func NewCustomFormRepository(db *sqlx.DB) port.CustomFormRepository {
	return &customFormRepo{db: db}
}

// Create a new custom form in tbl_custom_form.
func (r *customFormRepo) Create(ctx context.Context, form *domain.CustomForm) error {
	q := `INSERT INTO tbl_custom_form (
		id, clinic_id, name, description, form_type, status,
		calculation_method, default_payment_responsibility, created_by, created_at, updated_at, deleted_at
	) VALUES (
		:id, :clinic_id, :name, :description, :form_type, :status,
		:calculation_method, :default_payment_responsibility, :created_by, :created_at, :updated_at, :deleted_at
	)`
	_, err := r.db.NamedExecContext(ctx, q, form)
	return err
}

// Fetch a custom form by id from tbl_custom_form.
func (r *customFormRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomForm, error) {
	q := `
		SELECT
			id,
			clinic_id,
			name,
			description,
			form_type,
			status,
			calculation_method,
			default_payment_responsibility,
			created_by,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_custom_form
		WHERE id = $1 AND deleted_at IS NULL
	`
	var form domain.CustomForm
	err := r.db.GetContext(ctx, &form, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("custom form not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get custom form by id: %w", err)
	}
	return &form, nil
}

// Fetch all non-deleted custom forms for a clinic.
func (r *customFormRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomForm, error) {
	q := `
		SELECT
			id, clinic_id, name, description, form_type, status,
			calculation_method, default_payment_responsibility,
			created_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form
		WHERE clinic_id = $1 AND deleted_at IS NULL
		ORDER BY updated_at DESC
	`
	var forms []domain.CustomForm
	if err := r.db.SelectContext(ctx, &forms, q, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get custom forms: %w", err)
	}
	return forms, nil
}

// Fetch all published custom forms for a clinic.
func (r *customFormRepo) GetPublishedByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomForm, error) {
	q := `
		SELECT
			id, clinic_id, name, description, form_type, status,
			calculation_method, default_payment_responsibility,
			created_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form
		WHERE clinic_id = $1 AND status = $2 AND deleted_at IS NULL
		ORDER BY name
	`
	var forms []domain.CustomForm
	if err := r.db.SelectContext(ctx, &forms, q, clinicID, util.FormStatusPublished); err != nil {
		return nil, fmt.Errorf("failed to get published custom forms: %w", err)
	}
	return forms, nil
}

// Update allowed fields of a custom form.
func (r *customFormRepo) Update(ctx context.Context, form *domain.CustomForm) error {
	q := `
		UPDATE tbl_custom_form
		SET
			name = :name,
			description = :description,
			default_payment_responsibility = :default_payment_responsibility,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, q, form)
	return err
}

// Set custom form's status to 'PUBLISHED'
func (r *customFormRepo) Publish(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx,
		`UPDATE tbl_custom_form SET status = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`,
		util.FormStatusPublished, now, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

// Set custom form's status to 'DRAFT'
func (r *customFormRepo) Unpublish(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx,
		`UPDATE tbl_custom_form SET status = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`,
		util.FormStatusDraft, now, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

// Archive a form by status
func (r *customFormRepo) Archive(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE tbl_custom_form SET status = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`,
		util.FormStatusArchived, time.Now(), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

// Soft-delete a form by setting deleted_at.
func (r *customFormRepo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}
