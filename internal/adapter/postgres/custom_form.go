package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/iamarpitzala/aca-reca-backend/util"
	"github.com/jmoiron/sqlx"
)

type customFormRepo struct {
	db *sqlx.DB
}

func NewCustomFormRepository(db *sqlx.DB) port.CustomFormRepository {
	return &customFormRepo{db: db}
}

func (r *customFormRepo) Create(ctx context.Context, form *form.Form) error {
	q := `INSERT INTO tbl_custom_form (
		id, clinic_id, name, description, status,
		calculation_method, created_by, created_at, updated_at, deleted_at
	) VALUES (
		:id, :clinic_id, :name, :description, :status,
		:calculation_method, :created_by, :created_at, :updated_at, :deleted_at
	)`
	_, err := r.db.NamedExecContext(ctx, q, form)
	return err
}

func (r *customFormRepo) GetByID(ctx context.Context, id string) (*form.Form, error) {
	q := `
		SELECT
			id,
			clinic_id,
			name,
			description,
			status,
			calculation_method,
			created_by,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_custom_form
		WHERE id = $1 AND deleted_at IS NULL
	`
	var form form.Form
	err := r.db.GetContext(ctx, &form, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("custom form not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get custom form by id: %w", err)
	}
	return &form, nil
}

func (r *customFormRepo) GetByClinicID(ctx context.Context, clinicID string) ([]form.Form, error) {
	q := `
		SELECT
			id, clinic_id, name, description, status,
			calculation_method,
			created_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form
		WHERE clinic_id = $1 AND deleted_at IS NULL
		ORDER BY updated_at DESC
	`
	var forms []form.Form
	if err := r.db.SelectContext(ctx, &forms, q, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get custom forms: %w", err)
	}
	return forms, nil
}

func (r *customFormRepo) GetPublishedByClinicID(ctx context.Context, clinicID string) ([]form.Form, error) {
	q := `
		SELECT
			id, clinic_id, name, description, status,
			calculation_method,
			created_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form
		WHERE clinic_id = $1 AND status = $2 AND deleted_at IS NULL
		ORDER BY name
	`
	var forms []form.Form
	if err := r.db.SelectContext(ctx, &forms, q, clinicID, util.FormStatusPublished); err != nil {
		return nil, fmt.Errorf("failed to get published custom forms: %w", err)
	}
	return forms, nil
}

func (r *customFormRepo) Update(ctx context.Context, form *form.Form) error {
	q := `
		UPDATE tbl_custom_form
		SET
			name = :name,
			description = :description,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, q, form)
	return err
}

func (r *customFormRepo) Publish(ctx context.Context, id string) error {
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

func (r *customFormRepo) Unpublish(ctx context.Context, id string) error {
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

func (r *customFormRepo) Archive(ctx context.Context, id string) error {
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

func (r *customFormRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}
