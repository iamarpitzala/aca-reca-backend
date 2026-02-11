package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type clinicFinancialSettingsRepo struct {
	db *sqlx.DB
}

func NewClinicFinancialSettingsRepository(db *sqlx.DB) port.ClinicFinancialSettingsRepository {
	return &clinicFinancialSettingsRepo{db: db}
}

func (r *clinicFinancialSettingsRepo) Create(ctx context.Context, settings *domain.ClinicFinancialSettings) error {
	settings.ID = uuid.New()
	settings.CreatedAt = time.Now()
	settings.UpdatedAt = time.Now()
	
	query := `INSERT INTO tbl_clinic_financial_settings 
		(id, clinic_id, financial_year_start, accounting_method, gst_registered, gst_reporting_frequency, default_amount_mode, lock_date, gst_defaults, created_at, updated_at)
		VALUES (:id, :clinic_id, :financial_year_start, :accounting_method, :gst_registered, :gst_reporting_frequency, :default_amount_mode, :lock_date, :gst_defaults, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, settings)
	return err
}

func (r *clinicFinancialSettingsRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.ClinicFinancialSettings, error) {
	query := `SELECT id, clinic_id, financial_year_start, accounting_method, gst_registered, gst_reporting_frequency, default_amount_mode, lock_date, gst_defaults, created_at, updated_at, deleted_at
		FROM tbl_clinic_financial_settings 
		WHERE clinic_id = $1 AND deleted_at IS NULL`
	var settings domain.ClinicFinancialSettings
	err := r.db.GetContext(ctx, &settings, query, clinicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial settings not found")
		}
		return nil, errors.New("failed to get financial settings")
	}
	return &settings, nil
}

func (r *clinicFinancialSettingsRepo) Update(ctx context.Context, settings *domain.ClinicFinancialSettings) error {
	settings.UpdatedAt = time.Now()
	query := `UPDATE tbl_clinic_financial_settings SET
		financial_year_start = :financial_year_start,
		accounting_method = :accounting_method,
		gst_registered = :gst_registered,
		gst_reporting_frequency = :gst_reporting_frequency,
		default_amount_mode = :default_amount_mode,
		lock_date = :lock_date,
		gst_defaults = :gst_defaults,
		updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.NamedExecContext(ctx, query, settings)
	return err
}

func (r *clinicFinancialSettingsRepo) Delete(ctx context.Context, clinicID uuid.UUID) error {
	query := `UPDATE tbl_clinic_financial_settings SET deleted_at = $1 WHERE clinic_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), clinicID)
	return err
}
