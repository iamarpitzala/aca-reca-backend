package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/utils"
)

type QuarterService struct {
	financialSettingsRepo port.ClinicFinancialSettingsRepository
}

func NewQuarterService(financialSettingsRepo port.ClinicFinancialSettingsRepository) *QuarterService {
	return &QuarterService{
		financialSettingsRepo: financialSettingsRepo,
	}
}

// CalculateQuartersForClinic calculates quarters for a clinic based on its financial settings
// This is the new system-driven approach - quarters are calculated, not stored
func (s *QuarterService) CalculateQuartersForClinic(
	ctx context.Context,
	clinicID uuid.UUID,
	yearsBack int,
	yearsForward int,
) ([]utils.CalculatedQuarter, error) {
	// Get financial settings for the clinic
	settings, err := s.financialSettingsRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		// If settings don't exist, use defaults (July-June, no lock date)
		settings = &domain.ClinicFinancialSettings{
			FinancialYearStart: domain.FinancialYearStartJuly,
			LockDate:          nil,
		}
	}

	quarters := utils.GetQuartersForPeriod(
		settings.FinancialYearStart,
		settings.LockDate,
		yearsBack,
		yearsForward,
	)

	return quarters, nil
}

// GetQuarterForDate finds the quarter that contains a given date for a clinic
func (s *QuarterService) GetQuarterForDate(
	ctx context.Context,
	clinicID uuid.UUID,
	date time.Time,
) (*utils.CalculatedQuarter, error) {
	// Get financial settings for the clinic
	settings, err := s.financialSettingsRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		// If settings don't exist, use defaults
		settings = &domain.ClinicFinancialSettings{
			FinancialYearStart: domain.FinancialYearStartJuly,
			LockDate:          nil,
		}
	}

	quarter := utils.GetQuarterForDate(date, settings.FinancialYearStart, settings.LockDate)
	if quarter == nil {
		return nil, errors.New("quarter not found for date")
	}

	return quarter, nil
}
