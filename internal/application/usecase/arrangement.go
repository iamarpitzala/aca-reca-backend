package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/coa"
)

var ErrArrangementNotFound = errors.New("arrangement not found")

type ArrangementService struct {
	repo port.ArrangementRepository
}

func NewArrangementService(repo port.ArrangementRepository) *ArrangementService {
	return &ArrangementService{repo: repo}
}

func (s *ArrangementService) Create(ctx context.Context, userID string, req *coa.ArrangementRequest) (*coa.ArrangementResponse, error) {
	a := req.ToArrangement(userID)
	a.ID = uuid.New().String()
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	// Reload to get created_at/updated_at
	created, err := s.repo.GetByID(ctx, a.ID, a.UserID)
	if err != nil || created == nil {
		return a.ToResponse(), nil
	}
	return created.ToResponse(), nil
}

func (s *ArrangementService) GetByID(ctx context.Context, id, userID string) (*coa.ArrangementResponse, error) {
	a, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrArrangementNotFound
	}
	return a.ToResponse(), nil
}

func (s *ArrangementService) ListByUserID(ctx context.Context, userID string) ([]coa.ArrangementResponse, error) {
	list, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]coa.ArrangementResponse, 0, len(list))
	for i := range list {
		out = append(out, *list[i].ToResponse())
	}
	return out, nil
}

func (s *ArrangementService) Update(ctx context.Context, id, userID string, req *coa.ArrangementRequest) (*coa.ArrangementResponse, error) {
	existing, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrArrangementNotFound
	}
	a := req.ToArrangement(existing.UserID)
	a.ID = id
	a.CreatedAt = existing.CreatedAt
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	updated, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if updated != nil {
		return updated.ToResponse(), nil
	}
	return a.ToResponse(), nil
}

func (s *ArrangementService) Delete(ctx context.Context, id, userID string) error {
	existing, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrArrangementNotFound
	}
	return s.repo.Delete(ctx, id)
}
