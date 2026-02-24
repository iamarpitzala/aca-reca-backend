package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/jmoiron/sqlx"
)

type customFormVersionRepo struct {
	db *sqlx.DB
}

func NewCustomFormVersionRepository(db *sqlx.DB) port.CustomFormVersionRepository {
	return &customFormVersionRepo{db: db}
}

func (r *customFormVersionRepo) Create(ctx context.Context, version *form.FormVersion) error {
	q := `INSERT INTO tbl_custom_form_version (
		form_id, version, is_active, created_by, created_at
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING id`
	err := r.db.QueryRowContext(ctx, q,
		version.FormID, version.Version, version.IsActive, version.CreatedBy, version.CreatedAt,
	).Scan(&version.ID)
	return err
}

func (r *customFormVersionRepo) GetLatestByFormID(ctx context.Context, formID string) (*form.FormVersion, error) {
	q := `
		SELECT id, form_id, version, is_active, created_by, created_at
		FROM tbl_custom_form_version
		WHERE form_id = $1
		ORDER BY version DESC
		LIMIT 1
	`
	var version form.FormVersion
	err := r.db.GetContext(ctx, &version, q, formID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no version found for form: %w", err)
		}
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}
	return &version, nil
}

func (r *customFormVersionRepo) GetByFormID(ctx context.Context, formID string) ([]form.FormVersion, error) {
	q := `
		SELECT id, form_id, version, is_active, created_by, created_at
		FROM tbl_custom_form_version
		WHERE form_id = $1
		ORDER BY version DESC
	`
	var versions []form.FormVersion
	if err := r.db.SelectContext(ctx, &versions, q, formID); err != nil {
		return nil, fmt.Errorf("failed to get form versions: %w", err)
	}
	return versions, nil
}

func (r *customFormVersionRepo) SetActive(ctx context.Context, formID string, versionID int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Set all versions for this form to inactive
	_, err = tx.ExecContext(ctx, `UPDATE tbl_custom_form_version SET is_active = false WHERE form_id = $1`, formID)
	if err != nil {
		return err
	}

	// Set the specified version to active
	_, err = tx.ExecContext(ctx, `UPDATE tbl_custom_form_version SET is_active = true WHERE id = $1`, versionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
