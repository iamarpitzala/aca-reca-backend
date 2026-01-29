package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateFinancialForm(ctx context.Context, db *sqlx.DB, form *domain.FinancialForm) error {
	query := `INSERT INTO tbl_financial_form (id, clinic_id, quarter_id, name, calculation_method, configuration, is_active, created_at, updated_at)
		VALUES (:id, :clinic_id, :quarter_id,  :name, :calculation_method, :configuration, :is_active, :created_at, :updated_at)`

	_, err := db.NamedExecContext(ctx, query, form)
	if err != nil {
		return err
	}
	return nil
}

func GetFinancialFormByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.FinancialFormResponse, error) {
	query := `SELECT id, clinic_id, quarter_id, name, calculation_method, configuration, is_active, created_at, updated_at, deleted_at
		FROM tbl_financial_form WHERE id = $1 AND deleted_at IS NULL`

	var row domain.FinancialForm
	err := db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial form not found")
		}
		return nil, errors.New("failed to get financial form")
	}

	response, err := row.ToFinancialFormResponse()
	if err != nil {
		return &domain.FinancialFormResponse{}, err
	}

	return response, nil
}

func GetFinancialFormsByClinicID(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.FinancialFormResponse, error) {

	query := `
		SELECT
			id,
			clinic_id,
			quarter_id,
			name,
			calculation_method,
			configuration,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_financial_form
		WHERE clinic_id = $1
		  AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var rows []domain.FinancialForm
	if err := db.SelectContext(ctx, &rows, query, clinicID); err != nil {
		return nil, fmt.Errorf("failed to get financial forms for clinic %s: %w", clinicID, err)
	}

	responses := make([]domain.FinancialFormResponse, 0, len(rows))

	for _, row := range rows {
		resp, err := row.ToFinancialFormResponse()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to map financial form %s: %w",
				row.ID,
				err,
			)
		}
		responses = append(responses, *resp)
	}

	return responses, nil
}

func UpdateFinancialForm(ctx context.Context, db *sqlx.DB, form *domain.FinancialForm) error {

	query := `UPDATE tbl_financial_form SET name = :name, calculation_method = :calculation_method, 
		configuration = :configuration, is_active = :is_active, updated_at = :updated_at WHERE id = :id`

	_, err := db.NamedExecContext(ctx, query, form)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFinancialForm(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `UPDATE tbl_financial_form SET deleted_at = $1 WHERE id = $2`
	_, err := db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

// GST Repository
func GetGSTByID(ctx context.Context, db *sqlx.DB, id int) (*domain.GST, error) {
	query := `SELECT id, name, type, percentage FROM tbl_gst WHERE id = $1`
	var gst domain.GST
	err := db.GetContext(ctx, &gst, query, id)
	if err != nil {
		return nil, err
	}
	return &gst, nil
}
