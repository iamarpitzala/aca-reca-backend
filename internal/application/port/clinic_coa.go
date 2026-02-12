package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type ClinicCOARepository interface {
	Create(ctx context.Context, cc *domain.ClinicCOA) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ClinicCOA, error)
	ListByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.ClinicCOA, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, clinicID, coaID uuid.UUID) (bool, error)
}
