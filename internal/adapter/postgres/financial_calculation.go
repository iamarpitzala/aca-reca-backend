package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type financialCalculationRepo struct {
	db *sqlx.DB
}

func NewFinancialCalculationRepository(db *sqlx.DB) port.FinancialCalculationRepository {
	return &financialCalculationRepo{db: db}
}

func (r *financialCalculationRepo) Create(ctx context.Context, calc *domain.FinancialCalculation) error {
	inputJSON, err := json.Marshal(calc.InputData)
	if err != nil {
		return errors.New("failed to marshal input data")
	}
	calculatedJSON, err := json.Marshal(calc.CalculatedData)
	if err != nil {
		return errors.New("failed to marshal calculated data")
	}
	basJSON, err := json.Marshal(calc.BASMapping)
	if err != nil {
		return errors.New("failed to marshal BAS mapping")
	}
	query := `INSERT INTO tbl_financial_calculation (id, financial_form_id, input_data, calculated_data, bas_mapping, created_at, created_by)
		VALUES (:id, :financial_form_id, :input_data, :calculated_data, :bas_mapping, :created_at, :created_by)`
	args := map[string]interface{}{
		"id":                calc.ID,
		"financial_form_id": calc.FinancialFormID,
		"input_data":        string(inputJSON),
		"calculated_data":   string(calculatedJSON),
		"bas_mapping":       string(basJSON),
		"created_at":        calc.CreatedAt,
		"created_by":         calc.CreatedBy,
	}
	_, err = r.db.NamedExecContext(ctx, query, args)
	return err
}

func (r *financialCalculationRepo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]domain.FinancialCalculation, error) {
	type calcRow struct {
		ID              uuid.UUID       `db:"id"`
		FinancialFormID uuid.UUID       `db:"financial_form_id"`
		InputData       json.RawMessage `db:"input_data"`
		CalculatedData  json.RawMessage `db:"calculated_data"`
		BASMapping      json.RawMessage `db:"bas_mapping"`
		CreatedAt       time.Time       `db:"created_at"`
		CreatedBy       *uuid.UUID      `db:"created_by"`
	}
	query := `SELECT id, financial_form_id, input_data, calculated_data, bas_mapping, created_at, created_by FROM tbl_financial_calculation WHERE financial_form_id = $1 ORDER BY created_at DESC`
	var rows []calcRow
	if err := r.db.SelectContext(ctx, &rows, query, formID); err != nil {
		return nil, errors.New("failed to get calculations")
	}
	out := make([]domain.FinancialCalculation, len(rows))
	for i, row := range rows {
		var inputData, calculatedData, basMapping map[string]interface{}
		if err := json.Unmarshal(row.InputData, &inputData); err != nil {
			return nil, errors.New("failed to unmarshal input data")
		}
		if err := json.Unmarshal(row.CalculatedData, &calculatedData); err != nil {
			return nil, errors.New("failed to unmarshal calculated data")
		}
		if err := json.Unmarshal(row.BASMapping, &basMapping); err != nil {
			return nil, errors.New("failed to unmarshal BAS mapping")
		}
		out[i] = domain.FinancialCalculation{
			ID:              row.ID,
			FinancialFormID: row.FinancialFormID,
			InputData:       inputData,
			CalculatedData:  calculatedData,
			BASMapping:      basMapping,
			CreatedAt:       row.CreatedAt,
			CreatedBy:       row.CreatedBy,
		}
	}
	return out, nil
}
