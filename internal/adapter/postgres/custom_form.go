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
	"github.com/jmoiron/sqlx"
)

type customFormRepo struct {
	db *sqlx.DB
}

func NewCustomFormRepository(db *sqlx.DB) port.CustomFormRepository {
	return &customFormRepo{db: db}
}

func (r *customFormRepo) Create(ctx context.Context, form *domain.CustomForm) error {
	q := `INSERT INTO tbl_custom_form (id, clinic_id, name, description, calculation_method, form_type, status, fields, default_payment_responsibility, service_facility_fee_percent, outwork_enabled, outwork_rate_percent, version, created_by, created_at, updated_at)
		VALUES (:id, :clinic_id, :name, :description, :calculation_method, :form_type, :status, :fields, :default_payment_responsibility, :service_facility_fee_percent, :outwork_enabled, :outwork_rate_percent, :version, :created_by, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, q, form)
	return err
}

func (r *customFormRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomForm, error) {
	q := `SELECT id, clinic_id, name, description, calculation_method, form_type, status, fields, default_payment_responsibility, service_facility_fee_percent, outwork_enabled, outwork_rate_percent, version, created_by, created_at, updated_at, published_at, deleted_at FROM tbl_custom_form WHERE id = $1 AND deleted_at IS NULL`
	var form domain.CustomForm
	if err := r.db.GetContext(ctx, &form, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("custom form not found")
		}
		return nil, err
	}
	return &form, nil
}

func (r *customFormRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomForm, error) {
	q := `SELECT id, clinic_id, name, description, calculation_method, form_type, status, fields, default_payment_responsibility, service_facility_fee_percent, outwork_enabled, outwork_rate_percent, version, created_by, created_at, updated_at, published_at, deleted_at FROM tbl_custom_form WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY updated_at DESC`
	var rows []domain.CustomForm
	if err := r.db.SelectContext(ctx, &rows, q, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get custom forms: %w", err)
	}
	return rows, nil
}

func (r *customFormRepo) GetPublishedByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomForm, error) {
	q := `SELECT id, clinic_id, name, description, calculation_method, form_type, status, fields, default_payment_responsibility, service_facility_fee_percent, outwork_enabled, outwork_rate_percent, version, created_by, created_at, updated_at, published_at, deleted_at FROM tbl_custom_form WHERE clinic_id = $1 AND status = 'published' AND deleted_at IS NULL ORDER BY name`
	var rows []domain.CustomForm
	if err := r.db.SelectContext(ctx, &rows, q, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get published custom forms: %w", err)
	}
	return rows, nil
}

func (r *customFormRepo) Update(ctx context.Context, form *domain.CustomForm) error {
	q := `UPDATE tbl_custom_form SET name = :name, description = :description, fields = :fields, default_payment_responsibility = :default_payment_responsibility, service_facility_fee_percent = :service_facility_fee_percent, outwork_enabled = :outwork_enabled, outwork_rate_percent = :outwork_rate_percent, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, q, form)
	return err
}

func (r *customFormRepo) Publish(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form SET status = 'published', published_at = $1, updated_at = $1 WHERE id = $2 AND deleted_at IS NULL`, now, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

func (r *customFormRepo) Unpublish(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form SET status = 'draft', published_at = NULL, updated_at = $1 WHERE id = $2 AND deleted_at IS NULL`, now, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

func (r *customFormRepo) Archive(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form SET status = 'archived', updated_at = $1 WHERE id = $2 AND deleted_at IS NULL`, time.Now(), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

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

func (r *customFormRepo) CreateEntry(ctx context.Context, entry *domain.CustomFormEntry) error {
	// Get form to get calculation method
	form, err := r.GetByID(ctx, entry.FormID)
	if err != nil {
		return fmt.Errorf("failed to get form: %w", err)
	}

	// Use transaction to ensure data consistency
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := r.saveNormalizedEntry(ctx, tx, entry, form); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *customFormRepo) GetEntryByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormEntry, error) {
	// Use transaction for read operations to ensure consistency
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	normalized, err := r.loadNormalizedEntry(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	// Convert normalized to JSONB (doesn't need form, uses header's calculation method)
	entry, err := domain.ConvertNormalizedToJSONB(normalized)
	if err != nil {
		return nil, fmt.Errorf("failed to convert normalized entry: %w", err)
	}

	// For read operations, we don't need to commit, but we should explicitly handle the transaction
	// The defer rollback is fine for read-only transactions
	return entry, nil
}

func (r *customFormRepo) GetEntriesByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormEntry, error) {
	// Get form
	_, err := r.GetByID(ctx, formID)
	if err != nil {
		return nil, fmt.Errorf("failed to get form: %w", err)
	}

	// Get entry headers
	q := `SELECT id FROM tbl_entry_header WHERE form_id = $1 AND deleted_at IS NULL ORDER BY entry_date DESC, created_at DESC`
	var entryIDs []uuid.UUID
	if err := r.db.SelectContext(ctx, &entryIDs, q, formID); err != nil {
		return nil, fmt.Errorf("failed to get entry IDs: %w", err)
	}

	// Load each entry
	entries := make([]domain.CustomFormEntry, 0, len(entryIDs))
	for _, entryID := range entryIDs {
		entry, err := r.GetEntryByID(ctx, entryID)
		if err != nil {
			continue // Skip entries that fail to load
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func (r *customFormRepo) GetEntriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormEntry, error) {
	// Get entry headers
	q := `SELECT id FROM tbl_entry_header WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY entry_date DESC, created_at DESC`
	var entryIDs []uuid.UUID
	if err := r.db.SelectContext(ctx, &entryIDs, q, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get entry IDs: %w", err)
	}

	// Load each entry
	entries := make([]domain.CustomFormEntry, 0, len(entryIDs))
	for _, entryID := range entryIDs {
		entry, err := r.GetEntryByID(ctx, entryID)
		if err != nil {
			continue // Skip entries that fail to load
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func (r *customFormRepo) GetEntriesByQuarter(ctx context.Context, clinicID, quarterID uuid.UUID) ([]domain.CustomFormEntry, error) {
	// Get entry headers
	q := `SELECT id FROM tbl_entry_header WHERE clinic_id = $1 AND quarter_id = $2 AND deleted_at IS NULL ORDER BY entry_date DESC`
	var entryIDs []uuid.UUID
	if err := r.db.SelectContext(ctx, &entryIDs, q, clinicID, quarterID); err != nil {
		return nil, fmt.Errorf("failed to get entry IDs: %w", err)
	}

	// Load each entry
	entries := make([]domain.CustomFormEntry, 0, len(entryIDs))
	for _, entryID := range entryIDs {
		entry, err := r.GetEntryByID(ctx, entryID)
		if err != nil {
			continue // Skip entries that fail to load
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func (r *customFormRepo) UpdateEntry(ctx context.Context, entry *domain.CustomFormEntry) error {
	// Get form
	form, err := r.GetByID(ctx, entry.FormID)
	if err != nil {
		return fmt.Errorf("failed to get form: %w", err)
	}

	// Use transaction to ensure data consistency
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := r.updateNormalizedEntry(ctx, tx, entry, form); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *customFormRepo) DeleteEntry(ctx context.Context, id uuid.UUID) error {
	// Use transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if entry exists
	var exists bool
	if err := tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM tbl_entry_header WHERE id = $1 AND deleted_at IS NULL)`, id); err != nil {
		return err
	}
	if !exists {
		return errors.New("entry not found")
	}

	if err := r.deleteNormalizedEntry(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}
