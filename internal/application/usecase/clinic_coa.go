package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

var (
	ErrClinicCOAExists   = errors.New("this AOC is already assigned to the clinic")
	ErrClinicCOANotFound = errors.New("clinic AOC association not found")
	ErrClinicNotFound    = errors.New("clinic not found")
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

func (s *ClinicCOAService) GetClinicAOCs(ctx context.Context, clinicID uuid.UUID) ([]domain.ClinicCOAWithDetails, error) {
	if _, err := s.clinicRepo.GetByID(ctx, clinicID); err != nil {
		return nil, ErrClinicNotFound
	}
	return s.repo.ListByClinicIDWithDetails(ctx, clinicID)
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

// CreateCOAForClinic creates a new chart-of-accounts entry and assigns it to the clinic (tenant-level COA).
func (s *ClinicCOAService) CreateCOAForClinic(ctx context.Context, clinicID uuid.UUID, req *domain.CreateCOAForClinicRequest) (*domain.ClinicCOAResponse, error) {
	if _, err := s.clinicRepo.GetByID(ctx, clinicID); err != nil {
		return nil, errors.New("clinic not found")
	}
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errors.New("code and name are required")
	}
	if _, err := s.aocRepo.GetAccountTypeByID(ctx, req.AccountTypeID); err != nil {
		return nil, errors.New("account type not found")
	}
	if _, err := s.aocRepo.GetAccountTaxByID(ctx, req.AccountTaxID); err != nil {
		return nil, errors.New("account tax not found")
	}
	existing, _ := s.aocRepo.GetByCode(ctx, code)
	if existing != nil {
		return nil, errors.New("an account with this code already exists")
	}
	now := time.Now()
	aoc := &domain.AOC{
		ID:            uuid.New(),
		AccountTypeID: req.AccountTypeID,
		AccountTaxID:  req.AccountTaxID,
		Code:          code,
		Name:          name,
		Description:   req.Description,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
	}
	if err := s.aocRepo.Create(ctx, aoc); err != nil {
		return nil, err
	}
	return s.AddClinicAOC(ctx, clinicID, aoc.ID)
}
