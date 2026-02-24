package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
	"github.com/jmoiron/sqlx"
)

type clinicFinancialSettingRepo struct {
	db *sqlx.DB
}

func NewClinicFinancialSettingRepository(db *sqlx.DB) port.ClinicFinancialSettingRepository {
	return &clinicFinancialSettingRepo{db: db}
}

func (r *clinicFinancialSettingRepo) Create(ctx context.Context, settings *clinic.ClinicFinancialSetting) error {
	query := `INSERT INTO tbl_clinic_financial_settings
	(id, clinic_id, clinic_financial_year_id, accounting_method, gst_registered, gst_reporting_frequency, default_amount_mode, lock_date, created_at, updated_at, deleted_at)
	VALUES (:id, :clinic_id, :clinic_financial_year_id, :accounting_method, :gst_registered, :gst_reporting_frequency, :default_amount_mode, :lock_date, :created_at, :updated_at, :deleted_at)`
	_, err := r.db.NamedExecContext(ctx, query, settings)
	return err
}

func (r *clinicFinancialSettingRepo) GetByClinicID(ctx context.Context, clinicID string) (*clinic.ClinicFinancialSetting, error) {
	query := `SELECT s.id, s.clinic_id, s.clinic_financial_year_id, s.accounting_method, s.gst_registered, s.gst_reporting_frequency, s.default_amount_mode, s.lock_date, s.created_at, s.updated_at, s.deleted_at
FROM tbl_clinic_financial_settings s
JOIN tbl_clinic_financial_year cfy ON cfy.id = s.clinic_financial_year_id AND cfy.deleted_at IS NULL
WHERE cfy.clinic_id = $1 AND cfy.is_current = true AND s.deleted_at IS NULL
LIMIT 1`
	var settings clinic.ClinicFinancialSetting
	err := r.db.GetContext(ctx, &settings, query, clinicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial settings not found")
		}
		return nil, errors.New("failed to get financial settings")
	}
	return &settings, nil
}

func (r *clinicFinancialSettingRepo) GetByClinicFinancialYearID(ctx context.Context, clinicFinancialYearID int) (*clinic.ClinicFinancialSetting, error) {
	query := `SELECT id, clinic_id, clinic_financial_year_id, accounting_method, gst_registered, gst_reporting_frequency, default_amount_mode, lock_date, created_at, updated_at, deleted_at
FROM tbl_clinic_financial_settings WHERE clinic_financial_year_id = $1 AND deleted_at IS NULL`
	var settings clinic.ClinicFinancialSetting
	err := r.db.GetContext(ctx, &settings, query, clinicFinancialYearID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *clinicFinancialSettingRepo) Update(ctx context.Context, settings *clinic.ClinicFinancialSetting) error {
	query := `UPDATE tbl_clinic_financial_settings SET
	clinic_financial_year_id = :clinic_financial_year_id,
	accounting_method = :accounting_method,
	gst_registered = :gst_registered,
	gst_reporting_frequency = :gst_reporting_frequency,
	default_amount_mode = :default_amount_mode,
	lock_date = :lock_date,
	updated_at = :updated_at
	WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, settings)
	return err
}

func (r *clinicFinancialSettingRepo) Delete(ctx context.Context, clinicID string) error {
	query := `UPDATE tbl_clinic_financial_settings SET deleted_at = $1 WHERE clinic_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), clinicID)
	return err
}
