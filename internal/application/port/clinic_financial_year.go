package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
)

type ClinicFinancialYearRepository interface {
	Create(ctx context.Context, cfy *clinic.ClinicFinancialYear) (int, error)
	GetByID(ctx context.Context, id int) (*clinic.ClinicFinancialYear, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]clinic.ClinicFinancialYear, error)
	GetCurrentByClinicID(ctx context.Context, clinicID string) (*clinic.ClinicFinancialYear, error)
}
