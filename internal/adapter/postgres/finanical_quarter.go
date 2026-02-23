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

type financialQuarterRepo struct {
	db *sqlx.DB
}

func NewFinancialQuarterRepository(db *sqlx.DB) port.FinancialQuarterRepository {
	return &financialQuarterRepo{db: db}
}

// Create implements [port.FinancialQuarterRepository].
func (f *financialQuarterRepo) Create(ctx context.Context, q *clinic.FinancialQuarter) error {
	query := `INSERT INTO tbl_financial_quarter (id, financial_year_id, name, quarter_number, start_date, end_date, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at) VALUES (:id, :financial_year_id, :name, :quarter_number, :start_date, :end_date, :is_closed, :closed_at, :closed_by, :created_at, :updated_at, :deleted_at)`
	_, err := f.db.NamedExecContext(ctx, query, q)
	return err
}

// Delete implements [port.FinancialQuarterRepository].
func (f *financialQuarterRepo) Delete(ctx context.Context, id int) error {
	query := `UPDATE tbl_financial_quarter SET deleted_at = $1 WHERE id = $2`
	_, err := f.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// GetByFinancialYearID implements [port.FinancialQuarterRepository].
func (f *financialQuarterRepo) GetByFinancialYearID(ctx context.Context, financialYearID int) (*clinic.FinancialQuarter, error) {
	query := `SELECT id, financial_year_id, name, quarter_number, start_date, end_date, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at FROM tbl_financial_quarter WHERE financial_year_id = $1 AND deleted_at IS NULL ORDER BY quarter_number ASC LIMIT 1`
	var q clinic.FinancialQuarter
	err := f.db.GetContext(ctx, &q, query, financialYearID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial quarter not found")
		}
		return nil, errors.New("failed to get financial quarter by financial year id")
	}
	return &q, nil
}

// GetByID implements [port.FinancialQuarterRepository].
func (f *financialQuarterRepo) GetByID(ctx context.Context, id int) (*clinic.FinancialQuarter, error) {
	query := `SELECT id, financial_year_id, name, quarter_number, start_date, end_date, is_closed, closed_at, closed_by, created_at, updated_at, deleted_at FROM tbl_financial_quarter WHERE id = $1 AND deleted_at IS NULL`
	var q clinic.FinancialQuarter
	err := f.db.GetContext(ctx, &q, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial quarter not found")
		}
		return nil, errors.New("failed to get financial quarter by id")
	}
	return &q, nil
}

// Update implements [port.FinancialQuarterRepository].
func (f *financialQuarterRepo) Update(ctx context.Context, q *clinic.FinancialQuarter) error {
	query := `UPDATE tbl_financial_quarter SET financial_year_id = :financial_year_id, name = :name, quarter_number = :quarter_number, start_date = :start_date, end_date = :end_date, is_closed = :is_closed, closed_at = :closed_at, closed_by = :closed_by, updated_at = :updated_at WHERE id = :id`
	_, err := f.db.NamedExecContext(ctx, query, q)
	return err
}
