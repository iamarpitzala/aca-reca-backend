package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

var ErrFinancialSettingsNotFound = errors.New("financial settings not found")
var ErrFinancialSettingsLocked = errors.New("cannot modify financial settings: financial year start is locked once transactions exist")

type ClinicFinancialSettingsService struct {
	repo port.ClinicFinancialSettingsRepository
	clinicRepo port.ClinicRepository
}

func NewClinicFinancialSettingsService(repo port.ClinicFinancialSettingsRepository, clinicRepo port.ClinicRepository) *ClinicFinancialSettingsService {
	return &ClinicFinancialSettingsService{
		repo:      repo,
		clinicRepo: clinicRepo,
	}
}

// GetByClinicID retrieves financial settings for a clinic, creating defaults if none exist
func (s *ClinicFinancialSettingsService) GetByClinicID(ctx context.Context, clinicID uuid.UUID) (*domain.ClinicFinancialSettings, error) {
	// Verify clinic exists
	_, err := s.clinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	
	settings, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		// If not found, return defaults (caller can create if needed)
		if err.Error() == "financial settings not found" {
			return s.getDefaultSettings(clinicID), nil
		}
		return nil, err
	}
	return settings, nil
}

// CreateOrUpdate creates or updates financial settings for a clinic
func (s *ClinicFinancialSettingsService) CreateOrUpdate(ctx context.Context, clinicID uuid.UUID, req *domain.ClinicFinancialSettingsRequest) (*domain.ClinicFinancialSettings, error) {
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
		settings, err := req.ToClinicFinancialSettings(clinicID)
		if err != nil {
			return nil, err
		}
		if err := s.repo.Create(ctx, settings); err != nil {
			return nil, err
		}
		return settings, nil
	}
	
	// Update existing settings
	// Check if financial_year_start can be changed (should be locked if transactions exist)
	// For now, we'll allow updates but this should be checked against transaction table
	if req.FinancialYearStart != "" && req.FinancialYearStart != existing.FinancialYearStart {
		// TODO: Check if transactions exist - if so, lock financial_year_start
		// For MVP, we'll allow the change
		existing.FinancialYearStart = req.FinancialYearStart
	}
	if req.AccountingMethod != "" {
		existing.AccountingMethod = req.AccountingMethod
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
	if req.GSTDefaults != nil {
		if err := existing.SetGSTDefaultsMap(req.GSTDefaults); err != nil {
			return nil, err
		}
	}
	
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// getDefaultSettings returns default financial settings
func (s *ClinicFinancialSettingsService) getDefaultSettings(clinicID uuid.UUID) *domain.ClinicFinancialSettings {
	now := time.Now()
	// Default lock date: end of last financial year (assuming July-June)
	lockDate := time.Date(now.Year()-1, 6, 30, 0, 0, 0, 0, now.Location())
	if now.Month() >= 7 {
		lockDate = time.Date(now.Year(), 6, 30, 0, 0, 0, 0, now.Location())
	}
	
	gstDefaults := map[string]string{
		"patient_fees": "GST_FREE",
		"service_fees":  "GST_10",
		"lab_fees":      "GST_FREE",
	}
	
	settings := &domain.ClinicFinancialSettings{
		ClinicID:              clinicID,
		FinancialYearStart:    domain.FinancialYearStartJuly,
		AccountingMethod:      domain.AccountingMethodCash,
		GSTRegistered:         true,
		GSTReportingFrequency: domain.GSTReportingFrequencyQuarterly,
		DefaultAmountMode:     domain.DefaultAmountModeGSTInclusive,
		LockDate:              &lockDate,
	}
	settings.SetGSTDefaultsMap(gstDefaults)
	
	return settings
}
