package usecase

import (
	"context"
	"errors"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
)

var ErrFinancialYearNotFound = errors.New("financial year not found")

type FinancialYearService struct {
	repo port.FinancialYearRepository
}

func NewFinancialYearService(repo port.FinancialYearRepository) *FinancialYearService {
	return &FinancialYearService{repo: repo}
}

// Create implements [port.FinancialYearRepository].
func (f *FinancialYearService) Create(ctx context.Context, year *clinic.FinancialYear) error {
	return f.repo.Create(ctx, year)
}

// Delete implements [port.FinancialYearRepository].
func (f *FinancialYearService) Delete(ctx context.Context, id int) error {
	return f.repo.Delete(ctx, id)
}

// GetByClinicID implements [port.FinancialYearRepository].
func (f *FinancialYearService) GetByClinicID(ctx context.Context, clinicID string) (*clinic.FinancialYear, error) {
	return f.repo.GetByClinicID(ctx, clinicID)
}

// GetByID implements [port.FinancialYearRepository].
func (f *FinancialYearService) GetByID(ctx context.Context, id int) (*clinic.FinancialYear, error) {
	return f.repo.GetByID(ctx, id)
}

// Update implements [port.FinancialYearRepository].
func (f *FinancialYearService) Update(ctx context.Context, year *clinic.FinancialYear) error {
	return f.repo.Update(ctx, year)
}
