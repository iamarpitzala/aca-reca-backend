package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type QuarterRepository interface {
	Create(ctx context.Context, q *domain.Quarter) error
	Update(ctx context.Context, q *domain.Quarter) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]domain.Quarter, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Quarter, error)
}
