package usecase

import (
	"context"
	"errors"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
)

var ErrFinancialYearNotFound = errors.New("financial year not found")

type FinancialYearService struct {
	repo    port.FinancialYearRepository
	cfyRepo port.ClinicFinancialYearRepository
}

func NewFinancialYearService(repo port.FinancialYearRepository, cfyRepo port.ClinicFinancialYearRepository) *FinancialYearService {
	return &FinancialYearService{repo: repo, cfyRepo: cfyRepo}
}

// Create creates a master financial year (global).
func (f *FinancialYearService) Create(ctx context.Context, year *clinic.FinancialYear) error {
	return f.repo.Create(ctx, year)
}

// CreateLink links a clinic to a master financial year (creates tbl_clinic_financial_year row).
func (f *FinancialYearService) CreateLink(ctx context.Context, clinicID string, financialYearID int) (*clinic.ClinicFinancialYear, error) {
	fy, err := f.repo.GetByID(ctx, financialYearID)
	if err != nil || fy == nil {
		return nil, ErrFinancialYearNotFound
	}
	cfy := &clinic.ClinicFinancialYear{
		ClinicID:        clinicID,
		FinancialYearID: financialYearID,
		IsCurrent:       true,
		IsClosed:        false,
	}
	id, err := f.cfyRepo.Create(ctx, cfy)
	if err != nil {
		return nil, err
	}
	cfy.ID = id
	return cfy, nil
}

// Delete deletes a master financial year (soft delete).
func (f *FinancialYearService) Delete(ctx context.Context, id int) error {
	return f.repo.Delete(ctx, id)
}

// GetByClinicID returns the master financial year for the clinic's current financial year link.
func (f *FinancialYearService) GetByClinicID(ctx context.Context, clinicID string) (*clinic.FinancialYear, error) {
	cfy, err := f.cfyRepo.GetCurrentByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	return f.repo.GetByID(ctx, cfy.FinancialYearID)
}

// GetByID returns a master financial year by id.
func (f *FinancialYearService) GetByID(ctx context.Context, id int) (*clinic.FinancialYear, error) {
	return f.repo.GetByID(ctx, id)
}

// List returns all master financial years.
func (f *FinancialYearService) List(ctx context.Context) ([]clinic.FinancialYear, error) {
	return f.repo.List(ctx)
}

// Update updates a master financial year.
func (f *FinancialYearService) Update(ctx context.Context, year *clinic.FinancialYear) error {
	return f.repo.Update(ctx, year)
}
