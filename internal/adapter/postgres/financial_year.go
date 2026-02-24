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

func (f *financialYearRepo) Create(ctx context.Context, year *clinic.FinancialYear) error {
	query := `INSERT INTO tbl_financial_year (fy_label, start_date, end_date, is_active, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, now(), now(), NULL)
RETURNING id`
	err := f.db.GetContext(ctx, &year.ID, query, year.FYLabel, year.StartDate, year.EndDate, year.IsActive)
	return err
}

func (f *financialYearRepo) GetByID(ctx context.Context, id int) (*clinic.FinancialYear, error) {
	query := `SELECT id, fy_label, start_date, end_date, is_active, created_at, updated_at, deleted_at
FROM tbl_financial_year WHERE id = $1 AND deleted_at IS NULL`
	var year clinic.FinancialYear
	err := f.db.GetContext(ctx, &year, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial year not found")
		}
		return nil, err
	}
	return &year, nil
}

func (f *financialYearRepo) GetByFYLabel(ctx context.Context, fyLabel string) (*clinic.FinancialYear, error) {
	query := `SELECT id, fy_label, start_date, end_date, is_active, created_at, updated_at, deleted_at
FROM tbl_financial_year WHERE fy_label = $1 AND deleted_at IS NULL`
	var year clinic.FinancialYear
	err := f.db.GetContext(ctx, &year, query, fyLabel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &year, nil
}

func (f *financialYearRepo) List(ctx context.Context) ([]clinic.FinancialYear, error) {
	query := `SELECT id, fy_label, start_date, end_date, is_active, created_at, updated_at, deleted_at
FROM tbl_financial_year WHERE deleted_at IS NULL ORDER BY start_date DESC`
	var list []clinic.FinancialYear
	err := f.db.SelectContext(ctx, &list, query)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []clinic.FinancialYear{}
	}
	return list, nil
}

func (f *financialYearRepo) Update(ctx context.Context, year *clinic.FinancialYear) error {
	query := `UPDATE tbl_financial_year SET fy_label = $1, start_date = $2, end_date = $3, is_active = $4, updated_at = now() WHERE id = $5`
	_, err := f.db.ExecContext(ctx, query, year.FYLabel, year.StartDate, year.EndDate, year.IsActive, year.ID)
	return err
}

func (f *financialYearRepo) Delete(ctx context.Context, id int) error {
	query := `UPDATE tbl_financial_year SET deleted_at = $1 WHERE id = $2`
	_, err := f.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
