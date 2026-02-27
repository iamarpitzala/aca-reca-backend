package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
)

// CustomFormEntryRepository handles tbl_custom_form_entry and tbl_custom_form_entry_value.
type CustomFormEntryRepository interface {
	Create(ctx context.Context, entry *form.Entry, values []form.EntryValue) error
	GetByID(ctx context.Context, id string) (*form.Entry, []form.EntryValue, error)
	GetByFormVersionID(ctx context.Context, formVersionID int) ([]form.Entry, error)
	GetByFormID(ctx context.Context, formID string) ([]form.Entry, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]form.Entry, error)
	Update(ctx context.Context, entryID string, values []form.EntryValue) error
	Delete(ctx context.Context, id string) error
}
