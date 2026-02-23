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
	db                   *sqlx.DB
	financialYearRepo    port.FinancialYearRepository
	financialQuarterRepo port.FinancialQuarterRepository
}

func NewClinicFinancialSettingRepository(db *sqlx.DB, financialYearRepo port.FinancialYearRepository, financialQuarterRepo port.FinancialQuarterRepository) port.ClinicFinancialSettingRepository {
	return &clinicFinancialSettingRepo{db: db, financialYearRepo: financialYearRepo, financialQuarterRepo: financialQuarterRepo}
}

func (r *clinicFinancialSettingRepo) Create(ctx context.Context, settings *clinic.ClinicFinancialSetting) error {
	financialYear, err := r.financialYearRepo.GetByClinicID(ctx, settings.ClinicID)
	if err != nil {
		return err
	}

	financialQuarter, err := r.financialQuarterRepo.GetByFinancialYearID(ctx, financialYear.ID)
	if err != nil {
		return err
	}
	settings.FinancialQuarterID = financialQuarter.ID

	query := `INSERT INTO tbl_clinic_financial_settings 
		(id, clinic_id, financial_year_id, financial_quarter_id, accounting_method, gst_registered, gst_reporting_frequency, default_amount_mode, lock_date, created_at, updated_at)
		VALUES (:id, :clinic_id, :financial_year_id, :financial_quarter_id, :accounting_method, :gst_registered, :gst_reporting_frequency, :default_amount_mode, :lock_date, :created_at, :updated_at)`
	_, err = r.db.NamedExecContext(ctx, query, settings)
	return err
}

func (r *clinicFinancialSettingRepo) GetByClinicID(ctx context.Context, clinicID string) (*clinic.ClinicFinancialSetting, error) {
	query := `SELECT id, clinic_id, financial_year_id, financial_quarter_id, accounting_method, gst_registered, gst_reporting_frequency, default_amount_mode, lock_date, created_at, updated_at, deleted_at
		FROM tbl_clinic_financial_settings 
		WHERE clinic_id = $1 AND deleted_at IS NULL`
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

func (r *clinicFinancialSettingRepo) Update(ctx context.Context, settings *clinic.ClinicFinancialSetting) error {
	query := `UPDATE tbl_clinic_financial_settings SET
		financial_year_id = :financial_year_id,
		financial_quarter_id = :financial_quarter_id,
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
