package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type fieldEntryRepo struct {
	db *sqlx.DB
}

func NewFieldEntryRepository(db *sqlx.DB) port.FieldEntryRepository {
	return &fieldEntryRepo{db: db}
}

// Create implements [port.FieldEntryRepository].
func (f *fieldEntryRepo) Create(ctx context.Context, entry *domain.FieldEntry) error {
	q := `INSERT INTO tbl_custom_form_entry (
		id, form_id, form_version_id, tbl_custom_form_field_id, value, created_by, created_at, updated_at, deleted_at
	) VALUES (
		:id, :form_id, :form_version_id, :tbl_custom_form_field_id, :value, :created_by, :created_at, :updated_at, :deleted_at
	)`
	_, err := f.db.NamedExecContext(ctx, q, entry)
	return err
}

// Delete implements [port.FieldEntryRepository].
func (f *fieldEntryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	q := `UPDATE tbl_custom_form_entry SET deleted_at = $1 WHERE id = $2`
	_, err := f.db.ExecContext(ctx, q, time.Now(), id)
	return err
}

// GetByClinicID implements [port.FieldEntryRepository].
func (f *fieldEntryRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.FieldEntry, error) {
	q := `
		SELECT 
			cfe.id, 
			cfe.form_id,
			cfe.form_version_id,
			cfe.tbl_custom_form_field_id, 
			cfe.value, 
			cfe.created_by, 
			cfe.created_at, 
			cfe.updated_at, 
			cfe.deleted_at 
		FROM tbl_custom_form_entry cfe
		INNER JOIN tbl_custom_form cf ON cfe.form_id = cf.id
		WHERE cf.clinic_id = $1 AND cfe.deleted_at IS NULL
	`
	var entries []domain.FieldEntry
	err := f.db.SelectContext(ctx, &entries, q, clinicID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// GetByFormID implements [port.FieldEntryRepository].
func (f *fieldEntryRepo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]domain.FieldEntry, error) {
	q := `SELECT id, form_id, form_version_id, tbl_custom_form_field_id, value, created_by, created_at, updated_at, deleted_at FROM tbl_custom_form_entry WHERE form_id = $1 AND deleted_at IS NULL`
	var entries []domain.FieldEntry
	err := f.db.SelectContext(ctx, &entries, q, formID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// GetByID implements [port.FieldEntryRepository].
func (f *fieldEntryRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.FieldEntry, error) {
	q := `SELECT id, form_id, form_version_id, tbl_custom_form_field_id, value, created_by, created_at, updated_at, deleted_at FROM tbl_custom_form_entry WHERE id = $1 AND deleted_at IS NULL`
	var entry domain.FieldEntry
	err := f.db.GetContext(ctx, &entry, q, id)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// Update implements [port.FieldEntryRepository].
func (f *fieldEntryRepo) Update(ctx context.Context, entry *domain.FieldEntry) error {
	q := `UPDATE tbl_custom_form_entry SET value = $1, updated_at = $2 WHERE id = $3`
	_, err := f.db.ExecContext(ctx, q, entry.Value, time.Now(), entry.ID)
	return err
}
