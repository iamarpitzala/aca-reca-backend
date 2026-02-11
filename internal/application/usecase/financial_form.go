package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type FinancialFormService struct {
	repo       port.FinancialFormRepository
	clinicRepo port.ClinicRepository
}

func NewFinancialFormService(repo port.FinancialFormRepository, clinicRepo port.ClinicRepository) *FinancialFormService {
	return &FinancialFormService{repo: repo, clinicRepo: clinicRepo}
}

func (s *FinancialFormService) CreateFinancialForm(ctx context.Context, req *domain.FinancialFormRequest) error {
	form, err := req.ToRepo()
	if err != nil {
		return err
	}
	if _, err := s.clinicRepo.GetByID(ctx, form.ClinicID); err != nil {
		return errors.New("clinic not found")
	}
	if form.CalculationMethod != "net" && form.CalculationMethod != "gross" {
		return errors.New("calculation method must be 'net' or 'gross'")
	}
	form.ID = uuid.New()
	return s.repo.Create(ctx, form)
}

func (s *FinancialFormService) GetFinancialFormByID(ctx context.Context, id uuid.UUID) (*domain.FinancialFormResponse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FinancialFormService) GetFinancialFormsByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.FinancialFormResponse, error) {
	return s.repo.GetByClinicID(ctx, clinicID)
}

func (s *FinancialFormService) UpdateFinancialForm(ctx context.Context, form *domain.FinancialForm) error {
	existing, err := s.repo.GetByID(ctx, form.ID)
	if err != nil {
		return errors.New("financial form not found")
	}
	if existing.ClinicID != form.ClinicID {
		return errors.New("cannot change clinic for financial form")
	}
	if form.CalculationMethod != "net" && form.CalculationMethod != "gross" {
		return errors.New("calculation method must be 'net' or 'gross'")
	}
	form.UpdatedAt = time.Now()
	return s.repo.Update(ctx, form)
}

func (s *FinancialFormService) DeleteFinancialForm(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errors.New("financial form not found")
	}
	return s.repo.Delete(ctx, id)
}
