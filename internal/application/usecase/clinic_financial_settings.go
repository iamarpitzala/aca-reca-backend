package usecase

import (
	"context"
	"errors"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

var ErrFinancialSettingsNotFound = errors.New("financial settings not found")
var ErrFinancialSettingsLocked = errors.New("cannot modify financial settings: financial year start is locked once transactions exist")

type ClinicFinancialSettingsService struct {
	repo                 port.ClinicFinancialSettingRepository
	clinicRepo           port.ClinicRepository
	financialYearRepo    port.FinancialYearRepository
	financialQuarterRepo port.FinancialQuarterRepository
}

func NewClinicFinancialSettingsService(repo port.ClinicFinancialSettingRepository, clinicRepo port.ClinicRepository, financialYearRepo port.FinancialYearRepository, financialQuarterRepo port.FinancialQuarterRepository) *ClinicFinancialSettingsService {
	return &ClinicFinancialSettingsService{
		repo:                 repo,
		clinicRepo:           clinicRepo,
		financialYearRepo:    financialYearRepo,
		financialQuarterRepo: financialQuarterRepo,
	}
}

// GetByClinicID retrieves financial settings for a clinic, creating defaults if none exist
func (s *ClinicFinancialSettingsService) GetByClinicID(ctx context.Context, clinicID string) (*clinic.ClinicFinancialSetting, error) {
	// Verify clinic exists
	_, err := s.clinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	settings, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		// If not found, return defaults (caller can create if needed)
		if err.Error() == "financial settings not found" {
			return s.getDefaultSettings(ctx, clinicID), nil
		}
		return nil, err
	}
	return settings, nil
}

// CreateOrUpdate creates or updates financial settings for a clinic
func (s *ClinicFinancialSettingsService) CreateOrUpdate(ctx context.Context, clinicID string, req *clinic.ClinicFinancialSettingRequest) (*clinic.ClinicFinancialSetting, error) {
	// Verify clinic exists
	_, err := s.clinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil && err.Error() != "financial settings not found" {
		return nil, err
	}

	if existing == nil {
		// Create new settings
		settings, err := req.ToClinicFinancialSetting(clinicID)
		if err != nil {
			return nil, err
		}
		if err := s.repo.Create(ctx, settings); err != nil {
			return nil, err
		}
		return settings, nil
	}

	// Update existing settings
	if req.FinancialYearID != existing.FinancialYearID {
		existing.FinancialYearID = req.FinancialYearID
	}
	if req.FinancialQuarterID != existing.FinancialQuarterID {
		existing.FinancialQuarterID = req.FinancialQuarterID
	}
	if req.CalculationMethod != "" {
		existing.CalculationMethod = req.CalculationMethod
	}
	existing.GSTRegistered = req.GSTRegistered
	if req.GSTReportingFrequency != "" {
		existing.GSTReportingFrequency = req.GSTReportingFrequency
	}
	if req.DefaultAmountMode != "" {
		existing.DefaultAmountMode = req.DefaultAmountMode
	}
	if req.LockDate != nil {
		existing.LockDate = req.LockDate
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// getDefaultSettings returns default financial settings
func (s *ClinicFinancialSettingsService) getDefaultSettings(ctx context.Context, clinicID string) *clinic.ClinicFinancialSetting {
	financialYear, err := s.financialYearRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil
	}
	financialQuarter, err := s.financialQuarterRepo.GetByFinancialYearID(ctx, financialYear.ID)
	if err != nil {
		return nil
	}

	settings := &clinic.ClinicFinancialSetting{
		ClinicID:              clinicID,
		FinancialYearID:       financialYear.ID,
		FinancialQuarterID:    financialQuarter.ID,
		CalculationMethod:     util.AccountingMethodAccrual,
		GSTRegistered:         true,
		GSTReportingFrequency: util.PeriodQuarterly,
		DefaultAmountMode:     util.Inclusive,
		LockDate:              &financialQuarter.EndDate,
	}

	return settings
}
