package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type ClinicFinancialSettingsRepository interface {
	Create(ctx context.Context, settings *domain.ClinicFinancialSettings) error
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.ClinicFinancialSettings, error)
	Update(ctx context.Context, settings *domain.ClinicFinancialSettings) error
	Delete(ctx context.Context, clinicID uuid.UUID) error
}
