package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
)

type ClinicFinancialSettingRepository interface {
	Create(ctx context.Context, settings *clinic.ClinicFinancialSetting) error
	GetByClinicID(ctx context.Context, clinicID string) (*clinic.ClinicFinancialSetting, error)
	GetByClinicFinancialYearID(ctx context.Context, clinicFinancialYearID int) (*clinic.ClinicFinancialSetting, error)
	Update(ctx context.Context, settings *clinic.ClinicFinancialSetting) error
	Delete(ctx context.Context, clinicID string) error
}
