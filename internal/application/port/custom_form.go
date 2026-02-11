package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type CustomFormRepository interface {
	Create(ctx context.Context, form *domain.CustomForm) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomForm, error)
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomForm, error)
	GetPublishedByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomForm, error)
	Update(ctx context.Context, form *domain.CustomForm) error
	Publish(ctx context.Context, id uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateEntry(ctx context.Context, entry *domain.CustomFormEntry) error
	GetEntryByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormEntry, error)
	GetEntriesByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormEntry, error)
	GetEntriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormEntry, error)
	GetEntriesByQuarter(ctx context.Context, clinicID, quarterID uuid.UUID) ([]domain.CustomFormEntry, error)
	UpdateEntry(ctx context.Context, entry *domain.CustomFormEntry) error
	DeleteEntry(ctx context.Context, id uuid.UUID) error
}
