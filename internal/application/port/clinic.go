package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type ClinicRepository interface {
	Create(ctx context.Context, clinic *domain.Clinic) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Clinic, error)
	Update(ctx context.Context, clinic *domain.Clinic) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]domain.Clinic, error)
	GetByABN(ctx context.Context, abnNumber string) (*domain.Clinic, error)
}
