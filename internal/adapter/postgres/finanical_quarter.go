package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
	"github.com/jmoiron/sqlx"
)

type financialQuarterRepo struct {
	db *sqlx.DB
}

func NewFinancialQuarterRepository(db *sqlx.DB) port.FinancialQuarterRepository {
	return &financialQuarterRepo{db: db}
}

func (f *financialQuarterRepo) Create(ctx context.Context, q *clinic.FinancialQuarter) error {
	query := `INSERT INTO tbl_financial_quarter (financial_year_id, quarter_number, name, start_date, end_date, created_at, updated_at)
VALUES (:financial_year_id, :quarter_number, :name, :start_date, :end_date, now(), now())
RETURNING id`
	rows, err := f.db.NamedQueryContext(ctx, query, q)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		_ = rows.Scan(&q.ID)
	}
	return nil
}

func (f *financialQuarterRepo) CreateQuarter(ctx context.Context, financialYearID int, name string, quarterNumber int, startDate, endDate string) error {
	query := `INSERT INTO tbl_financial_quarter (financial_year_id, quarter_number, name, start_date, end_date, created_at, updated_at)
VALUES ($1, $2, $3, $4::date, $5::date, now(), now())`
	_, err := f.db.ExecContext(ctx, query, financialYearID, quarterNumber, name, startDate, endDate)
	return err
}

func (f *financialQuarterRepo) GetByID(ctx context.Context, id int) (*clinic.FinancialQuarter, error) {
	query := `SELECT id, financial_year_id, quarter_number, name, start_date, end_date, created_at, updated_at
FROM tbl_financial_quarter WHERE id = $1`
	var q clinic.FinancialQuarter
	err := f.db.GetContext(ctx, &q, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (f *financialQuarterRepo) GetByFinancialYearID(ctx context.Context, financialYearID int) (*clinic.FinancialQuarter, error) {
	query := `SELECT id, financial_year_id, quarter_number, name, start_date, end_date, created_at, updated_at
FROM tbl_financial_quarter WHERE financial_year_id = $1 ORDER BY quarter_number ASC LIMIT 1`
	var q clinic.FinancialQuarter
	err := f.db.GetContext(ctx, &q, query, financialYearID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial quarter not found")
		}
		return nil, err
	}
	return &q, nil
}

func (f *financialQuarterRepo) ListByFinancialYearID(ctx context.Context, financialYearID int) ([]clinic.FinancialQuarter, error) {
	query := `SELECT id, financial_year_id, quarter_number, name, start_date, end_date, created_at, updated_at
FROM tbl_financial_quarter WHERE financial_year_id = $1 ORDER BY quarter_number ASC`
	var list []clinic.FinancialQuarter
	err := f.db.SelectContext(ctx, &list, query, financialYearID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []clinic.FinancialQuarter{}
	}
	return list, nil
}

func (f *financialQuarterRepo) Update(ctx context.Context, q *clinic.FinancialQuarter) error {
	query := `UPDATE tbl_financial_quarter SET financial_year_id = :financial_year_id, quarter_number = :quarter_number, name = :name, start_date = :start_date, end_date = :end_date, updated_at = now() WHERE id = :id`
	_, err := f.db.NamedExecContext(ctx, query, q)
	return err
}

func (f *financialQuarterRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tbl_financial_quarter WHERE id = $1`
	_, err := f.db.ExecContext(ctx, query, id)
	return err
}
