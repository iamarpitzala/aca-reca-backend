package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type QuarterService struct {
	db *sqlx.DB
}

func NewQuarterService(db *sqlx.DB) *QuarterService {
	return &QuarterService{
		db: db,
	}
}

func (qs *QuarterService) CreateQuarter(ctx context.Context, form *domain.Quarter) error {

	form.ID = uuid.New()
	form.CreatedAt = time.Now()

	return repository.CreateQuarter(ctx, qs.db, form)
}

func (qs *QuarterService) GetQuaterByID(ctx context.Context, id uuid.UUID) (domain.Quarter, error) {
	return repository.GetQuarterByID(ctx, qs.db, id)
}

func (qs *QuarterService) UpdateQuarter(ctx context.Context, form *domain.Quarter) error {

	existing, err := repository.GetQuarterByID(ctx, qs.db, form.ID)
	if err != nil {
		return errors.New("quarter not found")
	}

	if existing.DeletedAt != nil {
		return errors.New("cannot update deleted quarter")
	}

	updatedAt := time.Now()

	form.CreatedAt = existing.CreatedAt
	form.UpdatedAt = &updatedAt

	return repository.UpdateQuarter(ctx, qs.db, form)
}

func (qs *QuarterService) DeleteQuarter(ctx context.Context, id uuid.UUID) error {
	_, err := repository.GetQuarterByID(ctx, qs.db, id)
	if err != nil {
		return errors.New("quarter not found")
	}

	return repository.DeleteQuarter(ctx, qs.db, id)
}

func (qs *QuarterService) ListQuarter(ctx context.Context) ([]domain.Quarter, error) {
	return repository.ListQuarter(ctx, qs.db)
}
