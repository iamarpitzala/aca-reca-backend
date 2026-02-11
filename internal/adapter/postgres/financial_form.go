package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type financialFormRepo struct {
	db *sqlx.DB
}

func NewFinancialFormRepository(db *sqlx.DB) port.FinancialFormRepository {
	return &financialFormRepo{db: db}
}

func (r *financialFormRepo) Create(ctx context.Context, form *domain.FinancialForm) error {
	query := `INSERT INTO tbl_financial_form (id, clinic_id, quarter_id, name, calculation_method, configuration, is_active, created_at, updated_at)
		VALUES (:id, :clinic_id, :quarter_id, :name, :calculation_method, :configuration, :is_active, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, form)
	return err
}

func (r *financialFormRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.FinancialFormResponse, error) {
	query := `SELECT id, clinic_id, quarter_id, name, calculation_method, configuration, is_active, created_at, updated_at, deleted_at FROM tbl_financial_form WHERE id = $1 AND deleted_at IS NULL`
	var row domain.FinancialForm
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial form not found")
		}
		return nil, errors.New("failed to get financial form")
	}
	return row.ToFinancialFormResponse()
}

func (r *financialFormRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.FinancialFormResponse, error) {
	query := `SELECT id, clinic_id, quarter_id, name, calculation_method, configuration, is_active, created_at, updated_at, deleted_at FROM tbl_financial_form WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	var rows []domain.FinancialForm
	if err := r.db.SelectContext(ctx, &rows, query, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get financial forms for clinic %s: %w", clinicID, err)
	}
	out := make([]domain.FinancialFormResponse, 0, len(rows))
	for i := range rows {
		resp, err := rows[i].ToFinancialFormResponse()
		if err != nil {
			return nil, fmt.Errorf("failed to map financial form %s: %w", rows[i].ID, err)
		}
		out = append(out, *resp)
	}
	return out, nil
}

func (r *financialFormRepo) Update(ctx context.Context, form *domain.FinancialForm) error {
	query := `UPDATE tbl_financial_form SET name = :name, calculation_method = :calculation_method, configuration = :configuration, is_active = :is_active, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, form)
	return err
}

func (r *financialFormRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tbl_financial_form SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *financialFormRepo) GetGSTByID(ctx context.Context, id int) (*domain.GST, error) {
	query := `SELECT id, name, type, percentage FROM tbl_gst WHERE id = $1`
	var gst domain.GST
	err := r.db.GetContext(ctx, &gst, query, id)
	if err != nil {
		return nil, err
	}
	return &gst, nil
}
