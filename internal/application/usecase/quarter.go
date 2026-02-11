package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type QuarterService struct {
	repo port.QuarterRepository
}

func NewQuarterService(repo port.QuarterRepository) *QuarterService {
	return &QuarterService{repo: repo}
}

func (s *QuarterService) CreateQuarter(ctx context.Context, q *domain.Quarter) error {
	q.ID = uuid.New()
	now := time.Now()
	q.CreatedAt = now
	q.UpdatedAt = &now
	q.DeletedAt = nil
	return s.repo.Create(ctx, q)
}

func (s *QuarterService) GetQuarterByID(ctx context.Context, id uuid.UUID) (domain.Quarter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *QuarterService) UpdateQuarter(ctx context.Context, q *domain.Quarter) error {
	existing, err := s.repo.GetByID(ctx, q.ID)
	if err != nil {
		return errors.New("quarter not found")
	}
	if existing.DeletedAt != nil {
		return errors.New("cannot update deleted quarter")
	}
	now := time.Now()
	q.CreatedAt = existing.CreatedAt
	q.UpdatedAt = &now
	return s.repo.Update(ctx, q)
}

func (s *QuarterService) DeleteQuarter(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errors.New("quarter not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *QuarterService) ListQuarter(ctx context.Context) ([]domain.Quarter, error) {
	return s.repo.List(ctx)
}
