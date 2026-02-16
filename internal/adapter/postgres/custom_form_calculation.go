package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type customFormCalculationRepo struct {
	db *sqlx.DB
}

func NewCustomFormCalculationRepository(db *sqlx.DB) port.CustomFormCalculationRepository {
	return &customFormCalculationRepo{db: db}
}

// GetByFormVersionID retrieves calculation settings for a specific form version
func (r *customFormCalculationRepo) GetByFormVersionID(ctx context.Context, formVersionID uuid.UUID) (*domain.CustomFormCalculation, error) {
	q := `
		SELECT
			id,
			form_version_id,
			form_id,
			calculation_method,
			default_payment_responsibility,
			service_facility_fee_percent,
			outwork_enabled,
			outwork_rate_percent,
			created_at
		FROM tbl_custom_form_calculation
		WHERE form_version_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var calc domain.CustomFormCalculation
	err := r.db.GetContext(ctx, &calc, q, formVersionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil if not found (not an error - fallback to form-level)
		}
		return nil, fmt.Errorf("failed to get calculation by form version id: %w", err)
	}
	return &calc, nil
}

// GetByFormID retrieves the latest calculation settings for a form (across all versions)
func (r *customFormCalculationRepo) GetByFormID(ctx context.Context, formID uuid.UUID) (*domain.CustomFormCalculation, error) {
	q := `
		SELECT
			id,
			form_version_id,
			form_id,
			calculation_method,
			default_payment_responsibility,
			service_facility_fee_percent,
			outwork_enabled,
			outwork_rate_percent,
			created_at
		FROM tbl_custom_form_calculation
		WHERE form_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var calc domain.CustomFormCalculation
	err := r.db.GetContext(ctx, &calc, q, formID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil if not found
		}
		return nil, fmt.Errorf("failed to get calculation by form id: %w", err)
	}
	return &calc, nil
}
