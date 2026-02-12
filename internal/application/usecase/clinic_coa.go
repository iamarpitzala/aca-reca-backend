package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

var (
	ErrClinicCOAExists   = errors.New("this AOC is already assigned to the clinic")
	ErrClinicCOANotFound = errors.New("clinic AOC association not found")
)

type ClinicCOAService struct {
	repo     port.ClinicCOARepository
	clinicRepo port.ClinicRepository
	aocRepo   port.AOCRepository
}

func NewClinicCOAService(repo port.ClinicCOARepository, clinicRepo port.ClinicRepository, aocRepo port.AOCRepository) *ClinicCOAService {
	return &ClinicCOAService{
		repo:       repo,
		clinicRepo: clinicRepo,
		aocRepo:    aocRepo,
	}
}

func (s *ClinicCOAService) AddClinicAOC(ctx context.Context, clinicID uuid.UUID, coaID uuid.UUID) (*domain.ClinicCOAResponse, error) {
	if _, err := s.clinicRepo.GetByID(ctx, clinicID); err != nil {
		return nil, errors.New("clinic not found")
	}
	if aoc, err := s.aocRepo.GetByID(ctx, coaID); err != nil || aoc == nil {
		return nil, errors.New("AOC not found")
	}
	exists, err := s.repo.Exists(ctx, clinicID, coaID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrClinicCOAExists
	}
	now := time.Now()
	cc := &domain.ClinicCOA{
		ID:        uuid.New(),
		ClinicID:  clinicID,
		COAID:     coaID,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}
	if err := s.repo.Create(ctx, cc); err != nil {
		return nil, err
	}
	return cc.ToResponse(), nil
}

func (s *ClinicCOAService) GetClinicAOCs(ctx context.Context, clinicID uuid.UUID) ([]domain.ClinicCOAResponse, error) {
	list, err := s.repo.ListByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.ClinicCOAResponse, 0, len(list))
	for i := range list {
		out = append(out, *list[i].ToResponse())
	}
	return out, nil
}

func (s *ClinicCOAService) GetClinicAOCByID(ctx context.Context, id uuid.UUID) (*domain.ClinicCOAResponse, error) {
	cc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cc == nil {
		return nil, ErrClinicCOANotFound
	}
	return cc.ToResponse(), nil
}

func (s *ClinicCOAService) RemoveClinicAOC(ctx context.Context, id uuid.UUID) error {
	cc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if cc == nil {
		return ErrClinicCOANotFound
	}
	return s.repo.Delete(ctx, id)
}
