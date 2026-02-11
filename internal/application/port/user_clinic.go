package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type UserClinicRepository interface {
	Create(ctx context.Context, uc *domain.UserClinic) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.UserClinic, error)
	GetByUserAndClinic(ctx context.Context, userID, clinicID uuid.UUID) (*domain.UserClinic, error)
	GetUserClinics(ctx context.Context, userID uuid.UUID) ([]domain.UserClinicWithClinic, error)
	GetClinicUsers(ctx context.Context, clinicID uuid.UUID) ([]domain.UserClinicWithUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
