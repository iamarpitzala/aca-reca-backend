package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
	"github.com/jmoiron/sqlx"
)

type clinicFinancialYearRepo struct {
	db *sqlx.DB
}

func NewClinicFinancialYearRepository(db *sqlx.DB) port.ClinicFinancialYearRepository {
	return &clinicFinancialYearRepo{db: db}
}

func (r *clinicFinancialYearRepo) Create(ctx context.Context, cfy *clinic.ClinicFinancialYear) (int, error) {
	query := `INSERT INTO tbl_clinic_financial_year (clinic_id, financial_year_id, is_current, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now(), NULL)
RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query,
		cfy.ClinicID, cfy.FinancialYearID, cfy.IsCurrent, cfy.IsClosed, cfy.ClosedAt, cfy.ClosedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *clinicFinancialYearRepo) GetByID(ctx context.Context, id int) (*clinic.ClinicFinancialYear, error) {
	query := `SELECT id, clinic_id, financial_year_id, is_current, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at
FROM tbl_clinic_financial_year WHERE id = $1 AND deleted_at IS NULL`
	var cfy clinic.ClinicFinancialYear
	err := r.db.GetContext(ctx, &cfy, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cfy, nil
}

func (r *clinicFinancialYearRepo) GetByClinicID(ctx context.Context, clinicID string) ([]clinic.ClinicFinancialYear, error) {
	query := `SELECT id, clinic_id, financial_year_id, is_current, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at
FROM tbl_clinic_financial_year WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	var list []clinic.ClinicFinancialYear
	err := r.db.SelectContext(ctx, &list, query, clinicID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []clinic.ClinicFinancialYear{}
	}
	return list, nil
}

func (r *clinicFinancialYearRepo) GetCurrentByClinicID(ctx context.Context, clinicID string) (*clinic.ClinicFinancialYear, error) {
	query := `SELECT id, clinic_id, financial_year_id, is_current, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at
FROM tbl_clinic_financial_year WHERE clinic_id = $1 AND is_current = true AND deleted_at IS NULL LIMIT 1`
	var cfy clinic.ClinicFinancialYear
	err := r.db.GetContext(ctx, &cfy, query, clinicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no current financial year for clinic")
		}
		return nil, err
	}
	return &cfy, nil
}
