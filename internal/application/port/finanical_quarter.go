package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
)

type FinancialQuarterRepository interface {
	Create(ctx context.Context, q *clinic.FinancialQuarter) error
	GetByID(ctx context.Context, id string) (*clinic.FinancialQuarter, error)
	GetByFinancialYearID(ctx context.Context, financialYearID int) (*clinic.FinancialQuarter, error)
	Update(ctx context.Context, q *clinic.FinancialQuarter) error
	Delete(ctx context.Context, id string) error
}
