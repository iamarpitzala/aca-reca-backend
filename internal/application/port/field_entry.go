package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
)

type FieldEntryRepository interface {
	Create(ctx context.Context, entry *form.FieldEntry) error
	GetByID(ctx context.Context, id string) (*form.FieldEntry, error)
	GetByFormID(ctx context.Context, formID string) ([]form.FieldEntry, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]form.FieldEntry, error)
	Update(ctx context.Context, entry *form.FieldEntry) error
	Delete(ctx context.Context, id string) error
}
