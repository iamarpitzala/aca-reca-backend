package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type FinancialCalculationRepository interface {
	Create(ctx context.Context, calc *domain.FinancialCalculation) error
	GetByFormID(ctx context.Context, formID uuid.UUID) ([]domain.FinancialCalculation, error)
}
