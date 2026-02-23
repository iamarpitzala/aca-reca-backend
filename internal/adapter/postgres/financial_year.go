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

type financialYearRepo struct {
	db *sqlx.DB
}

func NewFinancialYearRepository(db *sqlx.DB) port.FinancialYearRepository {
	return &financialYearRepo{db: db}
}

// Create implements [port.FinancialYearRepository].
func (f *financialYearRepo) Create(ctx context.Context, year *clinic.FinancialYear) error {
	query := `INSERT INTO tbl_financial_year (id, clinic_id, fy_label, start_date, end_date, is_current, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at) VALUES (:id, :clinic_id, :fy_label, :start_date, :end_date, :is_current, :is_closed, :closed_at, :closed_by, :created_at, :updated_at, :deleted_at)`
	_, err := f.db.NamedExecContext(ctx, query, year)
	return err
}

// Delete implements [port.FinancialYearRepository].
func (f *financialYearRepo) Delete(ctx context.Context, id int) error {
	query := `UPDATE tbl_financial_year SET deleted_at = $1 WHERE id = $2`
	_, err := f.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// GetByClinicID implements [port.FinancialYearRepository].
func (f *financialYearRepo) GetByClinicID(ctx context.Context, clinicID string) (*clinic.FinancialYear, error) {
	query := `SELECT id, clinic_id, fy_label, start_date, end_date, is_current, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at FROM tbl_financial_year WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY start_date DESC LIMIT 1`
	var year clinic.FinancialYear
	err := f.db.GetContext(ctx, &year, query, clinicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial year not found")
		}
		return nil, errors.New("failed to get financial year by clinic id")
	}
	return &year, nil
}

// GetByID implements [port.FinancialYearRepository].
func (f *financialYearRepo) GetByID(ctx context.Context, id int) (*clinic.FinancialYear, error) {
	query := `SELECT id, clinic_id, fy_label, start_date, end_date, is_current, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at FROM tbl_financial_year WHERE id = $1 AND deleted_at IS NULL`
	var year clinic.FinancialYear
	err := f.db.GetContext(ctx, &year, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial year not found")
		}
		return nil, errors.New("failed to get financial year by id")
	}
	return &year, nil
}

// Update implements [port.FinancialYearRepository].
func (f *financialYearRepo) Update(ctx context.Context, year *clinic.FinancialYear) error {
	query := `UPDATE tbl_financial_year SET clinic_id = :clinic_id, fy_label = :fy_label, start_date = :start_date, end_date = :end_date, is_current = :is_current, is_closed = :is_closed, closed_at = :closed_at, closed_by = :closed_by, updated_at = :updated_at WHERE id = :id`
	_, err := f.db.NamedExecContext(ctx, query, year)
	return err
}
