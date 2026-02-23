package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
)

type FinancialYearRepository interface {
	Create(ctx context.Context, year *clinic.FinancialYear) error
	GetByID(ctx context.Context, id int) (*clinic.FinancialYear, error)
	GetByClinicID(ctx context.Context, clinicID string) (*clinic.FinancialYear, error)
	Update(ctx context.Context, year *clinic.FinancialYear) error
	Delete(ctx context.Context, id int) error
}
