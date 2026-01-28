package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateFinancialForm(ctx context.Context, db *sqlx.DB, form *domain.FinancialForm) error {
	configJSON, err := json.Marshal(form.Configuration)
	if err != nil {
		return errors.New("failed to marshal configuration")
	}

	query := `INSERT INTO tbl_financial_form (id, clinic_id, quarter_id, gst_id, name, calculation_method, configuration, is_active, created_at, updated_at)
		VALUES (:id, :clinic_id, :quarter_id, :gst_id, :name, :calculation_method, :configuration, :is_active, :created_at, :updated_at)`

	args := map[string]interface{}{
		"id":                 form.ID,
		"clinic_id":          form.ClinicID,
		"quarter_id":         form.QuarterID,
		"gst_id":             form.GSTID,
		"name":               form.Name,
		"calculation_method": form.CalculationMethod,
		"configuration":      string(configJSON),
		"is_active":          form.IsActive,
		"created_at":         form.CreatedAt,
		"updated_at":         form.UpdatedAt,
	}

	_, err = db.NamedExecContext(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func GetFinancialFormByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*domain.FinancialForm, error) {
	type formRow struct {
		ID                uuid.UUID       `db:"id"`
		ClinicID          uuid.UUID       `db:"clinic_id"`
		QuarterID         uuid.UUID       `db:"quarter_id"`
		GstID             *int            `db:"gst_id"`
		Name              string          `db:"name"`
		CalculationMethod string          `db:"calculation_method"`
		Configuration     json.RawMessage `db:"configuration"`
		IsActive          bool            `db:"is_active"`
		CreatedAt         time.Time       `db:"created_at"`
		UpdatedAt         time.Time       `db:"updated_at"`
		DeletedAt         *time.Time      `db:"deleted_at"`
	}

	query := `SELECT id, clinic_id, quarter_id, gst_id, name, calculation_method, configuration, is_active, created_at, updated_at, deleted_at
		FROM tbl_financial_form WHERE id = $1 AND deleted_at IS NULL`

	var row formRow
	err := db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("financial form not found")
		}
		return nil, errors.New("failed to get financial form")
	}

	var config map[string]interface{}
	if err := json.Unmarshal(row.Configuration, &config); err != nil {
		return nil, errors.New("failed to unmarshal configuration")
	}

	return &domain.FinancialForm{
		ID:                row.ID,
		ClinicID:          row.ClinicID,
		QuarterID:         row.QuarterID,
		GSTID:             row.GstID,
		Name:              row.Name,
		CalculationMethod: row.CalculationMethod,
		Configuration:     config,
		IsActive:          row.IsActive,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
		DeletedAt:         row.DeletedAt,
	}, nil
}

func GetFinancialFormsByClinicID(ctx context.Context, db *sqlx.DB, clinicID uuid.UUID) ([]domain.FinancialForm, error) {
	type formRow struct {
		ID                uuid.UUID       `db:"id"`
		ClinicID          uuid.UUID       `db:"clinic_id"`
		QuarterID         uuid.UUID       `db:"quarter_id"`
		GSTID             *int            `db:"gst_id"`
		Name              string          `db:"name"`
		CalculationMethod string          `db:"calculation_method"`
		Configuration     json.RawMessage `db:"configuration"`
		IsActive          bool            `db:"is_active"`
		CreatedAt         time.Time       `db:"created_at"`
		UpdatedAt         time.Time       `db:"updated_at"`
		DeletedAt         *time.Time      `db:"deleted_at"`
	}

	query := `SELECT id, clinic_id, quarter_id, gst_id, name, calculation_method, configuration, is_active, created_at, updated_at, deleted_at
		FROM tbl_financial_form WHERE clinic_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`

	var rows []formRow
	err := db.SelectContext(ctx, &rows, query, clinicID)
	if err != nil {
		return nil, errors.New("failed to get financial forms")
	}

	forms := make([]domain.FinancialForm, len(rows))
	for i, row := range rows {
		var config map[string]interface{}
		if err := json.Unmarshal(row.Configuration, &config); err != nil {
			return nil, errors.New("failed to unmarshal configuration")
		}

		forms[i] = domain.FinancialForm{
			ID:                row.ID,
			ClinicID:          row.ClinicID,
			QuarterID:         row.QuarterID,
			GSTID:             row.GSTID,
			Name:              row.Name,
			CalculationMethod: row.CalculationMethod,
			Configuration:     config,
			IsActive:          row.IsActive,
			CreatedAt:         row.CreatedAt,
			UpdatedAt:         row.UpdatedAt,
			DeletedAt:         row.DeletedAt,
		}
	}

	return forms, nil
}

func UpdateFinancialForm(ctx context.Context, db *sqlx.DB, form *domain.FinancialForm) error {
	configJSON, err := json.Marshal(form.Configuration)
	if err != nil {
		return errors.New("failed to marshal configuration")
	}

	query := `UPDATE tbl_financial_form SET name = :name, calculation_method = :calculation_method, 
		configuration = :configuration, is_active = :is_active, updated_at = :updated_at WHERE id = :id`

	args := map[string]interface{}{
		"id":                 form.ID,
		"name":               form.Name,
		"calculation_method": form.CalculationMethod,
		"configuration":      string(configJSON),
		"is_active":          form.IsActive,
		"updated_at":         form.UpdatedAt,
	}

	_, err = db.NamedExecContext(ctx, query, args)
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
