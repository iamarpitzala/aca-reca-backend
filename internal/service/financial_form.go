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

type FinancialFormService struct {
	db *sqlx.DB
}

func NewFinancialFormService(db *sqlx.DB) *FinancialFormService {
	return &FinancialFormService{
		db: db,
	}
}

func (ffs *FinancialFormService) CreateFinancialForm(ctx context.Context, form *domain.FinancialForm) error {
	// Verify clinic exists
	_, err := repository.GetClinicByID(ctx, ffs.db, form.ClinicID)
	if err != nil {
		return errors.New("clinic not found")
	}

	// Validate calculation method
	if form.CalculationMethod != "net" && form.CalculationMethod != "gross" {
		return errors.New("calculation method must be 'net' or 'gross'")
	}

	form.ID = uuid.New()
	form.CreatedAt = time.Now()
	form.UpdatedAt = time.Now()

	return repository.CreateFinancialForm(ctx, ffs.db, form)
}

func (ffs *FinancialFormService) GetFinancialFormByID(ctx context.Context, id uuid.UUID) (*domain.FinancialForm, error) {
	return repository.GetFinancialFormByID(ctx, ffs.db, id)
}

func (ffs *FinancialFormService) GetFinancialFormsByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.FinancialForm, error) {
	return repository.GetFinancialFormsByClinicID(ctx, ffs.db, clinicID)
}

func (ffs *FinancialFormService) UpdateFinancialForm(ctx context.Context, form *domain.FinancialForm) error {
	// Verify form exists
	existing, err := repository.GetFinancialFormByID(ctx, ffs.db, form.ID)
	if err != nil {
		return errors.New("financial form not found")
	}

	// Verify clinic ownership hasn't changed
	if existing.ClinicID != form.ClinicID {
		return errors.New("cannot change clinic for financial form")
	}

	// Validate calculation method
	if form.CalculationMethod != "net" && form.CalculationMethod != "gross" {
		return errors.New("calculation method must be 'net' or 'gross'")
	}

	form.UpdatedAt = time.Now()
	return repository.UpdateFinancialForm(ctx, ffs.db, form)
}

func (ffs *FinancialFormService) DeleteFinancialForm(ctx context.Context, id uuid.UUID) error {
	_, err := repository.GetFinancialFormByID(ctx, ffs.db, id)
	if err != nil {
		return errors.New("financial form not found")
	}

	return repository.DeleteFinancialForm(ctx, ffs.db, id)
}
