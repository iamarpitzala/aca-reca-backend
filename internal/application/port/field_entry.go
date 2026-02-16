package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type FieldEntryRepository interface {
	Create(ctx context.Context, entry *domain.FieldEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.FieldEntry, error)
	GetByFormID(ctx context.Context, formID uuid.UUID) ([]domain.FieldEntry, error)
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.FieldEntry, error)
	Update(ctx context.Context, entry *domain.FieldEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
}
