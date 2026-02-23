package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
)

type CustomFormRepository interface {
	Create(ctx context.Context, form *form.Form) error
	GetByID(ctx context.Context, id string) (*form.Form, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]form.Form, error)
	GetPublishedByClinicID(ctx context.Context, clinicID string) ([]form.Form, error)
	Update(ctx context.Context, form *form.Form) error
	Publish(ctx context.Context, id string) error
	Unpublish(ctx context.Context, id string) error
	Archive(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type CustomFormFieldRepository interface {
	Create(ctx context.Context, field *form.Field) error
	CreateBatch(ctx context.Context, fields []*form.Field) error
	GetByID(ctx context.Context, id string) (*form.Field, error)
	GetByFormID(ctx context.Context, formID string) ([]form.Field, error)
	GetByFormVersionID(ctx context.Context, formVersionID string) ([]form.Field, error)
	DeleteByFormID(ctx context.Context, formID string) error
	DeleteByFormVersionID(ctx context.Context, formVersionID string) error
}

type CustomFormVersionRepository interface {
	Create(ctx context.Context, version *form.FormVersion) error
	GetLatestByFormID(ctx context.Context, formID string) (*form.FormVersion, error)
	GetByFormID(ctx context.Context, formID string) ([]form.FormVersion, error)
	SetActive(ctx context.Context, formID string, versionID string) error
}
