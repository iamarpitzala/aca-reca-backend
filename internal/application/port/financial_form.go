package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type FinancialFormRepository interface {
	Create(ctx context.Context, form *domain.FinancialForm) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.FinancialFormResponse, error)
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.FinancialFormResponse, error)
	Update(ctx context.Context, form *domain.FinancialForm) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetGSTByID(ctx context.Context, id int) (*domain.GST, error)
}
