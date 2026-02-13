package port

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type BASSnapshotRepository interface {
	Create(ctx context.Context, snapshot *domain.BASSnapshot) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BASSnapshot, error)
	GetByClinicIDAndPeriod(ctx context.Context, clinicID uuid.UUID, periodStart, periodEnd time.Time) (*domain.BASSnapshot, error)
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.BASSnapshot, error)
	GetFinalisedByClinicIDs(ctx context.Context, clinicIDs []uuid.UUID, periodStart, periodEnd time.Time) ([]domain.BASSnapshot, error)
	Update(ctx context.Context, snapshot *domain.BASSnapshot) error
	Delete(ctx context.Context, id uuid.UUID) error
}
