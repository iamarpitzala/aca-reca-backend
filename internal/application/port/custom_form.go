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
	Unpublish(ctx context.Context, id uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomFormFieldRepository interface {
	Create(ctx context.Context, field *domain.CustomFormField) error
	CreateBatch(ctx context.Context, fields []*domain.CustomFormField) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormField, error)
	GetByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormField, error)
	GetByFormVersionID(ctx context.Context, formVersionID uuid.UUID) ([]domain.CustomFormField, error)
	DeleteByFormID(ctx context.Context, formID uuid.UUID) error
	DeleteByFormVersionID(ctx context.Context, formVersionID uuid.UUID) error
}

type CustomFormVersionRepository interface {
	Create(ctx context.Context, version *domain.CustomFormVersion) error
	GetLatestByFormID(ctx context.Context, formID uuid.UUID) (*domain.CustomFormVersion, error)
	GetByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormVersion, error)
	SetActive(ctx context.Context, formID uuid.UUID, versionID uuid.UUID) error
}

type CustomFormCalculationRepository interface {
	GetByFormVersionID(ctx context.Context, formVersionID uuid.UUID) (*domain.CustomFormCalculation, error)
	GetByFormID(ctx context.Context, formID uuid.UUID) (*domain.CustomFormCalculation, error)
}
