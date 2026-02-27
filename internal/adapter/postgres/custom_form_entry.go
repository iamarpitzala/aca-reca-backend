package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/jmoiron/sqlx"
)

type customFormEntryRepo struct {
	db *sqlx.DB
}

func NewCustomFormEntryRepository(db *sqlx.DB) port.CustomFormEntryRepository {
	return &customFormEntryRepo{db: db}
}

func (r *customFormEntryRepo) Create(ctx context.Context, entry *form.Entry, values []form.EntryValue) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qEntry := `INSERT INTO tbl_custom_form_entry (id, form_version_id, submitted_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, qEntry,
		entry.ID, entry.FormVersionID, entry.SubmittedBy, entry.CreatedAt, entry.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert entry: %w", err)
	}

	for _, v := range values {
		qVal := `INSERT INTO tbl_custom_form_entry_value (entry_id, field_id, value, gst_amount, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.ExecContext(ctx, qVal,
			v.EntryID, v.FieldID, v.Value, v.GSTAmount, v.CreatedAt, v.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert entry value: %w", err)
		}
	}

	return tx.Commit()
}

func (r *customFormEntryRepo) GetByID(ctx context.Context, id string) (*form.Entry, []form.EntryValue, error) {
	qEntry := `SELECT id, form_version_id, submitted_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form_entry WHERE id = $1 AND deleted_at IS NULL`
	var entry form.Entry
	if err := r.db.GetContext(ctx, &entry, qEntry, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("entry not found: %w", err)
		}
		return nil, nil, err
	}

	qValues := `SELECT id, entry_id, field_id, value, gst_amount, created_at, updated_at, deleted_at
		FROM tbl_custom_form_entry_value WHERE entry_id = $1 AND deleted_at IS NULL`
	var values []form.EntryValue
	if err := r.db.SelectContext(ctx, &values, qValues, id); err != nil {
		return nil, nil, err
	}
	if values == nil {
		values = []form.EntryValue{}
	}
	return &entry, values, nil
}

func (r *customFormEntryRepo) getValuesByEntryIDs(ctx context.Context, entryIDs []string) (map[string][]form.EntryValue, error) {
	if len(entryIDs) == 0 {
		return map[string][]form.EntryValue{}, nil
	}
	q := `SELECT id, entry_id, field_id, value, gst_amount, created_at, updated_at, deleted_at
		FROM tbl_custom_form_entry_value WHERE entry_id = ANY($1) AND deleted_at IS NULL`
	var rows []struct {
		form.EntryValue
	}
	query, args, err := sqlx.In(q, entryIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	out := make(map[string][]form.EntryValue)
	for _, row := range rows {
		out[row.EntryID] = append(out[row.EntryID], row.EntryValue)
	}
	return out, nil
}

func (r *customFormEntryRepo) GetByFormVersionID(ctx context.Context, formVersionID int) ([]form.Entry, error) {
	q := `SELECT id, form_version_id, submitted_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form_entry WHERE form_version_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	var entries []form.Entry
	if err := r.db.SelectContext(ctx, &entries, q, formVersionID); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *customFormEntryRepo) GetByFormID(ctx context.Context, formID string) ([]form.Entry, error) {
	q := `SELECT e.id, e.form_version_id, e.submitted_by, e.created_at, e.updated_at, e.deleted_at
		FROM tbl_custom_form_entry e
		INNER JOIN tbl_custom_form_version v ON v.id = e.form_version_id
		WHERE v.form_id = $1 AND e.deleted_at IS NULL ORDER BY e.created_at DESC`
	var entries []form.Entry
	if err := r.db.SelectContext(ctx, &entries, q, formID); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *customFormEntryRepo) GetByClinicID(ctx context.Context, clinicID string) ([]form.Entry, error) {
	q := `SELECT e.id, e.form_version_id, e.submitted_by, e.created_at, e.updated_at, e.deleted_at
		FROM tbl_custom_form_entry e
		INNER JOIN tbl_custom_form_version v ON v.id = e.form_version_id
		INNER JOIN tbl_custom_form f ON f.id = v.form_id
		WHERE f.clinic_id = $1 AND e.deleted_at IS NULL AND f.deleted_at IS NULL ORDER BY e.created_at DESC`
	var entries []form.Entry
	if err := r.db.SelectContext(ctx, &entries, q, clinicID); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *customFormEntryRepo) Update(ctx context.Context, entryID string, values []form.EntryValue) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	// Soft-delete existing values then insert new set; or delete and re-insert. Schema has UNIQUE(entry_id, field_id).
	// Simpler: delete all values for entry and insert new (or upsert). For soft delete we'd set deleted_at.
	_, err = tx.ExecContext(ctx, `UPDATE tbl_custom_form_entry_value SET deleted_at = $1 WHERE entry_id = $2`, now, entryID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE tbl_custom_form_entry SET updated_at = $1 WHERE id = $2`, now, entryID)
	if err != nil {
		return err
	}

	for _, v := range values {
		v.EntryID = entryID
		v.UpdatedAt = now
		v.CreatedAt = now
		q := `INSERT INTO tbl_custom_form_entry_value (entry_id, field_id, value, gst_amount, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.ExecContext(ctx, q, v.EntryID, v.FieldID, v.Value, v.GSTAmount, v.CreatedAt, v.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *customFormEntryRepo) Delete(ctx context.Context, id string) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_entry SET deleted_at = $1 WHERE id = $2`, now, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("entry not found")
	}
	return nil
}
