package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/user"
)

type UserClinicRepository interface {
	Create(ctx context.Context, uc *user.UserClinic) error
	GetByID(ctx context.Context, id string) (*user.UserClinic, error)
	GetByUserAndClinic(ctx context.Context, userID, clinicID string) (*user.UserClinic, error)
	GetUserClinics(ctx context.Context, userID string) ([]user.UserClinicWithClinic, error)
	Delete(ctx context.Context, id string) error
}
