package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateFinancialCalculation(ctx context.Context, db *sqlx.DB, calculation *domain.FinancialCalculation) error {
	inputJSON, err := json.Marshal(calculation.InputData)
	if err != nil {
		return errors.New("failed to marshal input data")
	}

	calculatedJSON, err := json.Marshal(calculation.CalculatedData)
	if err != nil {
		return errors.New("failed to marshal calculated data")
	}

	basJSON, err := json.Marshal(calculation.BASMapping)
	if err != nil {
		return errors.New("failed to marshal BAS mapping")
	}

	query := `INSERT INTO tbl_financial_calculation (id, financial_form_id, input_data, calculated_data, bas_mapping, created_at, created_by)
		VALUES (:id, :financial_form_id, :input_data, :calculated_data, :bas_mapping, :created_at, :created_by)`

	args := map[string]interface{}{
		"id":                calculation.ID,
		"financial_form_id": calculation.FinancialFormID,
		"input_data":        string(inputJSON),
		"calculated_data":   string(calculatedJSON),
		"bas_mapping":       string(basJSON),
		"created_at":        calculation.CreatedAt,
		"created_by":        calculation.CreatedBy,
	}

	_, err = db.NamedExecContext(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func GetCalculationsByFormID(ctx context.Context, db *sqlx.DB, formID uuid.UUID) ([]domain.FinancialCalculation, error) {
	type calcRow struct {
		ID              uuid.UUID       `db:"id"`
		FinancialFormID uuid.UUID       `db:"financial_form_id"`
		InputData       json.RawMessage `db:"input_data"`
		CalculatedData  json.RawMessage `db:"calculated_data"`
		BASMapping      json.RawMessage `db:"bas_mapping"`
		CreatedAt       time.Time       `db:"created_at"`
		CreatedBy       *uuid.UUID      `db:"created_by"`
	}

	query := `SELECT id, financial_form_id, input_data, calculated_data, bas_mapping, created_at, created_by
		FROM tbl_financial_calculation WHERE financial_form_id = $1 ORDER BY created_at DESC`

	var rows []calcRow
	err := db.SelectContext(ctx, &rows, query, formID)
	if err != nil {
		return nil, errors.New("failed to get calculations")
	}

	calculations := make([]domain.FinancialCalculation, len(rows))
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

		calculations[i] = domain.FinancialCalculation{
			ID:              row.ID,
			FinancialFormID: row.FinancialFormID,
			InputData:       inputData,
			CalculatedData:  calculatedData,
			BASMapping:      basMapping,
			CreatedAt:       row.CreatedAt,
			CreatedBy:       row.CreatedBy,
		}
	}

	return calculations, nil
}
