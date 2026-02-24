package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

var ErrFinancialSettingsNotFound = errors.New("financial settings not found")
var ErrFinancialSettingsLocked = errors.New("cannot modify financial settings: financial year start is locked once transactions exist")

const defaultFYLabel = "2025-26"

type ClinicFinancialSettingsService struct {
	repo                    port.ClinicFinancialSettingRepository
	clinicRepo              port.ClinicRepository
	financialYearRepo      port.FinancialYearRepository
	financialQuarterRepo   port.FinancialQuarterRepository
	clinicFinancialYearRepo port.ClinicFinancialYearRepository
	quarterLockRepo        port.ClinicFinancialQuarterLockRepository
}

func NewClinicFinancialSettingsService(
	repo port.ClinicFinancialSettingRepository,
	clinicRepo port.ClinicRepository,
	financialYearRepo port.FinancialYearRepository,
	financialQuarterRepo port.FinancialQuarterRepository,
	clinicFinancialYearRepo port.ClinicFinancialYearRepository,
	quarterLockRepo port.ClinicFinancialQuarterLockRepository,
) *ClinicFinancialSettingsService {
	return &ClinicFinancialSettingsService{
		repo:                    repo,
		clinicRepo:              clinicRepo,
		financialYearRepo:      financialYearRepo,
		financialQuarterRepo:   financialQuarterRepo,
		clinicFinancialYearRepo: clinicFinancialYearRepo,
		quarterLockRepo:        quarterLockRepo,
	}
}

// GetByClinicID retrieves financial settings for the clinic's current financial year.
func (s *ClinicFinancialSettingsService) GetByClinicID(ctx context.Context, clinicID string) (*clinic.ClinicFinancialSetting, error) {
	_, err := s.clinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	settings, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		if err.Error() == "financial settings not found" {
			return s.getDefaultSettings(ctx, clinicID), nil
		}
		return nil, err
	}
	return settings, nil
}

// CreateOrUpdate creates or updates financial settings for a clinic (keyed by clinic_financial_year_id).
func (s *ClinicFinancialSettingsService) CreateOrUpdate(ctx context.Context, clinicID string, req *clinic.ClinicFinancialSettingRequest) (*clinic.ClinicFinancialSetting, error) {
	_, err := s.clinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	existing, _ := s.repo.GetByClinicFinancialYearID(ctx, req.ClinicFinancialYearID)
	if existing == nil {
		settings, err := req.ToClinicFinancialSetting(clinicID)
		if err != nil {
			return nil, err
		}
		if err := s.repo.Create(ctx, settings); err != nil {
			return nil, err
		}
		return settings, nil
	}
	if req.ClinicFinancialYearID != existing.ClinicFinancialYearID {
		existing.ClinicFinancialYearID = req.ClinicFinancialYearID
	}
	if req.CalculationMethod != "" {
		existing.AccountingMethod = req.CalculationMethod
	}
	existing.GSTRegistered = req.GSTRegistered
	existing.GSTReportingFrequency = req.GSTReportingFrequency
	existing.DefaultAmountMode = req.DefaultAmountMode
	if req.LockDate != nil {
		existing.LockDate = req.LockDate
	}
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// CreateDefaultFinancialYearForClinic sets up default 2025-26 for a new clinic: master FY, quarters, clinic link, quarter locks, and default settings.
func (s *ClinicFinancialSettingsService) CreateDefaultFinancialYearForClinic(ctx context.Context, clinicID string) error {
	// 1) Get or create master financial year 2025-26
	fy, err := s.financialYearRepo.GetByFYLabel(ctx, defaultFYLabel)
	if err != nil {
		return err
	}
	if fy == nil {
		startDate, _ := time.Parse("2006-01-02", "2025-07-01")
		endDate, _ := time.Parse("2006-01-02", "2026-06-30")
		fy = &clinic.FinancialYear{
			FYLabel:   defaultFYLabel,
			StartDate: startDate,
			EndDate:   endDate,
			IsActive:  true,
		}
		if err := s.financialYearRepo.Create(ctx, fy); err != nil {
			return err
		}
	}
	// 2) Ensure 4 quarters exist for that master FY
	quarters, err := s.financialQuarterRepo.ListByFinancialYearID(ctx, fy.ID)
	if err != nil {
		return err
	}
	if len(quarters) == 0 {
		for _, q := range []struct{ name string; num int; start, end string }{
			{"Q1", 1, "2025-07-01", "2025-09-30"},
			{"Q2", 2, "2025-10-01", "2025-12-31"},
			{"Q3", 3, "2026-01-01", "2026-03-31"},
			{"Q4", 4, "2026-04-01", "2026-06-30"},
		} {
			if err := s.financialQuarterRepo.CreateQuarter(ctx, fy.ID, q.name, q.num, q.start, q.end); err != nil {
				return err
			}
		}
		quarters, _ = s.financialQuarterRepo.ListByFinancialYearID(ctx, fy.ID)
	}
	// 3) Link clinic to this financial year (tbl_clinic_financial_year)
	cfy := &clinic.ClinicFinancialYear{
		ClinicID:        clinicID,
		FinancialYearID: fy.ID,
		IsCurrent:       true,
		IsClosed:        false,
	}
	cfyID, err := s.clinicFinancialYearRepo.Create(ctx, cfy)
	if err != nil {
		return err
	}
	// 4) Create quarter locks for this clinic financial year
	for _, q := range quarters {
		_ = s.quarterLockRepo.CreateLock(ctx, cfyID, q.ID)
	}
	// 5) Create default settings row for this clinic_financial_year
	now := time.Now()
	settings := &clinic.ClinicFinancialSetting{
		ID:                    uuid.New().String(),
		ClinicID:              clinicID,
		ClinicFinancialYearID: cfyID,
		AccountingMethod:      util.AccountingMethodAccrual,
		GSTRegistered:         true,
		GSTReportingFrequency: util.PeriodQuarterly,
		DefaultAmountMode:     util.Inclusive,
		LockDate:              nil,
		CreatedAt:             now,
		UpdatedAt:             now,
		DeletedAt:             nil,
	}
	return s.repo.Create(ctx, settings)
}

// getDefaultSettings returns default financial settings for the clinic's current financial year (for display when no row exists).
func (s *ClinicFinancialSettingsService) getDefaultSettings(ctx context.Context, clinicID string) *clinic.ClinicFinancialSetting {
	cfy, err := s.clinicFinancialYearRepo.GetCurrentByClinicID(ctx, clinicID)
	if err != nil || cfy == nil {
		return nil
	}
	now := time.Now()
	return &clinic.ClinicFinancialSetting{
		ClinicID:              clinicID,
		ClinicFinancialYearID: cfy.ID,
		AccountingMethod:      util.AccountingMethodAccrual,
		GSTRegistered:         true,
		GSTReportingFrequency: util.PeriodQuarterly,
		DefaultAmountMode:     util.Inclusive,
		LockDate:              nil,
		CreatedAt:             now,
		UpdatedAt:             now,
		DeletedAt:             nil,
	}
}
