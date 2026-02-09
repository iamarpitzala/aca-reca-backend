package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateCustomForm(ctx context.Context, db *sqlx.DB, form *domain.CustomForm) error {
	query := `INSERT INTO tbl_custom_form (
		id, clinic_id, name, description, calculation_method, form_type, status, fields,
		default_payment_responsibility, service_facility_fee_percent, version, created_by, created_at, updated_at
	) VALUES (
		:id, :clinic_id, :name, :description, :calculation_method, :form_type, :status, :fields,
		:default_payment_responsibility, :service_facility_fee_percent, :version, :created_by, :created_at, :updated_at
	)`
	_, err := db.NamedExecContext(ctx, query, form)
	return err
}

func GetCustomFormByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.CustomForm, error) {
	query := `SELECT id, clinic_id, name, description, calculation_method, form_type, status, fields,
		default_payment_responsibility, service_facility_fee_percent, version, created_by, created_at, updated_at, published_at, deleted_at
		FROM tbl_custom_form WHERE id = $1 AND deleted_at IS NULL`
	var form domain.CustomForm
	err := db.GetContext(ctx, &form, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("custom form not found")
		}
		return nil, err
	}
	return &form, nil
}

func GetCustomFormsByClinicID(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.CustomForm, error) {
	query := `SELECT id, clinic_id, name, description, calculation_method, form_type, status, fields,
		default_payment_responsibility, service_facility_fee_percent, version, created_by, created_at, updated_at, published_at, deleted_at
		FROM tbl_custom_form WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY updated_at DESC`
	var rows []domain.CustomForm
	if err := db.SelectContext(ctx, &rows, query, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get custom forms: %w", err)
	}
	return rows, nil
}

func GetPublishedCustomFormsByClinicID(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.CustomForm, error) {
	query := `SELECT id, clinic_id, name, description, calculation_method, form_type, status, fields,
		default_payment_responsibility, service_facility_fee_percent, version, created_by, created_at, updated_at, published_at, deleted_at
		FROM tbl_custom_form WHERE clinic_id = $1 AND status = 'published' AND deleted_at IS NULL ORDER BY name`
	var rows []domain.CustomForm
	if err := db.SelectContext(ctx, &rows, query, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get published custom forms: %w", err)
	}
	return rows, nil
}

func UpdateCustomForm(ctx context.Context, db *sqlx.DB, form *domain.CustomForm) error {
	query := `UPDATE tbl_custom_form SET name = :name, description = :description, fields = :fields,
		default_payment_responsibility = :default_payment_responsibility, service_facility_fee_percent = :service_facility_fee_percent,
		updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL`
	_, err := db.NamedExecContext(ctx, query, form)
	return err
}

func PublishCustomForm(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE tbl_custom_form SET status = 'published', published_at = $1, updated_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	res, err := db.ExecContext(ctx, query, now, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

func ArchiveCustomForm(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `UPDATE tbl_custom_form SET status = 'archived', updated_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	res, err := db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

func DeleteCustomForm(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `UPDATE tbl_custom_form SET deleted_at = $1 WHERE id = $2`
	res, err := db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("custom form not found")
	}
	return nil
}

// CreateCustomFormEntry
func CreateCustomFormEntry(ctx context.Context, db *sqlx.DB, entry *domain.CustomFormEntry) error {
	query := `INSERT INTO tbl_custom_form_entry (
		id, form_id, form_name, form_type, clinic_id, quarter_id, values, calculations, entry_date,
		description, remarks, payment_responsibility, deductions, created_by, created_at, updated_at
	) VALUES (
		:id, :form_id, :form_name, :form_type, :clinic_id, :quarter_id, :values, :calculations, :entry_date,
		:description, :remarks, :payment_responsibility, :deductions, :created_by, :created_at, :updated_at
	)`
	_, err := db.NamedExecContext(ctx, query, entry)
	return err
}

func GetCustomFormEntryByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.CustomFormEntry, error) {
	query := `SELECT id, form_id, form_name, form_type, clinic_id, quarter_id, values, calculations, entry_date,
		description, remarks, payment_responsibility, deductions, created_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form_entry WHERE id = $1 AND deleted_at IS NULL`
	var entry domain.CustomFormEntry
	err := db.GetContext(ctx, &entry, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("entry not found")
		}
		return nil, err
	}
	return &entry, nil
}

func GetCustomFormEntriesByFormID(ctx context.Context, db *sqlx.DB, formID uuid.UUID) ([]domain.CustomFormEntry, error) {
	query := `SELECT id, form_id, form_name, form_type, clinic_id, quarter_id, values, calculations, entry_date,
		description, remarks, payment_responsibility, deductions, created_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form_entry WHERE form_id = $1 AND deleted_at IS NULL ORDER BY entry_date DESC, created_at DESC`
	var rows []domain.CustomFormEntry
	if err := db.SelectContext(ctx, &rows, query, formID); err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}
	return rows, nil
}

func GetCustomFormEntriesByClinicID(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.CustomFormEntry, error) {
	query := `SELECT id, form_id, form_name, form_type, clinic_id, quarter_id, values, calculations, entry_date,
		description, remarks, payment_responsibility, deductions, created_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form_entry WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY entry_date DESC, created_at DESC`
	var rows []domain.CustomFormEntry
	if err := db.SelectContext(ctx, &rows, query, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}
	return rows, nil
}

func GetCustomFormEntriesByQuarter(ctx context.Context, db *sqlx.DB, clinicID, quarterID uuid.UUID) ([]domain.CustomFormEntry, error) {
	query := `SELECT id, form_id, form_name, form_type, clinic_id, quarter_id, values, calculations, entry_date,
		description, remarks, payment_responsibility, deductions, created_by, created_at, updated_at, deleted_at
		FROM tbl_custom_form_entry WHERE clinic_id = $1 AND quarter_id = $2 AND deleted_at IS NULL ORDER BY entry_date DESC`
	var rows []domain.CustomFormEntry
	if err := db.SelectContext(ctx, &rows, query, clinicID, quarterID); err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}
	return rows, nil
}

func UpdateCustomFormEntry(ctx context.Context, db *sqlx.DB, entry *domain.CustomFormEntry) error {
	query := `UPDATE tbl_custom_form_entry SET values = :values, calculations = :calculations, updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`
	_, err := db.NamedExecContext(ctx, query, entry)
	return err
}

func DeleteCustomFormEntry(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `UPDATE tbl_custom_form_entry SET deleted_at = $1 WHERE id = $2`
	res, err := db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("entry not found")
	}
	return nil
}
