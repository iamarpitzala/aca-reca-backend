package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
)

type ClinicRepository interface {
	Create(ctx context.Context, clinic *clinic.Clinic) error
	GetByID(ctx context.Context, id string) (*clinic.Clinic, error)
	Update(ctx context.Context, clinic *clinic.Clinic) error
	SetActive(ctx context.Context, id string, active bool) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]clinic.Clinic, error)
	ABNExists(ctx context.Context, abnNumber string) (bool, error)
	GetByABN(ctx context.Context, abnNumber string) (*clinic.Clinic, error)
}
