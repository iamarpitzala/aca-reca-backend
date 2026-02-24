package postgres

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/jmoiron/sqlx"
)

type clinicFinancialQuarterLockRepo struct {
	db *sqlx.DB
}

func NewClinicFinancialQuarterLockRepository(db *sqlx.DB) port.ClinicFinancialQuarterLockRepository {
	return &clinicFinancialQuarterLockRepo{db: db}
}

func (r *clinicFinancialQuarterLockRepo) CreateLock(ctx context.Context, clinicFinancialYearID, financialQuarterID int) error {
	query := `INSERT INTO tbl_clinic_financial_quarter_lock (clinic_financial_year_id, financial_quarter_id, is_closed, closed_at, closed_by)
VALUES ($1, $2, FALSE, NULL, NULL)
ON CONFLICT (clinic_financial_year_id, financial_quarter_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, clinicFinancialYearID, financialQuarterID)
	return err
}
